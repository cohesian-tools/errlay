# errlay

`errlay` is a small, dependency-free Go package for attaching a logical **layer** and a numeric **code** to errors. It wraps underlying causes and integrates with the standard library's `errors.Is`, `errors.As`, and `errors.Unwrap` conventions.

Rather than comparing messages, `errlay` lets you search an error chain using a **partial target**: match by layer, by code, or by both.

## Usage

Import the package using the module path declared in your `go.mod` (replace `YOUR_MODULE_PATH` below).

```go
package main

import (
    "errors"
    "fmt"

    errlay "YOUR_MODULE_PATH"
)

func main() {
    cause := errors.New("connection refused")
    err := errlay.New("database", 10, "query failed", cause)

    fmt.Println(err)
    // |layer database - 10| query failed connection refused

    layer := "database"
    byLayer, _ := errlay.Target(&layer, nil)
    fmt.Println(errors.Is(err, byLayer)) // true

    code := uint8(10)
    byCode, _ := errlay.Target(nil, &code)
    fmt.Println(errors.Is(err, byCode)) // true

    exact, _ := errlay.Target(&layer, &code)
    fmt.Println(errors.Is(err, exact)) // true

    fmt.Println(errors.Is(err, cause)) // true (underlying cause)

    var found *errlay.Error
    if errors.As(err, &found) {
        fmt.Println(found.Layer(), found.Code(), found.Message())
        // database 10 query failed
    }
}
```

## Matching rules

`Target(layer *string, code *uint8)` builds an error target for `errors.Is`. A `nil` pointer means **do not filter by that field**; a non-`nil` pointer means **require an exact match**.

| Target fields | Match condition |
| --- | --- |
| `layer != nil`, `code == nil` | Same layer |
| `layer == nil`, `code != nil` | Same code |
| `layer != nil`, `code != nil` | Same layer **and** same code |
| Both `nil` | Invalid target; `Target` returns an error |

The matching rule is applied to **each individual error** in the chain. If a layer matches one wrapped error and the code matches a different wrapped error, the combined target does **not** match.

A specified empty layer (`""`) or code `0` is a valid match criterion. Comparisons are case-sensitive. Messages are not part of matching.

### Wrapped errors

```go
cause := errors.New("connection refused")
err := errlay.New("handler", 1, "request failed",
    errlay.New("database", 10, "query failed", cause),
)

layer := "database"
code := uint8(10)
target, _ := errlay.Target(&layer, &code)

fmt.Println(errors.Is(err, target)) // true: matches the inner error
fmt.Println(errors.Is(err, cause))  // true: finds the root cause
```

`errors.Is` traverses wrapped causes through `Unwrap()`. The package's `Is(target error) bool` method recognizes targets constructed by `errlay.Target`. A regular `*errlay.Error` is not treated as a field-filter target.

### Capturing an error with `errors.As`

```go
var e *errlay.Error
if errors.As(err, &e) {
    fmt.Printf("layer=%s code=%d message=%s\n", e.Layer(), e.Code(), e.Message())
}
```

`errors.As` retrieves the first matching `*errlay.Error` in the chain. No custom `As` method is necessary.

## API

```go
func New(layer string, code uint8, msg string, err error) error
func Target(layer *string, code *uint8) (target error, err error)

type Error struct { /* unexported fields */ }

func (e *Error) Error() string
func (e *Error) Layer() string
func (e *Error) Code() uint8
func (e *Error) Message() string
func (e *Error) Unwrap() error
func (e *Error) Is(target error) bool
```

- **`New`** creates a layered error. Use `nil` as the cause to create a leaf error.
- **`Target`** creates a search criterion; it fails if neither layer nor code is provided.
- **`Error`** formats the error as `|layer <layer> - <code>| <message>`, appending the cause when present.
- **`Layer`, `Code`, `Message`** expose the associated metadata.
- **`Unwrap`** exposes the underlying cause for standard Go error traversal.
- **`Is`** checks whether the current error satisfies the specified target fields.

**Note:** `Target` retains the supplied pointers, so modifying the pointed-to variables later also changes the target's matching criteria. Keep those values unchanged while using the target.

## Tests

```sh
go test ./...
go test ./... -v -cover
```

## License

MIT. See [LICENSE](LICENSE).
