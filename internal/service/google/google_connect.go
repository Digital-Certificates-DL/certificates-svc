package google

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"log"
	"net/http"
	"os"
	"strings"
)

func getClient(config *oauth2.Config, path string, code string) *http.Client {
	tokFile := path
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config, code)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config, code string) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	tok, err := config.Exchange(context.TODO(), code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open("token.json")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", "token.json")
	f, err := os.OpenFile("token.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func Connect(path, code string) (*http.Client, bool) {
	b, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Unable to read client secret file: %v", err)
		log.Printf("Could you continue to work without google drive ?")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if strings.ToLower(text) == "y\n" {
			return nil, false
		}
		return nil, true
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		log.Printf("Unable to parse client secret file to config: %v", err)
		log.Printf("Could you continue to work without google drive ?")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		if strings.ToLower(text) == "y" {
			return nil, false
		}
		return nil, true
	}
	return getClient(config, path, code), true

}
