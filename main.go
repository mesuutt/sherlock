package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	"github.com/fatih/color"
)

func main() {
	var wg sync.WaitGroup
	resultCh := make(chan Check, len(sites))
	username := "mesuutt"

	for i := 0; len(sites) > i; i++ {
		wg.Add(1)
		check := &Check{
			username: username,
			site:     &sites[i],
		}
		go checkSite(check, &wg, resultCh)
	}
	wg.Wait()
	close(resultCh)

	cyan := color.New(color.FgCyan).SprintFunc()
	boldCyan := color.New(color.FgCyan).Add(color.Bold).SprintFunc()
	boldRed := color.New(color.FgRed).Add(color.Bold).SprintFunc()
	bolGreen := color.New(color.FgGreen).Add(color.Bold).SprintFunc()

	for c := range resultCh {
		if c.failed {
			fmt.Printf("[%s] %s: %s (%s)\n", boldCyan("?"), boldCyan(c.site.name), c.ProfileUrl(), cyan("Check failed"))
		} else {
			if c.found {
				fmt.Printf("[%s] %s: %s\n", bolGreen("+"), bolGreen(c.site.name), c.ProfileUrl())
			} else {
				fmt.Printf("[%s] %s: %s\n", boldRed("-"), boldRed(c.site.name), c.ProfileUrl())
			}
		}
	}
}

func checkSite(check *Check, wg *sync.WaitGroup, resultCh chan Check) {
	resp, err := http.Get(check.ProfileUrl())
	defer wg.Done()

	if err != nil {
		// Check failed
		check.failed = true
		resultCh <- *check
		return
	}

	defer resp.Body.Close()
	if check.site.checkBy == "status_code" {
		check.found = resp.StatusCode == 200
		resultCh <- *check
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		check.failed = true
		resultCh <- *check
		return
	}

	check.found = !strings.Contains(string(body), check.site.errorMsg)
	resultCh <- *check
}
