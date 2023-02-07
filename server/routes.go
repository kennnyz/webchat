package main

import (
	"github.com/bmizerany/pat"
	"github.com/kennnyz/webchat/handlers"
	"net/http"
)

func routes() http.Handler {
	app := pat.New()

	app.Get("/", http.HandlerFunc(handlers.Home))
	app.Get("/ws", http.HandlerFunc(handlers.WsEndpoint))

	fileServer := http.FileServer(http.Dir("./static/"))
	app.Get("/static/", http.StripPrefix("/static", fileServer))

	return app
}
