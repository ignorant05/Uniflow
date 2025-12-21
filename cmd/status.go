package cmd

import (
	"fmt"
	"strings"
	"time"

	gh "github.com/google/go-github/v57/github"
	"github.com/ignorant05/Uniflow/cmd/helpers"
	"github.com/ignorant05/Uniflow/configs/github"
	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	"github.com/spf13/cobra"
)

var (
	showAllRuns bool
	limitRuns   int
)

var statusCmd = &cobra.Command{
	Use:     "status [workflow]",
	Aliases: []string{"s"},
	Short:   "Shows workflows's status",
	Long: `Status displays the current state of running or completed workflows.
If no workflow name is provided, it shows status for all workflows.


Example:
	
	# Check the status of all workflows
	uniflow status 

	# Check the status of a specific workflow
	uniflow status my-workflow (eg. deploy.yaml) 

	# Check the status of a specific workflow with all it's runs 
	uniflow status my-workflow --all 

	# Show only a limited number of runs provided by you
	uniflow status my-workflow --limit number-of-runs-desired (default: 5 most recent) 

	# Activate verbose output
	uniflow s --verbose`,
	Args: cobra.MaximumNArgs(1),
	Run:  runStatusCmd,
}

func init() {
	statusCmd.Flags().BoolVarP(&showAllRuns, "all", "a", false, "Show all workflow runs (default: 5 most recent)")
	statusCmd.Flags().IntVarP(&limitRuns, "limit", "l", 5, "Number of runs to show")

	rootCmd.AddCommand(statusCmd)
}

func runStatusCmd(cmd *cobra.Command, args []string) {
	if verbose {
		fmt.Println("Running in verbose mode...")
	}

	client, err := github.NewClientFromConfig(profileName)
	if err != nil {
		errorhandling.HandleError(err)
	}

	owner, repo, err := client.GetDefaultRepository()
	if err != nil {
		errorhandling.HandleError(err)
	}

	if verbose {
		fmt.Printf("<.> Info: Repository: %s/%s\n", owner, repo)
	}

	if len(args) > 0 {
		workflowFile := args[0]
		fmt.Printf("❯❯❯ Checking status of workflow: %s\n\n", workflowFile)

		if err := showWorkflowStatus(client, owner, repo, workflowFile); err != nil {
			errorhandling.HandleError(err)
		}
		return
	}

	fmt.Println("<.> Info: No workflow provided...")
	fmt.Println("❯❯❯ Checking status of all workflows...")

	if err := showAllWorkflowsStatus(client, owner, repo); err != nil {
		errorhandling.HandleError(err)
	}
}

func showWorkflowStatus(client *github.Client, owner, repo, workflowFile string) error {
	workflows, err := client.ListWorkflows(owner, repo)
	if err != nil {
		return err
	}
	var (
		workflowID   int64
		workflowName string
	)

	found := false

	for _, wf := range workflows {
		if strings.HasSuffix(wf.GetPath(), workflowFile) {
			workflowID = wf.GetID()
			workflowName = wf.GetName()
			found = !found
			break
		}
	}

	if !found {
		fmt.Printf("<?> Error: No workflow found with filename: %s\n\n", workflowFile)
		fmt.Println("</> Info: Available workflows:")
		for _, wf := range workflows {
			path := wf.GetPath()
			filename := strings.TrimPrefix(path, ".github/workflows/")
			fmt.Printf("   - %s\n", filename)
		}
		return nil
	}

	runs, err := client.GetWorkflowRuns(owner, repo, workflowID)
	if err != nil {
		if strings.Contains(err.Error(), "rate limit") {
			fmt.Println("<?> Error: GitHub API rate limit exceeded")
			fmt.Println("<.> Info:  Rate limits reset every hour")
			fmt.Println("<.> Info:	Check remaining quota: https://api.github.com/rate_limit")

			fmt.Println("<.> Info: Waiting 60 seconds before retry...")
			time.Sleep(60 * time.Second)

			runs, err = client.GetWorkflowRuns(owner, repo, workflowID)
		}
		return err
	}

	if len(runs) == 0 {
		fmt.Println("<.> Info: No runs found for this workflow")
		return nil
	}

	fmt.Printf("> Workflow: %s\n", workflowName)
	fmt.Printf("> File: %s\n", workflowFile)
	fmt.Println(strings.Repeat("─", 80))

	limit := limitRuns
	if showAllRuns {
		limit = len(runs)
	}
	if limit > len(runs) {
		limit = len(runs)
	}

	fmt.Printf("\n	- Recent Runs (showing %d of %d):\n\n", limit, len(runs))

	for i := 0; i < limit; i++ {
		DisplayRun(runs[i])
	}

	return nil
}

