package ggenums

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Define marks a type as an enum. Use it like:
//
//	//go:generate ggenums
//	type Status struct {
//	    enum `values:"pending,active,completed"`
//	}

// Generator handles the enum code generation
type Generator struct {
	Enums   map[string][]string
	pkgName string
	pkgDir  string
	fset    *token.FileSet
}

func NewGenerator(pkgPath string) *Generator {
	return &Generator{
		Enums:  make(map[string][]string),
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
			g.parseFile(file)
		}
	}

	return nil
}

func (g *Generator) parseFile(file *ast.File) {
	ast.Inspect(file, func(n ast.Node) bool {
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		structType, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		// Look for struct with enum field
		for _, field := range structType.Fields.List {
			if ident, ok := field.Type.(*ast.Ident); ok && ident.Name == "enum" {
				if field.Tag != nil {
					tag := strings.Trim(field.Tag.Value, "`")
					values := parseEnumTag(tag)
					if len(values) > 0 {
						g.Enums[typeSpec.Name.Name] = values
					}
				}
			}
		}

		return true
	})
}

func parseEnumTag(tag string) []string {
	const valuesPrefix = `values:"`
	if idx := strings.Index(tag, valuesPrefix); idx >= 0 {
		tag = tag[idx+len(valuesPrefix):]
		if idx = strings.Index(tag, `"`); idx >= 0 {
			val := cases.Title(language.English, cases.Compact).String(tag[:idx])
			values := val
			return strings.Split(values, ",")
		}
	}
	return nil
}
