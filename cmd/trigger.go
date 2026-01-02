package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/ignorant05/Uniflow/cmd/helpers"
	"github.com/ignorant05/Uniflow/internal/config"
	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	"github.com/ignorant05/Uniflow/platforms"

	ghlogs "github.com/ignorant05/Uniflow/platforms/configurations/github/logs"
	"github.com/ignorant05/Uniflow/types"
	"github.com/spf13/cobra"
)

// trigger command flags
var (
	// --branch (-b)
	// UTILITY: specify branch
	branch string

	// --workflow (-w)
	// UTILITY: specify workflow
	workflowFile string

	// --inputs (-i)
	// UTILITY: with inputs
	inputs map[string]string

	// --profile (-p)
	// UTILITY: specify profile
	profileName string

	// --stream
	// UTILITY: stream logs in real time
	streamLogs bool

	// --verbose (-v)
	// UTILITY: verbose output
	triggerVerbose bool
)

var triggerCmd = &cobra.Command{
	Use:     "trigger [workflow]",
	Aliases: []string{"t"},
	Short:   "Trigger a workflow execution",
	Long: `Trigger starts the execution of a specified workflow.
You can pass the workflow name as an argument.

Example:
	uniflow trigger deploy.yml

	# Trigger on a specific branch
	uniflow trigger deploy.yml --branch develop

	# Trigger with inputs
	uniflow trigger deploy.yml --input environment=prod --input version=v1.0

	# Use a specific profile
	uniflow trigger deploy.yml --profile prod`,
	Args: cobra.MaximumNArgs(1),
	Run:  runTriggerCmd,
}

func init() {
	triggerCmd.Flags().StringVarP(&branch, "branch", "b", "main", "Branch to trigger the workflow on")
	triggerCmd.Flags().StringVarP(&workflowFile, "workflow", "w", "", "Workflow file name (if different from arg)")
	triggerCmd.Flags().StringToStringVarP(&inputs, "input", "i", nil, "Workflow inputs (key=value)")
	triggerCmd.Flags().StringVarP(&profileName, "profile", "p", "default", "Config profile to use")
	triggerCmd.Flags().BoolVarP(&streamLogs, "stream", "s", false, "Stream workflow logs in real time")
	triggerCmd.Flags().BoolVarP(&triggerVerbose, "verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(triggerCmd)
}

// trigger command main function
func runTriggerCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		errMsg := fmt.Errorf("<?> Error: Not enough arguments")
		errorhandling.HandleError(errMsg)
	}

	workflow := args[0]

	// if verbose mode active
	if triggerVerbose {
		fmt.Printf("<!> Info: Verbose mode enabled\n")
		fmt.Printf("   Workflow: %s\n", workflow)
		fmt.Printf("   Branch: %s\n", branch)
		fmt.Printf("   Profile: %s\n", profileName)
		if len(inputs) > 0 {
			fmt.Printf("   Inputs: %v\n", inputs)
		}
	}

	fmt.Printf("❯ Triggering workflow: %s\n", workflow)

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

	// if verbose mode active
	if triggerVerbose {
		fmt.Printf("</> Info: %s/%s\n", owner, repo)
	}

	// parsing workflow inputs
	workflowInputs := make(map[string]interface{})
	for key, val := range inputs {
		workflowInputs[key] = val
	}

	triggerReqBody := types.TriggerRequest{
		WorkflowName: workflowFile,
		Branch:       branch,
		Inputs:       workflowInputs,
	}
	// trigger workflow
	_, err = client.TriggerWorkflow(ctx, &triggerReqBody)
	if err != nil {
		fmt.Printf("<?> Error: Failed to trigger workflow.\n")
		fmt.Printf("<?> Error: %v\n\n", err)

		fmt.Println("<!> Common issues:")
		fmt.Println("   1. Workflow file must exist in .github/workflows/")
		fmt.Printf("      Expected location: .github/workflows/%s\n", workflow)
		fmt.Println("")
		fmt.Println("   2. Workflow must have 'workflow_dispatch' trigger:")
		fmt.Println("      on:")
		fmt.Println("        workflow_dispatch:")
		fmt.Println("")
		fmt.Println("   3. Check available workflows with:")
		fmt.Println("      orchestrator workflows")

		errorhandling.HandleError(err)
	}

	fmt.Println("✓ Workflow triggered successfully!")
	fmt.Printf("   Repository: %s/%s\n", owner, repo)
	fmt.Printf("   Workflow: %s\n", workflow)
	fmt.Printf("   Branch: %s\n", branch)

	// retrieves repository information
	repoInfo, err := client.GetRepositoryInfo(ctx)
	if err != nil {
		errMsg := fmt.Errorf("<?> Error: Failed to retrieve repository %s/%s info.\n<?> Error: %w", owner, repo, err)
		errorhandling.HandleError(errMsg)
	} else {
		fmt.Printf("   View at: %s/actions\n", repoInfo.HTMLURL)
	}
	if len(workflowInputs) > 0 {
		fmt.Println("\n❯ Inputs:")
		for key, val := range workflowInputs {
			fmt.Printf("❯   %s: %v\n", key, val)
		}
	}

	if streamLogs {
		fmt.Println("❯ Waiting for workflow to start...")

		var runStatus string
		// wait for 30 secs
		for range 30 {
			ListWorkflowRuns := types.ListWorkflowRunsRequest{
				RunID:        runID,
				WorkflowName: workflow,
				Branch:       branch,
				Limit:        1,
			}
			runs, err := client.ListWorkflowRuns(ctx, &ListWorkflowRuns)
			if err == nil && runs[0].Status != "queued" {
				runStatus = runs[0].Status
				break
			}

			time.Sleep(1 * time.Second)
		}

		// Extract GitHub client for streamer
		githubClient, err := helpers.ExtractGithubClient(client)
		if err != nil {
			errorhandling.HandleError(fmt.Errorf("streaming only supported for GitHub: %w", err))
			return
		}

		if runStatus == "in_progress" || runStatus == "queued" {
			// Stream logs
			streamer := ghlogs.NewStreamer(
				githubClient,
				owner,
				repo,
				runID,
				ghlogs.StreamerOptions{
					Follow:    true,
					TailLines: 0,
					Colorize:  true,
				})

			setupGracefulShutdown(streamer)
			err := streamer.Stream()
			if err != nil {
				errorhandling.HandleError(err)
			}

		} else {
			fmt.Println("<!> Warn:  Workflow didn't start within expected time.")
			fmt.Printf("   View logs later with: uniflow logs --run-id %d\n", runID)
		}
	} else {
		fmt.Printf("   View logs with: uniflow logs --run-id %d\n", runID)
		fmt.Printf("   Or stream with: uniflow logs --run-id %d --follow\n", runID)
	}
}
