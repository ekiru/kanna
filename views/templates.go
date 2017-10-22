package views

import (
	"bytes"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var templates map[string]*template.Template

func init() {
	templates = make(map[string]*template.Template)
	var partials []string
	views := map[string]string{}
	readTemplates("templates", "", &partials, views)
	files := make([]string, len(partials)+2)
	files[0] = "templates/layout.html"
	copy(files[2:], partials)
	for name, fileName := range views {
		files[1] = fileName
		templates[name] = template.Must(template.ParseFiles(files...))
	}
}

func readTemplates(dirname string, prefix string, partials *[]string, views map[string]string) {
	fs := readDir(dirname)
	for _, f := range fs {
		if dirname == "templates" && f.Name() == "layout.html" {
			continue
		}
		if f.Name()[0] == '.' {
			continue
		}
		fileName := filepath.Join(dirname, f.Name())
		name := prefix + f.Name()
		if f.IsDir() {
			readTemplates(fileName, name+"/", partials, views)
		} else if strings.Contains(name, ".partial.") {
			*partials = append(*partials, fileName)
		} else {
			views[name] = fileName
		}
	}
}

func readDir(name string) []os.FileInfo {
	f, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fs, err := f.Readdir(0)
	if err != nil {
		panic(err)
	}
	return fs
}

// HtmlTemplate views serve an HTML document from a template defined
// using the template/html package.
type HtmlTemplate string

// Render renders the template to the response writer, passing the
// supplied data to the template.
func (template HtmlTemplate) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	var output bytes.Buffer
	if err := templates[string(template)].ExecuteTemplate(&output, "layout.html", data); err != nil {
		panic(err)
	}
	sendHtml(w, output.Bytes())
}

// ServeHTTP on an HtmlTemplate calls Render with nil data.
func (template HtmlTemplate) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	template.Render(w, r, nil)
}
