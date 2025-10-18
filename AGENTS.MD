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

## Coding Guidelines
- When log messages, always attach the name of the function in front of the message. Example:
    ```go
    func FooBar() error {
        pfx := "FooBar:"
        log.Printf("%s Hello %s", pfx, "World")    
    }
    ```
- ONLY add log statements at the beginning of service functions and when error is not returned. AVOID adding log statements in other places. NEVER add log statement at the end of the functions!!!
- Include function prefix `pfx` when creating new errors, e.g. `fmt.Errorf("%s new error message: %w", pfx, err)`
- Do not print log messages if you create a new error. The upstream client code will capture the error and print the message properly.
- When writing loops for 0 .. N, use modern "range over" notation, Example:
    ```go
    // instead of:
    for i:=0; i < N; i++ {
        ...
    }
    // write this:
    for i := range N {
        ...
    }
    ```
- When adding new public function to services, add them to the respective interface. Often the specific implementation of service would have a respective interface. When you are not sure if a function should be public or not, ASK USER! 
- When you write comments, ALWAYS add "ai: " prefix to a comment. Example:
    ```go
    // ai: function "foo" etc..
    func foo() {
        ...
    }
    ```
- Use Templ for HTML templating
- Follow HTMX patterns for dynamic interactions. Avoid using JavaScript.
- Use the db package for all database interactions

## Testing
- AVOID writing tests for new functions and functionality before the user tells you so.
- Once you finish the main implementation, PROPOSE to create tests and explicitly specify what exact tests you would like to write. If the user approves that, carry on and create the tests.
- Pay attention to what the user says: 
    - If they say to create a single test for the specific function, DO ONLY THAT
    - If they say to create tests for all the functions in a file / for the capability, then you are free to create all the tests.
- If the user asks to create tests, ALWAYS explicitly specify what exact tests you are going to create and ask for the confirmation from the user before proceeding.
- ALWAYS USE TEST NAMING FORMAT: `Test_FunctionName_Scenario`
- Use `setupTestDB` from db_common_test.go when writing tests for `db` package
- Use "DB mocks" with `db.NoOpDB` from `db/noopdb.go` when writing tests for the packages other than `db` package. See an example in `services/asksService_test.go`

## Working mode
- When you asked to implement a feature, ALWAYS provide the implementation plan first. NEVER start coding right away!
- When preparing an implementation plan or following the plan, ALWAYS follow the top-down principle: FIRST write top level code such as usage of methods and interfaces and THEN write low level code such as implementation of methods. If the implementation passes several abstraction layers: FIRST implement changes on the top layer and then follow to the next lower one. Example: you have layered architecture a) HTTP handler b) Service c) DB Repository. You need to implement a CRUD operation for a new entity. You planning and implementation principle is the following: 1) Define the endpoint shape 2) Flesh out rest of the API, such as payload structure 3) Create DTO 4) Add call to Service layer 5) Add function to Service interface 6) Implement domain Model is needed 7) Implement Service function 8) Add call to Repository 9) Add function to Repository interface 10) Add implementation to Repository 11) Add DB migration to for a new entity etc.
- When following implementation plan, ALWAYS implement one step at a time and then ask confirmation to proceed. NEVER implement the whole plan at once! 
- When asking the confirmation to proceed, ALWAYS provide a summary of the next step in the implementation plan!

# Implementation validation
- Always validate your implementation by starting the application as explained in `README.ME` using `run.sh` and check that the solution implemented correctly using `chrome-devtools-mcp`.
   
## Reference Links
- [HTMX Documentation](https://htmx.org/)
- [Templ Guide](https://templ.guide/)
