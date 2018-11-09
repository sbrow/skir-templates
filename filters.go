package main

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
)

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
	case map[interface{}]interface{}:
		for key, value := range t {
			switch tk := key.(type) {
			case string:
				switch tv := value.(type) {
				case nil:
					fmt.Fprintf(out, "%s %s\n", strings.Repeat("#", h), tk)
				case string:
					fmt.Fprintf(out, "* **%s**: %s\n", tk, tv)
				default:
					fmt.Fprintf(out, "%s %s\n%s", strings.Repeat("#", h), tk, Markdown(value, h+1))
				}
			case int:
				fmt.Fprintf(out, "%d. %v\n", tk, value)
			}
		}
	default:
		fmt.Fprintln(out, t)
	}
	return out.String()
}

// Type returns the name of an interface's type.
func Type(v interface{}) string {
	return fmt.Sprint(reflect.TypeOf(v))
}
