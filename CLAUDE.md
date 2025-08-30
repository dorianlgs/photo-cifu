# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Architecture

PhotoCifu is a photo gallery application built with Go + PocketBase backend and SvelteKit frontend. It follows a clean architecture pattern with dependency injection and proper separation of concerns.

**Architecture Overview:**
```
├── main.go                 # Application entry point with DI container setup
├── pkg/                    # Core application packages
│   ├── config/            # Environment-based configuration management
│   ├── container/         # Dependency injection container & service interfaces
│   ├── errors/            # Centralized error handling with structured responses
│   ├── handlers/          # Clean HTTP handlers with validation
│   └── validation/        # Request validation with type safety
├── workflow/              # Workflow definitions and activities
├── tools/                 # Utility functions and helpers
├── ui/                    # SvelteKit frontend with TypeScript
├── pb_data/               # PocketBase database and file storage
└── pb_migrations/         # Database migration files
```

**Core Components:**

### Backend (Go)
- **Container**: Dependency injection system in `pkg/container/` managing all services
- **Handlers**: HTTP request handlers in `pkg/handlers/` with proper validation
- **Services**: Business logic with interfaces for testability
- **Workflows**: go-workflows for async image processing
- **Configuration**: Environment-based config in `pkg/config/`
- **Error Handling**: Structured errors in `pkg/errors/` with proper HTTP status codes

### Frontend (SvelteKit + TypeScript)
- **Static Site Generation**: SvelteKit with static adapter for optimal performance
- **Configuration**: Environment variable support in `ui/src/config.ts`
- **PocketBase Integration**: Enhanced client in `ui/src/lib/pocketbase.ts`
- **Components**: Reusable UI components in `ui/src/lib/components/`
- **Stores**: Svelte stores for state management in `ui/src/lib/stores/`

### Database & Storage
- **PocketBase**: SQLite database with collections for users, galleries, images
- **Workflow DB**: Separate SQLite database for go-workflows state
- **File Storage**: Automatic thumbnail generation and secure file handling
- **Migrations**: Auto-generated JavaScript migrations in `pb_migrations/`

## Development Workflow

### Prerequisites
- Go 1.24+
- Node.js/Yarn for frontend development
- SQLite (included with PocketBase)

### Initial Setup
```bash
# Install Go dependencies
go mod tidy

# Install frontend dependencies
cd ui/ && yarn install && cd ..
```

### Development Commands

**Go Backend:**
```bash
# Run development server with hot reload
go run . serve --dev

# Run with debugging enabled
go run . serve --dev --debug

# Update all dependencies
go get -u -t ./...
go mod edit -require=modernc.org/libc@v1.66.3  # Pin to avoid version warnings
go clean -modcache  # Clear cache if version warnings persist
go mod tidy

# Build for production
go generate ./...           # Embed frontend assets
go build -ldflags "-s -w"   # Build optimized binary
```

**Frontend (SvelteKit + TypeScript):**
```bash
cd ui/

# Development server with hot reload
yarn run dev

# Build static site for production
yarn run build

# Preview production build
yarn run preview

# Code quality and type checking
yarn run lint          # ESLint linting
yarn run format        # Prettier formatting
yarn run format_check  # Check formatting only
yarn run check         # SvelteKit/TypeScript checking

# Testing
yarn run test          # Run tests in watch mode
yarn run test_run      # Run tests once
```

### Debugging

**Backend Debugging:**
- PocketBase admin UI available at `http://localhost:8090/admin`
- Database files in `pb_data/` directory
- Structured logging via PocketBase's built-in system
- Workflow state stored in `pb_data/workflow.db`
- Error handling with detailed stack traces in development mode

**Frontend Debugging:**
- Browser dev tools with TypeScript source maps
- SvelteKit dev mode with detailed error pages
- Network tab for API request/response inspection
- Enhanced PocketBase client with error handling utilities

## Code Quality Standards

