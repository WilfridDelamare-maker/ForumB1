package handlers

import "net/http"

func GetDarkMode(r *http.Request) bool {
	cookie, err := r.Cookie("theme")
	if err != nil {
		return true
	}

	return cookie.Value == "dark"
}

func ThemeHandler(w http.ResponseWriter, r *http.Request) {
	theme := GetDarkMode(r)

	newtheme := "dark"
	if theme {
		newtheme = "light"
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "theme",
        Value: newtheme,
        Path:  "/",
        MaxAge: 3600 * 24 * 30,
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}