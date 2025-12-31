package constants

const (
	DEFAULT_PLATFORM   = "github"
	GITHUB_PLATFORM    = "github"
	JENKINS_PLATFORM   = "jenkins"
	GITLAB_PLATFORM    = "gitlab"
	CIRCLE_CI_PLATFORM = "circleci"
)

var (
	GITHUB_PATHS    = []string{".github/workflows/", ".github/workflows/*.yaml", ".github/workflows/*.yml"}
	JENKINS_PATHS   = []string{"Jenkinsfile", "Jenkinsfile.groovy", "Jenkinsfile.jenkins"}
	GITLAB_PATHS    = []string{".gitlab-ci.yml"}
	CIRCLE_CI_PATHS = []string{".circleci/config.yml", ".circleci/config.yaml"}
)
