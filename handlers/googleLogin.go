package handlers

import (
	"encoding/json"
	"forum/database"
	"forum/models"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"log"
	"fmt"
)

func GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	if clientID == "" {
		http.Error(w, "GOOGLE_CLIENT_ID manquant", http.StatusInternalServerError)
		return
	}

	// preparer les données à envoyer à google pour vérifier 
	// qu'on a bien droit d'accéder au site pour se Oauth
	values := url.Values{}
	values.Set("client_id", clientID)
	values.Set("redirect_uri", "http://localhost:8080/auth/google/callback")
	values.Set("response_type", "code")
	values.Set("scope", "openid email profile")
	values.Set("access_type", "online")

	redirectURL := "https://accounts.google.com/o/oauth2/v2/auth?" + values.Encode()

	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code Google manquant", http.StatusBadRequest)
		return
	}

	token, err := exchangeGoogleCode(code)
	if err != nil {
		http.Error(w, "Erreur échange token Google", http.StatusInternalServerError)
		return
	}

	googleUser, err := getGoogleUser(token.AccessToken)
	if err != nil {
		http.Error(w, "Erreur récupération utilisateur Google", http.StatusInternalServerError)
		return
	}

	user, err := database.FindOrCreateGoogleUser(googleUser)
	if err != nil {
		log.Println("Erreur utilisateur Google:", err)
		http.Error(w, "Erreur utilisateur", http.StatusInternalServerError)
		return
	}

	sessionID, err := database.CreateSession(user.ID)
	if err != nil {
		log.Println("Erreur session Google:", err)
		http.Error(w, "Erreur session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400,
		HttpOnly: true,
	})

	// puis créer ta session/cookie
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func exchangeGoogleCode(code string) (*models.GoogleTokenResponse, error) {
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("variables OAuth Google manquantes")
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("grant_type", "authorization_code")
	data.Set("redirect_uri", "http://localhost:8080/auth/google/callback")

	req, err := http.NewRequest(
		"POST",
		"https://oauth2.googleapis.com/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Google token error: %s", string(body))
	}

	var token models.GoogleTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func getGoogleUser(accessToken string) (*models.GoogleUser, error) {
	req, err := http.NewRequest(
		"GET",
		"https://openidconnect.googleapis.com/v1/userinfo",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user models.GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}