# Feature Description Document - F07

## Overview
A new search section will be added to the filter panel, allowing users to filter tasks by searching for text in both the title and description of tasks. The search will be implemented using HTMX's active search pattern, providing real-time filtering as the user types without requiring custom JavaScript.

## Requirements
### Functional Requirements
- Add a new "Search" section to the filter panel with the same style as other sections
- Implement a text input field that filters tasks as the user types
- Search should check both task title and task content (description)
- Results should update in real-time without page reloads
- Search should work alongside existing filters (tags, completion status, etc.)
- Use HTMX for implementation, following the active-search pattern
- No custom JavaScript should be used for the implementation

## User Stories
```
As a user
I want to search for tasks by entering text
So that I can quickly find specific tasks without scrolling through the entire list
```

```
As a user
I want search results to update as I type
So that I can refine my search in real-time without waiting for page reloads
```

## Technical Specifications
### Architecture
- Extend the existing filter panel component to include a search section
- Add a new handler for processing search requests
- Update the task service to support text-based filtering
- Implement using HTMX's active-search pattern

### Data Model
- Add a new search field to the TasksQuery struct in models/settings.go:
```go
type TasksQuery struct {
    // existing fields...
    SearchText string
}
```

### API Endpoints
- Add a new endpoint for handling search requests:
  - Path: `/filter/search`
  - Method: POST
  - Request: Form data with search text
  - Response: Updated task table HTML

## UI/UX Design
- Implement debouncing to prevent excessive requests during typing
- Show visual feedback during search (optional loading indicator)

## Implementation Details
### HTMX Implementation
The search will be implemented using HTMX's active-search pattern:
- Text input with `hx-post` to send search text to the server
- `hx-trigger="keyup changed delay:300ms"` to debounce requests
- `hx-target` to update the task table with search results
- Server-side filtering based on the search text

### Code Changes
1. Update FilterPanel component to include search section
2. Add search handler in userHandlers.go
3. Update TasksQuery model to include search text
4. Extend task filtering logic in services/tasksService.go

### Example HTMX Implementation
```html
<input
  type="search"
  name="search"
  placeholder="Search tasks..."
  hx-post="/filter/search"
  hx-trigger="keyup changed delay:300ms"
  hx-target="#task-table-container"
  hx-indicator=".search-indicator"
/>
<div class="search-indicator htmx-indicator">Searching...</div>
```

## Testing Requirements
- Test search functionality with various input patterns
- Verify that search works alongside other filters
- Ensure proper handling of special characters in search text
