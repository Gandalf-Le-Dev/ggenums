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
	g := NewGenerator(pkgPath)
	
	require.NotNil(t, g)
	assert.Equal(t, pkgPath, g.PackageDir())
	assert.Empty(t, g.Enums)
	assert.NotNil(t, g.fset)
}

func TestGenerator_PackageName(t *testing.T) {
	g := NewGenerator("/test/path")
	g.pkgName = "testpkg"
	
	assert.Equal(t, "testpkg", g.PackageName())
}

func TestGenerator_PackageDir(t *testing.T) {
	pkgPath := "/test/path"
	g := NewGenerator(pkgPath)
	
	assert.Equal(t, pkgPath, g.PackageDir())
}

func TestTransformToPascalCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single word",
			input:    "active",
			expected: "Active",
		},
		{
			name:     "snake_case to PascalCase",
			input:    "in_progress",
			expected: "InProgress",
		},
		{
			name:     "multiple underscores",
			input:    "multi_word_example",
			expected: "MultiWordExample",
		},
		{
			name:     "already uppercase",
			input:    "PENDING",
			expected: "Pending",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "single character",
			input:    "a",
			expected: "A",
		},
		{
			name:     "with numbers",
			input:    "status_1",
			expected: "Status1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformToPascalCase(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseEnumComment(t *testing.T) {
	tests := []struct {
		name        string
		comment     string
		expected    EnumDef
		expectError bool
	}{
		{
			name:    "valid enum comment",
			comment: "//enum:name=Status values=pending,active,completed",
			expected: EnumDef{
				Name: "Status",
				Values: []EnumValue{
					{ConstantName: "Pending", StringValue: "pending"},
					{ConstantName: "Active", StringValue: "active"},
					{ConstantName: "Completed", StringValue: "completed"},
				},
			},
			expectError: false,
		},
		{
			name:    "valid enum with snake_case values",
			comment: "//enum:name=TaskStatus values=in_progress,not_started,completed",
			expected: EnumDef{
				Name: "TaskStatus",
				Values: []EnumValue{
					{ConstantName: "InProgress", StringValue: "in_progress"},
					{ConstantName: "NotStarted", StringValue: "not_started"},
					{ConstantName: "Completed", StringValue: "completed"},
				},
			},
			expectError: false,
		},
		{
			name:    "different order of parameters",
			comment: "//enum:values=active,inactive name=State",
			expected: EnumDef{
				Name: "State",
				Values: []EnumValue{
					{ConstantName: "Active", StringValue: "active"},
					{ConstantName: "Inactive", StringValue: "inactive"},
				},
			},
			expectError: false,
		},
		{
			name:        "missing name parameter",
			comment:     "//enum:values=active,inactive",
			expected:    EnumDef{},
			expectError: true,
		},
		{
			name:        "missing values parameter",
			comment:     "//enum:name=Status",
			expected:    EnumDef{},
			expectError: true,
		},
		{
			name:        "empty name",
			comment:     "//enum:name= values=active,inactive",
			expected:    EnumDef{},
			expectError: true,
		},
		{
			name:        "empty values",
			comment:     "//enum:name=Status values=",
			expected:    EnumDef{},
			expectError: true,
		},
		{
			name:        "malformed comment",
			comment:     "//enum:invalid",
			expected:    EnumDef{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseEnumComment(tt.comment)
			
			if tt.expectError {
				require.Error(t, err)
				return
			}
			
			require.NoError(t, err)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, len(tt.expected.Values), len(result.Values))
			
			for i, expectedValue := range tt.expected.Values {
				assert.Equal(t, expectedValue.ConstantName, result.Values[i].ConstantName)
				assert.Equal(t, expectedValue.StringValue, result.Values[i].StringValue)
			}
		})
	}
}

func TestGenerator_Parse(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	// Create a test Go file with enum comment
	testFile := `package testpkg

// Some regular comment
//enum:name=Status values=pending,active,completed

type SomeStruct struct {
	Field string
}
`
	
	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(testFile), 0644)
	require.NoError(t, err)
	
	// Test parsing
	g := NewGenerator(tmpDir)
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

//enum:name=Status values=pending,active,completed

type SomeStruct struct {
	Field string
}
`
	
	// Create second test file
	testFile2 := `package testpkg

//enum:name=Priority values=low,medium,high

type AnotherStruct struct {
	Field string
}
`
	
	err := os.WriteFile(filepath.Join(tmpDir, "test1.go"), []byte(testFile1), 0644)
	require.NoError(t, err)
	
	err = os.WriteFile(filepath.Join(tmpDir, "test2.go"), []byte(testFile2), 0644)
	require.NoError(t, err)
	
	// Test parsing
	g := NewGenerator(tmpDir)
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

func TestGenerator_Parse_NoEnums(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	// Create a test Go file without enum comments
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
	
	// Test parsing
	g := NewGenerator(tmpDir)
	err = g.Parse()
	require.NoError(t, err)
	
	// Verify results
	assert.Equal(t, "testpkg", g.PackageName())
	assert.Empty(t, g.Enums)
}

func TestGenerator_Parse_InvalidDirectory(t *testing.T) {
	g := NewGenerator("/nonexistent/directory")
	err := g.Parse()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse package")
}

func TestGenerator_Parse_InvalidEnumComment(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir := t.TempDir()
	
	// Create a test Go file with invalid enum comment
	testFile := `package testpkg

//enum:name=Status
// Missing values parameter should cause error

type SomeStruct struct {
	Field string
}
`
	
	err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(testFile), 0644)
	require.NoError(t, err)
	
	// Test parsing
	g := NewGenerator(tmpDir)
	err = g.Parse()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "enum values not specified")
}