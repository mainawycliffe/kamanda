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
	"github.com/logrusorgru/aurora"
	"github.com/mainawycliffe/kamanda/configs"
	"github.com/mainawycliffe/kamanda/oauth/templates"
	"github.com/pkg/browser"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"
)

const (
	oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	port              = "8000"
	noPortURL         = "urn:ietf:wg:oauth:2.0:oob"
)

// googleOAuthScopes return a list of scopes used by kamanda
var googleOAuthScopes = []string{
	"https://www.googleapis.com/auth/userinfo.email",
	"https://www.googleapis.com/auth/cloud-platform",
	"https://www.googleapis.com/auth/firebase",
}

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

// getGoogleOAuthConfig create a OAuth Config object
func getGoogleOAuthConfig(port string) *oauth2.Config {
	redirectURL := fmt.Sprintf("http://localhost:%s", port)
	if port == "" {
		redirectURL = noPortURL
	}
	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     viper.GetString(configs.GoogleOAuthClientIDConfigKey),
		ClientSecret: viper.GetString(configs.GoogleOAuthClientSecretConfigKey),
		Scopes:       googleOAuthScopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://accounts.google.com/o/oauth2/auth",
			TokenURL:  "https://oauth2.googleapis.com/token",
			AuthStyle: oauth2.AuthStyleInParams,
		},
	}
}

// LoginWithLocalhost starts a server that can be used to capture OAUTH
// token from Google Auth Server
func LoginWithLocalhost(isCILogin bool) {
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
		templateData := map[string]string{
			"AppRepoURL": "https://github.com/mainawycliffe/kamanda",
			"AppName":    "Kamanda - Firebase CLI Companion Tool",
		}
		if r.FormValue("state") != oauthStateTracker {
			if err := writeHTMLOutput(w, templateData, templates.LoginFailureTemplate); err != nil {
				fmt.Printf("Error showing response: %s", err.Error())
			}
			return
		}
		data, err := getUserDataFromGoogle(googleOauthConfig, r.FormValue("code"))
		if err != nil {
			log.Println(err.Error())
			if err := writeHTMLOutput(w, templateData, templates.LoginFailureTemplate); err != nil {
				fmt.Printf("Error showing response: %s", err.Error())
			}
			return
		}
		if isCILogin {
			fmt.Fprint(os.Stdout, aurora.Sprintf(aurora.Green("\n\nSuccess! Use this token to login on a CI server:\n\n%s"), data.RefreshToken))
			fmt.Printf("\n\nExample: firebase deploy --token \"$FIREBASE_TOKEN\"\n")
		} else {
			if err = saveRefreshToken(data); err != nil {
				fmt.Printf("%s", err.Error())
				if err := writeHTMLOutput(w, templateData, templates.LoginFailureTemplate); err != nil {
					fmt.Printf("Error showing response: %s", err.Error())
				}
				return
			}
			fmt.Fprint(os.Stdout, aurora.Sprintf(aurora.Green("\n\nSuccess! Logged in as %s\n\n"), data.Email))
		}
		if err := writeHTMLOutput(w, templateData, templates.LoginSuccessTemplate); err != nil {
			fmt.Printf("Error showing response: %s", err.Error())
		}
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
func LoginWithoutLocalhost(isCILogin bool) error {
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
	if isCILogin {
		fmt.Fprint(os.Stdout, aurora.Sprintf(aurora.Green("\n\nSuccess! Use this token to login on a CI server:\n\n%s\n"), data.RefreshToken))
		return nil
	}
	if err = saveRefreshToken(data); err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, aurora.Sprintf(aurora.Green("\n\nSuccess! Logged in as %s\n\n"), data.Email))
	return nil
}

// RevokeRefreshToken revoke refresh token for google oauth when the user
// logouts
func RevokeRefreshToken() error {
	url := fmt.Sprintf("https://accounts.google.com/o/oauth2/revoke?token=%s", viper.GetString("FirebaseRefreshToken"))
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error sending http request: %w", err)
	}
	// remove configs locally whether successfull or not. May be not best
	// practice here
	err = configs.UnsetViperConfig(
		configs.FirebaseRefreshTokenViperConfigKey,
		configs.FirebaseLoggedInUserEmailViperConfigKey,
	)
	if err != nil {
		return fmt.Errorf("Error removing configs: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}
	var responseJSON map[string]string
	if err = json.Unmarshal(body, &responseJSON); err != nil {
		return fmt.Errorf("Error reading response body: %w", err)
	}
	var message string
	if responseJSON["error_description"] != "" {
		message = responseJSON["error_description"]
	} else {
		message = responseJSON["error"]
	}
	return fmt.Errorf("%s", message)
}

func saveRefreshToken(data *GetUserDataFromGoogleResponse) error {
	viper.Set(configs.FirebaseRefreshTokenViperConfigKey, data.RefreshToken)
	viper.Set(configs.FirebaseLoggedInUserEmailViperConfigKey, data.Email)
	err := viper.WriteConfig()
	if err != nil {
		return fmt.Errorf("An error occurred while saving refresh token: %w", err)
	}
	return nil
}

func writeHTMLOutput(w http.ResponseWriter, data interface{}, html string) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	errorTemplate, err := template.New("response.html").Parse(html)
	if err != nil {
		fmt.Fprintf(w, "Unable to load and parse failure template")
		return err
	}
	err = errorTemplate.Execute(w, data)
	if err != nil {
		fmt.Fprintf(w, "Unable to load and parse failure template")
		return err
	}
	return nil
}
