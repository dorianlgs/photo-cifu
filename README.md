# PhotoCifu

A photo gallery application with workflow-based image processing, built with Go + PocketBase backend and SvelteKit frontend.

## Features

- **Photo Gallery Management**: Create galleries by uploading ZIP archives of images
- **Workflow Engine**: Async image processing with go-workflows
- **User Authentication**: JWT-based auth with PocketBase
- **File Storage**: Automatic thumbnail generation and file management
- **Responsive Frontend**: SvelteKit with TypeScript and Tailwind CSS

## Architecture

PhotoCifu follows a clean architecture pattern with dependency injection:

### Backend Structure
```
├── main.go                 # Application entry point
├── pkg/
│   ├── config/            # Configuration management
│   ├── container/         # Dependency injection container
│   ├── errors/            # Centralized error handling
│   ├── handlers/          # HTTP request handlers
│   └── validation/        # Input validation
├── workflow/              # Workflow definitions and activities
├── tools/                 # Utility functions
└── ui/                    # SvelteKit frontend
```

### Key Components

- **Container**: Dependency injection system for services
- **Handlers**: Clean HTTP route handlers with validation
- **Services**: Business logic implementations with interfaces
- **Workflows**: Long-running async processes for image processing
- **Validation**: Request validation with structured error responses

## Development

### Prerequisites
- Go 1.24+
- Node.js/Yarn for frontend development
- SQLite (included with PocketBase)

### Quick Start

1. **Install dependencies**:
   ```bash
   go mod tidy
   cd ui/ && yarn install && cd ..
   ```

2. **Run development server**:
   ```bash
   go run . serve --dev
   ```

3. **Access the application**:
   - Frontend: http://localhost:8090
   - Admin UI: http://localhost:8090/admin

### Frontend Development

```bash
cd ui/

# Development server with hot reload
yarn run dev

# Build production frontend
yarn run build

# Code quality checks
yarn run lint
yarn run format
yarn run check
```

### Build for Production

```bash
# Generate embedded frontend assets
go generate ./...

# Build optimized binary
GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o photo-cifu

# Run production server
./photo-cifu serve
```

## Configuration

### Environment Variables

**Backend Configuration**:
- `WORKFLOW_DB_NAME`: Workflow database filename (default: "workflow.db")
- `GALLERY_MAX_FILE_SIZE`: Max gallery ZIP size in bytes (default: 100MB)
- `GALLERY_MAX_IMAGES`: Max images per gallery (default: 100)
- `WORKFLOW_DEFAULT_TIMEOUT`: Default workflow timeout in seconds (default: 300)

**Frontend Configuration**:
- `PUBLIC_POCKETBASE_URL`: PocketBase API URL
- `PUBLIC_WEBSITE_URL`: Public website URL

### Command Line Options

```bash
# Development mode with debugging
go run . serve --dev --debug

# Custom data directory
go run . serve --dir ./custom-data

# Custom public files directory
go run . serve --publicDir ./custom-public

# Disable auto-migration
go run . serve --automigrate=false
```

## API Endpoints

All custom APIs use the `/api/photocifu/` prefix:

- `POST /api/photocifu/gallery/create` - Create gallery with ZIP upload
- `POST /api/photocifu/workflow/create` - Start workflow instance
- `POST /api/photocifu/signal/send` - Send workflow signals
- `POST /api/photocifu/settings` - Update application settings

All endpoints require authentication via PocketBase JWT tokens.

## Database Schema

### Collections
- **users**: Authentication and user profiles
- **galleries**: Photo gallery metadata
- **images**: Individual image records with file references
- **messages**: System messaging/notifications

### File Storage
- Images stored in `pb_data/storage/`
- Automatic thumbnail generation (100x100)
- Workflow state in separate SQLite database (`workflow.db`)

## Workflow System

PhotoCifu uses go-workflows for async image processing:

### Workflow Types
- `gallery_process`: Process uploaded gallery images
- `image_enhancement`: Individual image processing
- `cleanup`: Background maintenance tasks

### Example Workflow Usage
```bash
# Create a gallery processing workflow
curl -X POST http://localhost:8090/api/photocifu/workflow/create \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "workflow_type": "gallery_process",
    "input": {
      "gallery_id": "abc123",
      "gallery_name": "My Gallery",
      "user_email": "user@example.com"
    }
  }'
```

## Development Commands

### Go Backend
```bash
# Run with live reload
go run . serve --dev

# Update dependencies
go get -u -t ./...
go mod tidy

# Run tests (when available)
go test ./...

# Build for specific platform
GOOS=windows GOARCH=amd64 go build -o photo-cifu.exe
```

### Database Management
```bash
# Reset development database
rm pb_data/data.db pb_data/workflow.db

# Backup database
cp pb_data/data.db pb_data/backups/backup-$(date +%Y%m%d).db
```

### Troubleshooting

**Common Issues**:
- **Port conflicts**: Change port with `--http=0.0.0.0:8091`
- **File upload limits**: Configure via environment variables
- **Workflow failures**: Check `pb_data/workflow.db` for state
- **Frontend build errors**: Run `yarn run check` for TypeScript issues

**Development Database Reset**:
```bash
# Stop server, backup if needed, then:
rm pb_data/data.db pb_data/workflow.db
# Restart server to recreate with migrations
```

## Contributing

1. Follow Go and TypeScript best practices
2. Use the dependency injection container for new services
3. Add proper validation for all inputs
4. Include error handling with structured responses
5. Update tests when adding new features

## License

MIT License - see LICENSE file for details.