### Go Code Style
- Follow standard Go formatting (`gofmt`)
- Use meaningful variable names and interfaces
- Handle errors explicitly with structured error types
- Add comments for exported functions and interfaces
- Keep functions focused and small (single responsibility)
- **NEW**: Use dependency injection via `pkg/container/`
- **NEW**: Validate all inputs using `pkg/validation/`
- **NEW**: Use structured error handling via `pkg/errors/`
- **NEW**: Follow clean architecture patterns

### TypeScript/Svelte Standards
- Use TypeScript strict mode (enabled in `tsconfig.json`)
- Follow ESLint configuration with consistent formatting
- Use Prettier for code formatting
- Prefer composition over inheritance
- Use proper TypeScript types, avoid `any`
- **NEW**: Use structured configuration from `config.ts`
- **NEW**: Implement proper error handling for API calls
- **NEW**: Use environment variables for configuration

### Testing Strategy
- **Frontend**: Vitest for unit testing (`ui/src/**/*.{test,spec}.{js,ts}`)
- **Backend**: Go testing framework (when implemented)
- Tests automatically discovered by respective frameworks
- Use `globals: true` configuration for Vitest test functions
- Focus on testing business logic, services, and complex components
- **NEW**: Test service interfaces for better coverage
- **NEW**: Mock dependencies using container interfaces

## Database Management

### PocketBase Collections
- **users**: Authentication and profile data
- **galleries**: Photo gallery metadata and settings
- **images**: Individual image records with file references
- **messages**: System messaging/notifications

### Migrations
- Auto-generated in `pb_migrations/`
- JavaScript migration files with timestamps
- Run automatically on server start
- Manual migration: Use PocketBase admin interface

### File Storage
- Images stored in `pb_data/storage/`
- Automatic thumbnail generation (100x100)
- File naming with random suffixes for uniqueness

## API Development

### Custom Endpoints
All custom APIs use `/api/photocifu/` prefix and are handled via the new architecture:

**Gallery Management:**
- `POST /api/photocifu/gallery/create` - Create gallery with ZIP upload
  - Handler: `pkg/handlers/handlers.go:CreateGallery()`
  - Service: `pkg/container/services.go:GalleryServiceImpl`
  - Validation: `pkg/validation/validation.go:GalleryCreateRequest`

**Workflow Management:**
- `POST /api/photocifu/workflow/create` - Start workflow instance
  - Handler: `pkg/handlers/handlers.go:CreateWorkflow()`
  - Service: `pkg/container/services.go:WorkflowServiceImpl`
  - Supported types: `gallery_process`, `image_enhancement`, `cleanup`

**Signal Handling:**
- `POST /api/photocifu/signal/send` - Send workflow signals
  - Handler: `pkg/handlers/handlers.go:SendSignal()`
  - Service: `pkg/container/services.go:SignalServiceImpl`
  
**Settings:**
- `POST /api/photocifu/settings` - Update application settings
  - Handler: `pkg/handlers/handlers.go:UpdateSettings()`
  - Service: `pkg/container/services.go:SettingsServiceImpl`

### Authentication & Security
- All custom endpoints require `apis.RequireAuth()` middleware
- Uses PocketBase's built-in JWT authentication
- **NEW**: Enhanced client-side auth in `ui/src/lib/pocketbase.ts`
- **NEW**: Structured error responses with proper HTTP status codes
- **NEW**: Input validation with type safety

## Workflow Management

### go-workflows Integration
- Separate SQLite database (`pb_data/workflow.db`) for workflow state
- Background worker automatically started in `pkg/container/container.go`
- Workflow definitions in `workflow/` directory with proper input types
- Signal-based communication for external triggers
- **NEW**: Structured workflow input types (e.g., `GalleryProcessingInput`)

### Available Workflows
1. **Gallery Processing** (`gallery_process`):
   - Input: `workflow.GalleryProcessingInput`
   - Activities: Image processing, thumbnail generation
   - Timeout: 5 minutes with email notification fallback

