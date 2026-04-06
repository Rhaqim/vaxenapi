package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	AccessTokenCookie  = "vaxen_access_token"
	RefreshTokenCookie = "vaxen_refresh_token"
	CSRFTokenHeader    = "X-CSRF-Token"
	CSRFTokenCookie    = "vaxen_csrf_token"
)

// SetAuthCookies sets httpOnly secure cookies for the access and CSRF tokens.
// The access token is httpOnly (invisible to JS). The CSRF token is readable
// by JS so the frontend can include it in request headers.
func SetAuthCookies(c *gin.Context, accessToken string, csrfToken string, secure bool, maxAgeSec int) {
	// Access token: httpOnly, secure, strict same-site — JS cannot read this
	c.SetSameSite(http.SameSiteStrictMode)
	c.SetCookie(
		AccessTokenCookie,
		accessToken,
		maxAgeSec,
		"/",
		"",    // domain — empty = current host
		secure,
		true, // httpOnly
	)

	// CSRF token: NOT httpOnly so JS can read it and send as a header
	c.SetCookie(
		CSRFTokenCookie,
		csrfToken,
		maxAgeSec,
		"/",
		"",
		secure,
		false, // not httpOnly — JS must read this
	)
}

// ClearAuthCookies removes all auth cookies (logout).
func ClearAuthCookies(c *gin.Context, secure bool) {
	c.SetCookie(AccessTokenCookie, "", -1, "/", "", secure, true)
	c.SetCookie(CSRFTokenCookie, "", -1, "/", "", secure, false)
	c.SetCookie(RefreshTokenCookie, "", -1, "/", "", secure, true)
}
