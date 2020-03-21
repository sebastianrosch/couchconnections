package logout

import (
	"net/http"
	"net/url"

	"github.com/sebastianrosch/livingroompresentations/internal/config"
)

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	domain := config.Get().Auth0Domain

	logoutUrl, err := url.Parse("https://" + domain)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logoutUrl.Path += "/v2/logout"
	parameters := url.Values{}

	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", config.Get().Auth0ClientID)
	logoutUrl.RawQuery = parameters.Encode()

	http.Redirect(w, r, logoutUrl.String(), http.StatusTemporaryRedirect)
}