2. **Image Enhancement** (`image_enhancement`): 
   - Future implementation for individual image processing

3. **Cleanup** (`cleanup`):
   - Future implementation for background maintenance

### Workflow Patterns
- **NEW**: Proper typed input structures instead of generic strings
- **NEW**: Timeout handling with notification fallbacks
- **NEW**: Structured error handling throughout workflow execution
- Long-running image processing tasks with progress tracking
- Batch operations on galleries with transaction safety

## Common Development Tasks

### Adding New API Endpoints (NEW ARCHITECTURE)
1. **Define Request/Response Types**: Add validation struct in `pkg/validation/`
2. **Create Service Interface**: Define interface in `pkg/container/container.go`
3. **Implement Service**: Add implementation in `pkg/container/services.go`
4. **Add Handler**: Create handler function in `pkg/handlers/handlers.go`
5. **Register Route**: Add route in `pkg/handlers/routes.go`
6. **Update Container**: Wire service in `pkg/container/container.go`

Example workflow:
```go
// 1. Add validation struct
type MyRequest struct {
    Name string `json:"name"`
}
func (r *MyRequest) Validate() error { /* validation logic */ }

// 2. Add service interface
type MyService interface {
    DoSomething(name string) error
}

// 3. Implement service
func (s *MyServiceImpl) DoSomething(name string) error { /* business logic */ }

// 4. Add handler
func (h *Handlers) MyEndpoint(e *core.RequestEvent) error { /* HTTP handling */ }

// 5. Register route
router.POST("/api/photocifu/my-endpoint", h.MyEndpoint).Bind(apis.RequireAuth())
```

### Adding New Frontend Routes
1. Create route file in appropriate `ui/src/routes/` subdirectory
2. Use SvelteKit's file-based routing with TypeScript
3. Add layout files for shared UI elements
4. Implement authentication guards using `ui/src/lib/pocketbase.ts`
5. **NEW**: Use structured configuration from `ui/src/config.ts`

### Database Schema Changes
1. Use PocketBase admin UI to modify collections
2. Migration files auto-generated in `pb_migrations/`
3. Test migrations on development data
4. Update TypeScript types if using generated types
5. **NEW**: Consider impact on service interfaces and validation

## Troubleshooting

### Common Issues
- **Port conflicts**: Default PocketBase port is 8090, change with `--http=0.0.0.0:8091`
- **File upload limits**: Configure via environment variables (`GALLERY_MAX_FILE_SIZE`)
- **CORS issues**: PocketBase handles CORS automatically for API calls
- **Build failures**: Ensure `go generate ./...` runs successfully
- **Frontend build errors**: Check TypeScript errors with `npm run check`
- **Container/Service issues**: Check dependency injection wiring in `pkg/container/`
- **Validation failures**: Verify request structures in `pkg/validation/`
- **modernc.org/libc version warnings**: Pin to v1.66.3 with `go mod edit -require=modernc.org/libc@v1.66.3` then run `go clean -modcache && go mod tidy`

### Development Database Reset
```bash
# Stop server, backup if needed, then:
rm pb_data/data.db pb_data/workflow.db

# Restart server to recreate with migrations
go run . serve --dev
```

### Performance Considerations
- **NEW**: Configuration-based limits for file sizes and gallery counts
- **NEW**: Structured error handling reduces debugging time
- **NEW**: Service interfaces enable better testing and optimization
- Image thumbnails generated automatically on upload
- SQLite performance suitable for small-medium galleries
- Static site generation provides optimal frontend performance
- Workflow engine handles async processing without blocking main thread

## Deployment

### Production Build
```bash
# Build frontend and embed into Go binary
go generate ./...

# Build optimized binary for target platform
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o photo-cifu

# Frontend is automatically embedded via go generate
# No separate frontend deployment needed
```

### Environment Configuration
**NEW**: The application now supports environment variables for configuration:

