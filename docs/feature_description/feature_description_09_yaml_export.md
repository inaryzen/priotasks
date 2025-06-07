# Feature Description Document - F09: Export Tasks to YAML

## Overview
Add functionality to export currently visible tasks in the task table to a YAML file format. This feature will allow users to download their tasks data in a structured, human-readable YAML format for backup purposes or integration with other tools.

## Requirements
### Functional Requirements
- Add "Export YAML" option to the Operations menu
- Export only the tasks that are currently visible in the task table (respecting current filters)
- Generate a properly formatted YAML file containing task data
- Automatically trigger file download in the browser
- Include all relevant task fields: ID, title, description, priority, tags, value, creation date, etc.

### Non-Functional Requirements
- File format should be standard YAML that can be parsed by any YAML parser
- Generated YAML must be human-readable with proper indentation
- Maximum number of tasks that can be exported at once is 1000
- When exporting to YAML, make sure to escape special symbols from title and description 

## User Stories
```
As a user
I want to export my visible tasks to a YAML file
So that I can backup my tasks or use them in other tools that support YAML format
```

## Technical Specifications
### Architecture
- Add new YAML export handler in the `handlers` package
- Add YAML export service in the `services` package
- Add YAML marshaling functionality using Go's `gopkg.in/yaml.v3` package
- Extend the Operations menu UI component to include the new export option

### Data Model
No changes to existing data models required.

### API Endpoints
- `GET /tasks/export/yaml`
  - Respects current filter parameters from the task table
  - Returns YAML file as attachment
  - Content-Type: application/x-yaml
  - Content-Disposition: attachment; filename="tasks.yaml"

### Dependencies
- `gopkg.in/yaml.v3` for YAML generation

## UI/UX Design
- Add "Export YAML" button to the Operations menu
- Use HTMX to trigger the export operation
- Download should start automatically when the export is complete

## Testing Requirements
- Unit test scenarios:
  - Test YAML marshaling of task data
  - Test export with various task field combinations
  - Test export with empty task list
  - Test export with filtered tasks
- Integration test cases:
  - Test the complete export flow from UI action to file download
  - Verify YAML file structure and content
- Acceptance criteria:
  - YAML file contains all visible tasks
  - YAML is properly formatted and valid
  - File downloads automatically in browser
  - Export respects current table filters

## Security Considerations
- Validate maximum number of tasks that can be exported at once
- Ensure proper escaping of special characters in YAML output
- Add appropriate request validation

## Resources
- YAML v3 package documentation: https://pkg.go.dev/gopkg.in/yaml.v3
