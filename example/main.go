// examples/basic/main.go
package main

import (
	"encoding/json"
	"fmt"

	"github.com/Gandalf-Le-Dev/ggenums/example/enums"
)

// Example struct using the generated enums
type User struct {
	ID     string           `json:"id"`
	Status enums.StatusEnum `json:"status"`
	Role   enums.RoleEnum   `json:"role"`
}

func main() {
	// JSON marshaling
	user := User{
		ID:     "123",
		Status: enums.StatusPending,
		Role:   enums.RoleAdmin,
	}

	jsonData, _ := json.Marshal(user)
	fmt.Printf("JSON: %s\n", jsonData)

	// JSON unmarshaling
	var parsed User
	json.Unmarshal(jsonData, &parsed)
	fmt.Printf("Parsed: %+v\n", parsed)

	// Validation
	fmt.Printf("Is valid status? %v\n", enums.StatusInProgress.IsValid())

	// Parsing from string
	status, _ := enums.ParseStatus("in_progress")
	fmt.Printf("Parsed status: %s\n", status)

	// Comparison
	fmt.Printf("Is status in progress? %v\n", status == enums.StatusInProgress)

	// Error handling
	_, err := enums.ParseStatus("error_status")
	fmt.Printf("Invalid status error: %v\n", err)
}
