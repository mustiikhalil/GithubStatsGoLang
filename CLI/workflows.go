package CLI

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

type WorkflowInput struct {
	Token      string
	WorkflowId int
	Repository string
	Owner      string
	Limit      int
	PerPage    int
}

func registerCLICommands(rootCmd *cobra.Command, params *WorkflowInput) {
	rootCmd.Flags().StringVarP(&params.Token, "token", "t", "", "Github access token")
	rootCmd.Flags().IntVar(&params.WorkflowId, "workflowId", -1, "Repository that contains the workflow")
	rootCmd.Flags().StringVarP(&params.Repository, "repository", "r", "", "Repository that contains the workflow")
	rootCmd.Flags().StringVarP(&params.Owner, "owner", "o", "", "Owner that contains the workflow")
	rootCmd.Flags().IntVar(&params.Limit, "limit", 0, "Owner that contains the workflow")
	rootCmd.Flags().IntVar(&params.PerPage, "perPage", 0, "Owner that contains the workflow")

	rootCmd.MarkFlagRequired("token")
	rootCmd.MarkFlagRequired("workflowId")
	rootCmd.MarkFlagRequired("repository")
	rootCmd.MarkFlagRequired("owner")
	rootCmd.MarkFlagRequired("limit")
	rootCmd.MarkFlagRequired("perPage")
}

func AppendWorkflowSubCommand(command *cobra.Command) {
	var params WorkflowInput
	subCommand := &cobra.Command{
		Use:   "workflow",
		Short: "Gets a summary of a Github worflow",
		Run: func(cmd *cobra.Command, args []string) {
			params.Run()
		},
	}
	registerCLICommands(subCommand, &params)
	command.AddCommand(subCommand)
}

func (params *WorkflowInput) Run() {
	numberOfRequests := params.Limit + 1
	workflows := []WorkFlow{}
	for i := 1; i < numberOfRequests; i++ {
		response, err := params.requestGithubData(i)
		if err != nil {
			fmt.Println(err)
			return
		}
		workflows = append(workflows, response.WorkflowRuns...)
	}
	println(len(workflows))
	totalRuns := 0.0
	statuses := make(map[Status]StatusStats)
	var totalTime time.Duration
	for _, workflow := range workflows {
		totalRuns += 1
		diff := workflow.UpdatedAt.Sub(workflow.RunStartedAt)
		totalTime += diff
		_, prs := statuses[workflow.Conclusion]
		if !prs {
			statuses[workflow.Conclusion] = StatusStats{
				Count: 1,
				Diff:  diff,
			}
		} else {
			c := statuses[workflow.Conclusion]
			c.Count += 1
			c.Diff += diff
			statuses[workflow.Conclusion] = c
		}
	}

	fmt.Println(workflows[0].Path)
	fmt.Println(len(workflows))
	fmt.Println(workflows[0].CreatedAt)
	fmt.Println(workflows[len(workflows)-1].CreatedAt)
	for k, i := range statuses {
		fmt.Printf("Status: %s duration: ", k)
		fmt.Println(time.Duration(i.Diff.Nanoseconds() / i.Count))
	}
	fmt.Print("Duration: ")
	fmt.Println(totalTime)
}

func (params *WorkflowInput) requestGithubData(pageNumber int) (*GithubActionResponse, error) {
	uri := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/actions/workflows/%d/runs?per_page=%d&page=%d",
		params.Owner,
		params.Repository,
		params.WorkflowId,
		params.PerPage,
		pageNumber)
	client := http.DefaultClient
	req, reqErr := http.NewRequest("GET", uri, nil)
	if reqErr != nil {
		fmt.Println(reqErr)
		return nil, reqErr
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", params.Token))
	req.Header.Add("Accept", "application/vnd.github+json")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(uri)
	var response GithubActionResponse
	json.Unmarshal(body, &response)
	return &response, nil
}

type Status string

const (
	Cancelled  Status = "cancelled"
	Success    Status = "success"
	Failure    Status = "failure"
	InProgress Status = "inProgress"
)

type WorkFlow struct {
	Name         string
	Path         string
	Conclusion   Status
	RunStartedAt time.Time `json:"run_started_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type GithubActionResponse struct {
	TotalCount   int        `json:"total_count"`
	WorkflowRuns []WorkFlow `json:"workflow_runs"`
}

type StatusStats struct {
	Count int64
	Diff  time.Duration
}
