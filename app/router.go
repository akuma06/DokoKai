package app

import (
    _ "github.com/akuma06/DokoKai/app/controllers/animes" // animescontroller
    _ "github.com/akuma06/DokoKai/app/controllers/static" // staticcontroller
	"net/http"

	"github.com/akuma06/DokoKai/app/controllers"
	"github.com/justinas/nosurf"
)

// CSRFRouter : CSRF protection for Router variable for exporting the route configuration
var CSRFRouter *nosurf.CSRFHandler

func init() {
	CSRFRouter = nosurf.New(controllers.Get())
	CSRFRouter.ExemptRegexp("/api(?:/.+)*")
	CSRFRouter.ExemptRegexp("/mod(?:/.+)*")
	CSRFRouter.ExemptPath("/upload")
	CSRFRouter.ExemptPath("/user/login")
	CSRFRouter.ExemptPath("/oauth2/token")
	CSRFRouter.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Invalid CSRF tokens", http.StatusBadRequest)
	}))
	CSRFRouter.SetBaseCookie(http.Cookie{
		Path:   "/",
		MaxAge: nosurf.MaxAge,
	})

}