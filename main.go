// Command skir-templates fills out liquid templates.
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"regexp"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/osteele/liquid"
)

// Wiki represents a wiki page.
type Wiki struct {
	Title    string              // Title holds the contents of Header 1.
	ToC      bool                // ToC is whether to include a Table of Contents.
	Imports  []map[string]string // Imports says whether or not to import data.
	Contents interface{}         // Contents holds the Page contents.
}

// Compile compiles a yaml document into an "order matters" yaml document.
func Compile(yaml []byte) []byte {
	ret := regexp.MustCompile(`(?m)^(\s*)([-])`).ReplaceAll(yaml, []byte(`  $1$2`))
	ret = regexp.MustCompile(`(?m)^(\s*)(\w)`).ReplaceAll(ret, []byte(`$1- $2`))
	return ret
}

// Import imports .yml files and maps them to Contents
// as laid out in w.Imports.
func (w *Wiki) Import() error {
	cont := []map[string]interface{}{}
	for _, v := range w.Imports {
		for mnt, file := range v {
			data, err := ioutil.ReadFile(file)
			HandleError(err)
			var obj interface{}
			HandleError(yaml.Unmarshal(Compile(data), &obj))
			cont = append(cont, map[string]interface{}{
				mnt: obj,
			})
		}
	}
	w.Contents = cont
	return nil
}

// Binding returns the Wiki converted into map[string]interface{}
// that liquid.ParseAndRenderString can understand.
func (w *Wiki) Binding() map[string]interface{} {
	return map[string]interface{}{
		"title":    w.Title,
		"toc":      w.ToC,
		"contents": w.Contents,
	}
}

// HandleError panics if e is not nil.
func HandleError(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Use Args[1] as our template.
	Input, err := ioutil.ReadFile(os.Args[1])
	HandleError(err)

	// ERROR: Assumes we are using wiki.md
	var bindings Wiki

	// Separate out Front Matter
	s := regexp.MustCompile(`---`).Split(string(Input), 3)
	data := strings.TrimSpace(s[1])

	// Unmarshal Front Matter.
	HandleError(yaml.Unmarshal([]byte(data), &bindings))

	// Import external YAML files.
	bindings.Import()

	// Parse and fill out our template.
	engine := liquid.NewEngine()
	engine.RegisterFilter("type", Type)
	engine.RegisterFilter("md", Markdown)

	tmpl := strings.TrimSpace(s[2])
	out, err := engine.ParseAndRenderString(tmpl, bindings.Binding())
	HandleError(err)

	// fmt.Println(strings.TrimSpace(out))
	fmt.Println(out)
}
