package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/Gandalf-Le-Dev/ggenums/templates"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ggenums",
	Short: "Generate code for enums from a JSON file",
	Long:  `Generate structs and code for enums from a JSON file`,
	Run:   GenerateEnums,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

type config struct {
	Package string `json:"package"`
	Enums   []enum `json:"enums"`
}

type enum struct {
	Name   string            `json:"name"`
	Plural string            `json:"plural"`
	Values map[string]string `json:"values"`
}

var cfg config

func init() {
	bytes, err := os.ReadFile("./enums.json")
	cobra.CheckErr(err)

	err = json.Unmarshal(bytes, &cfg)
	cobra.CheckErr(err)
}

type tmplValues struct {
	Values       map[string]string
	Type         string
	TypePlural   string
	TypePluralLC string
	Package      string
}

// GenerateEnums generates the code for the enums
func GenerateEnums(cmd *cobra.Command, args []string) {
	flag.Parse()

	for _, enum := range cfg.Enums {
		filename := fmt.Sprintf("%s_generated.go", strings.ToLower(enum.Name[0:1])+enum.Name[1:])

		file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		cobra.CheckErr(err)

		defer file.Close()

		t := tmplValues{
			Values:       enum.Values,
			Type:         enum.Name,
			TypePlural:   enum.Plural,
			TypePluralLC: strings.ToLower(enum.Plural[0:1]) + enum.Plural[1:], // Lowercase the first letter
			Package:      cfg.Package,
		}

		templ, err := template.New("enum_template").Parse(templates.EnumTemplate)
		if err != nil {
			log.Fatal(err)
		}

		err = templ.Execute(file, t)
		if err != nil {
			log.Fatal(err)
		}

		cmd := exec.Command("goimports", "-local", "github.com/Gandalf-Le-Dev", "-w", filename)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}
}
