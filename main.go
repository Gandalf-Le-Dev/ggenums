package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Gandalf-Le-Dev/ggenums/generator"
	"github.com/Gandalf-Le-Dev/ggenums/templates"
)

func main() {
	var dir string
	var types string
	flag.StringVar(&dir, "dir", ".", "directory to scan for enum definitions")
	flag.StringVar(&types, "type", "", "comma-separated list of type names to generate enums for")
	flag.Parse()

	if types == "" {
		log.Fatal("must specify -type flag with comma-separated type names")
	}

	typeNames := strings.Split(types, ",")
	for i, name := range typeNames {
		typeNames[i] = strings.TrimSpace(name)
	}

	g := generator.NewGenerator(dir, typeNames)
	if err := g.Parse(); err != nil {
		log.Fatal(err)
	}

	if err := generate(g); err != nil {
		log.Fatal(err)
	}
}

func generate(g *generator.Generator) error {
	tmpl, err := template.New("enum").Funcs(template.FuncMap{
		"ToLower": strings.ToLower,
	}).Parse(templates.ConstEnumTemplate)
	if err != nil {
		return err
	}

	for _, enumDef := range g.Enums {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, struct {
			Package string
			Type    string
			Values  []generator.EnumValue
		}{
			Package: g.PackageName(),
			Type:    enumDef.Name,
			Values:  enumDef.Values,
		}); err != nil {
			return err
		}

		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return err
		}

		filename := filepath.Join(g.PackageDir(), fmt.Sprintf("%s_enum_generated.go", strings.ToLower(enumDef.Name)))
		if err := os.WriteFile(filename, formatted, 0644); err != nil {
			return err
		}
	}

	return nil
}