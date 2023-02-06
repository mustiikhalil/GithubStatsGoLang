package CLI

import (
	"github.com/spf13/cobra"
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "github-worflow",
		Short: "Gets a summary of a Github worflow",
	}
	AppendWorkflowSubCommand(rootCmd)
	AppendActionsSubCommand(rootCmd)
	rootCmd.Execute()
}
