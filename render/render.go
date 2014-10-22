// Package render renders responses to clients
package render

import (
	"bytes"
	"html/template"
	"path/filepath"
	"sync"

	log "github.com/Sirupsen/logrus"
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

	// Renderer is an interface type that different renderers should implement
	Renderer interface {
		Render(args ...interface{}) ([]byte, error)
	}

	html struct{}

	templateRegistry struct {
		sync.RWMutex
		templates map[string]*template.Template
	}
)

const packagename = "webserver.render:"

// Settings provides exported access to runtime configuration
var Settings = Conventions{
	TemplateDirectory:  "web-src/html/",
	LogDebugMessages:   false,
	LogErrorMessages:   true,
	LogTemplateResults: false,
	CacheTemplates:     true,
}

var (
	tr templateRegistry
	// HTML renderer which implements the Render method
	HTML = html{}
)

func init() {
	tr.templates = make(map[string]*template.Template)
}

// Render executes a template returning the rendered byte array and error
// While this method supports the Renderer interface only one args is allowed.
func (r html) Render(view string, args ...interface{}) ([]byte, error) {
	log.WithFields(log.Fields{"event": packagename + "Render", "view": view}).Debug("Rendering view")

	file := Settings.TemplateDirectory + view + ".html"
	return executeTemplate(file, args[0])
}

// executeTemplate ensures templates are cached
// if caching is enabled.
func executeTemplate(file string, data interface{}) (body []byte, err error) {
	log.WithFields(log.Fields{"event": packagename + "Render", "file": file, "willCache": Settings.CacheTemplates}).Debug("Rendering template")

	// Place a read lock on our registry
	tr.RLock()
	t, present := tr.templates[file]
	tr.RUnlock()

	// If the view is not already present in the registry
	if !present {
		log.WithFields(log.Fields{"event": packagename + "Render", "file": file}).Debug("Parsing template")

		// Create a new template
		t = template.New(filepath.Base(file))

		// Enhance our template with custom format so we can reuse with JS libs?
		// t.Delims("<%=", ">")
		_, err = t.ParseFiles(file)
		if err != nil {
			log.WithFields(log.Fields{"event": packagename + "Render", "file": file, "error": err}).Error("Unable to parse template")
			return
		}

		if Settings.CacheTemplates {
			log.WithFields(log.Fields{"event": packagename + "Render", "file": file}).Debug("Caching rendered template")

			tr.Lock()
			tr.templates[file] = t
			tr.Unlock()
		}
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		log.WithFields(log.Fields{"event": packagename + "Render", "file": file, "error": err}).Debug("Unable to execute template")

		return
	}

	return buf.Bytes(), nil
}
