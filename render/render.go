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
		TemplateDirectory  string
		LogDebugMessages   bool
		LogErrorMessages   bool
		LogTemplateResults bool
		CacheTemplates     bool
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

const packagename = "[webserver.render]"

var Settings = Conventions{
	TemplateDirectory:  "web-src/html/",
	LogDebugMessages:   false,
	LogErrorMessages:   true,
	LogTemplateResults: false,
	CacheTemplates:     true,
}

var (
	tr   templateRegistry
	HTML = html{}
)

func init() {
	tr.templates = make(map[string]*template.Template)
}

func (r html) Render(view string, args ...interface{}) ([]byte, error) {
	if Settings.LogDebugMessages {
		log.Printf("%s Rendering view %s", packagename, view)
	}

	file := Settings.TemplateDirectory + view + ".html"
	return executeTemplate(file, args)

	//return content, err
}

// executeTemplate ensures templates are cached
// if caching is enabled.
func executeTemplate(file string, data interface{}) (body []byte, err error) {

	if Settings.LogDebugMessages {
		log.Printf("%s Rendering %s", packagename, file)
	}

	// Place a read lock on our registry
	tr.RLock()
	t, present := tr.templates[file]
	tr.RUnlock()

	// If the view is not already present in the registry
	if !present {
		if Settings.LogDebugMessages {
			log.Printf("%s Parsing template %s", packagename, file)
		}

		// Create a new template
		t = template.New(filepath.Base(file))

		// Enhance our template with custom format so we can reuse with JS libs?
		// t.Delims("<%=", ">")
		_, err = t.ParseFiles(file)
		if err != nil {
			log.Printf("%s Unable to render template due to the following error: %s", packagename, err)
			return
		}

		if Settings.CacheTemplates {
			if Settings.LogDebugMessages {
				log.Printf("%s Caching rendered template %s", packagename, file)
			}
			tr.Lock()
			tr.templates[file] = t
			tr.Unlock()
		}
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		if Settings.LogErrorMessages {
			log.Printf("%s Unable to render template `%s`. Error: %s", packagename, file, err)
		}
		return
	}

	return buf.Bytes(), nil
}
