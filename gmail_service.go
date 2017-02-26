package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/carlqt/guavatracker/pivotal"
	gmail "google.golang.org/api/gmail/v1"
)

func ListMessages(svc *gmail.Service) {
	acct := "me"
	msgService := svc.Users.Messages
	listMessages, err := msgService.List(acct).Q("label:unread from:notifications@typeform.com").Do()
	if err != nil {
		logger.Fatal(err)
	}

	for _, m := range listMessages.Messages {
		msg, _ := svc.Users.Messages.Get("me", m.Id).Format("full").Do()
		data, _ := base64.URLEncoding.DecodeString(msg.Payload.Body.Data)
		byteReader := bytes.NewReader(data)

		ticket := pivotal.NewTicket()
		ticket.Create()
		// pivotal.NewTicket
		// fmt.Println(body)
		body := parseHtml(byteReader)
		markAsRead(msg, msgService)
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

// Private functions

func parseHtml(b io.Reader) (body string) {
	var nodeAnswer string
	doc, err := goquery.NewDocumentFromReader(b)
	if err != nil {
		logger.Println(err)
	}

	doc.Find("li").Each(func(i int, s *goquery.Selection) {
		if em := s.Find("em"); em.Text() == "" {
			nodeQuestion := s.ChildrenFiltered("b").Text()

			if anchor := s.Find("a"); anchor.Text() != "" {
				nodeAnswer = linkURL(anchor)
			} else {
				nodeAnswer = strings.Replace(s.Text(), nodeQuestion, "", -1)
			}

			body += nodeQuestion + "\n" + nodeAnswer + "\n\n"
		}
	})

	return strings.TrimSpace(body)
}

func linkURL(s *goquery.Selection) (val string) {
	for _, node := range s.Nodes {
		for _, attr := range node.Attr {
			val = attr.Val
		}
	}

	return val
}

func markAsRead(m *gmail.Message, svc *gmail.UsersMessagesService) error {
	unread := []string{"UNREAD"}
	msgRequest := gmail.ModifyMessageRequest{RemoveLabelIds: unread}
	_, err := svc.Modify("me", m.Id, &msgRequest).Do()

	return err
}
