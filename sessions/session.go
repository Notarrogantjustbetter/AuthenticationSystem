package sessions

import (
	"crypto/rand"
	"io"
	"net/http"
	"github.com/gorilla/sessions"
)


var Store = sessions.NewCookieStore(generateRandomCookieKey(64))

func generateRandomCookieKey(length int) []byte {
	maker := make([]byte, length)
	if _, err := io.ReadFull(rand.Reader, maker); err != nil {
		return nil
	}
	return maker
}

func SetSession(w http.ResponseWriter, r *http.Request, userId int64) error {
	session, _ := Store.Get(r, "session")
	session.Values["user_id"] = userId
	return session.Save(r, w)
}

func DeleteSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := Store.Get(r, "session")
	delete(session.Values, "user_id")
	return session.Save(r, w)
}

