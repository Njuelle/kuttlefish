package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

// Comment is a Github comment on a PR or Issue
type Comment struct {
	Repo string
	User struct {
		Login string `json:"login"`
	} `json:"user"`
	Token      string
	ThreadID   int
	ThreadType int
	Body       json.RawMessage `json:"body"`
	ID         int             `json:"id"`
}

// NewComment is a constructor for Comment structure
func NewComment(repo string, threadID int, threadType int, token string) *Comment {
	return &Comment{
		Repo:       repo,
		ThreadID:   threadID,
		ThreadType: threadType,
		Token:      token,
	}
}

// AddBodyFromFile give a body to a Comment structure from a text file
func (c *Comment) AddBodyFromFile(fn string) (*Comment, error) {
	t, err := template.ParseFiles(fn)
	if err != nil {
		return nil, err
	}

	var tpl bytes.Buffer
	if err = t.Execute(&tpl, c); err != nil {
		return nil, err
	}

	replacer := strings.NewReplacer("\n", "", "\"", "'") // Minify HTML
	c.Body = []byte(fmt.Sprintf(`{"body": "%s"}`, replacer.Replace(tpl.String())))

	return c, nil
}

// SetUser configure user using Github Token
func (c *Comment) SetUser() (*Comment, error) {
	user, err := getUser(c.Token)
	if err != nil {
		return nil, err
	}

	c.User.Login = user

	return c, nil
}

// Comment allow to post or update a comment on a Github thread
func (c *Comment) Comment() error {
	if p, err := c.getPreviousPost(); p != nil && err == nil {
		err = UpdateComment(c, p)
		if err != nil {
			return err
		}

		return nil
	} else {
		err = CreateComment(c)
		if err != nil {
			return err
		}

		return nil
	}
}

// hasAlreadyPost return if given user already post on a PR
func (c *Comment) getPreviousPost() (*Comment, error) {

	comments, err := getAllCommentsFromPR(c.Token, c.Repo, c.ThreadID)
	if err != nil {
		return nil, err
	}

	var previousComment *Comment

	for i := range comments {
		if comments[i].User.Login == c.User.Login {
			previousComment = comments[i]
		}
	}

	return previousComment, nil
}
