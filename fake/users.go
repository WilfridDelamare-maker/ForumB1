package fake

import "net/http"

func GetCurrentUser(r *http.Request) (string, bool) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return "", false
	}

	if cookie.Value == "session_nbr" {
		return "Boss", true
	}

	return "", false
}