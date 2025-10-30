# Feature Description Document - 12

## Overview
Implement a task cloning feature that allows users to create a copy of an existing task when viewing/editing the task details. This feature will help users quickly create similar tasks without having to manually re-enter all the task information.

## Requirements
### Functional Requirements
- Add a "Clone" button to the task detail modal (when viewing/editing an existing task)
- The Clone button should create a new task with all the same properties as the original task except:
  - New unique ID
  - Title should be prefixed with "Copy of " 
  - Creation date should be set to current timestamp
  - Updated date should be set to current timestamp
  - Completion status should be reset to false (not completed)
- After cloning, the user should see the edit panel of the newly created task
- The Clone button should only be visible when viewing an existing task (not when creating a new task)

## User Stories
```
As a task manager
I want to clone an existing task
So that I can quickly create similar tasks without re-entering all the details

As a project manager
I want to duplicate recurring tasks
So that I can maintain consistency in task properties and save time on task creation
```

## Technical Specifications
### Architecture
- Add new HTTP endpoint for task cloning
- Extend TaskModal component to include Clone button
- Add new service method for task cloning logic

### Data Model
- No changes to existing data models required
- New task will use existing Task model structure

### API Endpoints
- `POST /tasks/{id}/clone` - Clone an existing task by ID
  - Request: Task ID in URL path
  - Behavior: Creates new task with copied properties and redirects to task list

### Dependencies
- Existing HTMX framework for UI interactions
- Existing Templ templating system
- Current database layer (SQLite)

## Implementation Plan
- [ ] Add Clone button to TaskModal component: Add Clone button in the form-buttons-left section next to Delete button
- [ ] Create new HTTP endpoint: Add POST /tasks/{id}/clone endpoint in tasksHandlers.go
- [ ] Add service method: Implement CloneTask method in tasksService.go and add to interface
- [ ] Add task cloning logic: Implement the core cloning logic with proper property copying and defaults
- [ ] Use the existing database methods to save a new task
- [ ] Validate implementation: Use chrome-devtools-mcp to test the complete clone workflow