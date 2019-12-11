package auth

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"
	"github.com/mainawycliffe/kamanda/firebase"
)

type NewUser struct {
	UID           string
	Email         string
	EmailVerified bool
	PhoneNumber   string
	Password      string
	DisplayName   string
	Disabled      bool
	PhotoURL      string
}

// NewFirebaseUser create a new firebase user using Email/Password Auth Provider
func NewFirebaseUser(ctx context.Context, user *NewUser) (*auth.UserRecord, error) {
	params := (&auth.UserToCreate{}).
		Email(user.Email).
		EmailVerified(user.EmailVerified).
		Password(user.Password).
		Disabled(user.Disabled)
	// incase you want to use a custom UID instead a random one
	if user.UID != "" {
		params = params.UID(user.UID)
	}
	if user.DisplayName != "" {
		params = params.DisplayName(user.DisplayName)
	}
	if user.PhoneNumber != "" {
		params = params.PhoneNumber(user.PhotoURL)
	}
	if user.PhotoURL != "" {
		params = params.PhotoURL(user.PhotoURL)
	}
	client, err := firebase.Auth(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("Error authenticating firebase account: %w", err)
	}
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		return nil, firebase.NewError(err)
	}
	return u, nil
}

func AddCustomClaimsToFirebaseUser(ctx context.Context, uid string, listOfClaims *[]map[string]interface{}) error {
	client, err := firebase.Auth(ctx, "")
	if err != nil {
		return err
	}
	for _, v := range *listOfClaims {
		err := client.SetCustomUserClaims(ctx, uid, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteFirebaseUser(ctx context.Context, uid string) error {
	if uid == "" {
		return fmt.Errorf("The UID of the user can not be empty")
	}
	client, err := firebase.Auth(ctx, "")
	if err != nil {
		return fmt.Errorf("Error authenticating firebase account: %w", err)
	}
	err = client.DeleteUser(ctx, uid)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return fmt.Errorf("User not found!")
		}
		return fmt.Errorf("An unnkown error: %w", err)
	}
	return nil
}

// ListAllFirebaseUsers get all users in firebase auth
func ListAllFirebaseUsers(ctx context.Context, maxResults uint32, nextPageToken string) error {
	panic("not implemented")
}
