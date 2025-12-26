package cmd

import (
	"fmt"

	"github.com/ignorant05/Uniflow/configs/github"
	errorhandling "github.com/ignorant05/Uniflow/internal/errorHandling"
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

	rootCmd.AddCommand(triggerCmd)
}

// trigger command main function
func runTriggerCmd(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		errMsg := fmt.Errorf("<?> Error: Not enough arguments.\n")
		errorhandling.HandleError(errMsg)
	}

	workflow := args[0]

	// if verbose mode active
	if verbose {
		fmt.Printf("<!> Info: Verbose mode enabled\n")
		fmt.Printf("   Workflow: %s\n", workflow)
		fmt.Printf("   Branch: %s\n", branch)
		fmt.Printf("   Profile: %s\n", profileName)
		if len(inputs) > 0 {
			fmt.Printf("   Inputs: %v\n", inputs)
		}
	}

	fmt.Printf("❯ Triggering workflow: %s\n", workflow)

	// create new client with profileName
	client, err := github.NewClientFromConfig(profileName)
	if err != nil {
		errMsg := fmt.Errorf("<?> Error: Field to create client.\n<?> Error: %w.\n", err)
		errorhandling.HandleError(errMsg)
	}

	// if verbose mode active
	if verbose {
		fmt.Printf("</> Info: Testing connection...\n")
		if err := client.TestConnection(); err != nil {
			errMsg := fmt.Errorf("<?> Error: Connection Testing failed...\n%w.\n", err)
			errorhandling.HandleError(errMsg)
		}
		fmt.Println("✓ Testing connection passed...")
	}

	// retrieves default repoditory
	owner, repo, err := client.GetDefaultRepository()
	if err != nil {
		errMsg := fmt.Errorf("<?> Error: Failed to get default repository...\n<?> Error: %w.\n", err)
		errorhandling.HandleError(errMsg)
	}

	// if verbose mode active
	if verbose {
		fmt.Printf("</> Info: %s/%s\n", owner, repo)
	}

	// parsing workflow inputs
	workflowInputs := make(map[string]interface{})
	for key, val := range inputs {
		workflowInputs[key] = val
	}

	// trigger workflow
	err = client.TriggerWorkflow(owner, repo, workflowFile, branch, workflowInputs)
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
	repoInfo, err := client.GetRepositoryInfo(owner, repo)
	if err != nil {
		errMsg := fmt.Errorf("<?> Error: Failed to retrieve repository %s/%s info.\n<?> Error: %w.\n", owner, repo, err)
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
}
