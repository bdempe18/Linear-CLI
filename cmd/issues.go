/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"linear/write"
	"strings"

	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
)

var stateFilter string

type IssuesResponse struct {
	Cycles struct {
		Nodes []struct {
			ID     string `json:"id"`
			Issues struct {
				Nodes []struct {
					ID     string `json:"id"`
					Number int    `json:"number"`
					Title  string `json:"title"`
					State  struct {
						Name string `json:"name"`
					} `json:"state"`
				} `json:"nodes"`
			} `json:"issues"`
		} `json:"nodes"`
	} `json:"cycles"`
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
		apiKey := GetApiKey()

		client := graphql.NewClient("https://api.linear.app/graphql")
		request := graphql.NewRequest(`
			query {
				cycles(first: 11) {
					nodes {
						id
						issues {
							nodes {
								id
								number
								title
								state {
									name
								}
							}
						}
					}
				}
			}
		`)

		request.Header.Set("Authorization", apiKey)
		var response IssuesResponse
		if err := client.Run(context.Background(), request, &response); err != nil {
			write.Std.Error("Error making request")
		}

		groupedIssues := make(map[string][]string)

		for _, cycle := range response.Cycles.Nodes {
			for _, issue := range cycle.Issues.Nodes {
				state := issue.State.Name
				fmtName := fmt.Sprintf("(#%d) %s", issue.Number, issue.Title)
				groupedIssues[state] = append(groupedIssues[state], fmtName)

			}
		}

		for state, titles := range groupedIssues {
			if stateFilter != "" && !strings.EqualFold(state, stateFilter) {
				continue
			}

			write.Std.Success(state)
			for _, title := range titles {
				write.Std.Infof(" - %s\n", title)
			}
			write.Std.Info("\n")
		}
		// for _, issue := range response.Cycles.Nodes.Issues.Nodes {
		// write.Std.Info(issue.Title)
		// }
	},
}

func init() {
	issuesCmd.Flags().StringVarP(&stateFilter, "state", "s", "", "Filter issues by state")
	rootCmd.AddCommand(issuesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// issuesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// issuesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
