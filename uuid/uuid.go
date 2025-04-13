package uuid

import (
	"encoding/hex"
	"encoding/json"

	"github.com/google/uuid"

	"github.com/sjdaws/pkg/errors"
)

// UUID type.
type UUID uuid.UUID //nolint:recvcheck // UnmarshalJSON requires a pointer while remaining methods require value

// Nil empty UUID.
var Nil UUID //nolint:gochecknoglobals // Emulating uuid.Nil from google/uuid

// MustParse parse a UUID string and panic if it fails.
func MustParse(value string) UUID {
	parsed := uuid.MustParse(value)

	return UUID(parsed)
}

// New generate new UUID.
func New() UUID {
	return UUID(uuid.New())
}

// Parse a UUID string and return error if it fails.
func Parse(value string) (UUID, error) {
	parsed, err := uuid.Parse(value)
	if err != nil {
		return Nil, errors.Wrap(err, "unable to parse UUID")
	}

	return UUID(parsed), nil
}

// MarshalJSON gracefully handle marshaling to JSON.
func (u UUID) MarshalJSON() ([]byte, error) {
	if u == Nil {
		return json.Marshal(nil) //nolint:wrapcheck // Marshaling nil will never throw an error
	}

	return json.Marshal(u.String()) //nolint:wrapcheck // Marshaling a string will never throw an error
}

// String convert UUID to string.
func (u UUID) String() string {
	var buffer [36]byte

	hex.Encode(buffer[:], u[:4])
	buffer[8] = '-'
	hex.Encode(buffer[9:13], u[4:6])
	buffer[13] = '-'
	hex.Encode(buffer[14:18], u[6:8])
	buffer[18] = '-'
	hex.Encode(buffer[19:23], u[8:10])
	buffer[23] = '-'
	hex.Encode(buffer[24:], u[10:])

	return string(buffer[:])
}

// UnmarshalJSON gracefully handle unmarshaling invalid UUIDs.
func (u *UUID) UnmarshalJSON(data []byte) error {
	id, err := uuid.Parse(string(data))
	if err != nil {
		id = uuid.Nil
	}

	*u = UUID(id)

	return nil
}
