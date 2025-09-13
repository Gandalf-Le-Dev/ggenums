# GGenums - Generate Golang Enums

[![Go Report Card](https://goreportcard.com/badge/github.com/Gandalf-Le-Dev/ggenums)](https://goreportcard.com/report/github.com/Gandalf-Le-Dev/ggenums)
[![GoDoc](https://godoc.org/github.com/Gandalf-Le-Dev/ggenums?status.svg)](https://godoc.org/github.com/Gandalf-Le-Dev/ggenums)

GGenums is a powerful code generation tool that creates type-safe enums in Go with rich functionality. It follows the **enumer-style approach**, parsing Go type declarations directly from your source code to generate methods for string conversion, JSON marshaling/unmarshaling, and validation.

## Features

- 🎯 **Enumer-compatible API** - Uses `-type` flag and parses Go declarations
- 🔒 **Type-safe enum values** with full IDE support
- 📦 **JSON marshaling/unmarshaling** support
- ✅ **Built-in validation** with `IsValid()` method
- 🔄 **String conversion** methods (`String()`, `ParseX()`)
- 📝 **Clean separation** - source files own types, generated files provide methods
- ⚡ **All values functions** - `AllXs()` to get all enum values
- 🛡️ **No redeclaration conflicts** - proper file separation

## Installation

```bash
go install github.com/Gandalf-Le-Dev/ggenums@latest
```

## Quick Start

1. **Define your enums using Go types and constants:**

    ```go
    //go:generate ggenums -type=Status,Priority

    type Status int
    const (
        StatusPending Status = iota
        StatusActive
        StatusInProgress
        StatusCompleted
    )

    type Priority int
    const (
        PriorityLow Priority = iota
        PriorityMedium
        PriorityHigh
        PriorityUrgent
    )
    ```

2. **Run the generator:**

    ```bash
    go generate
    ```

3. **Use your generated enums with rich functionality!**

## Complete Example

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
)

//go:generate ggenums -type=Status,Role,Priority

type Status int
const (
    StatusPending Status = iota
    StatusActive
    StatusInProgress
    StatusCompleted
    StatusCancelled
)

type Role int
const (
    RoleAdmin Role = iota
    RoleUser
    RoleGuest
    RoleModerator
)

type Priority int
const (
    PriorityLow Priority = iota
    PriorityMedium
    PriorityHigh
    PriorityUrgent
)

type User struct {
    ID       int      `json:"id"`
    Name     string   `json:"name"`
    Status   Status   `json:"status"`
    Role     Role     `json:"role"`
    Priority Priority `json:"priority"`
}

func main() {
    // Create a user with enum values
    user := User{
        ID:       1,
        Name:     "Alice",
        Status:   StatusInProgress,
        Role:     RoleAdmin,
        Priority: PriorityHigh,
    }

    fmt.Printf("User: %+v\n", user)

    // JSON marshaling
    jsonData, _ := json.Marshal(user)
    fmt.Printf("JSON: %s\n", jsonData)

    // JSON unmarshaling
    var parsed User
    json.Unmarshal(jsonData, &parsed)
    fmt.Printf("Parsed: %+v\n", parsed)

    // String conversion
    fmt.Printf("Status: %s\n", user.Status)
    fmt.Printf("Role: %s\n", user.Role)
    fmt.Printf("Priority: %s\n", user.Priority)

    // Parsing from string
    status, err := ParseStatus("pending")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Parsed status: %s\n", status)

    // Validation
    if status.IsValid() {
        fmt.Printf("Status '%s' is valid\n", status)
    }

    // Comparison
    fmt.Printf("Is status pending? %v\n", status == StatusPending)
    fmt.Printf("Is status in progress? %v\n", status == StatusInProgress)

    // Error handling
    _, err = ParseStatus("invalid_status")
    if err != nil {
        fmt.Printf("Error parsing invalid status: %v\n", err)
    }

    // All values functions
    fmt.Printf("All statuses: %v\n", AllStatuss())
    fmt.Printf("All roles: %v\n", AllRoles())
    fmt.Printf("All priorities: %v\n", AllPrioritys())

    // Enum comparison and switching
    switch user.Priority {
    case PriorityLow:
        fmt.Println("This can wait")
    case PriorityMedium:
        fmt.Println("Handle when convenient")
    case PriorityHigh:
        fmt.Println("Handle soon")
    case PriorityUrgent:
        fmt.Println("Handle immediately!")
    }
}
```

**Output:**

```
User: {ID:1 Name:Alice Status:in_progress Role:admin Priority:high}
JSON: {"id":1,"name":"Alice","status":"in_progress","role":"admin","priority":"high"}
Parsed: {ID:1 Name:Alice Status:in_progress Role:admin Priority:high}
Status: in_progress
Role: admin
Priority: high
Parsed status: pending
Status 'pending' is valid
Is status pending? true
Is status in progress? false
Error parsing invalid status: invalid Status: invalid_status
All statuses: [pending active in_progress completed cancelled]
All roles: [admin user guest moderator]
All priorities: [low medium high urgent]
Handle soon
```

## Generated Code Features

For each enum type, ggenums automatically generates:

### **Methods**
- `String() string` - Convert enum to string representation
- `IsValid() bool` - Validate if enum value is defined
- `MarshalJSON() ([]byte, error)` - JSON marshaling support
- `UnmarshalJSON([]byte) error` - JSON unmarshaling support

### **Functions**
- `ParseStatus(s string) (Status, error)` - Parse string to enum
- `AllStatuss() []Status` - Get all defined enum values

### **Variables**
- `allStatuss []Status` - Slice containing all enum values
- `statusNames map[Status]string` - Enum to string mapping
- `statusValues map[string]Status` - String to enum mapping

## Key Benefits

✅ **Type Safety**: Full compile-time type checking
✅ **IDE Support**: Complete autocomplete and refactoring support
✅ **JSON Ready**: Built-in JSON marshaling/unmarshaling
✅ **Validation**: Runtime validation with `IsValid()`
✅ **Performance**: Efficient lookup maps for conversions
✅ **Standard**: Follows Go naming conventions and patterns
✅ **Clean**: Proper separation between source and generated code

## Migration from v1.x

If you're upgrading from the comment-based approach, see [https://github.com/Gandalf-Le-Dev/ggenums/releases/tag/v2.0.0](https://github.com/Gandalf-Le-Dev/ggenums/releases/tag/v2.0.0) for a complete migration guide.

## Requirements

- Go 1.18+ (uses generics internally for better type safety)
- Enum constants must follow `TypeNameValue` pattern (e.g., `StatusPending`)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
