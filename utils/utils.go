package utils

	"fmt"
	"github.com/logrusorgru/aurora"
// ProcessCustomClaimInput take in the input from cmd flags which a map of strings
// and convert it to a map of interface
func ProcessCustomClaimInput(input map[string]string) map[string]interface{} {
	customClaims := make(map[string]interface{})
	for k, v := range input {
		// @todo try and determine the value type and return it natively
		customClaims[k] = v
	}
	return customClaims
}

// StdOutError print an error message to the standard out
func StdOutError(format string, a ...interface{}) {
	m := aurora.Sprintf(aurora.Red(format), a)
	fmt.Fprintf(os.Stdout, "%s\n", m)
}

// StdOutSuccess print a success message to the standard out
func StdOutSuccess(format string, a ...interface{}) {
	m := aurora.Sprintf(aurora.Green(format), a)
	fmt.Fprintf(os.Stdout, "%s\n", m)
}

