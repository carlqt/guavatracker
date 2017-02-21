package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"

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
		byteReader := bytes.NewReader(data)
		body := parseHtml(byteReader)
		fmt.Println(body)
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

func parseHtml(b io.Reader) (body string) {
	doc, err := goquery.NewDocumentFromReader(b)
	if err != nil {
		logger.Println(err)
	}

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		nodeQuestion := s.ChildrenFiltered("b").Text()
		nodeAnswer := strings.Replace(s.Text(), nodeQuestion, "", -1)

		body += nodeQuestion + "\n" + nodeAnswer + "\n\n"
	})

	return strings.TrimSpace(body)
}
