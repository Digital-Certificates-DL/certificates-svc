package google

import (
	"testing"
)

var (
	path = "./client_secret.json"
	name = "test.jpg"
)

func TestConnect(t *testing.T) {
	Connect("", path, name)
	//require.Equal(t, "\x18Bitcoin Signed Mssage:\n\x0Erandom message", message)
}
