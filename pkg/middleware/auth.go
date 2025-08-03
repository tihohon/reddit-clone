package middleware

import (
	"net/http"
	"redditclone/pkg/helpers"
	"redditclone/pkg/models/session"
	"regexp"
	"strings"
)

var (
	noAuthPatterns = []*regexp.Regexp{
		regexp.MustCompile(`^/api/posts/[^/]*$`), // /api/posts/{category} - exact match
		regexp.MustCompile(`^/api/posts/$`),      // /api/posts/ - exact match
		regexp.MustCompile(`^/api/register$`),    // /api/register - exact match
		regexp.MustCompile(`^/api/login$`),       // /api/login - exact match
		regexp.MustCompile(`^/static/.*$`),       // /static/anything
		regexp.MustCompile(`^/$`),                // root path
	}
)

func isUnAuthUrl(url string, method string) bool {
	if !strings.HasPrefix(url, "/api") {
		return true
	}

	if strings.HasPrefix(url, "/api/user") {
		return true
	}

	if strings.HasPrefix(url, "/api/post/") && method == "GET" {
		parts := strings.Split(url, "/")
		if len(parts) == 4 && parts[1] == "api" && parts[2] == "post" {
			return true
		}
	}

	for _, pattern := range noAuthPatterns {
		if pattern.MatchString(url) {
			return true
		}
	}

	return false
}

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isUnAuthUrl(r.RequestURI, r.Method) {
			next.ServeHTTP(w, r)
			return
		}
		sess, err := sm.Check(r)
		if err != nil {
			sm.Logger.Errorf("Auth error %e", err)
			helpers.WriteBadRequest(w, "Auth error")
			return
		}

		ctx := session.ContextWithSession(r.Context(), sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
