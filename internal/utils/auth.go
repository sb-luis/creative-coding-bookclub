package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"time"
)

// HashPassword creates a hash of the password using SHA256
func HashPassword(password string) string {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	return hex.EncodeToString(hasher.Sum(nil))
}

// VerifyPassword checks if the provided password matches the hash
func VerifyPassword(password, hash string) bool {
	return HashPassword(password) == hash
}

// AuthenticateMember authenticates a member with name and password
func AuthenticateMember(name, password string) (*Member, error) {
	db := GetDB()
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	member, err := db.GetMemberByName(name)
	if err != nil {
		return nil, errors.New("invalid name or password")
	}

	if !VerifyPassword(password, member.PasswordHash) {
		return nil, errors.New("invalid name or password")
	}

	return member, nil
}

// CreateMemberAccount creates a new member account
func CreateMemberAccount(name, password string) (*Member, error) {
	db := GetDB()
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	if name == "" || password == "" {
		return nil, errors.New("name and password are required")
	}

	if len(password) < 6 {
		return nil, errors.New("password must be at least 6 characters long")
	}

	passwordHash := HashPassword(password)
	return db.CreateMember(name, passwordHash)
}

// GetMemberFromSession retrieves the member associated with a session
func GetMemberFromSession(sessionID string) (*Member, error) {
	db := GetDB()
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	session, err := db.GetSession(sessionID)
	if err != nil {
		return nil, err
	}

	return db.GetMemberByID(session.MemberID)
}

// CreateMemberSession creates a new session for a member
func CreateMemberSession(memberID int) (*Session, error) {
	db := GetDB()
	if db == nil {
		return nil, errors.New("database not initialized")
	}

	return db.CreateSession(memberID)
}

// LogoutMember logs out a member by deleting their session
func LogoutMember(sessionID string) error {
	db := GetDB()
	if db == nil {
		return errors.New("database not initialized")
	}

	return db.DeleteSession(sessionID)
}

// SetSessionCookie sets the session cookie for a member
func SetSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(24 * time.Hour),
	}
	http.SetCookie(w, cookie)
}

// GetSessionFromRequest retrieves the session ID from the request cookie
func GetSessionFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// ClearSessionCookie clears the session cookie
func ClearSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Unix(0, 0), // Set to past time to delete
	}
	http.SetCookie(w, cookie)
}

// GetCurrentMember gets the currently authenticated member from the request
func GetCurrentMember(r *http.Request) (*Member, error) {
	sessionID, err := GetSessionFromRequest(r)
	if err != nil {
		return nil, err
	}

	return GetMemberFromSession(sessionID)
}

// RequireAuth is a middleware that requires authentication
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := GetCurrentMember(r)
		if err != nil {
			http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}
