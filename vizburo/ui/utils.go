package ui

import (
	"html/template"
	"net/http"
	"strings"
)

// RenderTemplate renders a template with the base layout
func RenderTemplate(w http.ResponseWriter, templateName string, data map[string]interface{}) error {
	// Define safe HTML function for templates
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"split": func(s string, sep string) []string {
			return strings.Split(s, sep)
		},
		"mod": func(a, b int) int {
			return a % b
		},
	}

	// Parse the base template and all dependencies with custom functions
	t := template.New("base.html").Funcs(funcMap)
	t, err := t.ParseFiles(
		"vizburo/ui/templates/layouts/base.html",
		"vizburo/ui/templates/"+templateName,
	)
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := t.Execute(w, data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}

// RenderPartial renders just the content portion of a template (for HTMX responses)
func RenderPartial(w http.ResponseWriter, templateName string, data map[string]interface{}) error {
	// Define safe HTML function for templates
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
		"split": func(s string, sep string) []string {
			return strings.Split(s, sep)
		},
		"mod": func(a, b int) int {
			return a % b
		},
	}

	// Parse just the template file without base layout
	t := template.New("partial").Funcs(funcMap)
	t, err := t.ParseFiles("vizburo/ui/templates/" + templateName)
	if err != nil {
		http.Error(w, "Error loading template: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// Execute the "content" template defined in the partial file
	if err := t.ExecuteTemplate(w, "content", data); err != nil {
		http.Error(w, "Error rendering template: "+err.Error(), http.StatusInternalServerError)
		return err
	}

	return nil
}
