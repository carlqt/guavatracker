package pivotal

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"

	"github.com/carlqt/guavatracker/appconfig"
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
	baseURL = "https://www.pivotaltracker.com/services/v5/projects"
)

func NewTicket(c *appconfig.Config) Ticket {
	u, _ := url.Parse(baseURL)
	u.Path = path.Join(u.Path, c.ProjectID, "stories")
	id, _ := strconv.Atoi(c.ProjectID)

	return Ticket{
		ProjectID:    id,
		APIToken:     c.PivotalToken,
		url:          u.String(),
		StoryType:    "feature",
		CurrentState: "unstarted",
	}
}

func (t *Ticket) Create() {
	jsonVal, err := json.Marshal(t)
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
	err = json.NewDecoder(resp.Body).Decode(t)
	if err != nil {
		log.Fatal(err)
	} else if t.ID == 0 {
		log.Fatal("Id is 0")
	}
	log.Println("Success!")
}

// helper function for debugger purposes
func showJSON(r *http.Response) {
	io.Copy(os.Stdout, r.Body)
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}
