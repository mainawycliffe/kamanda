package oauth

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	port              = "8000"
)

// GoogleAPIUserInfo the user data returned from Google API Services
type GoogleAPIUserInfo struct {
	Email         string `json:"email"`
	Hd            string `json:"hd"`
	ID            string `json:"id"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

// GetUserDataFromGoogleResponse capture the response we need when users data is
// returned from google api services
type GetUserDataFromGoogleResponse struct {
	Email        string
	RefreshToken string
}

// generateOauthStateTracker generates a random string to act as oauth state
// tracker
func generateOauthStateTracker() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, 10)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return base64.URLEncoding.EncodeToString([]byte(string(b)))
}

// getUserDataFromGoogle fetch user data and the token from google using code
func getUserDataFromGoogle(googleOauthConfig *oauth2.Config, code string) (*GetUserDataFromGoogleResponse, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("OAuth: Error exchanging codes: %w", err)
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("OAuth: Error fetching user data: %w", err)
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("OAuth: Error reading response body: %w", err)
	}
	var userInfo *GoogleAPIUserInfo
	if err = json.Unmarshal(contents, &userInfo); err != nil {
		return nil, fmt.Errorf("OAuth: Error reading data from google: %w", err)
	}
	return &GetUserDataFromGoogleResponse{
		Email:        userInfo.Email,
		RefreshToken: token.RefreshToken,
	}, nil
}

// StartLocalhostServer starts a server that can be used to capture OAUTH
// token from Google Auth Server
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
	oauthStateTracker := generateOauthStateTracker()
	u := googleOauthConfig.AuthCodeURL(oauthStateTracker)
	fmt.Printf("Visit this URL on any device to log in:\n\n%s\n\nWaiting for authentication...", u)
	// open browser now
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
