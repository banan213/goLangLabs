package images

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://api.unsplash.com"

type Photo struct {
	ID          string `json:"id"`
	Urls        Urls   `json:"urls"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	Likes       int    `json:"likes"`
	User        User   `json:"user"`
}

type Urls struct {
	Full string `json:"full"`
}

type User struct {
	Name string `json:"name"`
}

type photoSearchResponse struct {
	Total      int     `json:"total"`
	TotalPages int     `json:"total_pages"`
	Results    []Photo `json:"results"`
}

type Client struct {
	apiKey string
}

func New(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

func (c *Client) Random() (*Photo, error) {
	u := fmt.Sprintf("%s/photos/random?client_id=%s", baseURL, c.apiKey)

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("random request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unsplash returned status: %s", resp.Status)
	}

	var p Photo
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, fmt.Errorf("decode random: %w", err)
	}
	return &p, nil
}

func (c *Client) Search(query string) ([]Photo, error) {
	q := url.QueryEscape(query)
	u := fmt.Sprintf("%s/search/photos?query=%s&client_id=%s", baseURL, q, c.apiKey)

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("search request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unsplash returned status: %s", resp.Status)
	}

	var r photoSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("decode search: %w", err)
	}
	return r.Results, nil
}
