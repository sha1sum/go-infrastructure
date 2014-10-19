// Package render renders responses to clients
package render

import (
	"bytes"
	"html/template"
	"log"
	"path/filepath"
	"sync"
)

type (
	// Conventions organizes the default settings for the Webserver renderer
	Conventions struct {
		TemplateDirectory string
		Debug             bool
		CacheTemplates    bool
	}

	Renderer interface {
		Render(args ...interface{}) ([]byte, error)
	}

	html struct{}

	templateRegistry struct {
		sync.RWMutex
		templates map[string]*template.Template
	}
)

var Settings = Conventions{
	TemplateDirectory: "web-src/html/",
	Debug:             false,
	CacheTemplates:    true,
}

var (
	tr   templateRegistry
	HTML = html{}
)

func init() {
	tr.templates = make(map[string]*template.Template)
}

func (r html) Render(view string, args ...interface{}) ([]byte, error) {

	if Settings.Debug {
		log.Printf("Rendering view %s", view)
	}

	file := Settings.TemplateDirectory + view + ".html"
	return executeTemplate(file, args)

	//return content, err
}

// executeTemplate ensures templates are cached
// if caching is enabled.
func executeTemplate(file string, data interface{}) (body []byte, err error) {

	// Place a read lock on our registry
	tr.RLock()
	t, present := tr.templates[file]
	tr.RUnlock()

	// If the view is not already present in the registry
	if !present {
		if Settings.Debug {
			log.Printf("Parsing template %s", file)
		}

		// Create a new template
		t = template.New(filepath.Base(file))

		// Enhance our template with custom format so we can reuse with JS libs?
		// t.Delims("<%=", ">")
		_, err = t.ParseFiles(file)
		if err != nil {
			log.Printf("Unable to render template due to the following error: %s", err)
			return
		}

		if Settings.CacheTemplates {
			if Settings.Debug {
				log.Printf("Caching rendered template %s", file)
			}
			tr.Lock()
			tr.templates[file] = t
			tr.Unlock()
		}
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return
	}

	return buf.Bytes(), nil
}
