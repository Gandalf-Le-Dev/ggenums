package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Gandalf-Le-Dev/ggenums/generator"
)

func TestGenerate(t *testing.T) {
	// Create a temporary directory
	tmpDir := t.TempDir()

	// Create a test enum definition
	enumDef := `package test

//go:generate ggenums
//enum:name=Status values=pending,active,completed,in_progress
`
	err := os.WriteFile(filepath.Join(tmpDir, "status.go"), []byte(enumDef), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Run the generator
	g := generator.NewGenerator(tmpDir)
	if err := g.Parse(); err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	if err := generate(g); err != nil {
		t.Fatalf("Failed to generate: %v", err)
	}

	// Check if the generated file exists
	generatedFile := filepath.Join(tmpDir, "status_enum_generated.go")
	content, err := os.ReadFile(generatedFile)
	if err != nil {
		t.Fatalf("Failed to read generated file: %v", err)
	}

	// Basic checks on the generated content
	generatedContent := string(content)
	checks := []string{
		"package test",
		"type StatusEnum string",
		"StatusPending",
		"StatusActive",
		"StatusCompleted",
		"StatusInProgress",
		"func (e StatusEnum) IsValid() bool",
	}

	for _, check := range checks {
		if !strings.Contains(generatedContent, check) {
			t.Errorf("Generated content missing expected string: %s\n%s", check, generatedContent)
		}
	}
}
