package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"

	"github.com/pkg/browser"
)

var logger *log.Logger
var logFile *os.File

func init() {
	var err error
	logFile, err = os.OpenFile("server.log", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(logFile, os.Stdout)
	logger = log.New(mw, "", log.Ldate|log.Ltime)
}

func main() {
	defer logFile.Close()

	// Add function to check if access token is available in local ~/.credentials/*.json
	token, _ := authenticateAccount()

	logger.Println("the token is:" + token.AccessToken)
}

func authenticateAccount() (*oauth2.Token, error) {
	config, _ := readConfig()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	browser.OpenURL(authURL)

	var code string
	fmt.Printf("Please enter access code in the browser: ")
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}

	// TODO: Add function to store token in local
	return tok, err
}

func readConfig() (*oauth2.Config, error) {
	// parse client_id.json file
	b, err := ioutil.ReadFile("config/client_id.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}
	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/gmail-go-quickstart.json
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config, err
}
