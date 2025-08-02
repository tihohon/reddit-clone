package middleware

import (
	"net/http"
	"redditclone/pkg/models/session"
	"strings"
)

var (
	noAuthUrls = map[string]struct{}{"/api/posts/": {}, "/api/register": {}}
)

func Auth(sm *session.SessionsManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.String(), "/api") {
			next.ServeHTTP(w, r)
			return
		}
		if _, ok := noAuthUrls[r.URL.Path]; ok {
			next.ServeHTTP(w, r)
			return
		}
		_, err := sm.Check(r)
		if err != nil {
			sm.Logger.Errorf("Auth error %w", err)
		}

		next.ServeHTTP(w, r)
	})

}
