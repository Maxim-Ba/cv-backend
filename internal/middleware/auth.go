package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/http"
)

const sessionCookieName = "admin_session"

func makeToken(user, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(user))
	return fmt.Sprintf("%x", mac.Sum(nil))
}

// SetSessionCookie выставляет cookie после успешного логина.
func SetSessionCookie(w http.ResponseWriter, user, secret string) {
	token := makeToken(user, secret)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Path:     "/admin",
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearSessionCookie удаляет cookie сессии.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/admin",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// RequireAuth возвращает middleware, которое проверяет cookie сессии.
func RequireAuth(user, secret string) func(http.Handler) http.Handler {
	expected := makeToken(user, secret)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(sessionCookieName)
			if err != nil || !hmac.Equal([]byte(cookie.Value), []byte(expected)) {
				http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
