package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	pkgPath := "/test/path"
	typeNames := []string{"Status", "Priority"}
	g := NewGenerator(pkgPath, typeNames)

	require.NotNil(t, g)
	assert.Equal(t, pkgPath, g.PackageDir())
	assert.Equal(t, typeNames, g.typeNames)
	assert.Empty(t, g.Enums)
	assert.NotNil(t, g.fset)
}

func TestGenerator_PackageName(t *testing.T) {
	g := NewGenerator("/test/path", []string{"Status"})
	g.pkgName = "testpkg"

	assert.Equal(t, "testpkg", g.PackageName())
}

func TestGenerator_PackageDir(t *testing.T) {
	pkgPath := "/test/path"
	g := NewGenerator(pkgPath, []string{"Status"})

	assert.Equal(t, pkgPath, g.PackageDir())
}

func TestCamelToSnake(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single word",
			input:    "Active",
			expected: "active",
		},
		{
			name:     "CamelCase to snake_case",
			input:    "InProgress",
			expected: "in_progress",
		},
		{
			name:     "multiple words",
			input:    "MultiWordExample",
			expected: "multi_word_example",
		},
		{
			name:     "single character",
			input:    "A",
			expected: "a",
		},
		{
			name:     "with consecutive caps",
			input:    "XMLParser",
			expected: "x_m_l_parser",
		},
		{
			name:     "with numbers",
			input:    "Status1",
			expected: "status1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := camelToSnake(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}


func TestGenerator_Parse(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create a test Go file with enum type and const declarations
	testFile := `package testpkg

//go:generate ggenums -type=Status
type Status int
const (
	StatusPending Status = iota
	StatusActive
	StatusCompleted
)

type SomeStruct struct {
	Field string
}
`

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(testFile), 0644)
	require.NoError(t, err)

	// Test parsing
	g := NewGenerator(tmpDir, []string{"Status"})
	err = g.Parse()
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", g.PackageName())
	assert.Len(t, g.Enums, 1)

	enum := g.Enums[0]
	assert.Equal(t, "Status", enum.Name)
	assert.Len(t, enum.Values, 3)

	expectedValues := []EnumValue{
		{ConstantName: "Pending", StringValue: "pending"},
		{ConstantName: "Active", StringValue: "active"},
		{ConstantName: "Completed", StringValue: "completed"},
	}

	for i, expected := range expectedValues {
		assert.Equal(t, expected.ConstantName, enum.Values[i].ConstantName)
		assert.Equal(t, expected.StringValue, enum.Values[i].StringValue)
	}
}

func TestGenerator_Parse_MultipleFiles(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create first test file
	testFile1 := `package testpkg

//go:generate ggenums -type=Status
type Status int
const (
	StatusPending Status = iota
	StatusActive
	StatusCompleted
)

type SomeStruct struct {
	Field string
}
`

	// Create second test file
	testFile2 := `package testpkg

//go:generate ggenums -type=Priority
type Priority int
const (
	PriorityLow Priority = iota
	PriorityMedium
	PriorityHigh
)

type AnotherStruct struct {
	Field string
}
`

	err := os.WriteFile(filepath.Join(tmpDir, "test1.go"), []byte(testFile1), 0644)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(tmpDir, "test2.go"), []byte(testFile2), 0644)
	require.NoError(t, err)

	// Test parsing
	g := NewGenerator(tmpDir, []string{"Status", "Priority"})
	err = g.Parse()
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", g.PackageName())
	assert.Len(t, g.Enums, 2)

	// Find the enums (order may vary based on file system)
	var statusEnum, priorityEnum *EnumDef
	for i := range g.Enums {
		if g.Enums[i].Name == "Status" {
			statusEnum = &g.Enums[i]
		} else if g.Enums[i].Name == "Priority" {
			priorityEnum = &g.Enums[i]
		}
	}

	// Check Status enum
	require.NotNil(t, statusEnum)
	assert.Equal(t, "Status", statusEnum.Name)
	assert.Len(t, statusEnum.Values, 3)

	// Check Priority enum
	require.NotNil(t, priorityEnum)
	assert.Equal(t, "Priority", priorityEnum.Name)
	assert.Len(t, priorityEnum.Values, 3)
	assert.Equal(t, "Low", priorityEnum.Values[0].ConstantName)
	assert.Equal(t, "Medium", priorityEnum.Values[1].ConstantName)
	assert.Equal(t, "High", priorityEnum.Values[2].ConstantName)
}

func TestGenerator_Parse_NoMatchingTypes(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()

	// Create a test Go file without matching types
	testFile := `package testpkg

type SomeStruct struct {
	Field string
}

// Just a regular comment
func someFunction() {
	// Another regular comment
}
`

	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(testFile), 0644)
	require.NoError(t, err)

	// Test parsing - looking for types that don't exist
	g := NewGenerator(tmpDir, []string{"Status", "Priority"})
	err = g.Parse()
	require.NoError(t, err)

	// Verify results
	assert.Equal(t, "testpkg", g.PackageName())
	assert.Empty(t, g.Enums)
}

func TestGenerator_Parse_InvalidDirectory(t *testing.T) {
	g := NewGenerator("/nonexistent/directory", []string{"Status"})
	err := g.Parse()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse package")
}