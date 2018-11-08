// Command skir-templates fills out liquid templates.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

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

// Markdown returns the markdown representation of v.
//
// header determines which header to use at the base level (default 1).
func Markdown(v interface{}, header ...int) string {
	if len(header) == 0 || header[0] < 1 {
		header = []int{1}
	}
	h := header[0]
	out := new(bytes.Buffer)
	switch t := v.(type) {
	case []interface{}:
		for _, value := range t {
			fmt.Fprint(out, Markdown(value, h))
		}
	case map[string]interface{}:
		for key, value := range t {
			fmt.Fprintf(out, "%s %s\n%s", strings.Repeat("#", h), key, Markdown(value, h))
		}
	case map[interface{}]interface{}:
		var frmt string
		for key, value := range t {
			switch tv := value.(type) {
			case string:
				switch key.(type) {
				case string:
					frmt = "* **%s**: %s"
				case int:
					frmt = "%d. %s"
				}
				fmt.Fprintf(out, frmt+"\n", key, value)
			case []interface{}:
				fmt.Fprintf(out, "%s %s\n%s", strings.Repeat("#", h+1), key, Markdown(tv, h))
				// TODO: Closest, but still not good enough
				fmt.Fprintln(out)
			}
		}
	}
	return out.String()
}

// Type returns the name of an interface's type.
func Type(v interface{}) string {
	return fmt.Sprint(reflect.TypeOf(v))
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
