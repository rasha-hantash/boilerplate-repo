package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/rasha-hantash/boilerplate-repo/platform/api/lib"
)

// Repository handles database operations for todos
type Repository struct {
	db *sql.DB
}

// NewRepository creates a new repository instance
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// CreateTodo inserts a new todo into the database
func (r *Repository) CreateTodo(todo *lib.Todo) error {
	query := `
		INSERT INTO todos (id, title, description, completed, priority, category, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.Priority,
		todo.Category,
		todo.DueDate,
		todo.CreatedAt,
		todo.UpdatedAt,
	)

	return err
}

// GetTodo retrieves a todo by ID
func (r *Repository) GetTodo(id string) (*lib.Todo, error) {
	query := `
		SELECT id, title, description, completed, priority, category, due_date, created_at, updated_at
		FROM todos WHERE id = $1
	`

	var todo lib.Todo
	var dueDate sql.NullTime

	err := r.db.QueryRow(query, id).Scan(
		&todo.ID,
		&todo.Title,
		&todo.Description,
		&todo.Completed,
		&todo.Priority,
		&todo.Category,
		&dueDate,
		&todo.CreatedAt,
		&todo.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("todo not found: %s", id)
		}
		return nil, err
	}

	if dueDate.Valid {
		todo.DueDate = &dueDate.Time
	}

	return &todo, nil
}

// ListTodos retrieves todos with optional filtering
func (r *Repository) ListTodos(completed *bool, priority *lib.Priority, category string, limit int) ([]*lib.Todo, error) {
	query := `
		SELECT id, title, description, completed, priority, category, due_date, created_at, updated_at
		FROM todos
	`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if completed != nil {
		conditions = append(conditions, fmt.Sprintf("completed = $%d", argIndex))
		args = append(args, *completed)
		argIndex++
	}

	if priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, *priority)
		argIndex++
	}

	if category != "" {
		conditions = append(conditions, fmt.Sprintf("category = $%d", argIndex))
		args = append(args, category)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*lib.Todo
	for rows.Next() {
		var todo lib.Todo
		var dueDate sql.NullTime

		err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Completed,
			&todo.Priority,
			&todo.Category,
			&dueDate,
			&todo.CreatedAt,
			&todo.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			todo.DueDate = &dueDate.Time
		}

		todos = append(todos, &todo)
	}

	return todos, nil
}

// UpdateTodo updates an existing todo
func (r *Repository) UpdateTodo(todo *lib.Todo) error {
	query := `
		UPDATE todos 
		SET title = $2, description = $3, completed = $4, priority = $5, category = $6, due_date = $7, updated_at = $8
		WHERE id = $1
	`

	result, err := r.db.Exec(query,
		todo.ID,
		todo.Title,
		todo.Description,
		todo.Completed,
		todo.Priority,
		todo.Category,
		todo.DueDate,
		time.Now(),
	)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo not found: %s", todo.ID)
	}

	return nil
}

// DeleteTodo removes a todo by ID
func (r *Repository) DeleteTodo(id string) error {
	query := "DELETE FROM todos WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("todo not found: %s", id)
	}

	return nil
}

// InitDB creates the todos table if it doesn't exist
func (r *Repository) InitDB() error {
	query := `
		CREATE TABLE IF NOT EXISTS todos (
			id VARCHAR(36) PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			description TEXT,
			completed BOOLEAN DEFAULT FALSE,
			priority INTEGER DEFAULT 0,
			category VARCHAR(100),
			due_date TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		)
	`

	_, err := r.db.Exec(query)
	return err
}
