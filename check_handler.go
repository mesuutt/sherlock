package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type checkerFunc func(*Checker, *Check)

// statusChecker check username by response status code
func statusChecker(checker *Checker, check *Check) {
	res, err := http.Get(check.ProfileUrl())
	if err != nil {
		// Check failed
		check.errorMsg = "Request Error"
		check.failed = true
		checker.results <- check
		return
	}

	defer res.Body.Close()
	check.found = res.StatusCode == 200
	checker.results <- check
}

// bodyChecker check username by searching username in body content
func bodyChecker(searchText string) checkerFunc {

	return func(checker *Checker, check *Check) {
		res, err := http.Get(check.ProbeUrl())
		if err != nil {
			// Check failed
			check.errorMsg = "Request Error"
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
