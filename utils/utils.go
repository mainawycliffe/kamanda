package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/logrusorgru/aurora"
	"github.com/mainawycliffe/kamanda/configs"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// FormatTimestampToDate takes in a timestamp in milliseconds since epoch and
// converts to a date format
func FormatTimestampToDate(timestamp int64, format string) string {
	nanoSeconds := timestamp * int64(time.Millisecond)
	tm := time.Unix(0, nanoSeconds)
	return tm.Format(format)
}

// IsUserLoggedIn checks whether a firebase refresh token is set in the configurations
func IsUserLoggedIn() bool {
	if viper.IsSet(configs.FirebaseRefreshTokenViperConfigKey) && viper.GetString(configs.FirebaseLoggedInUserEmailViperConfigKey) != "" {
		// todo: probably check the format of the token to ensure its correct
		return true
	}
	return false
}

// ParseStringToActualValue takes a string and converts the string the actual
// type of the string i.e. "true" => true
func ParseStringToActualValueType(input string) interface{} {
	if v, err := strconv.Atoi(input); err == nil {
		return v
	}
	if v, err := strconv.ParseFloat(input, 64); err == nil {
		return v
	}
	if v, err := strconv.ParseBool(input); err == nil {
		return v
	}
	return input
}

// PasswordGenerator generate password of given length
func PasswordGenerator(passwordLength int) string {
	rand.Seed(time.Now().UnixNano())
	letterBytes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()?.|")
	var password strings.Builder
	for i := 0; i < passwordLength; i++ {
		password.WriteRune(letterBytes[rand.Intn(len(letterBytes))])
	}
	return password.String()
}

// ProcessCustomClaimInput take in the input from cmd flags which a map of strings
// and convert it to a map of interface
func ProcessCustomClaimInput(input map[string]string) map[string]interface{} {
	customClaims := make(map[string]interface{})
	for k, v := range input {
		// @todo try and determine the value type and return it natively
		customClaims[k] = ParseStringToActualValueType(v)
	}
	return customClaims
}

// StdOutError print an error message to the standard out
func StdOutError(w io.Writer, format string, a ...interface{}) {
	m := aurora.Sprintf(aurora.Red(format), a...)
	fmt.Fprintf(w, "%s", m)
}

// StdOutSuccess print a success message to the standard out
func StdOutSuccess(w io.Writer, format string, a ...interface{}) {
	m := aurora.Sprintf(aurora.Green(format), a...)
	fmt.Fprintf(w, "%s", m)
}

// UnmarshalFormatFile read and unmarshal either a json/yaml file into a struct
func UnmarshalFormatFile(path string, extension string, v interface{}) error {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Error reading %s: %w", path, err)
	}
	switch extension {
	case "yaml":
		err = yaml.Unmarshal(content, v)
		if err != nil {
			return fmt.Errorf("Error decoding yaml: %w", err)
		}
	case "json":
		err = json.Unmarshal(content, v)
		if err != nil {
			return fmt.Errorf("Error decoding json: %w", err)
		}
	default:
		return fmt.Errorf("Unsupported file type")
	}
	return nil
}
