package pivotal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

type Ticket struct {
	ProjectID    int    `json:"project_id"`
	Name         string `json:"name"`
	StoryType    string `json:"story_type"`
	url          string
	CurrentState string `json:"current_state"`
	Estimate     int    `json:"estimate"`
	Description  string `json:"description"`
	APIToken     string `json:"-"`
	ID           int    `json:"id"`
}

type errorResponse struct {
	code           string
	kind           string
	err            string
	generalProblem string
}

const (
	baseURL       = "https://www.pivotaltracker.com/services/v5/projects"
	apiToken      = "d1f26e67a5273ea05f4926a241f7355d"
	landingPageID = "1927121"
)

func NewTicket() Ticket {
	u, _ := url.Parse(baseURL)
	u.Path = path.Join(u.Path, landingPageID, "stories")

	return Ticket{
		ProjectID:    1927121,
		APIToken:     apiToken,
		url:          u.String(),
		StoryType:    "feature",
		CurrentState: "unstarted",
	}
}

func (t *Ticket) Create() {
	jsonVal, err := json.Marshal(t)
	fmt.Println(string(jsonVal[:]))
	if err != nil {
		log.Fatal(err)
	}

	client := new(http.Client)
	req, err := http.NewRequest("POST", t.url, bytes.NewBuffer(jsonVal))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-TrackerToken", t.APIToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	showJSON(resp)
	err = json.NewDecoder(resp.Body).Decode(t)
	if err != nil {
		log.Fatal(err)
	} else if t.ID == 0 {
		log.Fatal("Id is 0")
	}
	log.Println("Success!")
}

func (t *Ticket) Show() {
	t.APIToken = apiToken
	t.ProjectID = 1927121

	strProjectID := strconv.Itoa(t.ProjectID)

	req, err := http.NewRequest("GET", baseURL+strProjectID+"/stories/139780413", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("X-TrackerToken", t.APIToken)
	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	// showJSON(resp)
	err = json.NewDecoder(resp.Body).Decode(t)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", t)
}

func showJSON(r *http.Response) {
	io.Copy(os.Stdout, r.Body)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
