package cmd

import (
	"strings"
	"testing"
)

// Test workflows flags
func TestWorkflowsFlags(t *testing.T) {
	tests := []struct {
		name         string
		flagName     string
		defaultValue string
	}{
		{
			name:         "with-dispatch flag",
			flagName:     "with-dispatch",
			defaultValue: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := workflowsCmd.Flags().Lookup(tt.flagName)

			if flag == nil {
				t.Errorf("flag %s does not exist", tt.flagName)
				return
			}

			if flag.DefValue != tt.defaultValue {
				t.Errorf("flag %s flag = %v, want %v", tt.flagName, flag.DefValue, tt.defaultValue)
			}
		})
	}
}

// Test workflows aliases
func TestWorkflowsCmdAlias(t *testing.T) {
	expectedAliases := []string{"wf"}
	lenAl, lenExp := len(workflowsCmd.Aliases), len(expectedAliases)

	if lenAl != lenExp {
		t.Errorf("Number of aliases = %v, want %v", lenAl, lenExp)
	}

	for i, alias := range expectedAliases {
		if i == lenAl {
			break
		}

		if workflowsCmd.Aliases[i] != alias {
			t.Errorf("alias[%d] = %v, want %v", i, workflowsCmd.Aliases[i], alias)
		}
	}
}

// Test workflows use
func TestWorkflowsCmdUse(t *testing.T) {
	use := "workflows"
	actualUse := workflowsCmd.Use

	if use != actualUse {
		t.Errorf("Use = %v, want %v", use, actualUse)
	}
}

// Test workflows short
func TestWorkflowsCmdShort(t *testing.T) {
	subShort := "List available workflows in the repository"
	actualShort := workflowsCmd.Short

	if !strings.Contains(actualShort, subShort) {
		t.Errorf("Short Sippet = %v isn't in %v", subShort, actualShort)
	}
}