**Backend Environment Variables:**
- `WORKFLOW_DB_NAME`: Workflow database filename (default: "workflow.db")
- `GALLERY_MAX_FILE_SIZE`: Max ZIP file size in bytes (default: 100MB)
- `GALLERY_MAX_IMAGES`: Max images per gallery (default: 100)
- `WORKFLOW_DEFAULT_TIMEOUT`: Workflow timeout in seconds (default: 300)

**Frontend Environment Variables:**
- `PUBLIC_POCKETBASE_URL`: PocketBase API URL for production
- `PUBLIC_WEBSITE_URL`: Public website URL for metadata

### Deployment Structure
```
/
├── photo-cifu                    # Optimized binary with embedded frontend
├── pb_data/                      # Application data directory
│   ├── data.db                   # Main PocketBase database
│   ├── workflow.db               # Workflow state database
│   ├── storage/                  # File uploads and thumbnails
│   └── types.d.ts                # TypeScript definitions
└── pb_migrations/                # Database migrations (optional, embedded)
```

### Container Deployment (Docker)

**Quick Start with Docker:**
```bash
# Build and run with docker-compose (production)
docker-compose up --build

# Development mode with hot reload
docker-compose -f docker-compose.dev.yml up --build

# Build Docker image manually
docker build -t photo-cifu .

# Run with custom port and data persistence
docker run -d \
  --name photo-cifu \
  -p 8091:8090 \
  -v $(pwd)/pb_data:/app/pb_data \
  -v $(pwd)/pb_migrations:/app/pb_migrations \
  -e GALLERY_MAX_FILE_SIZE=209715200 \
  photo-cifu
```

**Docker Configuration:**
- **Multi-stage build**: Compiles Go application with embedded SvelteKit frontend
- **Pure Go compilation**: Uses `CGO_ENABLED=0` with modernc.org/sqlite for maximum compatibility
- **Persistent data**: Volumes for `pb_data/` (database and uploads) and `pb_migrations/`
- **Environment variables**: All configuration via environment variables
- **Health checks**: Built-in container health monitoring
- **Production ready**: Optimized Alpine Linux runtime (~50MB total image size)

**SQLite Driver Choice:**
- Uses pure Go SQLite port (modernc.org/sqlite) for compilation simplicity
- Trade-off: 2-3x slower than CGO-based mattn/go-sqlite3, but eliminates build dependencies
- Optimal for deployment reliability over peak SQLite performance

**Available Docker Compose configurations:**
- `docker-compose.yml`: Production deployment
- `docker-compose.dev.yml`: Development with volume mounts and dev mode

## Key Changes Summary

This refactoring introduced a clean architecture with the following major improvements:

### ✅ What Was Improved
1. **Dependency Injection**: `pkg/container/` system replaces global state
2. **Error Handling**: Structured errors in `pkg/errors/` with proper HTTP codes
3. **Input Validation**: Type-safe validation in `pkg/validation/`
4. **Route Handling**: Clean handlers in `pkg/handlers/` with separation of concerns
5. **Configuration**: Environment-based config in `pkg/config/`
6. **Service Layer**: Interface-based services for better testability
7. **Workflow Types**: Proper typed inputs instead of generic strings
8. **Frontend Config**: Enhanced configuration with environment variable support

### 📋 Development Guidelines
- **Direct DB calls**: Use service interfaces instead of direct PocketBase calls
- **Error handling**: Use `pkg/errors.HandleError()` for consistent responses
- **Route registration**: Add new routes in `pkg/handlers/routes.go`
- **Configuration**: Use `pkg/config.New()` instead of hardcoded values

### 🔧 Development Best Practices
- Always use the dependency injection container for services
- Validate all inputs using `pkg/validation/` structures
- Handle errors with structured `pkg/errors/` types
- Follow the handler → service → validation pattern for new endpoints
- Use environment variables for configuration instead of hardcoded values
- IMPORTANT: use only npm