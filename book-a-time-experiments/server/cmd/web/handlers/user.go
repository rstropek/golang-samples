package handlers

import (
	"log"
	"net/http"
	"text/template"
)

func (h Handlers) UserHandler(w http.ResponseWriter, r *http.Request) {

	session, err := h.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ts, err := template.ParseFiles("./ui/html/user.page.tmpl")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, session.Values["profile"])
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}
