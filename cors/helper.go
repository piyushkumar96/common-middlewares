package cors

import "regexp"

func MatchStringWithRegex(pattern, inputStr string) (bool, error) {
	// Compile the regex pattern
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return false, err
	}
	return regex.MatchString(inputStr), nil
}
