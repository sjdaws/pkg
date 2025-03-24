package connectionmock

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sjdaws/pkg/database"
)

// New mock database connection.
func New(t *testing.T) *database.Connection {
	t.Helper()

	// Create temporary database
	filename := t.TempDir() + "/test.db"

	_, err := os.Create(filename)
	require.NoError(t, err)

	connection, _ := database.Connect(false, "sqlite", "", filename, "", 0, "", "", "")

	return connection
}
