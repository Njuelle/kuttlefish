package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	// Flags and options definition
	var r = flag.String("r", "org/repo", "Repository where you need to post a comment")
	var f = flag.String("f", "index.html", "The filename where the comment body is stored (prefered HTML)")
	var token = flag.String("token", "12345678", "Github Token")
	var id = flag.Int("id", 0, "PR or Issue number")
	var t = flag.Int("t", 0, "Thread type where to put comment (0 for pull request and 1 for issues)")

	flag.Parse()

	c, err := NewComment(*r, *id, *t, *token).AddBodyFromFile(*f)
	if err != nil {
		fmt.Println("%: File not found or bad permission", *f)
		os.Exit(1)
	}

	c.SetUser()
	err = c.Comment()
	if err != nil {
		fmt.Println(err.Error())
	}
}
