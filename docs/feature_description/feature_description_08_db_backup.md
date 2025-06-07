# Feature Description Document - F08: SQLite Database Backup System

## Overview
Replace the current CSV export functionality with an automated SQLite database backup system. The system will create a backup of the database file at application startup and maintain a maximum of two backup files, automatically removing the oldest backup when this limit is reached.

## Requirements
### Functional Requirements
- Remove existing CSV export functionality completely (respective command flags, classes, etc)
- Create an automatic backup of the SQLite database file at application startup
- Maintain a maximum of two backup files at any time
- Implement a cleanup mechanism to remove the oldest backup when exceeding the limit
- Use a consistent naming convention for backup files (e.g., `priotasks_db_backup_YYYYMMDD_HHMMSS.db`)

### Non-Functional Requirements
- Minimal impact on application startup time
- Zero data loss during backup process

## User Stories
```
As a user
I want my database to be automatically backed up when starting the application
So that I can recover my data in case of corruption or accidental changes

As a user
I want old backups to be automatically managed
So that I don't need to manually clean up backup files
```

## Technical Specifications
### Architecture
- Remove CSV-related components (`csv` package)
- Add new backup service in `services` package
- Modify application startup flow to include backup creation

## Testing Requirements
- Unit test scenarios:
  - Test backup file creation
  - Test backup rotation (removing oldest file)
  - Test backup naming convention
  - Test error handling during file operations
- Integration test cases:
  - Test backup process during application startup
  - Test backup limit enforcement
  - Test backup file integrity

## Security Considerations
- Ensure backup files have appropriate file permissions
- Handle backup failures gracefully without compromising the main database

## Resources
### Related Code
- Package to remove: `csv/`
- Files to modify:
  - `main.go` (add backup initialization)
  - Create new `services/backupService.go`
