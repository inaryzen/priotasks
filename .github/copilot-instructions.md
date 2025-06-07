# Github Copilot Instructions for PrioTasks Project

## Project Overview
- Go-based task management application

## Technical Stack
- Backend: Go
- Frontend: HTMX + Templ. You can find documentation for HTMX in `docs/htmx_doc.md` file.
- Database: SQLite with modernc.org/sqlite

## Project Structure
- `assets/`: Static files (CSS, JS, favicon)
- `common/`: Shared utilities
- `components/`: Templ-based UI components
- `consts/`: Constants and enums
- `csv/`: CSV-related functionality
- `db/`: Database operations
- `docs/`: Feature documentation
- `handlers/`: HTTP request handlers
- `models/`: Data models
- `services/`: Business logic

## Testing
All the new functions and functionality MUST BE supplied with the tests. Once you finish the main implementation add the tests accordingly. Run the tests at the end of you work and verify that that the code changes are correct.

- Use `setupTestDB` from db_common_test.go when writing tests for `db` package
- Use "DB mocks" with `db.NoOpDB` from `db/noopdb.go` when writing tests for the packages other than `db` package. See an example in `services/asksService_test.go`
- Do not write tests unless explicitly asked to write tests.
- ALWAYS USE TEST NAMING FORMAT: `Test_FunctionName_Scenario`

## Best Practices
1. Component Development:
   - Use Templ for HTML templating
   - Follow HTMX patterns for dynamic interactions. Avoid using JavaScript.

2. Database Operations:
   - Use the db package for all database interactions

## Reference Links
- [HTMX Documentation](https://htmx.org/)
- [Templ Guide](https://templ.guide/)