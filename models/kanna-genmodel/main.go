package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
	"text/template"
)

var fileTemplate = template.Must(template.New("model.go").Parse(`
package {{.Package}}

import (
	{{ if .Table }}
	"context"
	"database/sql"
	{{ end }}
	"net/url"
	{{ if .Table }}
	"github.com/ekiru/kanna/db"
	{{ end }}
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

{{ if .Table -}}
func {{.Name}}ById(ctx context.Context, id string) (*{{.Name}}, error) {
	var model {{.Name}}
	rows, err := db.DB(ctx).QueryContext(ctx, "select id, type {{- range .Properties -}}
		, {{ .ColumnName }}
	{{- end }} from {{.Table}} where id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

	err = rows.Scan(
		db.URLScanner{ &model.id },
		&model.typ,
{{- range .Properties -}}
{{- if eq .Type "*url.URL" }}
		db.URLScanner{ &model.{{.FieldName}} },
{{- else }}
		&model.{{.FieldName}},
{{- end -}}
{{- end }}
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
{{- end }}
`))

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		panic("Usage: kanna-genmodel MODELJSON")
	}
	models := parseFile(flag.Arg(0))
	for _, model := range models {
		var output bytes.Buffer
		if err := fileTemplate.Execute(&output, model); err != nil {
			panic(err)
		}
		if err := ioutil.WriteFile(model.File, output.Bytes(), 0664); err != nil {
			panic(err)
		}
	}
	//debugPrint(model)
}

type rawModel struct {
	Package string
	Types   []struct {
		File       string
		Name       string
		Table      string
		Properties map[string]string
	}
}

type Model struct {
	Package    string
	File       string
	Name       string
	Table      string
	Properties []Property
}

type Property struct {
	FieldName  string
	Name       string
	ColumnName string
	Type       string
}

func parseFile(filename string) []*Model {
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

func processModel(raws rawModel) []*Model {
	var models []*Model
	for _, raw := range raws.Types {
		props := make([]Property, 0, len(raw.Properties))
		for prop, typ := range raw.Properties {
			props = append(props, Property{
				FieldName:  strings.Title(prop),
				Name:       prop,
				ColumnName: prop,
				Type:       typ,
			})
		}
		sort.Slice(props, func(i, j int) bool {
			return props[i].Name < props[j].Name
		})
		models = append(models, &Model{
			Package:    raws.Package,
			File:       raw.File,
			Name:       raw.Name,
			Table:      raw.Table,
			Properties: props,
		})
	}
	return models
}

func debugPrint(m *Model) {
	fmt.Println("Properties:")
	for _, prop := range m.Properties {
		fmt.Printf("\t%q of type %s\n", prop.Name, prop.Type)
	}
}
