package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/negroni"
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/prometheus/common/version"
	"github.com/sebastianrosch/livingroompresentations/app"
	"github.com/sebastianrosch/livingroompresentations/internal/config"
	"github.com/sebastianrosch/livingroompresentations/internal/store"
	"github.com/sebastianrosch/livingroompresentations/pkg/log"
	"github.com/sebastianrosch/livingroompresentations/routes/callback"
	"github.com/sebastianrosch/livingroompresentations/routes/home"
	"github.com/sebastianrosch/livingroompresentations/routes/login"
	"github.com/sebastianrosch/livingroompresentations/routes/logout"
	"github.com/sebastianrosch/livingroompresentations/routes/middlewares"
	"github.com/sebastianrosch/livingroompresentations/routes/newevent"
	"github.com/sebastianrosch/livingroompresentations/routes/user"
)

func main() {
	logger := log.NewDefaultLogger()
	logger.Info("Starting Living Room API",
		"version", version.Info(),
		"build_context", version.BuildContext())

	s, err := store.NewMongoStore(
		config.Get().DatabaseURI,
		config.Get().DatabaseName,
		config.Get().DatabaseUsername,
		config.Get().DatabasePassword,
	)
	if err != nil {
		logger.Error(err, "couldn't create MongoDB store")
		os.Exit(2)
	}

	s.CreateEvent("How viruses spread", "Epidemologist talks about how viruses spread")

	events, _ := s.GetAllEvents()
	for _, event := range events {
		fmt.Print(event)
	}

	app.Init()
	startServer(logger)
}

func startServer(logger logr.Logger) {
	r := mux.NewRouter()

	r.HandleFunc("/", home.HomeHandler)
	r.HandleFunc("/login", login.LoginHandler)
	r.HandleFunc("/logout", logout.LogoutHandler)
	r.HandleFunc("/callback", callback.CallbackHandler)
	r.Handle("/newevent", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(newevent.NewEventHandler)),
	))
	r.Handle("/user", negroni.New(
		negroni.HandlerFunc(middlewares.IsAuthenticated),
		negroni.Wrap(http.HandlerFunc(user.UserHandler)),
	))
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	http.Handle("/", r)
	logger.Info(fmt.Sprintf("Server listening on http://%s:%s/", config.Get().Host, config.Get().Port))
	logger.Error(http.ListenAndServe(fmt.Sprintf("%s:%s", config.Get().Host, config.Get().Port), nil), "terminated")
}
