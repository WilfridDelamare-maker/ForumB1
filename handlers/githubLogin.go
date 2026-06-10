package handlers

import (
	"encoding/json"
	"fmt"
	"forum/models"
	"forum/database"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://github.com/login/oauth/authorize" + "?client_id=" +
		os.Getenv("GITHUB_CLIENT_ID") +
		"&redirect_uri=http://localhost:8080/auth/github/callback" +
		"&scope=user:email"

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Code github manquant", http.StatusBadRequest)
		return
	}

	token, err := GetGitHubToken(code)
	if err != nil {
		log.Println("Erreur token GitHub:", err)
		http.Error(w, "Impossible de récupérer le token Github", http.StatusInternalServerError)
		return
	}

	GHuser, err := GetGitHubUser(token.AccessToken)
	if err != nil {
		log.Println("Erreur utilisateur GitHub:", err)
		http.Error(w, "Impossible de récupérer l'utilisateur", http.StatusInternalServerError)
		return
	}

	user, err := database.FindOrCreateGitHubUser(GHuser.ID, GHuser.Login, GHuser.Email)
	if err != nil {
		log.Println("Erreur user GitHub BDD:", err)
		http.Error(w, "Erreur utilisateur GitHub", http.StatusInternalServerError)
		return
    }

	sessionID, err := database.CreateSession(user.ID)
	if err != nil {
		log.Println("Erreur session GitHub:", err)
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func GetGitHubToken(code string) (*models.GithubTokenResponse, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	data := url.Values{} // equivaut à une map de string []string
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:8080/auth/github/callback")

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode())) // data.encode transforme le texte en urlencoded format
	// io.reader c'est une interface qui lit progressivement ce qu'on lui donne.
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var client *http.Client = &http.Client{} //objet qui permet d'envoyer req, recevoir resp et autres... il sert à envoyer la requete et recevoir reponse.

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token models.GithubTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("github token vide")
	}

	return &token, nil
}

func GetGitHubUser(token string) (models.GitHubUser, error) {
	var githubUser models.GitHubUser

	req, err := http.NewRequest( // body conforme a ce que github attend
		"GET",
		"https://api.github.com/user",
		nil,
	)
	if err != nil {
		return githubUser, err
	}

	req.Header.Set("Authorization", "Bearer "+token) // header conforme à ce que github attend

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return githubUser, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&githubUser)
	if err != nil {
		return githubUser, err
	}
	if githubUser.ID == 0 {
		return githubUser, fmt.Errorf("utilisateur github vide")
	}

	return githubUser, nil
}
