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

var fileTemplate = template.Must(template.New("model.go").Option("missingkey=error").Parse(`
{{- with $model := . -}}
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
	rows, err := db.DB(ctx).QueryContext(ctx, "select {{ $model.Table }}.id, {{ $model.Table }}.type {{- range .Properties -}}
		, {{ $model.Table }}.{{ .ColumnName }}
	{{- end -}}
	{{- range $join := .Joins -}}
		, {{ $join.Model.Table }}.type
		{{- range $join.Model.Properties -}}
			, {{ $join.Model.Table }}.{{ .ColumnName }}
		{{- end -}}
	{{- end }} from {{.Table}} {{- range .Joins -}}
		{{- ""}} join {{ .Model.Table }} on {{ $model.Table }}.{{ .LinkColumn }} = {{ .Model.Table }}.id
	{{- end }} where {{ .Table }}.id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}

{{- range .Joins }}
	model.{{ .LinkField }} = new({{ .Model.Name }})
{{- end }}

	err = rows.Scan(
		db.URLScanner{ &model.id },
		&model.typ,

{{- range .Properties -}}
{{- if .LinksTo }}
		db.URLScanner{ &model.{{.FieldName}}.id },
{{- else if eq .Type "*url.URL" }}
		db.URLScanner{ &model.{{.FieldName}} },
{{- else }}
		&model.{{.FieldName}},
{{- end -}}
{{- end -}}

{{- range $join := .Joins }}
		&model.{{$join.LinkField}}.typ,
{{- range $join.Model.Properties -}}
{{- if eq .Type "*url.URL" }}
		db.URLScanner{ &model.{{$join.LinkField}}.{{.FieldName}} },
{{- else }}
		&model.{{$join.LinkField}}.{{.FieldName}},
{{- end -}}
{{- end -}}
{{- end }}
	)
	if err != nil {
		return nil, err
	}

	return &model, nil
}
{{- end }}

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
}

type rawModel struct {
	Package string
	Types   []struct {
		File       string
		Name       string
		Table      string
		Properties map[string]interface{}
	}
}

type Model struct {
	Package    string
	File       string
	Name       string
	Table      string
	Properties []Property
	Joins      []ModelJoin
}

type Property struct {
	FieldName  string
	Name       string
	ColumnName string
	Type       string
	LinksTo    string
}

type ModelJoin struct {
	LinkColumn string
	LinkField  string
	Model      *Model
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
	modelMap := make(map[string]*Model)
	for _, raw := range raws.Types {
		props := make([]Property, 0, len(raw.Properties))
		for name, desc := range raw.Properties {
			var prop Property
			prop.Name = name
			prop.ColumnName = prop.Name
			prop.FieldName = strings.Title(prop.Name)
			switch desc := desc.(type) {
			case string:
				prop.FieldName = strings.Title(name)
				prop.Name = name
				prop.ColumnName = name
				prop.Type = desc
			case map[string]interface{}:
				if col, found := desc["column_name"]; found {
					prop.ColumnName = col.(string)
				}
				if linksTo, found := desc["links_to"]; found {
					prop.LinksTo = linksTo.(string)
				}
				prop.Type = desc["type"].(string)
			default:
				panic("invalid property descriptor")
			}
			props = append(props, prop)
		}
		sort.Slice(props, func(i, j int) bool {
			return props[i].Name < props[j].Name
		})
		model := &Model{
			Package:    raws.Package,
			File:       raw.File,
			Name:       raw.Name,
			Table:      raw.Table,
			Properties: props,
			Joins:      nil,
		}
		modelMap[model.Name] = model
		models = append(models, model)
	}
	for _, model := range models {
		for _, prop := range model.Properties {
			if prop.LinksTo != "" {
				linkedModel := modelMap[prop.LinksTo]
				if linkedModel == nil {
					panic("missing linked model " + prop.LinksTo)
				}
				model.Joins = append(model.Joins, ModelJoin{
					LinkColumn: prop.ColumnName,
					LinkField:  prop.FieldName,
					Model:      linkedModel,
				})
			}
		}
	}
	return models
}

func debugPrint(m *Model) {
	fmt.Println("Properties:")
	for _, prop := range m.Properties {
		fmt.Printf("\t%q of type %s\n", prop.Name, prop.Type)
	}
}
