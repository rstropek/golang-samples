package handlers

import (
	"net/http"
	"text/template"
)

func (h Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	ts, err := template.ParseFiles("./ui/html/home.page.tmpl")
	if err != nil {
		h.ErrorLog.Printf("Error parsing template (%s)", err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		h.ErrorLog.Printf("Error executing template (%s)", err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}
