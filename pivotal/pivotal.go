package pivotal

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Ticket struct {
	ProjectID    int    `json:"project_id"`
	Name         string `json:"name"`
	StoryType    string `json:"story_type"`
	URL          string `json:"url"`
	CurrentState string `json:"current_state"`
	Estimate     int    `json:"estimate"`
	Description  string `json:"description"`
	APIToken     string `json:"-"`
	Kind         string `json:"kind"`
	ID           int    `json:"id"`
}

const (
	baseURL  = "https://www.pivotaltracker.com/services/v5/projects/"
	apiToken = "d1f26e67a5273ea05f4926a241f7355d"
)

func NewTicket() {

}

func (t *Ticket) Create() {

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
