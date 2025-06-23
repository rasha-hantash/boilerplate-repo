package handler

import (
	"context"
	"time"

	"connectrpc.com/connect"
	todov1 "github.com/rasha-hantash/boilerplate-repo/platform/api/gen/proto/todo/v1"
	"github.com/rasha-hantash/boilerplate-repo/platform/api/lib"
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
	reqData := req.Msg

	var dueDate *time.Time
	if reqData.DueDate != nil {
		t := reqData.DueDate.AsTime()
		dueDate = &t
	}

	todo := lib.NewTodo(
		reqData.Title,
		reqData.Description,
		lib.Priority(reqData.Priority),
		reqData.Category,
		dueDate,
	)

	err := h.repo.CreateTodo(todo)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	response := &todov1.CreateTodoResponse{
		Todo: todo.ToProto(),
	}

	return connect.NewResponse(response), nil
}

// GetTodo retrieves a todo by ID
func (h *TodoHandler) GetTodo(ctx context.Context, req *connect.Request[todov1.GetTodoRequest]) (*connect.Response[todov1.GetTodoResponse], error) {
	todo, err := h.repo.GetTodo(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	response := &todov1.GetTodoResponse{
		Todo: todo.ToProto(),
	}

	return connect.NewResponse(response), nil
}

// ListTodos retrieves a list of todos with optional filtering
func (h *TodoHandler) ListTodos(ctx context.Context, req *connect.Request[todov1.ListTodosRequest]) (*connect.Response[todov1.ListTodosResponse], error) {
	reqData := req.Msg

	// Handle optional filters
	var completed *bool
	var priority *lib.Priority

	// Check if completed filter is set (protobuf uses zero values, so we need to check if it's explicitly set)
	// For simplicity, we'll assume if the field is present in the request, it should be used
	// In a real implementation, you might want to use field masks or optional wrappers

	if reqData.Priority != todov1.Priority_PRIORITY_UNSPECIFIED {
		p := lib.Priority(reqData.Priority)
		priority = &p
	}

	limit := int(reqData.PageSize)
	if limit <= 0 {
		limit = 50 // Default limit
	}

	todos, err := h.repo.ListTodos(completed, priority, reqData.Category, limit)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	protoTodos := make([]*todov1.Todo, len(todos))
	for i, todo := range todos {
		protoTodos[i] = todo.ToProto()
	}

	response := &todov1.ListTodosResponse{
		Todos: protoTodos,
		// NextPageToken can be implemented later for pagination
		NextPageToken: "",
	}

	return connect.NewResponse(response), nil
}

// UpdateTodo updates an existing todo
func (h *TodoHandler) UpdateTodo(ctx context.Context, req *connect.Request[todov1.UpdateTodoRequest]) (*connect.Response[todov1.UpdateTodoResponse], error) {
	reqData := req.Msg

	// First get the existing todo
	existingTodo, err := h.repo.GetTodo(reqData.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	// Update fields
	existingTodo.Title = reqData.Title
	existingTodo.Description = reqData.Description
	existingTodo.Completed = reqData.Completed
	existingTodo.Priority = lib.Priority(reqData.Priority)
	existingTodo.Category = reqData.Category

	if reqData.DueDate != nil {
		t := reqData.DueDate.AsTime()
		existingTodo.DueDate = &t
	} else {
		existingTodo.DueDate = nil
	}

	err = h.repo.UpdateTodo(existingTodo)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Get the updated todo to return
	updatedTodo, err := h.repo.GetTodo(reqData.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	response := &todov1.UpdateTodoResponse{
		Todo: updatedTodo.ToProto(),
	}

	return connect.NewResponse(response), nil
}

// DeleteTodo deletes a todo by ID
func (h *TodoHandler) DeleteTodo(ctx context.Context, req *connect.Request[todov1.DeleteTodoRequest]) (*connect.Response[todov1.DeleteTodoResponse], error) {
	err := h.repo.DeleteTodo(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	response := &todov1.DeleteTodoResponse{}
	return connect.NewResponse(response), nil
}
