package auth

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"
	"github.com/mainawycliffe/kamanda/firebase"
)

type FirebaseUser struct {
	UID                       string                     `json:"uid" yaml:"uid"`
	Email                     string                     `json:"email" yaml:"email"`
	ShouldUpdateEmailVerified bool                       `json:"-" yaml:"-"`
	EmailVerified             bool                       `json:"email_verified" yaml:"email_verified"`
	PhoneNumber               string                     `json:"phone" yaml:"phone"`
	Password                  string                     `json:"password" yaml:"password"`
	DisplayName               string                     `json:"name" yaml:"name" `
	ShouldUpdateDisabled      bool                       `json:"-" yaml:"-"`
	Disabled                  bool                       `json:"disabled" yaml:"disabled"`
	PhotoURL                  string                     `json:"photo_url" yaml:"photo_url"`
	CustomClaims              []FirebaseUserCustomClaims `json:"custom_claims" yaml:"custom_claims"`
}

type FirebaseUserCustomClaims struct {
	Key   string      `json:"key" yaml:"key"`
	Value interface{} `json:"value" yaml:"value"`
}

// NewFirebaseUser create a new firebase user using Email/Password Auth Provider
func NewFirebaseUser(ctx context.Context, user *FirebaseUser) (*auth.UserRecord, error) {
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
	client, err := firebase.Auth(ctx, "", "")
	if err != nil {
		return nil, fmt.Errorf("Error authenticating firebase account: %w", err)
	}
	u, err := client.CreateUser(ctx, params)
	if err != nil {
		return nil, firebase.NewError(err)
	}
	return u, nil
}

// UpdateFirebaseUser update a user details on firebase.
func UpdateFirebaseUser(ctx context.Context, UID string, user *FirebaseUser) (*auth.UserRecord, error) {
	params := &auth.UserToUpdate{}
	// incase you want to use a custom UID instead a random one
	if user.Email != "" {
		params = params.Email(user.Email)
	}
	if user.Password != "" {
		params = params.Password(user.Password)
	}
	if user.DisplayName != "" {
		params = params.DisplayName(user.DisplayName)
	}
	if user.PhoneNumber != "" {
		params = params.PhoneNumber(user.PhoneNumber)
	}
	if user.PhotoURL != "" {
		params = params.PhotoURL(user.PhotoURL)
	}
	if user.ShouldUpdateDisabled {
		params = params.Disabled(user.Disabled)
	}
	if user.ShouldUpdateEmailVerified {
		params = params.EmailVerified(user.EmailVerified)
	}
	client, err := firebase.Auth(ctx, "", "")
	if err != nil {
		return nil, fmt.Errorf("Error authenticating firebase account: %w", err)
	}
	u, err := client.UpdateUser(ctx, UID, params)
	if err != nil {
		return nil, firebase.NewError(err)
	}
	return u, nil
}

func AddCustomClaimToFirebaseUser(ctx context.Context, uid string, customClaims map[string]interface{}) error {
	client, err := firebase.Auth(ctx, "", "")
	if err != nil {
		return err
	}
	if err := client.SetCustomUserClaims(ctx, uid, customClaims); err != nil {
		return firebase.NewError(err)
	}
	return nil
}

// RemoveCustomClaimFromUser remove a custom claim or custom claims from a users
// account
func RemoveCustomClaimFromUser(ctx context.Context, UID string, keys []string) (*auth.UserRecord, error) {
	customClaimsToUnset := map[string]interface{}{}
	for _, v := range keys {
		customClaimsToUnset[v] = nil
	}
	client, err := firebase.Auth(ctx, "", "")
	if err != nil {
		return nil, fmt.Errorf("Error authenticating firebase account: %w", err)
	}
	params := &auth.UserToUpdate{}
	params.CustomClaims(customClaimsToUnset)
	user, err := client.UpdateUser(ctx, UID, params)
	if err != nil {
		return nil, firebase.NewError(err)
	}
	return user, nil
}

func DeleteFirebaseUser(ctx context.Context, uid string) error {
	if uid == "" {
		return fmt.Errorf("The UID of the user can not be empty")
	}
	client, err := firebase.Auth(ctx, "", "")
	if err != nil {
		return fmt.Errorf("Error authenticating firebase account: %w", err)
	}
	err = client.DeleteUser(ctx, uid)
	if err != nil {
		if auth.IsUserNotFound(err) {
			return fmt.Errorf("User not found!")
		}
		return fmt.Errorf("An unknown error: %w", err)
	}
	return nil
}
