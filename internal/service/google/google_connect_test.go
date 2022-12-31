package google

import (
	"testing"
)

var (
	path     = "./client_secret.json"
	login    = "mark.cherepovskyi@nure.ua"
	password = "Markys1256"
)

func TestConnect(t *testing.T) {
	Connect(path, login, password)
	//require.Equal(t, "\x18Bitcoin Signed Mssage:\n\x0Erandom message", message)
}
