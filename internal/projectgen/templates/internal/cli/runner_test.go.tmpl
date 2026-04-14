package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConsoleWriterDisablesColorForNonInteractiveOutput(t *testing.T) {
	var stderr bytes.Buffer

	writer := newConsoleWriter(&stderr)

	assert.True(t, writer.NoColor)
}

func TestNewConsoleWriterDisablesColorForNonTerminalFile(t *testing.T) {
	file, err := os.CreateTemp(t.TempDir(), "logs-*")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, file.Close())
	}()

	writer := newConsoleWriter(file)

	assert.True(t, writer.NoColor)
}
