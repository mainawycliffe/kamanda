package firebase

import "firebase.google.com/go/auth"

// Error custom error for firebase
type Error struct {
	OriginalError error
}

// NewError return our error object with the original error inside
func NewError(err error) Error {
	return Error{
		OriginalError: err,
	}
}

// Error parse the original error and return a coherent message
// @todo: expand this to parse all Firebase errors and return a readable message
func (err Error) Error() string {
	if auth.IsEmailAlreadyExists(err.OriginalError) {
		return "Email already exists"
	}
	if auth.IsPhoneNumberAlreadyExists(err.OriginalError) {
		return "Phone number already exists"
	}
	if auth.IsUIDAlreadyExists(err.OriginalError) {
		return "UID already exists"
	}
	if auth.IsUserNotFound(err.OriginalError) {
		return "User not found"
	}
	if auth.IsUnknown(err.OriginalError) {
		return "Unknown error"
	}
	return err.OriginalError.Error()
}

//IsUserNotFound check if an error is of type UserNotFound
func IsUserNotFound(err error) bool {
	return auth.IsUserNotFound(err)
}
