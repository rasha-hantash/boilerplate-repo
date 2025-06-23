package lib

import (
	"time"

	"github.com/google/uuid"
	todov1 "github.com/rasha-hantash/boilerplate-repo/platform/api/gen/proto/todo/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Priority represents the priority level of a todo
type Priority int

const (
	PriorityUnspecified Priority = iota
	PriorityLow
	PriorityMedium
	PriorityHigh
)

// Todo represents a todo item in the database
type Todo struct {
	ID          string     `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Completed   bool       `json:"completed" db:"completed"`
	Priority    Priority   `json:"priority" db:"priority"`
	Category    string     `json:"category" db:"category"`
	DueDate     *time.Time `json:"due_date" db:"due_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// NewTodo creates a new Todo with default values
func NewTodo(title, description string, priority Priority, category string, dueDate *time.Time) *Todo {
	now := time.Now()
	return &Todo{
		ID:          uuid.New().String(),
		Title:       title,
		Description: description,
		Completed:   false,
		Priority:    priority,
		Category:    category,
		DueDate:     dueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// ToProto converts a Todo to its protobuf representation
func (t *Todo) ToProto() *todov1.Todo {
	todo := &todov1.Todo{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		Priority:    priorityToProto(t.Priority),
		Category:    t.Category,
		CreatedAt:   timestamppb.New(t.CreatedAt),
		UpdatedAt:   timestamppb.New(t.UpdatedAt),
	}

	if t.DueDate != nil {
		todo.DueDate = timestamppb.New(*t.DueDate)
	}

	return todo
}

// FromProto creates a Todo from its protobuf representation
func FromProto(todo *todov1.Todo) *Todo {
	t := &Todo{
		ID:          todo.Id,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		Priority:    priorityFromProto(todo.Priority),
		Category:    todo.Category,
	}

	if todo.CreatedAt != nil {
		t.CreatedAt = todo.CreatedAt.AsTime()
	}

	if todo.UpdatedAt != nil {
		t.UpdatedAt = todo.UpdatedAt.AsTime()
	}

	if todo.DueDate != nil {
		dueDate := todo.DueDate.AsTime()
		t.DueDate = &dueDate
	}

	return t
}

// priorityToProto converts Priority to protobuf Priority
func priorityToProto(p Priority) todov1.Priority {
	switch p {
	case PriorityLow:
		return todov1.Priority_PRIORITY_LOW
	case PriorityMedium:
		return todov1.Priority_PRIORITY_MEDIUM
	case PriorityHigh:
		return todov1.Priority_PRIORITY_HIGH
	default:
		return todov1.Priority_PRIORITY_UNSPECIFIED
	}
}

// priorityFromProto converts protobuf Priority to Priority
func priorityFromProto(p todov1.Priority) Priority {
	switch p {
	case todov1.Priority_PRIORITY_LOW:
		return PriorityLow
	case todov1.Priority_PRIORITY_MEDIUM:
		return PriorityMedium
	case todov1.Priority_PRIORITY_HIGH:
		return PriorityHigh
	default:
		return PriorityUnspecified
	}
}
