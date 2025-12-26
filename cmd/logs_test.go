package cmd

import (
	"strings"
	"testing"
)

// Test logs cmd flags
func TestLogsCmdFlag(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
		wantType string
	}{
		{
			name:     "run-id flag exists",
			flagName: "run-id",
			wantType: "int64",
		},
		{
			name:     "job flag exists",
			flagName: "job",
			wantType: "string",
		},
		{
			name:     "follow flag exists",
			flagName: "follow",
			wantType: "bool",
		},
		{
			name:     "tail flag exists",
			flagName: "tail",
			wantType: "int",
		},
		{
			name:     "no-color flag exists",
			flagName: "no-color",
			wantType: "bool",
		},
		{
			name:     "platform flag exists",
			flagName: "platform",
			wantType: "string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := logsCmd.Flags().Lookup(tt.flagName)

			if flag == nil {
				t.Errorf("flag %s does not exist", tt.flagName)
				return
			}

			if flag.Value.Type() != tt.wantType {
				t.Errorf("flag %s type = %v, want %v", tt.flagName, flag.Value.Type(), tt.wantType)
			}
		})
	}
}

// Test logs aliases
func TestLogsCmdAlias(t *testing.T) {
	expectedAliases := []string{"l"}
	lenAl, lenExp := len(logsCmd.Aliases), len(expectedAliases)

	if lenAl != lenExp {
		t.Errorf("Number of aliases = %v, want %v", lenAl, lenExp)
	}

	for i, alias := range expectedAliases {
		if i == lenAl {
			break
		}

		if logsCmd.Aliases[i] != alias {
			t.Errorf("alias[%d] = %v, want %v", i, logsCmd.Aliases[i], alias)
		}
	}
}

// Test logs use
func TestLogsCmdUse(t *testing.T) {
	use := "logs [workflow]"
	actualUse := logsCmd.Use

	if use != actualUse {
		t.Errorf("Use = %v, want %v", use, actualUse)
	}
}

// Test logs short
func TestLogsCmdShort(t *testing.T) {
	subShort := "View workflow execution logs"
	actualShort := logsCmd.Short

	if !strings.Contains(actualShort, subShort) {
		t.Errorf("Short Sippet = %v isn't in %v", subShort, actualShort)
	}
}

// Testing run func is sat
func TestLogsCmdRunFunction(t *testing.T) {
	if logsCmd.Run == nil && logsCmd.RunE == nil {
		t.Error("logs command has no Run or RunE function")
	}
}

// Test logs cmd args
func TestLogsCmdArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantError bool
	}{
		{
			name:      "no args",
			args:      []string{},
			wantError: false,
		},
		{
			name:      "one arg",
			args:      []string{"deploy.yml"},
			wantError: false,
		},
		{
			name:      "two args",
			args:      []string{"deploy.yml", "extra"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := logsCmd.Args(logsCmd, tt.args)

			if tt.wantError && err == nil {
				t.Error("expected error for invalid args, got nil")
			}

			if !tt.wantError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// Test logs flags
func TestLogsFlags(t *testing.T) {
	tests := []struct {
		name         string
		flagName     string
		defaultValue string
	}{
		{
			name:         "run-id default",
			flagName:     "run-id",
			defaultValue: "0",
		},
		{
			name:         "job default",
			flagName:     "job",
			defaultValue: "",
		},
		{
			name:         "follow default",
			flagName:     "follow",
			defaultValue: "false",
		},
		{
			name:         "tail default",
			flagName:     "tail",
			defaultValue: "0",
		},
		{
			name:         "no-color default",
			flagName:     "no-color",
			defaultValue: "false",
		},
		{
			name:         "platform default",
			flagName:     "platform",
			defaultValue: "github",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := logsCmd.Flags().Lookup(tt.flagName)

			if flag == nil {
				t.Errorf("flag %s does not exist", tt.flagName)
				return
			}

			if flag.DefValue != tt.defaultValue {
				t.Errorf("flag %s default = %v, want %v", tt.flagName, flag.DefValue, tt.defaultValue)
			}
		})
	}
}

// Test log validation
func TestLogsCmdValidation(t *testing.T) {
	tests := []struct {
		name    string
		check   func() bool
		wantErr bool
	}{
		{
			name: "command has use string",
			check: func() bool {
				return logsCmd.Use != ""
			},
			wantErr: false,
		},
		{
			name: "command has short description",
			check: func() bool {
				return logsCmd.Short != ""
			},
			wantErr: false,
		},
		{
			name: "command has long description",
			check: func() bool {
				return logsCmd.Long != ""
			},
			wantErr: false,
		},
		{
			name: "command has run function",
			check: func() bool {
				return logsCmd.Run != nil || logsCmd.RunE != nil
			},
			wantErr: false,
		},
		{
			name: "command has at least one alias",
			check: func() bool {
				return len(logsCmd.Aliases) > 0
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.check()
			if !result {
				t.Errorf("%s check failed", tt.name)
			}
		})
	}
}
