package main

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func main() {
	var filterOnlyFounds bool
	var proxyStr string
	var verbose bool
	var rootCmd = &cobra.Command{
		Use:     "sherlock USERNAME",
		Short:   "Find usernames across social networks",
		Args:    cobra.MinimumNArgs(1),
		Example: "sherlock user123",
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			proxyURL, err := url.Parse(proxyStr)

			if proxyStr != "" && (err != nil || proxyURL.Scheme == "" || proxyURL.Host == "") {
				log.Fatal("Proxy URL not valid")
			}

			showBanner()
			checker := newChecker(username, &sites, proxyURL, verbose)
			go checker.Check()

			red := color.New(color.FgRed).SprintFunc()
			boldRed := color.New(color.FgRed).Add(color.Bold).SprintFunc()
			boldGreen := color.New(color.FgGreen).Add(color.Bold).SprintFunc()
			boldWhite := color.New(color.FgWhite).Add(color.Bold).SprintFunc()
			boldYellow := color.New(color.FgYellow).Add(color.Bold).SprintFunc()

			fmt.Printf(
				"%s%s%s %s %s %s\n",
				boldGreen("["),
				boldWhite("*"),
				boldGreen("]"),
				boldGreen("Checking username"),
				boldWhite(username),
				boldGreen("on:"),
			)

			for c := range checker.Results() {
				if filterOnlyFounds && (c.failed || !c.found) {
					continue
				}

				if c.failed {
					if checker.conf.Verbose {
						fmt.Printf("[%s] %s: %s (%s\n", boldRed("?"), boldRed(c.site.name), c.ProfileURL(), red(fmt.Sprintf("Error Message: %v", c.errorMsg)))

					} else {
						fmt.Printf("[%s] %s: %s (%s)\n", boldRed("?"), boldRed(c.site.name), c.ProfileURL(), red("Check failed. Use -v flag for verbose mode"))
					}
				} else {
					if c.found {
						fmt.Printf("[%s] %s: %s\n", boldGreen("+"), boldGreen(c.site.name), c.ProfileURL())
					} else {
						if c.errorMsg != "" {
							fmt.Printf("[%s] %s: %s\n", boldRed("-"), boldGreen(c.site.name), boldYellow(c.errorMsg))
						} else {
							fmt.Printf("[%s] %s: %s\n", boldRed("-"), boldGreen(c.site.name), boldYellow("Not Found!"))
						}
					}
				}
			}
		},
	}

	rootCmd.Flags().BoolVarP(&filterOnlyFounds, "only-found", "i", false, "Prints only found messages. Errors, and invalid username errors will not appear.")

	rootCmd.Flags().StringVarP(&proxyStr, "proxy", "p", "", "Make requests over a proxy.")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose mode for detailed error messages.")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func showBanner() {
	banner := `
                                              ."""-.
                                             /      \
 ____  _               _            _        |  _..--'-.
/ ___|| |__   ___ _ __| | ___   ___| |__    >.` + "`" + `__.-""\;"` + "`" + `
\___ \| '_ \ / _ \ '__| |/ _ \ / __| |/ /   / /(     ^\
 ___) | | | |  __/ |  | | (_) | (__|   <    '-` + "`" + `)     =|-.
|____/|_| |_|\___|_|  |_|\___/ \___|_|\_\    /` + "`" + `--.'--'   \ .-.
                                           .'` + "`" + `-._ ` + "`" + `.\    | J /`

	fmt.Printf("%v\n\n", banner)
}
