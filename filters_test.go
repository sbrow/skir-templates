package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/go-yaml/yaml"
)

func TestMarkdown(t *testing.T) {
	type args struct {
		v      interface{}
		header []int
	}
	type test struct {
		name string
		args args
		want string
	}

	yml := []byte(`
string: &str Hello, World!

multi-string:
- *str
- *str

int: &in 12

float: &flo 12.3456789

\[\]interface{}:
- *str
- *in
- *flo

header:
  Heading 1:

list-header:
- Heading 1:

sub-header:
  Heading 1:
    Heading 2:

list-sub-header:
- Heading 1:
  - Heading 2:

header-str:
  Heading 1:
  - Hello, World!

list-item:
  Heading 1:
  - Hello: World

ordered-item:
  Heading 1:
  - 1: Hello, World
`)
	var ts map[string]interface{}
	if err := yaml.Unmarshal(yml, &ts); err != nil {
		panic(err)
	}
	tests := []test{}
	Add := func(arg string, h []int, want string) {
		t := test{arg, args{ts[arg], h}, want}
		tests = append(tests, t)
	}
	Add("string", nil, fmt.Sprintln(ts["string"]))
	Add("multi-string", nil, strings.Repeat(tests[0].want, 2))
	Add("int", nil, fmt.Sprintln(ts["int"]))
	Add("float", nil, fmt.Sprintln(ts["float"]))
	Add(`\[\]interface{}`, nil, fmt.Sprintf("%s\n%d\n%0.7f\n", ts["string"], ts["int"], ts["float"]))
	Add("header", nil, "# Heading 1\n")
	Add("list-header", nil, "# Heading 1\n")
	Add("header", []int{0}, "# Heading 1\n")
	Add("header", []int{1}, "# Heading 1\n")
	Add("header", []int{2}, "## Heading 1\n")
	Add("sub-header", nil, "# Heading 1\n## Heading 2\n")
	Add("list-sub-header", nil, "# Heading 1\n## Heading 2\n")
	Add("sub-header", []int{2}, "## Heading 1\n### Heading 2\n")
	Add("header-str", nil, "# Heading 1\nHello, World!\n")
	Add("list-item", nil, "# Heading 1\n* **Hello**: World\n")
	Add("ordered-item", nil, "# Heading 1\n1. Hello, World\n")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Markdown(tt.args.v, tt.args.header...); got != tt.want {
				t.Errorf("Markdown() = \r\n\"%v\", want \r\n\"%v\"", got, tt.want)
			}
		})
	}
}
