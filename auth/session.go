package auth

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/aliBazrkar/go-chatapp/db"
	"gorm.io/gorm"
)

var (
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInvalidSession = errors.New("invalid session token")
	ErrInvalidCSRF    = errors.New("invalid csrf token")
	ErrExpiredSession = errors.New("session expired")
	ErrInvalidOrigin  = errors.New("invalid origin")
)

const (
	// DEV NOTES:
	// make sure to update origin later in deployment
	// this parameter is passed to chat package too.
	// Secure : True will block localhost connection
	// only secure it when you actually have deployed
	// it and have a cert to set on the http start.
	Origin            = "http://localhost:3000"
	SessionCookieName = "session_token"
	CSRFTokenName     = "csrf_token"
	CSRFHeaderName    = "X-CSRF-Token"
	SessionDuration   = 15 * time.Minute
	TokenLength       = 32
	Secure            = false
)

type SessionManager struct {
	db *db.Database
}

func NewSessionManager(database *db.Database) *SessionManager {
	return &SessionManager{db: database}
}

func (sm *SessionManager) CreateSession(userID uint32, length int) (sessionToken string, csrfToken string, expiresAt time.Time, err error) {
	sessionToken = generateToken(length)
	csrfToken = generateToken(length)
	expiresAt = time.Now().Add(SessionDuration)

	_, err = sm.db.CreateToken(userID, sessionToken, csrfToken, expiresAt)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return sessionToken, csrfToken, expiresAt, nil
}

func (sm *SessionManager) ValidateSession(w http.ResponseWriter, r *http.Request) (*db.Token, error) {

	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil, ErrInvalidSession
	}

	csrfToken := r.Header.Get(CSRFHeaderName)
	if csrfToken == "" {
		return nil, ErrInvalidCSRF
	}

	// since we are passing the value from browser
	// header right to find token from database,
	// if the result is found, ST is evaluated already.
	token, err := sm.db.GetTokenBySession(cookie.Value)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if token.CSRFToken != csrfToken {
		return nil, ErrInvalidCSRF
	}

	if time.Now().After(token.ExpiresAt) {
		sm.DeleteSession(cookie.Value)
		return nil, ErrExpiredSession
	}

	if w != nil {
		if time.Until(token.ExpiresAt) < SessionDuration/2 {
			if err := sm.ExtendSession(w, cookie.Value); err != nil {
				log.Println("Session Extending Error:", err)
			}
		}
	}
	return token, nil
}

// the reason why we use CSRF tokens is to make sure the
// request is coming from same origins, there are scenarios
// where csrf tokens are not needed, for example page
// refresh or websocket connection. for WS we need simple
// session cookie or JWT to just validate it's actually us.
// this method validates your session without CSRF token.
func (sm *SessionManager) CheckAuth(w http.ResponseWriter, r *http.Request) (*db.Token, error) {

	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || cookie.Value == "" {
		return nil, ErrInvalidSession
	}

	token, err := sm.db.GetTokenBySession(cookie.Value)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUnauthorized
		}
		return nil, err
	}

	if time.Now().After(token.ExpiresAt) {
		sm.DeleteSession(cookie.Value)
		return nil, ErrExpiredSession
	}

	if w != nil {
		if time.Until(token.ExpiresAt) < SessionDuration/2 {
			if err := sm.ExtendSession(w, cookie.Value); err != nil {
				log.Println("Session Extending Error:", err)
			}
		}
	}
	return token, nil
}

func (sm *SessionManager) ExtendSession(w http.ResponseWriter, sessionToken string) error {
	newExpiresAt := time.Now().Add(SessionDuration)
	if err := sm.db.ExtendSessionExpiration(sessionToken, newExpiresAt); err != nil {
		return err
	}
	SetSessionCookie(w, sessionToken, newExpiresAt)
	return nil
}

func (sm *SessionManager) DeleteSession(sessionToken string) error {
	return sm.db.DeleteToken(sessionToken)
}

func (sm *SessionManager) CleanupExpiredSessions() error {
	return sm.db.CleanupExpiredTokens()
}

func SetSessionCookie(w http.ResponseWriter, sessionToken string, expiresAt time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    sessionToken,
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   Secure,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: false,
		Secure:   Secure,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})
}
