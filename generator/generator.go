package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
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
	Enums     []EnumDef
	pkgName   string
	pkgDir    string
	fset      *token.FileSet
	typeNames []string
}

func NewGenerator(pkgPath string, typeNames []string) *Generator {
	return &Generator{
		Enums:     []EnumDef{},
		fset:      token.NewFileSet(),
		pkgDir:    pkgPath,
		typeNames: typeNames,
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
	// Find type declarations for our target types
	typeDecls := make(map[string]*ast.TypeSpec)
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					for _, typeName := range g.typeNames {
						if typeSpec.Name.Name == typeName {
							typeDecls[typeName] = typeSpec
						}
					}
				}
			}
		}
	}

	// Find const declarations that follow the pattern
	for _, decl := range file.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for typeName, typeSpec := range typeDecls {
				enumDef, err := g.parseConstBlock(genDecl, typeName, typeSpec)
				if err != nil {
					return err
				}
				if enumDef != nil {
					g.Enums = append(g.Enums, *enumDef)
				}
			}
		}
	}

	return nil
}

func (g *Generator) parseConstBlock(genDecl *ast.GenDecl, typeName string, typeSpec *ast.TypeSpec) (*EnumDef, error) {
	var values []EnumValue

	for _, spec := range genDecl.Specs {
		if valueSpec, ok := spec.(*ast.ValueSpec); ok {
			// Check if this const belongs to our type
			if valueSpec.Type != nil {
				if ident, ok := valueSpec.Type.(*ast.Ident); ok && ident.Name == typeName {
					// This const is explicitly typed with our enum type
					for _, name := range valueSpec.Names {
						constName := name.Name
						if strings.HasPrefix(constName, typeName) {
							// Extract the suffix after the type name
							suffix := strings.TrimPrefix(constName, typeName)
							if suffix != "" {
								stringValue := camelToSnake(suffix)
								values = append(values, EnumValue{
									ConstantName: suffix,
									StringValue:  stringValue,
								})
							}
						}
					}
				}
			} else {
				// Check if this is an iota-style const that might belong to our type
				for _, name := range valueSpec.Names {
					constName := name.Name
					if strings.HasPrefix(constName, typeName) {
						// Extract the suffix after the type name
						suffix := strings.TrimPrefix(constName, typeName)
						if suffix != "" {
							stringValue := camelToSnake(suffix)
							values = append(values, EnumValue{
								ConstantName: suffix,
								StringValue:  stringValue,
							})
						}
					}
				}
			}
		}
	}

	if len(values) == 0 {
		return nil, nil
	}

	return &EnumDef{
		Name:   typeName,
		Values: values,
	}, nil
}

// camelToSnake converts CamelCase to snake_case
func camelToSnake(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		if r >= 'A' && r <= 'Z' {
			result = append(result, r-'A'+'a')
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
