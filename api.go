package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	New  StoriesType = "new"
	Best StoriesType = "best"
	Top  StoriesType = "top"
	Ask  StoriesType = "ask"
	Job  StoriesType = "job"
	Poll StoriesType = "poll"
)

var (
	ErrCode          = errors.New("invalid HTTP status code")
	ErrCreateRequest = errors.New("create HTTP request")
	ErrSendRequest   = errors.New("send HTTP request")
)

type StoriesType string

type Story struct {
	ID          int         `json:"id"`
	Type        StoriesType `json:"type"`
	By          string      `json:"by"`
	Descendants int         `json:"descendants"`
	Kids        []int       `json:"kids"`
	Score       int         `json:"score"`
	Time        int64       `json:"time"`
	Title       string      `json:"title"`
	URL         string      `json:"url"`
}

func Request[T any](client *http.Client, method, URL string) (t T, err error) {
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		return t, fmt.Errorf("%w: %v", ErrCreateRequest, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return t, fmt.Errorf("%w: %v", ErrSendRequest, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return t, fmt.Errorf("%w: %d", ErrCode, resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	return
}

type ClientAPI struct {
	endpoint string
	client   *http.Client
}

func NewClientAPI() *ClientAPI {
	return &ClientAPI{
		endpoint: "https://hacker-news.firebaseio.com/v0/",
		client:   &http.Client{},
	}
}

func (c *ClientAPI) GetStory(ID int) (s Story, err error) {
	s, err = Request[Story](
		c.client,
		http.MethodGet,
		fmt.Sprintf("%s/item/%d.json", c.endpoint, ID),
	)
	if err != nil {
		return s, err
	}

	// log.Printf("Get story: %s (%d)", s.Title, s.ID)

	return
}

func (c *ClientAPI) Stories(t StoriesType) (stories []int, err error) {
	stories, err = Request[[]int](
		c.client,
		http.MethodGet,
		fmt.Sprintf("%s/%sstories.json", c.endpoint, t),
	)
	if err != nil {
		return nil, err
	}

	return
}
