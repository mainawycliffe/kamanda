package firebase

import "firebase.google.com/go/auth"

// NewError return our error object with the original error inside
func NewError(err error) FirebaseError {
	return FirebaseError{
		FirebaseError: err,
	}
}

// FirebaseError custom error for firebase
type FirebaseError struct {
	FirebaseError error
}

// Error parse the original error and return a coherent message
// @todo: expand this to parse all Firebase errors and return a readable message
func (err FirebaseError) Error() string {
	if auth.IsEmailAlreadyExists(err.FirebaseError) {
		return "Email already exists"
	}
	if auth.IsPhoneNumberAlreadyExists(err.FirebaseError) {
		return "Phone number already exists"
	}
	if auth.IsUIDAlreadyExists(err.FirebaseError) {
		return "UID already exists"
	}
	if auth.IsUserNotFound(err.FirebaseError) {
		return "User was not found"
	}
	if auth.IsUnknown(err.FirebaseError) {
		return "Unknown error"
	}
	return err.FirebaseError.Error()
}

//IsUserNotFound check if an error is of type UserNotFound
func IsUserNotFound(err error) bool {
	return auth.IsUserNotFound(err)
}
