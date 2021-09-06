package main

import (
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const viewsDir = "views/"

var views = template.New("")

func init() {
	parseViewsDir("*.html")
}

func parseViewsDir(path string) {
	vs, err := filepath.Glob(viewsDir + path)
	if err != nil {
		panic("Glob Views:" + err.Error())
	}

	for _, v := range vs {
		name := strings.TrimPrefix(v, viewsDir)
		name = strings.TrimSuffix(name, ".html")
		log.Printf("parsing: %s", name)

		c, err := os.ReadFile(v)
		if err != nil {
			panic(err)
		}

		views.AddParseTree(name, template.Must(template.New(name).Parse(string(c))).Tree)
	}
}
