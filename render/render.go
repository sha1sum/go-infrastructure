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
		log.WithFields(log.Fields{"event": packagename + "Render", "view": view}).Debug("Rendering view")
	}

	file := Settings.TemplateDirectory + view + ".html"
	return executeTemplate(file, args)

	//return content, err
}

// executeTemplate ensures templates are cached
// if caching is enabled.
func executeTemplate(file string, data interface{}) (body []byte, err error) {

	if Settings.LogDebugMessages {
		log.WithFields(log.Fields{"event": packagename + "Render", "file": file}).Debug("Rendering template")
	}

	// Place a read lock on our registry
	tr.RLock()
	t, present := tr.templates[file]
	tr.RUnlock()

	// If the view is not already present in the registry
	if !present {
		if Settings.LogDebugMessages {
			log.WithFields(log.Fields{"event": packagename + "Render", "file": file}).Debug("Parsing template")
		}

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
			if Settings.LogDebugMessages {
				log.WithFields(log.Fields{"event": packagename + "Render", "file": file}).Debug("Caching rendered template")
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
			log.WithFields(log.Fields{"event": packagename + "Render", "file": file, "error": err}).Debug("Unable to execute template")
		}
		return
	}

	return buf.Bytes(), nil
}
