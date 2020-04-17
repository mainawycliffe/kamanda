package auth

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"
	"github.com/mainawycliffe/kamanda/firebase"
)

type FindUserCriteria int

const (
	ByUserUIDCriteria   FindUserCriteria = 0
	ByUserEmailCriteria FindUserCriteria = 1
	ByUserPhoneCriteria FindUserCriteria = 2
)

//IsValid check if user criteria is valid i.e by uid, email or phone
func (c FindUserCriteria) IsValid() bool {
	if c != ByUserUIDCriteria && c != ByUserEmailCriteria && c != ByUserPhoneCriteria {
		return false
	}
	return true
}

// GetUser find a user by either uid, email or phone number
func GetUser(ctx context.Context, query string, criteria FindUserCriteria) (*auth.UserRecord, error) {
	if isValid := criteria.IsValid(); !isValid {
		return nil, fmt.Errorf("Invalid find user criteria.")
	}
	client, err := firebase.Auth(ctx, "", "")
	if err != nil {
		return nil, fmt.Errorf("Error authenticating firebase account: %w", err)
	}
	if criteria == ByUserEmailCriteria {
		return client.GetUserByEmail(ctx, query)
	}
	if criteria == ByUserPhoneCriteria {
		return client.GetUserByPhoneNumber(ctx, query)
	}
	return client.GetUser(ctx, query)
}
