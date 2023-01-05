package signature

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateMagicMessage(t *testing.T) {
	message := CreateMagicMessage("random message")
	require.Equal(t, "\x18Bitcoin Signed Message:\n\x0Erandom message", message)
}
