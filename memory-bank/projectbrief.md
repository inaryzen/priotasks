# Project Brief: PrioTasks

## Project Overview
PrioTasks is a task management application. It provides a web-based interface for managing prioritized tasks with various attributes and filtering capabilities.

## Core Requirements

### Task Management
- Create, read, update, and delete tasks
- Assign priorities, impact levels, cost estimates, and "fun" ratings to tasks
- Mark tasks as completed, work-in-progress, or planned
- Tag tasks for organization and filtering
- Calculate task value based on priority, impact, cost, and fun factors

### User Interface
- Display tasks in a sortable table view
- Filter tasks by various attributes and tags
- Provide modal views for task creation and editing
- Support for responsive design

### Technical Requirements
- Built with Go backend
- Use Templ for HTML templating
- Implement HTMX for dynamic UI updates without full page reloads
- SQLite database for data persistence

## Non-Goals
1. Multi-user support (single user application)
2. Mobile application (web-based only)
4. Integration with third-party services
