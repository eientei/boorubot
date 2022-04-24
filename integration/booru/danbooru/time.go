package danbooru

import (
	"strings"
	"time"
)

// Time with danbooru timestamp format (RFC3339)
type Time time.Time

// MarshalJSON returns danbooru-formatted time
func (t Time) MarshalJSON() ([]byte, error) {
	return ([]byte)(`"` + time.Time(t).Format(time.RFC3339) + `"`), nil
}

// UnmarshalJSON parses danbooru-formatted time
func (t *Time) UnmarshalJSON(bs []byte) error {
	if t == nil {
		return nil
	}

	s := string(bs)

	*t = Time{}

	if !strings.HasPrefix(s, `"`) || !strings.HasSuffix(s, `"`) {
		return nil
	}

	s = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)

	ts, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return err
	}

	*t = Time(ts)

	return nil
}
