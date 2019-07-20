package main

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"sync"
)

type resultChan chan *Check

type Check struct {
	site     *Site
	username string
	found    bool // Keeps username is found or not on the website
	failed   bool // Keeps check is success or not
	errorMsg string
}

type CheckConfig struct {
	Verbose  bool
	ProxyURL *url.URL
}

type Checker struct {
	username string
	sites    []Site
	results  resultChan
	wg       *sync.WaitGroup
	conf     *CheckConfig
}

// ProfileUrl return profile url of username
func (c *Check) ProfileUrl() string {
	return fmt.Sprintf(c.site.profileUrl, c.username)
}

// ProbeUrl return page which using for check existance of username
func (c *Check) ProbeUrl() string {
	if c.site.probeUrl != "" {
		return fmt.Sprintf(c.site.probeUrl, c.username)
	}
	return fmt.Sprintf(c.site.profileUrl, c.username)

}

func newChecker(username string, sites *[]Site, proxyURL *url.URL, verbose bool) *Checker {
	return &Checker{
		username: username,
		results:  make(resultChan),
		sites:    *sites,
		wg:       &sync.WaitGroup{},
		conf: &CheckConfig{
			Verbose:  verbose,
			ProxyURL: proxyURL,
		},
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
	defer c.wg.Done()

	if check.site.regexCheck != "" {
		match, _ := regexp.MatchString(check.site.regexCheck, c.username)
		if !match {
			check.errorMsg = "Illegal username format!"
			c.results <- check
			return
		}
	}

	check.site.checkerFn(c, check)
}

func (checker *Checker) CreateClient() *http.Client {
	var client = &http.Client{}
	if *checker.conf.ProxyURL != (url.URL{}) {
		fmt.Printf("Proxy kullanilacak")
		//adding the proxy settings to the Transport object
		transport := &http.Transport{
			Proxy: http.ProxyURL(checker.conf.ProxyURL),
		}

		//adding the Transport object to the http Client
		client = &http.Client{
			Transport: transport,
		}
	}

	return client
}
