# Technical Context: PrioTasks

## Technologies Used

### Backend
- **Go**: Core programming language for the backend
- **Standard Library HTTP**: For HTTP server and routing
- **SQLite**: Embedded database for data storage
- **CSV**: For data import/export functionality

### Frontend
- **Templ**: Type-safe HTML templating for Go
- **HTMX**: For dynamic UI updates without heavy JavaScript
- **CSS**: Custom styling for the application interface
- **Minimal JavaScript**: Only where necessary for enhanced interactions

## Development Setup

### Prerequisites
- Go installed (with `/home/{user}/go/bin` in PATH)
- Templ CLI tool installed (`go install github.com/a-h/templ/cmd/templ@latest`)

### Running the Application
- Execute `./run.sh` script to start the application
- The application will be available at `http://localhost:{port}` (port defined in configuration)

### Testing
- Run tests with `go test ./...`
- Tests are organized by package with naming convention `Test_FunctionName_Scenario`
- Table-driven tests are used for testing functions with multiple scenarios

### Version Management
- Use `./bump-version.sh` to increment the application version

## Technical Constraints

### Database
- SQLite is used for simplicity and portability
- Database schema is managed through code
- Migration support for schema changes

### UI Rendering
- Server-side rendering with Templ
- Dynamic updates via HTMX
- No client-side JavaScript framework

### Performance Considerations
- Efficient database queries for task filtering and sorting
- Minimal data transfer for UI updates (HTMX partial updates)
- Embedded database for reduced latency

## Dependencies

### Core Dependencies
- `github.com/a-h/templ`: HTML templating
- `github.com/google/uuid`: For generating unique IDs
- `github.com/mattn/go-sqlite3`: SQLite driver for Go

### Project Structure
```
priotasks/
├── assets/         # Static assets (CSS, JS, favicons)
├── common/         # Common utilities and configuration
├── components/     # Templ UI components
├── consts/         # Constants used throughout the application
├── csv/            # CSV import/export functionality
├── db/             # Database access layer
├── docs/           # Documentation and feature descriptions
├── handlers/       # HTTP request handlers
├── models/         # Data models
├── services/       # Business logic services
└── static/         # Static file serving
```

## Tool Usage Patterns

### Database Operations
- Repository pattern for database access
- Transactions for multi-step operations
- Prepared statements for query efficiency

### HTTP Request Handling
- Handler functions follow the standard Go HTTP handler signature
- URL patterns defined in main.go
- HTMX-specific response headers for dynamic updates

### Templ Component Usage
- Components are defined in `.templ` files
- Components are compiled to Go code with the Templ CLI
- Components are composed to create complex views

### Testing Approach
- Unit tests for core functionality
- Mock implementations for database testing
- Test utilities in `db/db_common_test.go`

## Configuration Management

### Application Configuration
- Configuration loaded at startup in `common.InitConfig()`
- Environment-based configuration options
- Debug mode toggle for development

### Build and Deployment
- Simple build process with Go standard tools
- Embedded assets for simplified deployment
- Single binary output
