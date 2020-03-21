package newevent

import (
	"fmt"
	"net/http"

	"github.com/sebastianrosch/livingroompresentations/app"
	"github.com/sebastianrosch/livingroompresentations/routes/templates"
)

func NewEventHandler(w http.ResponseWriter, r *http.Request) {
	session, err := app.Store.Get(r, "auth-session")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	profile := session.Values["profile"]
	fmt.Println(profile)

	templates.RenderTemplate(w, "newevent", profile)
}
