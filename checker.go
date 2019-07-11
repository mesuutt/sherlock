package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

type resultChan chan *Check

type Check struct {
	site     *Site
	username string
	found    bool // Keeps username is found or not on the website
	failed   bool // Keeps check is success or not
}

func (c *Check) ProfileUrl() string {
	return fmt.Sprintf(c.site.profileUrl, c.username)
}

type Checker struct {
	username string
	sites    []Site
	results  resultChan
	wg       *sync.WaitGroup
}

func newChecker(username string, sites *[]Site) *Checker {
	return &Checker{
		username: username,
		results:  make(resultChan),
		sites:    *sites,
		wg:       &sync.WaitGroup{},
	}
}

func (c *Checker) Check() {
	for i := 0; len(c.sites) > i; i++ {
		c.wg.Add(1)
		check := &Check{
			username: c.username,
			site:     &sites[i],
		}

		go c.checkSite(check)
	}
	c.wg.Wait()
	close(c.results)
}

func (c *Checker) Results() resultChan {
	return c.results
}

func (c *Checker) checkSite(check *Check) {
	resp, err := http.Get(check.ProfileUrl())
	defer c.wg.Done()
	if err != nil {
		// Check failed
		check.failed = true
		c.results <- check
		return
	}

	defer resp.Body.Close()
	if check.site.checkBy == "status_code" {
		check.found = resp.StatusCode == 200
		c.results <- check
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		check.failed = true
		c.results <- check
		return
	}

	check.found = !strings.Contains(string(body), check.site.errorMsg)
	c.results <- check
}
