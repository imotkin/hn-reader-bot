package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

const apiURL = "https://hacker-news.firebaseio.com/v0"

type Story struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int64  `json:"time"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

type storiesType string

const (
	New  storiesType = "new"
	Best storiesType = "best"
	Top  storiesType = "top"
)

func SendRequest[T any](method, URL string) (t T, err error) {
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		return t, fmt.Errorf("create HTTP request: %v", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return t, fmt.Errorf("send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return t, fmt.Errorf("invalid HTTP status code: %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&t)
	return
}

func GetStory(id int) (s Story, err error) {
	s, err = SendRequest[Story](http.MethodGet, fmt.Sprintf("%s/item/%d.json", apiURL, id))
	if err != nil {
		return s, fmt.Errorf("send API request: %v", err)
	}

	log.Printf("Get story: %s (%d)", s.Title, s.ID)

	return
}

func Stories(t storiesType) (stories []int, err error) {
	stories, err = SendRequest[[]int](http.MethodGet, fmt.Sprintf("%s/%sstories.json", apiURL, t))
	if err != nil {
		return nil, fmt.Errorf("send API request: %v", err)
	}

	return
}
