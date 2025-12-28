package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	"github.com/ignorant05/Uniflow/internal/logs"
	"github.com/ignorant05/Uniflow/platforms/github"
	"github.com/spf13/cobra"
)

// Config command flags representatives
var (
	// --run-id flag
	// UTILITY: specify workflow run id
	runID int64

	// --job flag
	// UTILITY: real-time log streaming
	jobName string

	// --follow flag
	// UTILITY: how many lines to show
	followLogs bool

	// --tail flag
	// UTILITY: how many lines to show (recent)
	tailLines int

	// --download-only flag
	// UTILITY: download only logs option
	downloadOnly bool

	// --no-color flag
	// UTILITY: no colored output
	noColor bool

	// --platform flag
	// UTILITY: specify platform
	platformFlag string
)

// Command: logs (or l)
//
// Example usage:
//   - uniflow logs deploy.yaml
var logsCmd = &cobra.Command{
	Use:     "logs [workflow]",
	Aliases: []string{"l"},
	Short:   "View workflow execution logs",
	Long: `Logs displays the execution logs for a specified workflow.
You can use this to debug issues or monitor workflow progress.


Features:
	• Real-time log streaming with --follow
	• Colored output for different log levels
	• Timestamps for each log line
	• Tail support to limit output
	• Graceful handling of Ctrl+C
	• Auto-detection of run completion

Example:
	# Latest run logs
	uniflow logs deploy.yml

	# Specific run with streaming
	uniflow logs deploy.yml --run-id 123456 --follow

	# Last 100 lines only
	uniflow logs deploy.yml --tail 100

	# Without colors
	uniflow logs deploy.yml --no-color

	# Specific job
	uniflow logs deploy.yml --job build`,
	Args: cobra.MaximumNArgs(1),
	Run:  runLogsCmd,
}

// Commands and subcommnds declaration
func init() {
	// Flags declaration
	logsCmd.Flags().Int64Var(&runID, "run-id", 0, "Specific run ID")
	logsCmd.Flags().StringVarP(&jobName, "job", "j", "", "Specific job name")
	logsCmd.Flags().BoolVarP(&followLogs, "follow", "f", false, "Follow logs (not implemented)")
	logsCmd.Flags().IntVarP(&tailLines, "tail", "t", 0, "Show last N lines (0 = all)")
	logsCmd.Flags().BoolVarP(&downloadOnly, "download-only", "d", false, "Just show download URL")
	logsCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	logsCmd.Flags().StringVarP(&platformFlag, "platform", "p", "github", "Platform (github, jenkins, gitlab, circleci). The default is github")

	// root command
	rootCmd.AddCommand(logsCmd)
}

// runLogsCmd
func runLogsCmd(cmd *cobra.Command, args []string) {
	// if verbose mode is active
	if verbose {
		fmt.Println("<!> Info: Verbose mode enabled")
	}

	// only works for github for now
	if platformFlag != "github" {
		fmt.Printf("<?> Error: Platform '%s' not yet supported. Currently only 'github' is available.\n", platformFlag)
		fmt.Println("<!> Warning: Jenkins, GitLab, CircleCI aren't supported yet.\n</> Info: Feel free to contribute if you want them in please.")
		return
	}

	// creating new client using profile name
	client, err := github.NewClientFromConfig(profileName)
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	// retrieving default repository field
	owner, repo, err := client.GetDefaultRepository()
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	// if verbose mode is active
	if verbose {
		fmt.Printf("</> Info: Repository: %s/%s\n", owner, repo)
		fmt.Printf("</> Info: Platform: %s\n", platformFlag)
	}

	// getting workflow runID (and name)
	targetRunID, workflowName, err := resolveRunID(client, owner, repo, args)
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	// Verify workflowName variable value and print output as needed
	if workflowName != "" {
		fmt.Printf("❯ Fetching logs for workflow: %s (Run #%d)\n\n", workflowName, targetRunID)
	} else {
		fmt.Printf("❯ Fetching logs for run ID: %d\n\n", targetRunID)
	}

	// creating new streamer
	streamer := logs.NewStreamer(
		client,
		owner,
		repo,
		targetRunID,
		logs.StreamerOptions{
			Follow:    followLogs,
			TailLines: tailLines,
			Colorize:  !noColor,
		})

	// Graceful shutdown
	setupGracefulShutdown(streamer)

	if err := streamer.Stream(); err != nil {
		errorhandling.HandleError(err)
		return
	}

}

// resolveRunID retrieves workflow runID (and name)
//
// Parameters:
//   - client: github client
//   - owner: owner name
//   - repo: repository name
//   - args: arguments from command
//
// Errors possible causes:
//   - invalid runID and workflow name
//   - cannot retrieve workflow (either deosn't exist or internal problem)
func resolveRunID(client *github.Client, owner, repo string, args []string) (int64, string, error) {
	if runID != 0 {
		return runID, "", nil
	}

	if len(args) > 0 {
		var (
			workflowName string
			workflowID   int64
		)

		workflowFile = args[0]

		workflows, err := client.ListWorkflows(owner, repo)
		if err != nil {
			return 0, "", err
		}

		for _, wf := range workflows {
			if strings.HasSuffix(wf.GetPath(), workflowFile) {
				workflowID, workflowName = wf.GetID(), wf.GetName()
				break
			}
		}

		if workflowID == 0 {
			return 0, "", fmt.Errorf("<?> Error: Invalid workflow.\n")
		}

		logsURL, err := client.GetWorkflowRunLogs(owner, repo, workflowID)
		if err != nil {
			return 0, "", err
		}

		if downloadOnly {
			logs.DownloadLogs(logsURL)
		}

		runs, err := client.GetWorkflowRuns(owner, repo, workflowID)
		if err != nil || len(runs) <= 0 {
			return 0, "", fmt.Errorf("<?> Error: cannot retrieve workflow runs.\n<?> Workflow: %s\n", workflowName)
		}

		return runs[0].GetID(), workflowName, nil
	}

	return 0, "", fmt.Errorf("<?> Error: Please specify either a workflow name or use --run-id")
}

func setupGracefulShutdown(s *logs.Streamer) {
	// making new channel with buffer
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\n<!> Warning: Received interrupt signal...")
		s.Stop()
	}()
}
