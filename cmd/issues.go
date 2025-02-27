/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"linear/internal"
	"strings"

	"github.com/spf13/cobra"
)

var stateFilter string

type Status struct {
	Name string `json:"name"`
}

type Issue struct {
	ID     string `json:"id"`
	Number int    `json:"number"`
	Title  string `json:"title"`
	State  Status `json:"state"`
}

type Cycle struct {
	ID     string `json:"id"`
	Issues struct {
		Nodes []Issue `json:"nodes"`
	} `json:"issues"`
}

type IssuesResponse struct {
	Cycles struct {
		Nodes []Cycle `json:"nodes"`
	} `json:"cycles"`
}

func (r IssuesResponse) Show() {
	groupedIssues := make(map[string][]string)

	for _, cycle := range r.Cycles.Nodes {
		for _, issue := range cycle.Issues.Nodes {
			state := issue.State.Name
			fmtName := fmt.Sprintf("(#%d) %s", issue.Number, issue.Title)
			groupedIssues[state] = append(groupedIssues[state], fmtName)
		}
	}

	// TODO the filters should be a struct or slice argument
	for state, titles := range groupedIssues {
		if stateFilter != "" && !strings.EqualFold(state, stateFilter) {
			continue
		}

		internal.Std.Success(state)
		for _, title := range titles {
			internal.Std.Infof(" - %s\n", title)
		}
		internal.Std.Info("\n")
	}
}

// issuesCmd represents the issues command
var issuesCmd = &cobra.Command{
	Use:   "issues",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		// cycleString := getCycleQLFilter()

		queryStructure := internal.GenNode("cycles",
			internal.GenField("id"),
			internal.GenNode("issues",
				internal.GenField("id"),
				internal.GenField("number"),
				internal.GenField("title"),
				internal.GenField("state",
					internal.GenField("name"),
				),
			),
		)

		apiKey := GetApiKey()
		var response IssuesResponse
		err := internal.Request(&response, queryStructure, apiKey)

		if err != nil {
			internal.Std.Errorf("Request could not be completed: %v", err)
			return
		}

		response.Show()
	},
}

func init() {
	issuesCmd.Flags().StringVarP(&stateFilter, "state", "s", "", "Filter issues by state")
	// issuesCmd.Flags().StringVarP(&cycle, "cycle", "c", "", "Cycle number")
	rootCmd.AddCommand(issuesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issuesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// issuesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
