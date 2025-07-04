# GGenums - Generate Golang Enums

[![Go Report Card](https://goreportcard.com/badge/github.com/Gandalf-Le-Dev/ggenums)](https://goreportcard.com/report/github.com/Gandalf-Le-Dev/ggenums)
[![GoDoc](https://godoc.org/github.com/Gandalf-Le-Dev/ggenums?status.svg)](https://godoc.org/github.com/Gandalf-Le-Dev/ggenums)

GGenums is a code generation tool that ease use of enums in Go. It automatically generates type-safe enums with built-in string conversion, JSON marshaling/unmarshaling, and validation.

## Features

- 🚀 Simple enum declaration using comments
- 🔒 Type-safe enum values
- 📦 JSON marshaling/unmarshaling support
- ✅ Built-in validation with IsValid()
- 🔄 String conversion methods
- 📝 Easy to read generated code

## Installation

```bash
go get -tool github.com/Gandalf-Le-Dev/ggenums@latest
```

## Quick Start

1. Define your enum using a comment:

    ```go
    //enum:name=Status values=pending,active,in_progress,completed
    ```

2. Run the generator:

    ```go
	//go:generate go tool ggenums
    ```

3. Use your generated enum:

    See the example below

## Example

```go
package main

import (
    "encoding/json"
    "fmt"
)

//go:generate go tool ggenums
//enum:name=Status values=pending,active,in_progress,completed
//enum:name=Role values=admin,user,guest

type User struct {
    ID     string `json:"id"`
    Status Status `json:"status"`
    Role   Role   `json:"role"`
}

func main() {
	var err error
	// JSON marshaling
	user := User{
		ID:     "123",
		Status: StatusPending,
		Role:   RoleAdmin,
	}

	jsonData, _ := json.Marshal(user)
	fmt.Printf("JSON: %s\n", jsonData)

	// JSON unmarshaling
	var parsed User
	json.Unmarshal(jsonData, &parsed)
	fmt.Printf("Parsed: %+v\n", parsed)

	// Parsing from string
	status, _ := ParseStatus("in_progress")
	fmt.Printf("Parsed status: %s\n", status)

	// Validation
	if status.IsValid() {
		fmt.Printf("Status '%s' is valid\n", status)
	} else {
		fmt.Printf("Status '%s' is invalid\n", status)
	}
	// Comparison
	fmt.Printf("Is status in progress? %v\n", status == StatusInProgress)

	// Error handling
	_, err = ParseStatus("error_status")
	fmt.Printf("Invalid status error: %v\n", err)
}
```

Output:

```bash
JSON: {"id":"123","status":"pending","role":"admin"}
Parsed: {ID:123 Status:pending Role:admin}
Parsed status: in_progress
Status 'in_progress' is valid
Is status in progress? true
Invalid status error: invalid Status: error_status
```

## Generated Code Features

The generated code includes:

- Type-safe enum values
- String conversion methods
- JSON marshaling/unmarshaling
- IsValid() validation method
- Parsing functions
- All[EnumName]s() function to get all values
- Lookup maps for efficient conversion

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
