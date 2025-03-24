package ormmock

import (
	"database/sql/driver"
	"regexp"
	"time"
)

// StringArg mock argument for variable string.
type StringArg struct{}

// TimeArg mock argument for variable time.
type TimeArg struct{}

// UUIDArg mock argument for variable uuid.
type UUIDArg struct{}

// Match determine if a string argument is string compatible.
func (a StringArg) Match(v driver.Value) bool {
	_, ok := v.(string)

	return ok
}

// Match determine if a time argument is time compatible.
func (a TimeArg) Match(v driver.Value) bool {
	_, ok := v.(time.Time)

	return ok
}

// Match determine if a uuid argument is uuid compatible.
func (u UUIDArg) Match(v driver.Value) bool {
	regex := `^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`

	if value, ok := v.(string); ok {
		return regexp.MustCompile(regex).MatchString(value)
	}

	return false
}
