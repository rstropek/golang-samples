package handlers

import (
	"log"

	"github.com/gorilla/sessions"
)

type Handlers struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	Store    *sessions.CookieStore
}
