package CLI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/spf13/cobra"
)

type ActionInput struct {
	token      string
	repository string
	owner      string
	limit      int
	perPage    int
}

func registerCLIActionsCommands(rootCmd *cobra.Command, params *ActionInput) {
	rootCmd.Flags().StringVarP(&params.token, "token", "t", "", "Github access token")
	rootCmd.Flags().StringVarP(&params.repository, "repository", "r", "", "Repository that contains the workflow")
	rootCmd.Flags().StringVarP(&params.owner, "owner", "o", "", "Owner that contains the workflow")
	rootCmd.Flags().IntVar(&params.limit, "limit", 0, "Owner that contains the workflow")
	rootCmd.Flags().IntVar(&params.perPage, "perPage", 0, "Owner that contains the workflow")

	rootCmd.MarkFlagRequired("token")
	rootCmd.MarkFlagRequired("workflowId")
	rootCmd.MarkFlagRequired("repository")
	rootCmd.MarkFlagRequired("owner")
	rootCmd.MarkFlagRequired("limit")
	rootCmd.MarkFlagRequired("perPage")
}

func AppendActionsSubCommand(command *cobra.Command) {
	var params ActionInput
	subCommand := &cobra.Command{
		Use:   "actions",
		Short: "Gets a summary of a Github worflow",
		Run: func(cmd *cobra.Command, args []string) {
			params.run()
		},
	}
	registerCLIActionsCommands(subCommand, &params)
	command.AddCommand(subCommand)
}

func (params *ActionInput) run() {
	uri := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/workflows", params.owner, params.repository)
	client := http.DefaultClient
	req, reqErr := http.NewRequest("GET", uri, nil)
	if reqErr != nil {
		fmt.Println(reqErr)
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", params.token))
	req.Header.Add("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(uri)
	var response Actions
	json.Unmarshal(body, &response)
	fmt.Println(response)
	for _, r := range response.Workflows {
		input := WorkflowInput{
			Token:      params.token,
			WorkflowId: r.Id,
			Repository: params.repository,
			Owner:      params.owner,
			Limit:      params.limit,
			PerPage:    params.perPage,
		}
		input.Run()
	}
}

type ActionsWorkflow struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Actions struct {
	Workflows []ActionsWorkflow `json:"workflows"`
}
