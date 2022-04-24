package danbooru

import (
	"strconv"
)

// ID type
type ID uint64

// IsZero returns true if ID is not defined
func (t ID) IsZero() bool {
	return t == 0
}

// String returns string representation of the ID
func (t ID) String() string {
	return strconv.FormatUint(uint64(t), 10)
}

// MarshalJSON serialized ID value to JSON using null to represent zero value
func (t ID) MarshalJSON() ([]byte, error) {
	if t == 0 {
		return []byte("null"), nil
	}

	return ([]byte)(strconv.FormatUint(uint64(t), 10)), nil
}

// UnmarshalJSON unmarshals ID value from JSON, mapping null to 0
func (t *ID) UnmarshalJSON(bs []byte) error {
	if t == nil {
		return nil
	}

	s := string(bs)

	if s == "null" || len(bs) == 0 {
		*t = 0

		return nil
	}

	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}

	*t = ID(v)

	return nil
}