func DisplayRun(run *gh.WorkflowRun) {
	fmt.Printf("  Run #%d\n", run.GetRunNumber())
	fmt.Printf("    Status:     %s\n", helpers.FormatStatus(run.GetStatus()))
	fmt.Printf("    Conclusion: %s\n", helpers.FormatConclusion(run.GetConclusion()))
	fmt.Printf("    Branch:     %s\n", run.GetHeadBranch())
	fmt.Printf("    Triggered:  %s\n", helpers.FormatTime(run.GetCreatedAt().Time))

	if verbose {
		fmt.Printf("    Run ID:     %d\n", run.GetID())
		fmt.Printf("    Commit:     %.7s\n", run.GetHeadSHA())
		fmt.Printf("    Actor:      %s\n", run.GetActor().GetLogin())
		fmt.Printf("    Event:      %s\n", run.GetEvent())
		fmt.Printf("    Updated:    %s\n", helpers.FormatTime(run.GetUpdatedAt().Time))
		fmt.Printf("    URL:        %s\n", run.GetHTMLURL())
	}

	fmt.Println()
}

func showAllWorkflowsStatus(client *github.Client, owner, repo string) error {
	totalRuns := 0

	workflows, err := client.ListWorkflows(owner, repo)
	if err != nil {
		return err
	}

	if len(workflows) == 0 {
		fmt.Println("<.> Info: No workflows found in this repository")
		return nil
	}

	for _, wf := range workflows {
		runs, err := client.GetWorkflowRuns(owner, repo, wf.GetID())
		if err != nil {
			if strings.Contains(err.Error(), "rate limit") {
				fmt.Println("<?> Error: GitHub API rate limit exceeded")
				fmt.Println("<.> Info:  Rate limits reset every hour")
				fmt.Println("<.> Info:	Check remaining quota: https://api.github.com/rate_limit")

				fmt.Println("<.> Info: Waiting 60 seconds before retry...")
				time.Sleep(60 * time.Second)

				runs, err = client.GetWorkflowRuns(owner, repo, wf.GetID())
			}

			if verbose {
				fmt.Printf("<?> Warning: Failed to get runs for %s: %v\n", wf.GetName(), err)
			}

			continue
		}

		if len(runs) == 0 {
			continue
		}

		fmt.Printf("   Workflow: %s\n", wf.GetName())
		fmt.Printf("   File: %s\n", strings.TrimPrefix(wf.GetPath(), ".github/workflows/"))
		fmt.Println(strings.Repeat("─", 80))

		limit := 1
		if showAllRuns {
			limit = helpers.Min(limitRuns, len(runs))
		}

		fmt.Printf("\n	- Recent Runs (showing %d of %d):\n\n", limit, len(runs))

		for i := 0; i < limit; i++ {
			DisplayRun(runs[i])
		}

		fmt.Println()
		totalRuns++
	}

	if totalRuns == 0 {
		fmt.Println("<.> Info: No workflow runs found")
	} else {
		fmt.Printf("<✓> Displayed status for %d workflow(s)\n", totalRuns)
		fmt.Println("\n<.> Info: To see more runs for a specific workflow:")
		fmt.Println("   uniflow status <workflow-file> --limit 10")
	}

	return nil
}
