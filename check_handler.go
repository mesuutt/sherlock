package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type checkerFunc func(*Checker, *Check)

var userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/44.0.2403.157 Safari/537.36"

// statusChecker check username by response status code
func statusChecker(checker *Checker, check *Check) {
	// res, err := http.Get(check.ProbeUrl())
	var client = checker.CreateClient()

	req, err := http.NewRequest("GET", check.ProbeURL(), nil)
	req.Header.Set("User-Agent", userAgent)
	res, err := client.Do(req)

	if err != nil {
		// Check failed
		if checker.conf.Verbose {
			check.errorMsg = err.Error()
		} else {
			check.errorMsg = "Request Error"
		}
		check.failed = true
		checker.results <- check
		return
	}

	defer res.Body.Close()
	check.found = res.StatusCode == 200
	checker.results <- check
}

// bodyChecker unsure username by searching given text in page content
func bodyChecker(searchText string) checkerFunc {

	return func(checker *Checker, check *Check) {
		// res, err := http.Get(check.ProbeUrl())
		var client = checker.CreateClient()
		req, err := http.NewRequest("GET", check.ProbeURL(), nil)
		req.Header.Set("User-Agent", userAgent)
		res, err := client.Do(req)

		if err != nil {
			// Check failed
			if checker.conf.Verbose {
				check.errorMsg = err.Error()
			} else {
				check.errorMsg = "Request Error"
			}
			check.failed = true
			checker.results <- check
			return
		}

		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)

		if err != nil {
			check.failed = true
			checker.results <- check
			return
		}

		check.found = !strings.Contains(string(body), searchText)
		checker.results <- check
	}
}

// catchByRedirectUrl check username by redirected page url
func redirectChecker(redirectUrlTemplate string) checkerFunc {
	return func(checker *Checker, check *Check) {
		var client = checker.CreateClient()
		client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}

		req, err := http.NewRequest("GET", check.ProbeURL(), nil)
		req.Header.Set("User-Agent", userAgent)
		res, err := client.Do(req)

		if err != nil {
			// Check failed
			if checker.conf.Verbose {
				check.errorMsg = err.Error()
			} else {
				check.errorMsg = "Request Error"
			}
			// check.errorMsg = "Request Error"
			check.failed = true
			checker.results <- check
			return
		}
		var errorURL string
		if strings.Contains(redirectUrlTemplate, "%s") {
			errorURL = fmt.Sprintf(redirectUrlTemplate, check.username)
		} else {
			errorURL = redirectUrlTemplate
		}

		check.found = !strings.Contains(res.Header.Get("Location"), errorURL)

		checker.results <- check
	}
}
