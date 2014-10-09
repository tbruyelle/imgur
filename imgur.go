package imgur

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.imgur.com/3/"
)

type Client struct {
	clientId string
	baseUrl  *url.URL
	client   *http.Client
}

type Response struct {
	Success bool `json:"success"`
	Status  int  `json:"status"`
}

type SearchResponse struct {
	Response
	Data []Image `json:"data"`
}

type Image struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DateTime    int    `json:"datetime"`
	Type        string `json:"type"`
	Animated    bool   `json:"animated"`
	Nsfw        bool   `json:"nsfw"`
	Link        string `json:"link"`
}

func NewClient(clientId string) *Client {
	baseUrl, _ := url.Parse(defaultBaseURL)

	c := &Client{
		clientId: clientId,
		client:   http.DefaultClient,
		baseUrl:  baseUrl,
	}
	return c
}

func (c *Client) NewRequest(method, uri string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}

	u := c.baseUrl.ResolveReference(rel)
	buf := new(bytes.Buffer)
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Client-Id "+c.clientId)
	return req, nil
}

func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if c := resp.StatusCode; c < 200 || c > 299 {
		return resp, fmt.Errorf("Server returns status %d", c)
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

func (c *Client) Search(q string) (*SearchResponse, *http.Response, error) {
	uri := fmt.Sprintf("gallery/search/viral/all/1?q=%s", q)
	req, err := c.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, nil, err
	}

	searchResp := new(SearchResponse)
	resp, err := c.Do(req, searchResp)
	if err != nil {
		return nil, resp, err
	}
	return searchResp, resp, nil
}
