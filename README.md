# Todo App - Full Stack Boilerplate

A modern full-stack TODO application built with:

## Backend
- **Go** with ConnectRPC (gRPC-web compatible)
- **PostgreSQL** database
- **ConnectRPC** for type-safe API communication

## Frontend
- **React** with TypeScript
- **TanStack Query** for server state management
- **ConnectRPC** web client
- **TailwindCSS** with **shadcn/ui** components
- **Vite** for fast development

## Features

- ✅ Create, Read, Update, Delete (CRUD) operations
- ✅ Mark todos as complete/incomplete
- ✅ Due dates, priorities, categories
- ✅ Filtering by status, priority, and category
- ✅ Modern, responsive UI
- ✅ Type-safe API communication
- ✅ Docker containerization

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Go 1.24+ (for local development)
- Node.js 18+ (for local development)

### Using Docker (Recommended)

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd boilerplate-repo
   ```

2. **Start the application**
   ```bash
   docker-compose up --build
   ```

3. **Access the application**
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080
   - PostgreSQL: localhost:5432

### Local Development

#### Backend Setup

1. **Navigate to the backend directory**
   ```bash
   cd platform/api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Generate protobuf code**
   ```bash
   # Install buf if not already installed
   go install github.com/bufbuild/buf/cmd/buf@latest
   
   # Generate code
   buf generate
   ```

4. **Set up PostgreSQL**
   ```bash
   # Using Docker
   docker run --name todo-postgres -e POSTGRES_DB=todos -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -p 5432:5432 -d postgres:15-alpine
   ```

5. **Run the backend**
   ```bash
   go run main.go
   ```

#### Frontend Setup

1. **Navigate to the frontend directory**
   ```bash
   cd frontend
   ```

2. **Install dependencies**
   ```bash
   npm install
   ```

3. **Generate protobuf code**
   ```bash
   # From the root directory
   cd platform/api
   buf generate
   ```

4. **Run the frontend**
   ```bash
   npm run dev
   ```

## Project Structure

```
boilerplate-repo/
├── platform/
│   ├── api/                    # Go backend
│   │   ├── config/            # Database configuration
│   │   ├── handler/           # ConnectRPC handlers
│   │   ├── lib/               # Data models and utilities
│   │   ├── services/          # Database repository
│   │   ├── proto/             # Protobuf definitions
│   │   ├── gen/               # Generated Go code
│   │   ├── main.go            # Application entry point
│   │   └── go.mod             # Go dependencies
│   └── sql/                   # Database migrations
├── frontend/                   # React frontend
│   ├── src/
│   │   ├── components/        # React components
│   │   │   ├── ui/           # shadcn/ui components
│   │   │   ├── TodoList.tsx  # Todo list component
│   │   │   └── CreateTodoForm.tsx
│   │   ├── hooks/            # React Query hooks
│   │   ├── lib/              # Utilities and API client
│   │   ├── gen/              # Generated TypeScript code
│   │   ├── App.tsx           # Main app component
│   │   └── main.tsx          # Entry point
│   ├── package.json          # Node.js dependencies
│   └── vite.config.ts        # Vite configuration
├── docker/                    # Docker configuration
│   ├── Dockerfile.backend    # Backend Dockerfile
│   └── Dockerfile.frontend   # Frontend Dockerfile
├── docker-compose.yml        # Docker Compose setup
└── README.md                 # This file
```

## API Endpoints

The backend provides the following ConnectRPC endpoints:

- `CreateTodo` - Create a new todo
- `GetTodo` - Retrieve a single todo by ID
- `ListTodos` - List todos with optional filtering
- `UpdateTodo` - Update an existing todo
- `DeleteTodo` - Delete a todo by ID

## Database Schema

```sql
CREATE TABLE todos (
    id VARCHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    priority INTEGER DEFAULT 0,
    category VARCHAR(100),
    due_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);
```

## Environment Variables

### Backend
- `DB_HOST` - PostgreSQL host (default: localhost)
- `DB_PORT` - PostgreSQL port (default: 5432)
- `DB_USER` - PostgreSQL user (default: postgres)
- `DB_PASSWORD` - PostgreSQL password (default: password)
- `DB_NAME` - Database name (default: todos)
- `DB_SSLMODE` - SSL mode (default: disable)

## Development

### Adding New Features

1. **Update Protobuf definitions** in `platform/api/proto/todo/v1/todo.proto`
2. **Generate code** with `buf generate`
3. **Implement backend handlers** in `platform/api/handler/`
4. **Update frontend hooks** in `frontend/src/hooks/`
5. **Add UI components** as needed

### Code Generation

The project uses [Buf](https://buf.build/) for protobuf code generation:

```bash
# Generate Go code
cd platform/api
buf generate

# Generate TypeScript code (also generates Go code)
buf generate
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

This project is licensed under the MIT License.
