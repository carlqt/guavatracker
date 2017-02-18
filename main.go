package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	ctx := context.Background()

	// Add function to check if access token is available in local ~/.credentials/*.json
	client := getClient(ctx)

	srv, err := gmail.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve gmail Client %v", err)
	}

	user := "me"
	r, err := srv.Users.Labels.List(user).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels. %v", err)
	}
	if len(r.Labels) > 0 {
		fmt.Print("Labels:\n")
		for _, l := range r.Labels {
			fmt.Printf("- %s\n", l.Name)
		}
	} else {
		fmt.Print("No labels found.")
	}

	// logger.Println("the token is:" + token.AccessToken)
}

func getClient(ctx context.Context) *http.Client {
	config, _ := readConfig()
	token, err := authenticateAccount(config)
	if err != nil {
		log.Fatal(err)
	}

	return config.Client(ctx, token)
}

func authenticateAccount(config *oauth2.Config) (*oauth2.Token, error) {
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
