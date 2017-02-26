package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	gmail "google.golang.org/api/gmail/v1"

	"github.com/carlqt/guavatracker/pivotal"
	"github.com/pkg/browser"
)

var logger *log.Logger
var logFile *os.File

func getClient(ctx context.Context) *http.Client {
	config, _ := readConfig()
	token, err := authenticateAccount(config)
	if err != nil {
		logger.Fatalln(err)
	}

	return config.Client(ctx, token)
}

func authenticateAccount(config *oauth2.Config) (token *oauth2.Token, err error) {
	token, success := fetchTokenFromFile(TokenPath()) //If token is already saved in local

	// TODO: Need to handle expired tokens. If token is expire, launch browser
	if !success {
		logger.Println("failed to get token from file")
		//If failed to get token from file, get token from web
		token, err = fetchTokenFromWeb(config)
		saveToken(TokenPath(), token)
	}

	return token, err
}

func readConfig() (*oauth2.Config, error) {
	// parse client_id.json file
	b, err := ioutil.ReadFile("config/client_id.json")
	if err != nil {
		logger.Fatalf("Unable to read client secret file: %v", err)
	}
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope, gmail.MailGoogleComScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	return config, err
}

func TokenPath() string {
	user, _ := user.Current()
	tokenDir := filepath.Join(user.HomeDir, ".credentials")
	os.MkdirAll(tokenDir, 0700)

	return tokenDir + "guavatracker.json"
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func fetchTokenFromFile(path string) (*oauth2.Token, bool) {
	// Check if file exists, if YES, return token
	f, err := os.Open(path)
	if err != nil {
		return nil, false
	}

	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, true
}

func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func fetchTokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	browser.OpenURL(authURL)

	var code string
	fmt.Printf("Please enter access code in the browser: ")
	if _, err := fmt.Scan(&code); err != nil {
		logger.Fatalf("Unable to read authorization code %v", err)
	}

	token, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		logger.Fatalf("Unable to retrieve token from web %v", err)
	}
	return token, err
}

func init() {
	var err error
	logFile, err = os.OpenFile("server.log", os.O_RDWR|os.O_APPEND, 0660)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(logFile, os.Stdout)
	logger = log.New(mw, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func main() {
	defer logFile.Close()
	// ctx := context.Background()

	// Add function to check if access token is available in local ~/.credentials/*.json
	// client := getClient(ctx)

	// srv, err := gmail.New(client)
	// if err != nil {
	// 	logger.Fatalf("Unable to retrieve gmail Client %v", err)
	// }

	// ListMessages(srv)

	ticket := pivotal.NewTicket()
	ticket.Name = "Automation Test"
	ticket.Description = "This is created through carlqt bot automation"
	ticket.Create()

	logger.Println("Done")
}
