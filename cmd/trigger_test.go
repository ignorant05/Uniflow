package cmd

import (
	"strings"
	"testing"
)

// Test trigger flags
func TestTriggerFlags(t *testing.T) {
	tests := []struct {
		name         string
		flagName     string
		defaultValue string
	}{
		{
			name:         "branch flag",
			flagName:     "branch",
			defaultValue: "main",
		},
		{
			name:         "workflow flag",
			flagName:     "workflow",
			defaultValue: "",
		},
		{
			name:         "input flag",
			flagName:     "input",
			defaultValue: "[]",
		},
		{
			name:         "profile flag",
			flagName:     "profile",
			defaultValue: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := triggerCmd.Flags().Lookup(tt.flagName)

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

// Test trigger aliases
func TestTriggerCmdAlias(t *testing.T) {
	expectedAliases := []string{"t"}
	lenAl, lenExp := len(triggerCmd.Aliases), len(expectedAliases)

	if lenAl != lenExp {
		t.Errorf("Number of aliases = %v, want %v", lenAl, lenExp)
	}

	for i, alias := range expectedAliases {
		if i == lenAl {
			break
		}

		if triggerCmd.Aliases[i] != alias {
			t.Errorf("alias[%d] = %v, want %v", i, triggerCmd.Aliases[i], alias)
		}
	}
}

// Test trigger use
func TestTriggerCmdUse(t *testing.T) {
	use := "trigger [workflow]"
	actualUse := triggerCmd.Use

	if use != actualUse {
		t.Errorf("Use = %v, want %v", use, actualUse)
	}
}

// Test trigger short
func TestTriggerCmdShort(t *testing.T) {
	subShort := "Trigger a workflow execution"
	actualShort := triggerCmd.Short

	if !strings.Contains(actualShort, subShort) {
		t.Errorf("Short Sippet = %v isn't in %v", subShort, actualShort)
	}
}

