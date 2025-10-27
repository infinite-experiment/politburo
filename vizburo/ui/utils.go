package ui

import (
	"html/template"
	"net/http"
)

// RenderTemplate parses and executes a template with the given data
func RenderTemplate(w http.ResponseWriter, templateName string, data map[string]interface{}) error {
	// Define safe HTML function for templates
	funcMap := template.FuncMap{
		"safeHTML": func(s string) template.HTML {
			return template.HTML(s)
		},
	}

	// Parse the base template and all dependencies with custom functions
	t := template.New("base.html").Funcs(funcMap)
	t, err := t.ParseFiles(
		"vizburo/ui/templates/base.html",
		"vizburo/ui/templates/"+templateName,
		"vizburo/ui/templates/showcase.html",
		"vizburo/ui/templates/components/navbar.html",
		"vizburo/ui/templates/components/sidebar.html",
		"vizburo/ui/templates/components/theme-switcher.html",
		"vizburo/ui/templates/components/footer.html",
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

// getThemeFromRequest extracts the theme preference from the request (cookie or default)
func getThemeFromRequest(r *http.Request) string {
	cookie, err := r.Cookie("theme_preference")
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}
	return "light" // default theme
}
