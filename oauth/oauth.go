// Package oauth get refresh token and save for future use
package oauth

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	port              = "8000"
	noPortURL         = "urn:ietf:wg:oauth:2.0:oob"
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

// googleOAuthScopes return a list of scopes used by kamanda
func googleOAuthScopes() []string {
	return []string{
		"https://www.googleapis.com/auth/userinfo.email",
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/firebase",
	}
}

// getGoogleOAuthConfig create a OAuth Config object
func getGoogleOAuthConfig(port string) *oauth2.Config {
	redirectURL := fmt.Sprintf("http://localhost:%s", port)
	if port == "" {
		redirectURL = noPortURL
	}
	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     viper.GetString("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: viper.GetString("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes:       googleOAuthScopes(),
		Endpoint:     google.Endpoint,
	}
}

// LoginWithLocalhost starts a server that can be used to capture OAUTH
// token from Google Auth Server
func LoginWithLocalhost() {
	googleOauthConfig := getGoogleOAuthConfig(port)
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

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		defer cancel()
		if r.FormValue("state") != oauthStateTracker {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			t, err := template.ParseFiles("templates/loginFailure.html")
			if err != nil {
				fmt.Fprintf(w, "Unable to load and parse failure template")
				return
			}
			err = t.Execute(w, nil)
			if err != nil {
				fmt.Fprintf(w, "Unable to load and parse failure template")
				return
			}
			return
		}
		data, err := getUserDataFromGoogle(googleOauthConfig, r.FormValue("code"))
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			t, err := template.ParseFiles("templates/loginFailure.html")
			if err != nil {
				fmt.Fprintf(w, "Unable to load and parse failure template")
				return
			}
			err = t.Execute(w, nil)
			if err != nil {
				fmt.Fprintf(w, "Unable to load and parse failure template")
				return
			}
			return
		}
		refreshTokenObject := RefreshToken{
			ClientID:     viper.GetString("GOOGLE_OAUTH_CLIENT_ID"),
			ClientSecret: viper.GetString("GOOGLE_OAUTH_CLIENT_SECRET"),
			RefreshToken: data.RefreshToken,
			Type:         "authorized_user",
		}
		err = SaveRefreshToken(refreshTokenObject)
		if err != nil {
			log.Println(err.Error())
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			t, err := template.ParseFiles("templates/loginFailure.html")
			if err != nil {
				fmt.Fprintf(w, "Unable to load and parse failure template")
				return
			}
			err = t.Execute(w, nil)
			if err != nil {
				fmt.Fprintf(w, "Unable to load and parse failure template")
				return
			}
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		t, err := template.ParseFiles("templates/loginSuccess.html")
		if err != nil {
			fmt.Fprintf(w, "Unable to load and parse success template %v", err)
			return
		}
		err = t.Execute(w, nil)
		if err != nil {
			fmt.Fprintf(w, "Unable to load and parse success template: %v", err)
			return
		}
		fmt.Printf("\n\nLogin Successful!\n\n")
		cancel()
	})
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	// Shutdown the server when the context is canceled
	if err := httpServer.Shutdown(context.Background()); err != nil {
		log.Fatal(err)
	}
}

// LoginWithoutLocalhost login without localhost localhost server, suitable in
// environments without a GUI. The user enters the authorization code manually
func LoginWithoutLocalhost() error {
	googleOauthConfig := getGoogleOAuthConfig("")
	oauthStateTracker := generateOauthStateTracker()
	u := googleOauthConfig.AuthCodeURL(oauthStateTracker)
	fmt.Printf("Visit this URL on any device to log in:\n\n%s\n\nWaiting for authentication...", u)
	// open browser now
	_ = browser.OpenURL(u)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\nEnter the authorization code here: ")
	code, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("An error occurred while reading user input: %w", err)
	}
	data, err := getUserDataFromGoogle(googleOauthConfig, code)
	if err != nil {
		return fmt.Errorf("An error occurred while exchanging code with token: %w", err)
	}
	refreshTokenObject := RefreshToken{
		ClientID:     viper.GetString("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: viper.GetString("GOOGLE_OAUTH_CLIENT_SECRET"),
		RefreshToken: data.RefreshToken,
		Type:         "authorized_user",
	}
	err = SaveRefreshToken(refreshTokenObject)
	if err != nil {
		return fmt.Errorf("An error occurred while saving refresh token: %w", err)
	}
	fmt.Fprintf(os.Stdout, "\n\nSuccess! Logged in as %s\n\n", data.Email)
	return nil
}
