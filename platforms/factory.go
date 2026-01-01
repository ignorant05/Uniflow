package platforms

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/ignorant05/Uniflow/internal/config"
	platforms "github.com/ignorant05/Uniflow/platforms/adapters"
	"github.com/ignorant05/Uniflow/platforms/constants"
	"github.com/ignorant05/Uniflow/platforms/github"
	"github.com/ignorant05/Uniflow/types"
)

type Factory struct {
	Config *config.Config
}

type PlatformInfo struct {
	Platform   string
	ConfigPath string
	Confidence int
}

// NewFactory creates Factory object from configuration
//
// Parameters:
//   - cfg: configuration
//
// Example:
// factory := platforms.NewFactory(cfg)
func NewFactory(cfg *config.Config) *Factory {
	return &Factory{Config: cfg}
}

// CreateClientForProfile Creates a client with/or without the platform name and the profile name
// NOTE: if not provided, it automatically uses default variables
//
// Parameters:
//   - ctx: the context variable
//   - platform: user selected platform name (default: "github")
//   - profileName: user selected profile name (default: "default")
//
// Example:
// client, err := f.CreateClientForProfile(ctx, "github", "mine")
func (f *Factory) CreateClientForProfile(ctx context.Context, platform, profileName string) (PlatformClient, error) {
	if platform == "" {
		platform = f.Config.DefaultPlatform
	}

	if profileName == "" {
		profileName = "default"
	}

	profile, err := f.Config.GetProfile(profileName)
	if err != nil {
		return nil, err
	}

	switch platform {
	case constants.GITHUB_PLATFORM:
		return f.CreateGithubClient(ctx, profile)
	default:
		return nil, &types.PlatformError{
			Code:    "unsupported_platform",
			Message: fmt.Sprintf("Platform %s is not supported.", platform),
		}
	}
}

// CreateClientAutoDetectPlatform creates a client with/or without the profile name, but it scans the dir to generate a corresponding config
// NOTE: if not provided, it automatically uses default variables
//
// This allows commands to work without --platform flag.
//
// Detection strategy:
//  1. Check for .github/workflows/ directory -> GitHub Actions
//  2. Check for Jenkinsfile -> Jenkins
//  3. Check for .gitlab-ci.yml -> GitLab CI
//  4. Check for .circleci/config.yml -> CircleCI
//  5. Fall back to default platform from config
//
// Parameters:
//   - ctx: the context variable
//   - profileName: user selected profile name (default: "default")
//
// Example:
// client, err := f.CreateClientAutoDetectPlatform(ctx, "mine")
func (f *Factory) CreateClientAutoDetectPlatform(ctx context.Context, profileName string) (PlatformClient, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	platformInfo, err := f.detectPlatformDirectory(cwd)
	if err != nil {
		return nil, err
	}

	if platformInfo.Confidence > 0 {
		return f.CreateClientForProfile(ctx, platformInfo.Platform, profileName)
	}

	return f.CreateClientAutoDetectPlatform(ctx, profileName)
}

// detectPlatformDirectory detects platforms directory and returns it's information if existed
//
// Parameters:
//   - dir: current working directory (aka the directory that you're already in it)
//
// Example:
// info, err := f.detectPlatformDirectory("~/Uniflow")
func (f *Factory) detectPlatformDirectory(dir string) (*PlatformInfo, error) {
	detectors := []struct {
		Platform   string
		Paths      []string
		Confidence int
	}{
		{
			Platform:   constants.GITHUB_PLATFORM,
			Paths:      constants.GITHUB_PATHS,
			Confidence: 100,
		},
		{
			Platform:   constants.JENKINS_PLATFORM,
			Paths:      constants.JENKINS_PATHS,
			Confidence: 100,
		}, {
			Platform:   constants.GITLAB_PLATFORM,
			Paths:      constants.GITLAB_PATHS,
			Confidence: 100,
		}, {
			Platform:   constants.CIRCLE_CI_PLATFORM,
			Paths:      constants.CIRCLE_CI_PATHS,
			Confidence: 100,
		},
	}

	for _, detector := range detectors {
		for _, path := range detector.Paths {
			fullFilePattern := filepath.Join(dir, string(path))

			if strings.HasSuffix(fullFilePattern, "/") {
				if info, err := os.Stat(fullFilePattern); err == nil && info.IsDir() {
					return &PlatformInfo{
						Platform:   detector.Platform,
						ConfigPath: fullFilePattern,
						Confidence: detector.Confidence,
					}, nil
				}
			} else {
				if matches, err := filepath.Glob(fullFilePattern); err == nil && len(matches) > 0 {
					return &PlatformInfo{
						Platform:   detector.Platform,
						ConfigPath: matches[0],
						Confidence: detector.Confidence,
					}, nil
				}
			}
		}

	}

	return &PlatformInfo{
		Platform:   constants.DEFAULT_PLATFORM,
		Confidence: 0,
	}, nil
}

// CreateClientForProfile Creates a client from a config profile
//
// Parameters:
//   - ctx: the context variable
//   - profile: profile configuration
//
// Example:
// client, err := f.CreateGithubClient(ctx, profile)
func (f *Factory) CreateGithubClient(ctx context.Context, profile *config.Profile) (PlatformClient, error) {
	if profile.Github == nil {
		return nil, &types.PlatformError{
			Code:     "not_configured",
			Message:  "Github is not configured for this profile",
			Platform: "github",
		}
	}

	client, err := github.NewClientFromProfile(ctx, profile)
	if err != nil {
		return nil, err
	}

	return platforms.NewGithubAdapter(client)
}

// ListSupportedPlatforms list all platforms supported by uniflow
//
// Parameters:
//   - None
//
// Example:
// supportedPlatforms := ListSupportedPlatforms()
func ListSupportedPlatforms() []string {
	return []string{
		constants.GITHUB_PLATFORM,
		constants.JENKINS_PLATFORM,
		constants.GITLAB_PLATFORM,
		constants.CIRCLE_CI_PLATFORM,
	}
}

// IsPlatformSupported checks if the platform is supported or not
//
// Parameters:
//   - platform: platform name
//
// Example:
// supported := IsPlatformSupported("github")
func IsPlatformSupported(platform string) bool {
	platform = strings.ToLower(platform)
	supportedPlatforms := ListSupportedPlatforms()

	return slices.Contains(supportedPlatforms, platform)
}
