package gopodcast

import (
	"strings"
)

// FlexBool is an alias for `bool` which supports unmarshalling
// "yes", "no" and similar values
type FlexBool bool

func (b *FlexBool) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case "true", "t", "y", "yes":
		*b = true
	case "false", "f", "n", "no":
		*b = false
	default:
		*b = false
	}
	return nil
}

func (b FlexBool) MarshalText() ([]byte, error) {
	if b {
		return []byte("true"), nil
	} else {
		return []byte("false"), nil
	}
}
