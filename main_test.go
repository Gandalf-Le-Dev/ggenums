package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gandalf-Le-Dev/ggenums/ggenums"
)

func TestGenerate(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a test enum definition
	enumDef := `package test

type Status struct {
    enum ` + "`values:\"pending,active,completed\"`" + `
}
`
	err := os.WriteFile(filepath.Join(tmpDir, "status.go"), []byte(enumDef), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Run the generator
	g := ggenums.NewGenerator(tmpDir)
	if err := g.Parse(); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if err := generate(g); err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Check if the generated file exists
	generatedFile := filepath.Join(tmpDir, "status_generated.go")
	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Basic checks on the generated content
	generatedContent := string(content)
	checks := []string{
		"package test",
		"type Status string",
		"StatusPending",
		"StatusActive",
		"StatusCompleted",
		"func (e Status) IsValid() bool",
	}

	for _, check := range checks {
		if !strings.Contains(generatedContent, check) {
			t.Errorf("Generated content missing expected string: %s\n%s", check, generatedContent)
		}
	}
}
