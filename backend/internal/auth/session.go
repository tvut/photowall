package auth

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

var Sessions *scs.SessionManager

func Init() {
	Sessions = scs.New()
	Sessions.Cookie.HttpOnly = true
	Sessions.Cookie.SameSite = http.SameSiteLaxMode
	Sessions.Cookie.Secure = false // true in prod (HTTPS)
}
