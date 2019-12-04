package oauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/pkg/browser"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/firebase?access_token="
	port              = "8000"
)

func generateOauthStateTracker() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return base64.URLEncoding.EncodeToString([]byte(string(b)))
}

func getUserDataFromGoogle(googleOauthConfig *oauth2.Config, code string) ([]byte, error) {
	// Use code to get token and get user info from Google.
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func StartLocalhostServer() {
	scopes := []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/firebase",
	}
	googleOauthConfig := &oauth2.Config{
		RedirectURL:  fmt.Sprintf("http://localhost:%s", port),
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
	// open browser now
	oauthStateTracker := generateOauthStateTracker()
	u := googleOauthConfig.AuthCodeURL(oauthStateTracker)
	fmt.Printf("Visit this URL on any device to log in:\n\n%s\n\nWaiting for authentication...", u)
	_ = browser.OpenURL(u)
	router := chi.NewRouter()
	httpServer := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router,
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("state") != oauthStateTracker {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		data, err := getUserDataFromGoogle(googleOauthConfig, r.FormValue("code"))
		if err != nil {
			log.Println(err.Error())
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		// todo: save credentials
		// todo: return success message to the user
		fmt.Fprintf(w, "UserInfo: %s\n", data)
		// use subroutine to shutdown server
		cancel()

	})
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	// Shutdown the server when the context is canceled
	fmt.Printf("\n\nLogin Successful!\n\n")
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}
