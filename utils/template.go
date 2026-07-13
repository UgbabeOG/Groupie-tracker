package utils

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

// templates holds the pre-parsed template collection loaded at init.
var templates *template.Template

func init() {
    var err error
    patterns := []string{
        filepath.Join("templates", "*.html"),
        filepath.Join("..", "templates", "*.html"),
    }

    for _, pattern := range patterns {
        templates, err = template.ParseGlob(pattern)
        if err == nil {
            return
        }
    }

    log.Fatalf("Failed to parse templates from %v: %v", patterns, err)
}

// RenderTemplate executes the named template into an internal bytes.Buffer
// before writing any HTTP headers or body content. This ensures that a
// malformed template never results in a half-rendered page or application
// panic — if execution fails the caller receives an error and can redirect
// to the error page instead.
func RenderTemplate(w http.ResponseWriter, tmpl string, data any) error {
	var buf bytes.Buffer

	if err := templates.ExecuteTemplate(&buf, tmpl, data); err != nil {
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := buf.WriteTo(w)
	return err
}

// RenderError renders the error.html template with the given status code
// and message. If the error template itself fails, a plain-text fallback
// is written directly to the response.
func RenderError(w http.ResponseWriter, code int, message string) {
	data := ErrorData{Code: code, Message: message}

	var buf bytes.Buffer
	if err := templates.ExecuteTemplate(&buf, "error.html", data); err != nil {
		log.Printf("[template] ERROR rendering error page: %v", err)
		http.Error(w, message, code)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(code)
	buf.WriteTo(w)
}
