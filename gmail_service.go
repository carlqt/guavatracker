package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/net/html"

	gmail "google.golang.org/api/gmail/v1"
)

func ListMessages(svc *gmail.Service) {
	acct := "me"
	listMessages, err := svc.Users.Messages.List(acct).Q("label:unread from:notifications@typeform.com").Do()
	if err != nil {
		logger.Fatal(err)
	}

	for _, m := range listMessages.Messages {
		msg, _ := svc.Users.Messages.Get("me", m.Id).Format("full").Do()
		data, _ := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		byteData := bytes.NewReader(data)
		displayHtml(byteData)
	}
}
func ListLabels(srv *gmail.Service) {
	acct := "me"
	r, err := srv.Users.Labels.List(acct).Do()
	if err != nil {
		logger.Fatalf("Unable to retrieve labels. %v", err)
	}
	if len(r.Labels) > 0 {
		fmt.Print("Labels:\n")
		for _, l := range r.Labels {
			fmt.Printf("- %s\n", l.Name)
		}
	} else {
		fmt.Print("No labels found.")
	}

}

func PrettyPrint(in *gmail.Message) {
	fmt.Printf("%+v\n", in)
}

func displayHtml(b io.Reader) {
	doc := html.NewTokenizer(b)

	for {
		token := doc.Next()
		switch {
		case token == html.ErrorToken:
			return
		case token == html.TextToken:
			t := doc.Token()
			fmt.Println(t.Data)
		}
	}
}
