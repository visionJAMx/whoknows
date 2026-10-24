package server

import (
	"html/template"
)

var appTemplates = make(map[string]*template.Template)

func InitTemplates() {
	pages := []string{"search.html", "about.html", "login.html", "register.html"}

	for _, page := range pages {

		appTemplates[page] = template.Must(template.ParseFiles(
			"templates/layout.html",
			"templates/"+page,
		))
	}
}
