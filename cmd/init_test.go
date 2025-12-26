package cmd

import (
	"strings"
	"testing"
)

// Test init flags
func TestInitFlags(t *testing.T) {
	tests := []struct {
		name         string
		flagName     string
		defaultValue string
	}{
		{
			name:         "force flag",
			flagName:     "force",
			defaultValue: "false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := initCmd.Flags().Lookup(tt.flagName)

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

// Test init aliases
func TestInitCmdAlias(t *testing.T) {
	expectedAliases := []string{"i"}
	lenAl, lenExp := len(initCmd.Aliases), len(expectedAliases)

	if lenAl != lenExp {
		t.Errorf("Number of aliases = %v, want %v", lenAl, lenExp)
	}

	for i, alias := range expectedAliases {
		if i == lenAl {
			break
		}

		if initCmd.Aliases[i] != alias {
			t.Errorf("alias[%d] = %v, want %v", i, initCmd.Aliases[i], alias)
		}
	}
}

// Test init use
func TestInitCmdUse(t *testing.T) {
	use := "init"
	actualUse := initCmd.Use

	if use != actualUse {
		t.Errorf("Use = %v, want %v", use, actualUse)
	}
}

// Test init short
func TestInitCmdShort(t *testing.T) {
	subShort := "Initialize uniflow configuration"
	actualShort := initCmd.Short

	if !strings.Contains(actualShort, subShort) {
		t.Errorf("Short Sippet = %v isn't in %v", subShort, actualShort)
	}
}
