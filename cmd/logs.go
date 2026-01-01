package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/ignorant05/Uniflow/cmd/helpers"
	"github.com/ignorant05/Uniflow/internal/config"
	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	"github.com/ignorant05/Uniflow/platforms"
	ghlogs "github.com/ignorant05/Uniflow/platforms/github/logs"
	"github.com/ignorant05/Uniflow/types"
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

	// --output (-o) flag
	// UTILITY: output file name
	output string

	// --verbose flag
	// UTILITY: logsVerbose output
	logsVerbose bool
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
	logsCmd.Flags().StringVarP(&output, "output", "o", "", "download file name for logs")
	logsCmd.Flags().BoolVarP(&followLogs, "follow", "f", false, "Follow logs (not implemented)")
	logsCmd.Flags().BoolVarP(&logsVerbose, "verbose", "v", false, "verbose output")
	logsCmd.Flags().IntVarP(&tailLines, "tail", "t", 0, "Show last N lines (0 = all)")
	logsCmd.Flags().BoolVarP(&downloadOnly, "download-only", "d", false, "Just show download URL")
	logsCmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	logsCmd.Flags().StringVarP(&platformFlag, "platform", "p", "github", "Platform (github, jenkins, gitlab, circleci). The default is github")

	// root command
	rootCmd.AddCommand(logsCmd)
}

// runLogsCmd
func runLogsCmd(cmd *cobra.Command, args []string) {
	// if Verbose mode is active
	if logsVerbose {
		fmt.Println("<!> Info: logsVerbose mode enabled")
	}

	// only works for github for now
	if platformFlag != "github" {
		fmt.Printf("<?> Error: Platform '%s' not yet supported. Currently only 'github' is available.\n", platformFlag)
		fmt.Println("<!> Warning: Jenkins, GitLab, CircleCI aren't supported yet.\n</> Info: Feel free to contribute if you want them in please.")
		return
	}

	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		errorhandling.HandleError(err)
	}

	factory := platforms.NewFactory(cfg)

	// create new client with profileName
	client, err := factory.CreateClientAutoDetectPlatform(ctx, profileName)
	if err != nil {
		errMsg := fmt.Errorf("<?> Error: Field to create client.\n<?> Error: %w", err)
		errorhandling.HandleError(errMsg)
	}

	owner, repo := client.GetRepository(ctx)

	// if Verbose mode is active
	if logsVerbose {
		fmt.Printf("</> Info: Repository: %s/%s", owner, repo)
		fmt.Printf("</> Info: Platform: %s", platformFlag)
	}

	// getting workflow runID (and name)
	targetRunID, workflowName, err := resolveRunID(ctx, client, owner, repo, args)
	if err != nil {
		errorhandling.HandleError(err)
		return
	}

	// For --follow, verify workflow is still running
	if followLogs {
		statusReq := types.StatusRequest{}

		status, err := client.GetStatus(ctx, &statusReq)
		if err != nil {
			errorhandling.HandleError(fmt.Errorf("failed to get workflow run: %w", err))
			return
		}

		if status.Status == "completed" {
			fmt.Println("<!> Info: Workflow already completed. Use without --follow to view logs.")
			followLogs = false
		} else if status.Status != "in_progress" && status.Status != "queued" {
			fmt.Printf("<!> Info: Workflow status is '%s'. Cannot follow.\n", status.Status)
			followLogs = false
		}
	}

	// Print header
	if workflowName != "" {
		fmt.Printf("❯ %s workflow: %s (Run #%d)\n\n",
			map[bool]string{true: "Streaming", false: "Viewing"}[followLogs],
			workflowName, targetRunID)
	} else {
		fmt.Printf("❯ %s run ID: %d\n\n",
			map[bool]string{true: "Streaming", false: "Viewing"}[followLogs],
			targetRunID)
	}

	// Extract GitHub client for streamer
	githubClient, err := helpers.ExtractGithubClient(client)
	if err != nil {
		errorhandling.HandleError(fmt.Errorf("streaming only supported for GitHub: %w", err))
		return
	}

	// Create streamer
	streamer := ghlogs.NewStreamer(
		githubClient,
		owner,
		repo,
		targetRunID,
		ghlogs.StreamerOptions{
			Follow:    followLogs,
			TailLines: tailLines,
			Colorize:  !noColor,
		})

	// Handle graceful shutdown
	setupGracefulShutdown(streamer)

	// Stream or display logs
	if err := streamer.Stream(); err != nil {
		errorhandling.HandleError(err)
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
func resolveRunID(ctx context.Context, client platforms.PlatformClient, owner, repo string, args []string) (int64, string, error) {
	if runID != 0 {
		return runID, "", nil
	}

	if len(args) > 0 {
		var (
			workflowName string
			workflowID   int64
		)

		workflowFile = args[0]

		workflows, err := client.ListWorkflows(ctx, &types.ListWorkflowsRequest{WithDispatch: wfWithDispatch})
		if err != nil {
			return 0, "", err
		}

		for _, wf := range workflows {
			if strings.HasSuffix(wf.Path, workflowFile) {
				workflowID, workflowName = wf.ID, wf.Name
				break
			}
		}

		if workflowID == 0 {
			return 0, "", fmt.Errorf("<?> Error: Invalid workflow")
		}

		workflowRunLogsReq := types.LogsRequest{
			RunID: workflowID,
			Tail:  tailLines,
		}

		logsURL, err := client.ListWorkflowRunLogs(ctx, &workflowRunLogsReq)
		if err != nil {
			return 0, "", err
		}

		if downloadOnly {
			err := ghlogs.DownloadLogs(logsURL.URL, output)
			if err != nil {
				return 0, "", err
			}
		}

		listWorkflowRunsReq := types.ListWorkflowRunsRequest{
			RunID:        workflowID,
			WorkflowName: workflowFile,
			Branch:       branch,
			Limit:        tailLines,
		}

		runs, err := client.ListWorkflowRuns(ctx, &listWorkflowRunsReq)
		if err != nil || len(runs) <= 0 {
			return 0, "", fmt.Errorf("<?> Error: cannot retrieve workflow runs.\n<?> Workflow: %s", workflowName)
		}

		return runs[0].RunID, workflowName, nil
	}

	return 0, "", fmt.Errorf("<?> Error: Please specify either a workflow name or use --run-id")
}

func setupGracefulShutdown(s *ghlogs.Streamer) {
	// making new channel with buffer
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\n\n<!> Warning: Received interrupt signal...")
		s.Stop()
	}()
}
