package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// getUser return user associated to Github Token
func getUser(token string) (string, error) {

	body, err := httpClient(token, "GET", "https://api.github.com/user")
	if err != nil {
		return "", err
	}

	var response struct {
		Login string `json:"login"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	return response.Login, nil
}

// getAllCommentsFromPR return all comments from PR
func getAllCommentsFromPR(token string, repo string, prID int) (*map[int]*Comment, error) {
	body, err := httpClient(token, "GET", fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/comments", repo, prID))
	if err != nil {
		return nil, err
	}

	comments := make(map[int]*Comment)

	err = json.Unmarshal(body, &comments)
	if err != nil {
		return nil, err
	}

	return &comments, nil

}

// UpdateComment update previous comment on a Github PR
func UpdateComment(nc *Comment, oc *Comment) error {
	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://api.github.com/repos/%s/pulls/comments/%d", nc.Repo, oc.ID), bytes.NewBuffer(nc.Body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(nc.Token, "x-oauth-basic")

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

// CreateComment create comment on a Github PR
func CreateComment(nc *Comment) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/comments", nc.Repo, nc.ThreadID), bytes.NewBuffer(nc.Body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(nc.Token, "x-oauth-basic")

	client := http.Client{}
	_, err = client.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func httpClient(token string, method string, route string) (json.RawMessage, error) {
	req, err := http.NewRequest(method, route, nil)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(token, "x-oauth-basic")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// read body
	var body json.RawMessage
	body, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
