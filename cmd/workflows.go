package cmd

import (
	"context"
	"fmt"

	"github.com/ignorant05/Uniflow/internal/config"
	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
	"github.com/ignorant05/Uniflow/platforms"
	"github.com/ignorant05/Uniflow/types"
	"github.com/spf13/cobra"
)

// workflows command flags
var (
	// --with-dispatch (-w)
	// UTILITY: displays only workflows with workflow_dispatch trigger
	wfWithDispatch bool

	// --verbose (-v)
	// UTILITY: verbose output
	workflowsVerbose bool
)

var workflowsCmd = &cobra.Command{
	Use:     "workflows",
	Aliases: []string{"wf"},
	Short:   "List available workflows in the repository",
	Long: `List all GitHub Actions workflows in the configured repository.

This helps you see which workflows are available to trigger.

Examples:
	# Show all workflows for the configured repo
	uniflow workflows

	# Show the workflows related to a specific profile 
	uniflow workflows --profile my-profile (eg. prod)

	# Activate verbose output
	uniflow wf -v`,
	RunE: runWorkflows,
}

func init() {
	workflowsCmd.Flags().BoolVarP(&wfWithDispatch, "with-dispatch", "w", false, "Show only workflows with 'workflow_dispatch' trigger")
	workflowsCmd.Flags().BoolVarP(&workflowsVerbose, "verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(workflowsCmd)
}

// runWorkflows is the main function for status command
func runWorkflows(cmd *cobra.Command, args []string) error {
	// if verbose mode is active
	if workflowsVerbose {
		fmt.Println("<!> Info: Verbose mode enabled")
		fmt.Printf("   Profile: %s\n", profileName)
	}

	fmt.Println("❯ Listing available workflows...")

	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	factory := platforms.NewFactory(cfg)

	// create new client with profileName
	client, err := factory.CreateClientAutoDetectPlatform(ctx, profileName)
	if err != nil {
		errMsg := fmt.Errorf("<?> Error: Field to create client.\n<?> Error: %w.\n", err)
		errorhandling.HandleError(errMsg)
	}

	owner, repo := client.GetRepository(ctx)

	// if verbose mode is active
	if workflowsVerbose {
		fmt.Printf("</> Info: Repository: %s/%s\n", owner, repo)
	}

	listWorkflowsReq := types.ListWorkflowsRequest{
		WithDispatch: wfWithDispatch,
	}

	workflows, err := client.ListWorkflows(ctx, &listWorkflowsReq)
	if err != nil {
		return err
	}

	if len(workflows) == 0 {
		fmt.Println("<?> No workflows found in this repository.")
		fmt.Println("")
		fmt.Println("</> Info: To add a workflow:")
		fmt.Printf("   1. Create .github/workflows/ directory in %s/%s\n", owner, repo)
		fmt.Println("   2. Add a workflow file (e.g., deploy.yml)")
		fmt.Println("   3. Include 'workflow_dispatch:' trigger")
		return nil
	}

	fmt.Printf("\n✓ Found %d workflow(s):\n\n", len(workflows))

	for idx, wf := range workflows {
		fmt.Printf("%d - workflow: %s\n", idx+1, wf.Name)
		fmt.Printf("   - File: %s\n", wf.Path)
		fmt.Printf("   - State: %s\n", wf.State)

		hasDispatch := false

		// if verbose mode is active
		if workflowsVerbose {
			fmt.Printf("   - ID: %d\n", wf.ID)
			fmt.Printf("   - URL: %s\n", wf.URL)
		}

		// if verbose mode is active
		if !hasDispatch && workflowsVerbose {
			fmt.Println("<!> Warning: Check if this workflow has 'workflow_dispatch' trigger")
		}

		fmt.Println()
	}

	fmt.Println("</> Info: Trigger a workflow with:")
	if len(workflows) > 0 {
		firstWorkflow := workflows[0].Path
		fileName := firstWorkflow

		if len(firstWorkflow) > len(".github/workflows/") {
			fileName = firstWorkflow[len(".github/workflows/"):]
		}

		fmt.Printf("   uniflow trigger %s\n", fileName)
	}

	return nil
}
