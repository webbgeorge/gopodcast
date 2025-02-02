package gopodcast

import (
	"fmt"
	"strings"
	"time"
)

// Bool is an alias for `bool` which supports unmarshalling
// "yes", "no" and similar values
type Bool bool

func (b *Bool) UnmarshalText(text []byte) error {
	*b = Bool(unmarshalBoolLike(text))
	return nil
}

func (b Bool) MarshalText() ([]byte, error) {
	if b {
		return []byte("true"), nil
	} else {
		return []byte("false"), nil
	}
}

// YesNo is an alias for `bool` which unmarshals bool-like
// values, and marshals to `yes` or `no`
type YesNo bool

func (b *YesNo) UnmarshalText(text []byte) error {
	*b = YesNo(unmarshalBoolLike(text))
	return nil
}

func (b YesNo) MarshalText() ([]byte, error) {
	if b {
		return []byte("yes"), nil
	} else {
		return []byte("no"), nil
	}
}

func unmarshalBoolLike(text []byte) bool {
	switch strings.ToLower(string(text)) {
	case "true", "t", "y", "yes":
		return true
	case "false", "f", "n", "no":
		return false
	default:
		return false
	}
}

type Time time.Time

func (t *Time) UnmarshalText(text []byte) error {
	formats := []string{
		time.RFC822,
		time.RFC822Z,
		time.RFC1123,
		time.RFC1123Z,
		"Mon, _2 Jan 2006 15:04:05 MST",   // 1123, with _2
		"Mon, _2 Jan 2006 15:04:05 -0700", // 1123Z, with _2
	}
	var tt time.Time
	var err error
	for _, f := range formats {
		tt, err = time.Parse(f, string(text))
		if err == nil {
			*t = Time(tt)
			return nil
		}
	}
	return fmt.Errorf("failed to parse time '%s'", string(text))
}

func (t Time) MarshalText() ([]byte, error) {
	return []byte(time.Time(t).Format(time.RFC1123)), nil
}
