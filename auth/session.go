package auth

import (
	"errors"
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
)

const (
	SessionCookieName = "session_token"
	CSRFHeaderName    = "X-CSRF-Token"
	SessionDuration   = 60 * time.Minute
	TokenLength       = 32
	Secure            = false // fix -> True later
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

func (sm *SessionManager) ValidateSession(r *http.Request) (*db.User, error) {

	// here we get both session token and CSRF token values from browser header
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

	return token.User, nil
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
	})
}

func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Expires:  time.Unix(0, 0),
		HttpOnly: true,
		Secure:   Secure,
		SameSite: http.SameSiteStrictMode,
	})
}
