package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

func main() {
	var filterOnlyFounds bool

	var rootCmd = &cobra.Command{
		Use:     "sharlock",
		Short:   "Find usernames across social networks",
		Args:    cobra.MinimumNArgs(1),
		Example: "sherlock mesuutt",
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			checker := newChecker(username, &sites)
			go checker.Check()

			cyan := color.New(color.FgCyan).SprintFunc()
			boldCyan := color.New(color.FgCyan).Add(color.Bold).SprintFunc()
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
					fmt.Printf("[%s] %s: %s (%s)\n", boldCyan("?"), boldCyan(c.site.name), c.ProfileUrl(), cyan("Check failed"))
				} else {
					if c.found {
						fmt.Printf("[%s] %s: %s\n", boldGreen("+"), boldGreen(c.site.name), c.ProfileUrl())
					} else {
						fmt.Printf("[%s] %s: %s\n", boldRed("-"), boldRed(c.site.name), boldYellow("Not Found!"))
					}
				}
			}
		},
	}

	showBanner()
	rootCmd.Flags().BoolVarP(&filterOnlyFounds, "only-found", "i", false, "Prints only found messages. Errors, and invalid username errors will not appear.")
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
