package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func main() {
	// var wg sync.WaitGroup
	// username := *flag.String("username", "", "One or more usernames to check with social networks.")
	/* onlyFound := *flag.Bool("found", false, "Prints only found messages. Errors, and invalid username errors will not appear.")
	flag.Parse()

	if flag.NArg() == 0 {
		printUsage()
		os.Exit(1)
	}

	username := flag.Args()[0]
	*/

	onlyFound := false
	var rootCmd = &cobra.Command{
		Use:     "sharlock",
		Short:   "Find usernames across social networks",
		Args:    cobra.MinimumNArgs(1),
		Example: "sherlock mesuutt",
		Run: func(cmd *cobra.Command, args []string) {
			checker := newChecker(args[0], &sites)
			go checker.Check()

			cyan := color.New(color.FgCyan).SprintFunc()
			boldCyan := color.New(color.FgCyan).Add(color.Bold).SprintFunc()
			boldRed := color.New(color.FgRed).Add(color.Bold).SprintFunc()
			bolGreen := color.New(color.FgGreen).Add(color.Bold).SprintFunc()

			for c := range checker.Results() {
				if onlyFound && (c.failed || !c.found) {
					continue
				}

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

		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
