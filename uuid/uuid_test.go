package uuid_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/uuid"
)

func TestNew(t *testing.T) {
	t.Parallel()

	actual := uuid.New()

	assert.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", actual.String())
}

func TestMustParse(t *testing.T) {
	t.Parallel()

	actual := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	assert.Equal(t, uuid.MustParse("00000000-0000-0000-0000-000000000001"), actual)
}

func TestMustParse_InvalidUUID(t *testing.T) {
	t.Parallel()

	defer func() {
		recovery := recover()
		require.Equal(t, "uuid: Parse(1): invalid UUID length: 1", recovery)
	}()

	actual := uuid.MustParse("1")

	assert.Equal(t, uuid.Nil, actual)
}

func TestParse(t *testing.T) {
	t.Parallel()

	actual, err := uuid.Parse("00000000-0000-0000-0000-000000000001")
	require.NoError(t, err)

	assert.Equal(t, uuid.MustParse("00000000-0000-0000-0000-000000000001"), actual)
}

func TestParse_InvalidUUID(t *testing.T) {
	t.Parallel()

	actual, err := uuid.Parse("1")
	require.Error(t, err)

	require.EqualError(t, err, "unable to parse UUID: invalid UUID length: 1")
	assert.Equal(t, uuid.Nil, actual)
}

func TestUUID_MarshalJSON(t *testing.T) {
	t.Parallel()

	type Test struct {
		Test uuid.UUID `json:"uuid"`
	}

	target := Test{
		Test: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	}

	actual, err := json.Marshal(target)
	require.NoError(t, err)

	assert.JSONEq(t, `{"uuid":"00000000-0000-0000-0000-000000000001"}`, string(actual))
}

func TestUUID_MarshalJSON_NilUUID(t *testing.T) {
	t.Parallel()

	type Test struct {
		Test uuid.UUID `json:"uuid"`
	}

	target := Test{
		Test: uuid.MustParse("00000000-0000-0000-0000-000000000000"),
	}

	actual, err := json.Marshal(target)
	require.NoError(t, err)

	assert.JSONEq(t, `{"uuid":null}`, string(actual))
}

func TestUUID_UnmarshalJSON(t *testing.T) {
	t.Parallel()

	type Test struct {
		Test uuid.UUID `json:"uuid"`
	}

	payload := `{"uuid":"00000000-0000-0000-0000-000000000001"}`

	var target Test

	err := json.Unmarshal([]byte(payload), &target)
	require.NoError(t, err)

	expected := Test{
		Test: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
	}

	assert.Equal(t, expected, target)
}

func TestUUID_UnmarshalJSON_InvalidUUID(t *testing.T) {
	t.Parallel()

	type Test struct {
		Test uuid.UUID `json:"uuid"`
	}

	payload := `{"uuid":"1"}`

	var target Test

	err := json.Unmarshal([]byte(payload), &target)
	require.NoError(t, err)

	expected := Test{
		Test: uuid.Nil,
	}

	assert.Equal(t, expected, target)
}

func TestUUID_String(t *testing.T) {
	t.Parallel()

	actual := uuid.New()

	assert.Regexp(t, "[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}", actual.String())
}

func TestUUID_Scan(t *testing.T) {
	t.Parallel()

	testcases := map[string]struct {
		error    string
		expected uuid.UUID
		value    any
	}{
		"16 byte": {
			expected: uuid.MustParse("30303030-3030-3030-3030-303030303031"),
			value:    []byte("0000000000000001"),
		},
		"byte": {
			expected: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			value:    []byte("00000000000000000000000000000001"),
		},
		"empty bytes": {
			expected: uuid.Nil,
			value:    []byte{},
		},
		"empty string": {
			expected: uuid.Nil,
			value:    "",
		},
		"invalid": {
			error: "unable to parse uuid string '00000000-0000-0000-0000-0000000001': invalid UUID length: 34",
			value: "00000000-0000-0000-0000-0000000001",
		},
		"nil": {
			expected: uuid.Nil,
			value:    nil,
		},
		"string": {
			expected: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			value:    "00000000-0000-0000-0000-000000000001",
		},
		"wrong type": {
			error: "invalid uuid type 'int'",
			value: 5,
		},
	}

	for name, testcase := range testcases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual := uuid.UUID{}
			err := (&actual).Scan(testcase.value)

			if testcase.error != "" {
				require.Error(t, err)
				require.EqualError(t, err, testcase.error)

				return
			}

			require.NoError(t, err)
			require.Equal(t, testcase.expected, actual)
		})
	}
}

func TestUUID_Value(t *testing.T) {
	t.Parallel()

	value := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	actual, err := value.Value()
	require.NoError(t, err)

	assert.Equal(t, value.String(), actual)
}
