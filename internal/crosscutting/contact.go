package crosscutting

import "regexp"

var emailRegexPattern = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

type Email string

func (value Email) String() string {
	return string(value)
}

func (value Email) IsValid() bool {
	return emailRegexPattern.MatchString(value.String())
}
