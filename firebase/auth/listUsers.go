package auth

import (
	"context"
	"fmt"

	"firebase.google.com/go/auth"
	"github.com/mainawycliffe/kamanda/firebase"
	"google.golang.org/api/iterator"
)

type ListUsersResponse struct {
	Users         []*auth.ExportedUserRecord
	nextPageToken string
}

// ListUsers get all users in firebase auth
func ListUsers(ctx context.Context, maxSize int, nextPageToken string) (ListUsersResponse, error) {
	client, err := firebase.Auth(ctx, "")
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("Error authenticating firebase account: %w", err)
	}
	usersIterator := client.Users(ctx, nextPageToken)
	users := make([]*auth.ExportedUserRecord, 0)
	for {
		user, err := usersIterator.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return ListUsersResponse{}, err
		}
		users = append(users, user)
	}
	response := ListUsersResponse{
		Users:         users,
		nextPageToken: usersIterator.PageInfo().Token,
	}
	return response, nil
}
