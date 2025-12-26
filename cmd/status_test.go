package cmd

import (
	"strings"
	"testing"
)

// Test status flags
func TestStatusFlags(t *testing.T) {
	tests := []struct {
		name         string
		flagName     string
		defaultValue string
	}{
		{
			name:         "all flag",
			flagName:     "all",
			defaultValue: "false",
		},
		{
			name:         "limit flag",
			flagName:     "limit",
			defaultValue: "5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := statusCmd.Flags().Lookup(tt.flagName)

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

// Test status aliases
func TestStatusCmdAlias(t *testing.T) {
	expectedAliases := []string{"s"}
	lenAl, lenExp := len(statusCmd.Aliases), len(expectedAliases)

	if lenAl != lenExp {
		t.Errorf("Number of aliases = %v, want %v", lenAl, lenExp)
	}

	for i, alias := range expectedAliases {
		if i == lenAl {
			break
		}

		if statusCmd.Aliases[i] != alias {
			t.Errorf("alias[%d] = %v, want %v", i, statusCmd.Aliases[i], alias)
		}
	}
}

// Test status use
func TestStatusCmdUse(t *testing.T) {
	use := "status [workflow]"
	actualUse := statusCmd.Use

	if use != actualUse {
		t.Errorf("Use = %v, want %v", use, actualUse)
	}
}

// Test status short
func TestStatusCmdShort(t *testing.T) {
	subShort := "Shows workflows's status"
	actualShort := statusCmd.Short

	if !strings.Contains(actualShort, subShort) {
		t.Errorf("Short Sippet = %v isn't in %v", subShort, actualShort)
	}
}
