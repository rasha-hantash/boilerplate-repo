import { createFileRoute } from '@tanstack/react-router'
import { useState } from 'react'
import { useQuery, useMutation } from '@connectrpc/connect-query'
import { useQueryClient } from '@tanstack/react-query'
import { listTodos, createTodo, deleteTodo } from '~/gen/proto/todo/v1/todo-TodoService_connectquery'
import { Priority } from '~/gen/proto/todo/v1/todo_pb'

export const Route = createFileRoute('/todos')({
  component: TodosComponent,
})

function TodosComponent() {
  const [newTodoTitle, setNewTodoTitle] = useState('')
  const [newTodoDescription, setNewTodoDescription] = useState('')
  const queryClient = useQueryClient()

  // Fetch todos using Connect-Query
  const { data: todosResponse, isLoading, error } = useQuery(listTodos, {
    pageSize: 20,
    completed: false,
  })

  // Create todo mutation
  const createTodoMutation = useMutation(createTodo, {
    onSuccess: () => {
      // Invalidate and refetch todos after creating a new one
      queryClient.invalidateQueries({ 
        queryKey: ['connect-query']
      })
      setNewTodoTitle('')
      setNewTodoDescription('')
    },
  })

  // Delete todo mutation
  const deleteTodoMutation = useMutation(deleteTodo, {
    onSuccess: () => {
      // Invalidate and refetch todos after deleting
      queryClient.invalidateQueries({ 
        queryKey: ['connect-query']
      })
    },
  })

  const handleCreateTodo = (e: React.FormEvent) => {
    e.preventDefault()
    if (newTodoTitle.trim()) {
      createTodoMutation.mutate({
        title: newTodoTitle,
        description: newTodoDescription,
        priority: Priority.MEDIUM,
        category: 'general',
      })
    }
  }

  const handleDeleteTodo = (todoId: string) => {
    deleteTodoMutation.mutate({ id: todoId })
  }

  if (isLoading) {
    return (
      <div className="p-6">
        <h1 className="text-2xl font-bold mb-4">Todos</h1>
        <div>Loading todos...</div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="p-6">
        <h1 className="text-2xl font-bold mb-4">Todos</h1>
        <div className="text-red-500">
          Error loading todos: {error.message}
        </div>
      </div>
    )
  }

  return (
    <div className="p-6 max-w-4xl mx-auto">
      <h1 className="text-3xl font-bold mb-6">Todos</h1>
      
      {/* Create Todo Form */}
      <form onSubmit={handleCreateTodo} className="mb-8 p-4 border rounded-lg bg-gray-50">
        <h2 className="text-xl font-semibold mb-4">Create New Todo</h2>
        <div className="mb-4">
          <label htmlFor="title" className="block text-sm font-medium mb-2">
            Title
          </label>
          <input
            id="title"
            type="text"
            value={newTodoTitle}
            onChange={(e) => setNewTodoTitle(e.target.value)}
            className="w-full p-2 border rounded focus:ring-2 focus:ring-blue-500"
            placeholder="Enter todo title"
            required
          />
        </div>
        <div className="mb-4">
          <label htmlFor="description" className="block text-sm font-medium mb-2">
            Description
          </label>
          <textarea
            id="description"
            value={newTodoDescription}
            onChange={(e) => setNewTodoDescription(e.target.value)}
            className="w-full p-2 border rounded focus:ring-2 focus:ring-blue-500"
            placeholder="Enter todo description"
            rows={3}
          />
        </div>
        <button
          type="submit"
          disabled={createTodoMutation.isPending}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600 disabled:opacity-50"
        >
          {createTodoMutation.isPending ? 'Creating...' : 'Create Todo'}
        </button>
      </form>

      {/* Todos List */}
      <div className="space-y-4">
        <h2 className="text-xl font-semibold">Your Todos ({todosResponse?.todos?.length || 0})</h2>
        
        {!todosResponse?.todos || todosResponse.todos.length === 0 ? (
          <div className="text-gray-500 text-center py-8">
            No todos found. Create your first todo above!
          </div>
        ) : (
          todosResponse.todos.map((todo) => (
            <div
              key={todo.id}
              className="p-4 border rounded-lg bg-white shadow-sm hover:shadow-md transition-shadow"
            >
              <div className="flex justify-between items-start">
                <div className="flex-1">
                  <h3 className="font-semibold text-lg">{todo.title}</h3>
                  {todo.description && (
                    <p className="text-gray-600 mt-1">{todo.description}</p>
                  )}
                  <div className="flex gap-4 mt-2 text-sm text-gray-500">
                    <span>Priority: {Priority[todo.priority]}</span>
                    {todo.category && <span>Category: {todo.category}</span>}
                    <span>Status: {todo.completed ? 'Completed' : 'Pending'}</span>
                  </div>
                </div>
                <button
                  onClick={() => handleDeleteTodo(todo.id)}
                  disabled={deleteTodoMutation.isPending}
                  className="ml-4 px-3 py-1 bg-red-500 text-white text-sm rounded hover:bg-red-600 disabled:opacity-50"
                >
                  {deleteTodoMutation.isPending ? 'Deleting...' : 'Delete'}
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  )
} 