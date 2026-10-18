# SIS Project - Checkpoint 1 (Core Go)

Domain: **Text Analyzer**. A task is one unit of work: analyze a piece of text and produce statistics.

- Input: a text string (leading and trailing whitespace is trimmed).
- Processing: count words, runes (Unicode characters) and bytes.
- Result: `Result{Words, Runes, Bytes}` stored inside the task.

Example: `"Hello Go world"` gives 3 words, 14 runes, 14 bytes. `"Привет мир"` gives 2 words, 10 runes, 19 bytes.

## Validation

Empty or whitespace-only text is rejected with `ErrInvalidInput`.

## Task lifecycle

```
Pending --> Processing --> Completed
                |
                +--------> Failed
```

| From       | Allowed next states   |
|------------|-----------------------|
| Pending    | Processing            |
| Processing | Completed, Failed     |
| Completed  | none (final)          |
| Failed     | none (final)          |

Any other transition returns `ErrInvalidTransition` and leaves the task unchanged. For example, processing a task that is already `Completed` or `Failed` is rejected.

## Errors

- `ErrInvalidInput` - invalid task input.
- `ErrTaskNotFound` - no task with the given ID.
- `ErrInvalidTransition` - operation not allowed in the current state.

Errors are wrapped with `%w`, so callers use `errors.Is`.

## Structure

- `main.go` - small demo.
- `internal/task/task.go` - Task, State, Result, state transitions.
- `internal/task/processor.go` - `Processor` interface, `TextAnalyzer`, input validation.
- `internal/task/manager.go` - in-memory storage (`map[int]*Task`), Create / Get / List / Process.
- `internal/task/*_test.go` - unit tests.

## Run and check

```
go run .
gofmt -l .
go vet ./...
go test ./...
```
