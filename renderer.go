// Package webserver is responsible for parses templates and rendering them at
// at runtime. This package makes it easy to render HTML for the applications
// that use the Wrecker Labs infrastructure.
package webserver

import (
	"bytes"
	"html/template"
	"net/http"
	"path/filepath"
	"sync"
)

// Registry serves as a instance cache for our application. As new templates
// are encountered at runtime they are parsed and cached within the registry
// so we can avoid repeated compilation as well as subsequent reads from
// the disk. When enabled, a template is read and compiled once per instance
// of the application.
type registry struct {
	sync.RWMutex
	templates map[string]*template.Template
}

var r registry

// Init constructs a map to cache templates by name within our registry.
func init() {
	// Create a new map with a string key and a template value
	r.templates = make(map[string]*template.Template)
}

// RenderLayoutHTML renders a layout and embeds a view within that layout.
func RenderLayoutHTML(
	resp http.ResponseWriter,
	req *http.Request,
	layout string,
	view string,
	data interface{}) (err error) {

	// Attempt to render the view using our private render method
	body, err := render(req, layout, view, data)
	if err != nil {
		// Render a standard error from the HTTP library
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the response
	resp.Write(body)

	return
}

// RenderErrorHTML renders an exception
func RenderErrorHTML(
	resp http.ResponseWriter,
	req *http.Request,
	message string,
	code int) (err error) {

	// Delegate to the render method and render our error
	body, err := render(req, "app", "error", map[string]string{"Error": message})
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(code)
	resp.Write(body)
	return
}

func render(
	req *http.Request,
	layout string,
	name string,
	data interface{}) (body []byte, err error) {

	// Relative path to the actual view file
	file := "web-src/html/view/" + name + ".html"
	// Parse our view
	view, err := parse(req, file, data)
	if err != nil {
		return
	}

	// Parse the layout
	body, err = parse(req, "web-src/html/layout/"+layout+".html", map[string]template.HTML{
		"Content": template.HTML(view),
	})
	if err != nil {
		return
	}
	return
}

// Parse a template file to be rendered. This method ensures templates are cached
// if caching is enabled.
func parse(
	req *http.Request,
	file string,
	data interface{}) (body []byte, err error) {

	// Place a read lock on our registry
	r.RLock()
	t, present := r.templates[file]
	r.RUnlock()

	// If the view is not already present in the registry
	if !present {
		// Create a new template
		t = template.New(filepath.Base(file))

		// Enhance our template with custom format so we can reuse with
		// underscore
		// t.Delims("<%=", ">")
		_, err = t.ParseFiles(file)
		if err != nil {
			return
		}

		// TODO Enable via feature flag
		//
		// Total lock to write to the registry
		//r.Lock()
		//r.templates[file] = t
		//r.Unlock()
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}
