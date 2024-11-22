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

	"github.com/Gandalf-Le-Dev/ggenums/ggenums"
	"github.com/Gandalf-Le-Dev/ggenums/templates"
)

func main() {
	var dir string
	flag.StringVar(&dir, "dir", ".", "directory to scan for enum definitions")
	flag.Parse()

	g := ggenums.NewGenerator(dir)
	if err := g.Parse(); err != nil {
		log.Fatal(err)
	}

	if err := generate(g); err != nil {
		log.Fatal(err)
	}
}

func generate(g *ggenums.Generator) error {
	tmpl, err := template.New("enum").Parse(templates.EnumTemplate)
	if err != nil {
		return err
	}

	for typeName, values := range g.Enums {
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, struct {
			Package string
			Type    string
			Values  []string
		}{
			Package: g.PackageName(),
			Type:    typeName,
			Values:  values,
		}); err != nil {
			return err
		}

		formatted, err := format.Source(buf.Bytes())
		if err != nil {
			return err
		}

		filename := filepath.Join(g.PackageDir(), fmt.Sprintf("%s_generated.go", strings.ToLower(typeName)))
		if err := os.WriteFile(filename, formatted, 0644); err != nil {
			return err
		}
	}

	return nil
}
