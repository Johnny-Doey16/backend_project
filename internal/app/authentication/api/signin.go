package api

import (
	"log"
	"net/http"

	"github.com/steve-mir/diivix_backend/internal/app/authentication/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var ssogolang *oauth2.Config
var RandomStr, _ = utils.GenerateUniqueToken(15)

func init() {
	ssogolang = &oauth2.Config{
		Endpoint:     google.Endpoint,
	}
}

func Signin(w http.ResponseWriter, r *http.Request) {
	url := ssogolang.AuthCodeURL(RandomStr)
	log.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
