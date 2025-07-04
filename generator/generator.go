package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type EnumValue struct {
	ConstantName string // The PascalCase name used for the constant (e.g., "InProgress")
	StringValue  string // The original string value (e.g., "in_progress")
}

type EnumDef struct {
	Name   string
	Values []EnumValue
}

// Generator handles the enum code generation
type Generator struct {
	Enums   []EnumDef
	pkgName string
	pkgDir  string
	fset    *token.FileSet
}

func NewGenerator(pkgPath string) *Generator {
	return &Generator{
		Enums:  []EnumDef{},
		fset:   token.NewFileSet(),
		pkgDir: pkgPath,
	}
}

func (g *Generator) PackageName() string {
	return g.pkgName
}

func (g *Generator) PackageDir() string {
	return g.pkgDir
}

func (g *Generator) Parse() error {
	pkgs, err := parser.ParseDir(g.fset, g.pkgDir, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse package: %w", err)
	}

	for pkgName, pkg := range pkgs {
		g.pkgName = pkgName
		for _, file := range pkg.Files {
			err := g.parseFile(file)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *Generator) parseFile(file *ast.File) error {
	// Look for //enum: comments in the file
	for _, commentGroup := range file.Comments {
		for _, comment := range commentGroup.List {
			if strings.HasPrefix(comment.Text, "//enum:") {
				def, err := parseEnumComment(comment.Text)
				if err != nil {
					return err
				}
				g.Enums = append(g.Enums, def)
			}
		}
	}
	return nil
}

// transformToPascalCase converts a snake_case string to PascalCase
func transformToPascalCase(input string) string {
	words := strings.Split(input, "_")
	for i, word := range words {
		words[i] = cases.Title(language.English, cases.Compact).String(word)
	}
	return strings.Join(words, "")
}

func parseEnumComment(comment string) (EnumDef, error) {
	// Remove //enum: prefix
	content := strings.TrimPrefix(comment, "//enum:")

	// Parse name and values
	parts := strings.Split(content, " ")
	var name, valuesStr string

	for _, part := range parts {
		if after, ok := strings.CutPrefix(part, "name="); ok {
			name = after
		} else if after, ok := strings.CutPrefix(part, "values="); ok {
			valuesStr = after
		}
	}

	if name == "" {
		return EnumDef{}, fmt.Errorf("enum name not specified")
	}
	if valuesStr == "" {
		return EnumDef{}, fmt.Errorf("enum values not specified")
	}

	// Split values and create enumValue structs
	valuesList := strings.Split(valuesStr, ",")
	values := make([]EnumValue, len(valuesList))
	for i, value := range valuesList {
		values[i] = EnumValue{
			ConstantName: transformToPascalCase(value),
			StringValue:  value, // Keep the original string value
		}
	}

	return EnumDef{
		Name:   name,
		Values: values,
	}, nil
}
