package controllers

import "html/template"

var templates *template.Template

func SetTemplates(t *template.Template) {
	templates = t
}
