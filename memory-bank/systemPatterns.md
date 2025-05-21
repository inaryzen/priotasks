# System Patterns: PrioTasks

## System Architecture

PrioTasks follows a layered architecture pattern with clear separation of concerns:

```mermaid
flowchart TD
    A["HTTP Layer(handlers/*.go)"] --> B["Service Layer(services/*.go)"]
    B --> C["Data Layer(db/*.go)"]
    C --> D["Database(SQLite)"]
```

### UI Components (Templ)

```mermaid
flowchart TD
    A[TasksView] --> B[FilterPanel]
    A --> C[TaskTable]
    D[TaskModal]
```

## Design Patterns

### Repository Pattern
- Implemented in the `db` package
- Abstracts database operations from business logic
- Provides a clean interface for data access

### Service Layer Pattern
- Implemented in the `services` package
- Encapsulates business logic
- Coordinates between handlers and data access

### Component-Based UI
- UI elements organized as reusable components
- Each component has a single responsibility
- Components can be composed to create complex views

### Model-View Pattern
- Models defined in the `models` package
- Views implemented as Templ components
- Clear separation between data and presentation

## Key concepts and models
- `models/tasks.go` - The primary entity of the project, contain `Task` struct
- `models/settings.go` - Supporting logic, primary user settings. Contains all the preferences and the table filtration settings.
- Every entity such as `Task` or `Settings` are persisted in the database. Whenever you change or introduce new fields to these structs you must also update their SQL tables accordingly. Typically you can find their table definitions in their `db` file, for instance `go/dbtasks.go` or `go/dbsettings.go` respectively. Remember to always introduce migration for the DB changes. See `initSettings()` and `settingsTableAddTagsColumn()` for examples. 

## Critical Implementation Paths

### Task Creation Flow

```mermaid
sequenceDiagram
    participant User
    participant Handler as handlers.PostTaskHandler
    participant Service as services.SaveNewTask
    participant DB as db.SaveTask
    
    User->>Handler: Submit form (HTMX POST)
    Handler->>Service: Process request data
    Service->>DB: Persist task
    DB-->>Service: Confirmation
    Service-->>Handler: Success/Error
    Handler-->>User: Updated task table (HTMX response)
```

### Task Filtering Flow

```mermaid
sequenceDiagram
    participant User
    participant Handler as Filter Handler
    participant Settings as User Settings
    participant Service as Task Service
    participant DB as Database
    
    User->>Handler: Select filter criteria (HTMX request)
    Handler->>Settings: Update filter settings
    Handler->>Service: Request filtered tasks
    Service->>DB: Query with filter criteria
    DB-->>Service: Return filtered tasks
    Service-->>Handler: Filtered task list
    Handler-->>User: Updated task table (HTMX response)
```

### Task Value Calculation

```mermaid
flowchart LR
    A[Task Attributes Set] --> B["Task.CalculateValue()"]
    B --> C[Value Stored with Task]
    C --> D[Value Displayed in UI]
```

## Component Relationships

### Handler Dependencies

```mermaid
flowchart TD
    A[Handlers] --> B[Services]
    A --> C[Templ Components]
    B --> D[Database Layer]
```

### Service Dependencies
- Services depend on the database layer for data access
- Services implement business logic independent of UI concerns

### Database Layer
- Provides a unified interface for data operations
- Abstracts the underlying database implementation
- Supports different database backends (currently SQLite)

### Component Hierarchy

```mermaid
classDiagram
    TasksView "1" --> "1" FilterPanel
    TasksView "1" --> "1" TaskTable
    TaskTable "1" --> "*" SortableHeader
    class TasksView {
        +Render()
    }
    class FilterPanel {
        +Render()
    }
    class TaskTable {
        +Render()
    }
    class SortableHeader {
        +Render()
    }
    class TaskModal {
        +Render()
    }
```
