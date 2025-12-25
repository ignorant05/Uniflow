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

// Status command flags
var (
	// --all (-a) flag
	// UTILITY: show all runs for workflow
	showAllRuns bool

	// --limit (-l) flag
	// UTILITY: limits output (sorted)
	limitRuns int
)

// status command declaration
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

// runStatusCmd is the main status command function
func runStatusCmd(cmd *cobra.Command, args []string) {
	// if verbose mode is active
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

	// if verbose mode is active
	if verbose {
		fmt.Printf("</> Info: Repository: %s/%s\n", owner, repo)
	}

	if len(args) > 0 {
		workflowFile := args[0]
		fmt.Printf("❯ Checking status of workflow: %s\n\n", workflowFile)

		if err := showWorkflowStatus(client, owner, repo, workflowFile); err != nil {
			errorhandling.HandleError(err)
		}
		return
	}

	fmt.Println("</> Info: No workflow provided...")
	fmt.Println("❯ Checking status of all workflows...")

	if err := showAllWorkflowsStatus(client, owner, repo); err != nil {
		errorhandling.HandleError(err)
	}
}

// showWorkflowStatus shows target workflow run status
//
// Parameters:
//   - client: github client
//   - owner: owner name
//   - repo: repository name
//   - workflowFile: target workflow file
//
// Errors possible causes:
//   - invalid workflow name
//   - rate limit exceeded
//   - no runs for this workflow
func showWorkflowStatus(client *github.Client, owner, repo, workflowFile string) error {
	found := false

	var (
		workflowID   int64
		workflowName string
	)

	// gettiing all workflows for owner/repo
	workflows, err := client.ListWorkflows(owner, repo)
	if err != nil {
		return err
	}

	// looking for a specific workflow with name: workflowFile
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

	// If it exists, getting workflowFile's runs
	runs, err := client.GetWorkflowRuns(owner, repo, workflowID)
	if err != nil {
		if strings.Contains(err.Error(), "rate limit") {
			fmt.Println("<?> Error: GitHub API rate limit exceeded")
			fmt.Println("</> Info:  Rate limits reset every hour")
			fmt.Println("</> Info:	Check remaining quota: https://api.github.com/rate_limit")

			fmt.Println("</> Info: Waiting 60 seconds before retry...")
			time.Sleep(60 * time.Second)

			runs, err = client.GetWorkflowRuns(owner, repo, workflowID)
		}
		return err
	}

	if len(runs) == 0 {
		fmt.Println("</> Info: No runs found for this workflow")
		return nil
	}

	fmt.Printf("> Workflow: %s\n", workflowName)
	fmt.Printf("> File: %s\n", workflowFile)
	fmt.Println(strings.Repeat("─", 80))

	// limiting (if present)
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

// DisplayRun run information
//
// Parameters:
//   - run: workflow run target
func DisplayRun(run *gh.WorkflowRun) {
	fmt.Printf("  Run #%d\n", run.GetRunNumber())
	fmt.Printf("    Status:     %s\n", helpers.FormatStatus(run.GetStatus()))
	fmt.Printf("    Conclusion: %s\n", helpers.FormatConclusion(run.GetConclusion()))
	fmt.Printf("    Branch:     %s\n", run.GetHeadBranch())
	fmt.Printf("    Triggered:  %s\n", helpers.FormatTime(run.GetCreatedAt().Time))

	// if verbose mode is active
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

// showAllWorkflowsStatus shows all workflow status
//
// Parameters:
//   - client: github client
//   - owner: owner name
//   - repo: repository name
//
// Errors Possible causes:
//   - no workflows
//   - rate limit exceeded
func showAllWorkflowsStatus(client *github.Client, owner, repo string) error {
	totalRuns := 0

	// retrieving all workflows for owner/repo
	workflows, err := client.ListWorkflows(owner, repo)
	if err != nil {
		return err
	}

	// if no workflows then it quits
	if len(workflows) == 0 {
		fmt.Println("</> Info: No workflows found in this repository")
		return nil
	}

	// getting all workflow runs
	for _, wf := range workflows {
		runs, err := client.GetWorkflowRuns(owner, repo, wf.GetID())
		if err != nil {
			if strings.Contains(err.Error(), "rate limit") {
				fmt.Println("<?> Error: GitHub API rate limit exceeded")
				fmt.Println("</> Info:  Rate limits reset every hour")
				fmt.Println("</> Info:	Check remaining quota: https://api.github.com/rate_limit")

				fmt.Println("</> Info: Waiting 60 seconds before retry...")
				time.Sleep(60 * time.Second)

				runs, err = client.GetWorkflowRuns(owner, repo, wf.GetID())
			}

			// if verbose mode is active
			if verbose {
				fmt.Printf("<?> Warning: Failed to get runs for %s: %v\n", wf.GetName(), err)
			}

			continue
		}

		// if it has no runs, then print nothing and continue
		// no need to print anything for this workflow
		if len(runs) == 0 {
			continue
		}

		fmt.Printf("   Workflow: %s\n", wf.GetName())
		fmt.Printf("   File: %s\n", strings.TrimPrefix(wf.GetPath(), ".github/workflows/"))
		fmt.Println(strings.Repeat("─", 80))

		// limit = 1 just in case
		limit := 1
		if showAllRuns {
			limit = helpers.Min(limitRuns, len(runs))
		}

		fmt.Printf("\n	- Recent Runs (showing %d of %d):\n\n", limit, len(runs))

		// limiting results
		for i := 0; i < limit; i++ {
			DisplayRun(runs[i])
		}

		fmt.Println()
		totalRuns++
	}

	// Print result conveniently
	if totalRuns == 0 {
		fmt.Println("</> Info: No workflow runs found")
	} else {
		fmt.Printf("✓  Displayed status for %d workflow(s)\n", totalRuns)
		fmt.Println("   To see more runs for a specific workflow:")
		fmt.Println("   uniflow status <workflow-file> --limit 10")
	}

	return nil
}
