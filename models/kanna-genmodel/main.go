package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
	"text/template"
)

var fileTemplate = template.Must(template.New("model.go").Parse(`
package {{.Package}}

import (
	"net/url"
)

type {{.Name}} struct {
	id *url.URL
	typ string
{{ range .Properties -}}
{{""}}	{{ .FieldName }} {{ .Type }}
{{ end -}}
}

func (model *{{.Name}}) ID() *url.URL {
	return model.id
}

func (model *{{.Name}}) Types() []string {
	return []string{ model.typ }
}

func (model *{{.Name}}) HasType(t string) bool {
	return t == model.typ
}

func (model *{{.Name}}) Props() []string {
	return []string{ "id", "type", {{ range .Properties -}}
		{{printf "%q" .Name }}, 
	{{- end }} }
}

func (model *{{.Name}}) GetProp(prop string) (interface{}, bool) {
	switch prop {
	case "id":
		return model.id, true
	case "type":
		return model.typ, true
{{ range .Properties -}}
{{""}}	case {{ printf "%q" .Name }}:
		return model.{{ .FieldName }}, true
{{ end -}}
{{""}}	default:
		return nil, false
	}
}
`))

func main() {
	outfile := flag.String("output", "", "an optional file in which to save the outpu")
	flag.Parse()
	if flag.NArg() != 1 {
		panic("Usage: kanna-genmodel MODELJSON")
	}
	model := parseFile(flag.Arg(0))
	var output bytes.Buffer
	if err := fileTemplate.Execute(&output, model); err != nil {
		panic(err)
	}
	if *outfile != "" {
		if err := ioutil.WriteFile(*outfile, output.Bytes(), 0664); err != nil {
			panic(err)
		}
	} else {
		fmt.Print(output.String())
	}
	//debugPrint(model)
}

type rawModel struct {
	Package    string
	Name       string
	Properties map[string]string
}

type Model struct {
	Package    string
	Name       string
	Properties []Property
}

type Property struct {
	FieldName string
	Name      string
	Type      string
}

func parseFile(filename string) *Model {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	var raw rawModel
	if err := json.Unmarshal(buf, &raw); err != nil {
		panic(err)
	}
	return processModel(raw)
}

func processModel(raw rawModel) *Model {
	props := make([]Property, 0, len(raw.Properties))
	for prop, typ := range raw.Properties {
		props = append(props, Property{
			FieldName: strings.Title(prop),
			Name:      prop,
			Type:      typ,
		})
	}
	return &Model{
		Package:    raw.Package,
		Name:       raw.Name,
		Properties: props,
	}
}

func debugPrint(m *Model) {
	fmt.Println("Properties:")
	for _, prop := range m.Properties {
		fmt.Printf("\t%q of type %s\n", prop.Name, prop.Type)
	}
}
