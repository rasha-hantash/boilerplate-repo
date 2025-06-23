package handler

import (
	"context"
	"connectrpc.com/connect"
	todov1 "github.com/rasha-hantash/boilerplate-repo/platform/api/gen/proto/todo/v1"
	// "github.com/rasha-hantash/boilerplate-repo/platform/api/lib"
	services "github.com/rasha-hantash/boilerplate-repo/platform/api/services/todo"
)

// TodoHandler implements the TodoService
type TodoHandler struct {
	repo *services.Repository
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(repo *services.Repository) *TodoHandler {
	return &TodoHandler{repo: repo}
}

// CreateTodo creates a new todo
func (h *TodoHandler) CreateTodo(ctx context.Context, req *connect.Request[todov1.CreateTodoRequest]) (*connect.Response[todov1.CreateTodoResponse], error) {
	panic("unimplemented")
}

// GetTodo retrieves a todo by ID
func (h *TodoHandler) GetTodo(ctx context.Context, req *connect.Request[todov1.GetTodoRequest]) (*connect.Response[todov1.GetTodoResponse], error) {
	panic("unimplemented")
}

// ListTodos retrieves a list of todos with optional filtering
func (h *TodoHandler) ListTodos(ctx context.Context, req *connect.Request[todov1.ListTodosRequest]) (*connect.Response[todov1.ListTodosResponse], error) {
	panic("unimplemented")
}

// UpdateTodo updates an existing todo
func (h *TodoHandler) UpdateTodo(ctx context.Context, req *connect.Request[todov1.UpdateTodoRequest]) (*connect.Response[todov1.UpdateTodoResponse], error) {
	panic("unimplemented")
}

// DeleteTodo deletes a todo by ID
func (h *TodoHandler) DeleteTodo(ctx context.Context, req *connect.Request[todov1.DeleteTodoRequest]) (*connect.Response[todov1.DeleteTodoResponse], error) {
	panic("unimplemented")
}
