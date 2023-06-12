package handlers

import (
	"github.com/CloudyKit/jet/v6"
	"log"
	"net/http"
)

// renderPage is used to render a page using the jet templating engine
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {
	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Println(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
