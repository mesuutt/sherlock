package main

import "fmt"

type Site struct {
	name       string
	errorMsg   string
	checkBy    string
	rank       int
	profileUrl string
	mainUrl    string
}

type Check struct {
	site     *Site
	username string
	found    bool // Keeps username is found or not on the website
	failed   bool // Keeps check is success or not
}

func (c *Check) ProfileUrl() string {
	return fmt.Sprintf(c.site.profileUrl, c.username)
}
