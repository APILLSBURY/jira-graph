package graph

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type issue struct {
	Key           string `json:"key"`
	Type          string `json:"type"`
	Summary       string `json:"summary"`
	Status        string `json:"status"`
	Assignee      string `json:"assignee"`
	AssigneeImage string `json:"assigneeImage"`
	Estimate      int    `json:"estimate"` // note that this doesn't differentiate between '0' and unset
	blockedByKeys []string
}

type epic issue

func (e epic) IsActive() bool {
	return e.Status == "Development Active"
}

type jiraClient struct {
	host          string
	user          string
	pass          string
	estimateField string
}

func (j jiraClient) Get(path string, q url.Values) (*http.Response, error) {
	baseURL := url.URL{
		Scheme: "https",
		Host:   j.host,
		Path:   path,
	}
	req, err := http.NewRequest("GET", baseURL.String(), nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(j.user, j.pass)
	req.URL.RawQuery = q.Encode()

	client := &http.Client{}
	return client.Do(req)
}

func (j jiraClient) Search(jql string, fields []string, startAt int) ([]byte, error) {
	q := url.Values{
		"jql":     []string{jql},
		"fields":  fields,
		"startAt": []string{strconv.Itoa(startAt)},
	}
	resp, err := j.Get("/rest/api/2/search", q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
