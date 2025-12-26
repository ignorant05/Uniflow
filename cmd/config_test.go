package cmd

import (
	"fmt"
	"strings"
	"testing"
)

// Test config flags
func TestConfigListFlags(t *testing.T) {
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
		{
			name:         "secrets flag",
			flagName:     "show-secrets",
			defaultValue: "false",
		},
		{
			name:         "profile flag",
			flagName:     "profile",
			defaultValue: "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := configListCmd.Flags().Lookup(tt.flagName)

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

// Test config aliases
func TestConfigCmdAlias(t *testing.T) {
	expectedAliases := []string{"c"}
	lenAl, lenExp := len(configCmd.Aliases), len(expectedAliases)

	if lenAl != lenExp {
		t.Errorf("Number of aliases = %v, want %v", lenAl, lenExp)
	}

	for i, alias := range expectedAliases {
		if i == lenAl {
			break
		}

		if configCmd.Aliases[i] != alias {
			t.Errorf("alias[%d] = %v, want %v", i, configCmd.Aliases[i], alias)
		}
	}
}

// Test config use
func TestConfigCmdUse(t *testing.T) {
	use := "config"
	actualUse := configCmd.Use

	if use != actualUse {
		t.Errorf("Use = %v, want %v", use, actualUse)
	}
}

// Test config short
func TestConfigCmdShort(t *testing.T) {
	subShort := "Manage Uniflow configuration"
	actualShort := configCmd.Short

	if !strings.Contains(actualShort, subShort) {
		t.Errorf("Short Sippet = %v isn't in %v", subShort, actualShort)
	}
}

// Testing configCmd subcommands
func TestConfigCmdSubcommands(t *testing.T) {
	tests := []struct {
		command    string
		subCommand string
	}{
		{
			command:    "config",
			subCommand: "list",
		},
		{
			command:    "config",
			subCommand: "set",
		},
		{
			command:    "config",
			subCommand: "get",
		},
		{
			command:    "config",
			subCommand: "validate",
		},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_%s", tt.command, tt.subCommand), func(t *testing.T) {
			got := configCmd.Commands()
			found := false
			for _, cmd := range got {
				if tt.subCommand == cmd.Name() {
					found = !found
					break
				}
			}

			if !found {
				t.Errorf("Subcommand %v, not found for %v", tt.subCommand, configCmd.Name())
			}
		})
	}
}
