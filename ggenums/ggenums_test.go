package ggenums

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGenerator_Parse(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create test file
	testFile := `package testpkg

type Status struct {
    enum ` + "`values:\"pending,active,completed\"`" + `
}

type Role struct {
    enum ` + "`values:\"admin,user,guest\"`" + `
}

// Should be ignored
type NotAnEnum struct {
    Field string
}

// Should be ignored - invalid tag
type InvalidEnum struct {
    enum ` + "`invalid:\"tag\"`" + `
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "enums.go"), []byte(testFile), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	g := NewGenerator(tmpDir)
	if err := g.Parse(); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	tests := []struct {
		name        string
		enumName    string
		wantValues  []string
		shouldExist bool
	}{
		{
			name:        "Status enum",
			enumName:    "Status",
			wantValues:  []string{"Pending", "Active", "Completed"},
			shouldExist: true,
		},
		{
			name:        "Role enum",
			enumName:    "Role",
			wantValues:  []string{"Admin", "User", "Guest"},
			shouldExist: true,
		},
		{
			name:        "Not an enum",
			enumName:    "NotAnEnum",
			shouldExist: false,
		},
		{
			name:        "Invalid enum",
			enumName:    "InvalidEnum",
			shouldExist: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			values, exists := g.Enums[tt.enumName]
			if exists != tt.shouldExist {
				t.Errorf("enum %s existence = %v, want %v", tt.enumName, exists, tt.shouldExist)
				return
			}
			if !tt.shouldExist {
				return
			}
			if !reflect.DeepEqual(values, tt.wantValues) {
				t.Errorf("enum %s values = %v, want %v", tt.enumName, values, tt.wantValues)
			}
		})
	}
}

func TestParseEnumTag(t *testing.T) {
	tests := []struct {
		name    string
		tag     string
		want    []string
		wantNil bool
	}{
		{
			name: "Simple values",
			tag:  `values:"one,two,three"`,
			want: []string{"one", "two", "three"},
		},
		{
			name: "Single value",
			tag:  `values:"single"`,
			want: []string{"single"},
		},
		{
			name:    "Invalid tag",
			tag:     `invalid:"tag"`,
			wantNil: true,
		},
		{
			name:    "Empty tag",
			tag:     ``,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseEnumTag(tt.tag)
			if tt.wantNil {
				if got != nil {
					t.Errorf("parseEnumTag() = %v, want nil", got)
				}
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseEnumTag() = %v, want %v", got, tt.want)
			}
		})
	}
}
