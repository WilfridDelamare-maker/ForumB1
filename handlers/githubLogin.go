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

// simple fonction pour rediriger vers github avec le client_id. 
// github va vérifier qui demande une OAuth avec la client_id
func GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://github.com/login/oauth/authorize" + "?client_id=" +
		os.Getenv("GITHUB_CLIENT_ID") +
		"&redirect_uri=http://localhost:8085/auth/github/callback" +
		"&scope=user:email"

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func GitHubCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code") // on récupere le code de github pour lui renvoyer
	if code == "" {
		http.Error(w, "Code github manquant", http.StatusBadRequest)
		return
	}

	token, err := GetGitHubToken(code) // grace au code on peut récupérer un token github
	if err != nil {
		log.Println("Erreur token GitHub:", err)
		http.Error(w, "Impossible de récupérer le token Github", http.StatusInternalServerError)
		return
	}

	GHuser, err := GetGitHubUser(token.AccessToken) // on peut enfin récupérer les données de l'user
	// grace au token
	if err != nil {
		log.Println("Erreur utilisateur GitHub:", err)
		http.Error(w, "Impossible de récupérer l'utilisateur", http.StatusInternalServerError)
		return
	}
	if GHuser.Email == "" { // si github masque l'email, on essaye de le récupérer explicitement (merci à scope=user:email)
		email, err := GetGitHubPrimaryEmail(token.AccessToken)
		if err != nil {
			log.Println("Erreur email non récupéré: ", err)
			http.Error(w, "Impossible de récupérer un email", http.StatusInternalServerError)
			return
		}

		GHuser.Email = email
	}

	// créer ou retrouver l'user depuis la db
	user, err := database.FindOrCreateGitHubUser(GHuser.ID, GHuser.Login, GHuser.Email)
	if err != nil {
		log.Println("Erreur user GitHub BDD:", err)
		http.Error(w, "Erreur utilisateur GitHub", http.StatusInternalServerError)
		return
    }

	// ensuite on obtient le sessionid et on crée les cookies de connexion
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

/* permet de récupérer un token de github grace au code valide recu.
il faut envoyer le client_id et client_secret en plus du code pour vérifier 
l'identité de notre requete 
*/
func GetGitHubToken(code string) (*models.GithubTokenResponse, error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	data := url.Values{} // equivaut à une map de string []string
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", "http://localhost:8085/auth/github/callback")

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode())) // data.encode transforme le texte en urlencoded format
	// io.reader c'est une interface qui lit progressivement ce qu'on lui donne.
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json") // on explique qu'on veut du JSON en format
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //on explique le format du body (url encoded)

	var client *http.Client = &http.Client{} //objet qui permet d'envoyer req, recevoir resp et autres... il sert à envoyer la requete et recevoir reponse.

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token models.GithubTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&token) // récupere la data de la réponse et le stocke dans token
	if err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("github token vide")
	}

	return &token, nil
}

// fonction pour récupérer l'id, le login et l'email depuis github
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

// fonction pour essayer de récupérer l'email si pas envoyé lors de la réponse d'origine
// cherche parmi les emails de l'user son mail primaire et verifié
func GetGitHubPrimaryEmail(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []models.GitHubEmail

	err = json.NewDecoder(resp.Body).Decode(&emails)
	if err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	for _, m := range emails {
		if m.Verified {
			return m.Email, nil
		}
	}

	return "", fmt.Errorf("aucun email github vérifié...")
}
