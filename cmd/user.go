/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"linear/internal"

	"github.com/machinebox/graphql"
	"github.com/spf13/cobra"
)

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		printUserInfo()
	},
}

type UserResponse struct {
	Viewer struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Workspace struct {
			ID     string `json:"id"`
			Name   string `json:"name"`
			URLKey string `json:"urlKey"`
			Teams  struct {
				Nodes []struct {
					ID   string `json:"id"`
					Name string `json:"name"`
					Key  string `json:"key"`
				} `json:"nodes"`
			} `json:"teams"`
		} `json:"organization"`
	} `json:"viewer"`
}

func printUserInfo() {
	apiKey := GetApiKey()

	client := graphql.NewClient("https://api.linear.app/graphql")
	request := graphql.NewRequest(`
		query {
			viewer {
				id
				name
				email
				organization {
					id
					name
					urlKey
					teams {
						nodes {
							id
							name
							key
						}
					}
				}
			}
		}
	`)

	request.Header.Set("Authorization", apiKey)
	var response UserResponse
	if err := client.Run(context.Background(), request, &response); err != nil {
		// TODO need Successf
		internal.Std.Error("Error making request")
	}

	viewer := response.Viewer
	// Print user information
	fmt.Printf("User Information:\n")
	fmt.Printf("  Name: %s\n", viewer.Name)
	fmt.Printf("  Email: %s\n", viewer.Email)

	// Print workspace information
	fmt.Printf("\nWorkspace Information:\n")
	fmt.Printf("  Name: %s\n", viewer.Workspace.Name)
	fmt.Printf("  URL Key: %s\n", viewer.Workspace.URLKey)

	// Print teams information
	fmt.Printf("\nTeams:\n")
	for _, team := range viewer.Workspace.Teams.Nodes {
		fmt.Printf("  - %s (Key: %s)\n", team.Name, team.Key)
	}
}

// Here you will define your flags and configuration settings.

// Cobra supports Persistent Flags which will work for this command
// and all subcommands, e.g.:
// userCmd.PersistentFlags().String("foo", "", "A help for foo")

// Cobra supports local flags which will only run when this command
// is called directly, e.g.:
// userCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

func init() {
	rootCmd.AddCommand(userCmd)
}
