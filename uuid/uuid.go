package uuid

import (
	"database/sql/driver"
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

// Scan implements sql.Scanner so UUIDs can be read from databases transparently.
func (uuid *UUID) Scan(value interface{}) error {
	switch source := value.(type) {
	case nil:
		return nil

	case string:
		// if an empty UUID comes from a table, we return a null UUID
		if source == "" {
			return nil
		}

		// see Parse for required string format
		parsed, err := Parse(source)
		if err != nil {
			return errors.Wrap(err, "unable to parse uuid string '%s'", source)
		}

		*uuid = parsed

	case []byte:
		// if an empty UUID comes from a table, we return a null UUID
		if len(source) == 0 {
			return nil
		}

		// assumes a simple slice of bytes if 16 bytes otherwise attempts to parse
		const simpleSize = 16
		if len(source) != simpleSize {
			return uuid.Scan(string(source))
		}

		copy((*uuid)[:], source)

	default:
		return errors.New("invalid uuid type '%T'", source)
	}

	return nil
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

// Value implements sql.Valuer so that UUIDs can be written to databases transparently.
//
//nolint:ireturn // return value is determined by database driver interface.
func (uuid UUID) Value() (driver.Value, error) {
	return uuid.String(), nil
}
