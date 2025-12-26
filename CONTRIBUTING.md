# Contributing to Uniflow

Thank you for your interest in contributing to Uniflow! 🎉

This document provides guidelines and instructions for contributing.

---
## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Pull Request Process](#pull-request-process)
- [Testing](#testing)
- [Documentation](#documentation)
- [Community](#community)

---
## 📜 Code of Conduct

### Our Pledge

We pledge to make participation in our project a harassment-free experience for everyone, regardless of age, body size, disability, ethnicity, gender identity, level of experience, nationality, personal appearance, race, religion, or sexual identity and orientation.

### Our Standards

**Positive behavior includes:**

- Using welcoming and inclusive language
- Being respectful of differing viewpoints
- Gracefully accepting constructive criticism
- Focusing on what is best for the community
- Showing empathy towards other community members

**Unacceptable behavior includes:**

- Trolling, insulting/derogatory comments, and personal attacks
- Public or private harassment
- Publishing others' private information without permission
- Other conduct which could reasonably be considered inappropriate

---
## 🚀 Getting Started

### Ways to Contribute

- 🐛 **Report bugs** - Found a bug? Let us know!
- ✨ **Suggest features** - Have an idea? Share it!
- 📝 **Improve documentation** - Help others understand
- 💻 **Write code** - Fix bugs or add features
- 🧪 **Write tests** - Improve code coverage

### First Time Contributors

Look for issues labeled:

- `good first issue` - Great for beginners
- `help wanted` - We need help with these
- `documentation` - Documentation improvements

---
## 🛠️ Development Setup

### Prerequisites

- Go 1.24.4 or higher
- Git
- GitHub account
- IDE 

### Fork and Clone

```bash
# 1. Fork the repository on GitHub
# Click the "Fork" button at https://github.com/ignorant05/uniflow

# 2. Clone your fork
git clone https://github.com/YOUR_USERNAME/uniflow.git
cd uniflow

# 3. Add upstream remote
git remote add upstream https://github.com/ignorant05/uniflow.git

# 4. Verify remotes
git remote -vv
```

### Install Dependencies

```bash
# Download dependencies
make install

# Verify everything works
make build
./uniflow --version
```

### Run Tests

```bash
# Run all tests
make test

# Run with coverage
make cov-test
```

### Development Workflow

```bash
# 1. Create a new branch
git checkout -b feature/my-feature

# 2. Make your changes
# ... edit files ...

# 3. Test your changes
make test
make build
./uniflow --help

# 4. Commit your changes
git add .
git commit -m "Feature: description"

# 5. Push to your fork
git push origin feature/my-feature

# 6. Create Pull Request on GitHub
```

## 🔄 Making Changes

### Branch Naming

Use descriptive branch names:

```bash
feature/add-gitlab-support
enhacement/enhanced-github-client
fix/config-validation-bug
docs/improve-readme
refactor/simplify-logger
test/add-status-tests
```

### Commit Messages

Follow conventional commits:

```
type(scope): subject

body (optional)

footer (optional)
```
**Notes:** there's no concrete structure for commit messages yet.

**Types:**

- `feature` - New feature
- `enhacement` - Enhancement for existing feature... etc
- `fix` - Bug fix
- `docs` - Documentation only
- `style` - Code style (formatting, etc.)
- `refactor` - Code refactoring
- `test` - Adding tests
-
**Examples:**

```bash
feat(jenkins): add Jenkins platform support

Implements basic Jenkins integration including:
- API client setup
- Job triggering
- Status checking

Closes (or Solves/or Fixes) #42

---

fix(config): validate repository format

Fixes bug where invalid repo format wasn't caught
during config validation.

Fixes #38

---

docs(readme): add installation instructions

Adds detailed installation steps for Linux, macOS, and Windows.
```

**Rules:**

- Use `gofmt` to format code
- Use meaningful variable names
- Add comments for exported functions
- Keep functions focused and small
- Handle errors properly

### File Organization

```
Uniflow/
.
├── cmd  # commands with their sources of truth, helpers and tests
│   ├── config.go
│   ├── config_test.go
│   ├── constants
│   │   └── config_constants.go
│   ├── helpers
│   │   ├── config_helper.go
│   │   └── status_helper.go
│   ├── init.go
│   ├── init_test.go
│   ├── logs.go
│   ├── logs_test.go
│   ├── root.go
│   ├── status.go
│   ├── status_test.go
│   ├── trigger.go
│   ├── trigger_test.go
│   ├── workflows.go
│   └── workflows_test.go
├── configs   # platform's client configurations
│   ├── github
│   │   ├── constants
│   │   │   └── github.go
│   │   ├── github_client.go
│   │   ├── github_factory.go
│   │   ├── github_workflows.go
│   │   └── helpers
│   │       └── github_client_helper.go
│   └── tests       # platforms's unit tests 
│       └── unit
│           └── github
│               ├── clientCreation_test.go
│               ├── getDefaultRepo_test.go
│               ├── getWorkflowRun_test.go
│               ├── listWorkflows_test.go
│               ├── mock_server.go
│               ├── triggerWorkflow_test.go
│               └── workflowRunJobs_test.go
├── CONTRIBUTING.md        # contribution guidelines
├── doc
│   ├── commands.md        # commands documentation
│   └── github-actions.md  # github-actions documentation
├── go.mod
├── go.sum
├── install.sh        # Installation script
├── internal          # internal configurations
│   ├── config
│   │   ├── config.go
│   │   ├── loader.go
│   │   ├── platforms.go
│   │   └── validator.go
│   ├── constants
│   │   ├── config
│   │   │   ├── github_constants.go
│   │   │   ├── loader_constants.go
│   │   │   └── validator_constants.go
│   │   ├── credentials
│   │   │   └── keyring_constants.go
│   │   └── logs
│   │       └── streamer_constants.go
│   ├── credentials
│   │   └── keyring.go
│   ├── errorHandling
│   │   └── errorHandling.go
│   ├── helpers
│   │   ├── loader_helpers.go
│   │   ├── logs_helper.go
│   │   └── validator_helpers.go
│   └── logs   
│       ├── downloader.go
│       ├── streamer.go
│       └── streamer_test.go
├── LICENSE
├── main.go           
├── makefile          # makefile for dev experience
├── README.md
└── replacements.txt  # Add your tokens here and filtre them with git filtre repo
					  # Example: 
					  # git filter-repo --force --replace-text replacements.txt
```

---
## 📥 Pull Request Process

### Before Submitting

- ✅ Tests pass: `make test`
- ✅ Code is formatted: `make fmt`
- ✅ No linting errors: `make lint`
- ✅ Documentation updated

### PR Template

When creating a PR, include:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring

## Testing
- [ ] Unit tests added/updated
- [ ] Manual testing completed
- [ ] All tests pass

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
- [ ] Commit messages follow conventions

## Related Issues
Closes #42
Related to #38
```

### PR Review Process

1. **Automated checks** run (tests, linting)
2. **Maintainer review** (1-3 days)
3. **Feedback addressed** (if needed)
4. **Approval** from maintainer
5. **Merge** to main branch

### After Merge

```bash
# Update your local main
git checkout main
git pull upstream main

# Delete your feature branch
git branch -d feature/my-feature
git push origin --delete feature/my-feature
```
### README Updates

If you change functionality, update:

- `README.md` - Quick start section
- `docs/commands.md` - Command reference
- `docs/configuration.md` - Config options
- `docs/github-actions.md` - GitHub specifics

---
## 🐛 Reporting Bugs

### Bug Report Template

```markdown
**Describe the bug**
A clear description of what the bug is.

**To Reproduce**
Steps to reproduce:
1. Run command '...'
2. See error '...'

**Expected behavior**
What you expected to happen.

**Actual behavior**
What actually happened.

**Environment**
- OS: [e.g., Ubuntu 22.04]
- Go version: [e.g., 1.21.5]
- Uniflow version: [e.g., 1.0.0]

**Additional context**
Any other information about the problem.
```

---
## ✨ Requesting Features

### Feature Request Template

```markdown
**Is your feature request related to a problem?**
Description of the problem.

**Describe the solution you'd like**
Clear description of what you want to happen.

**Describe alternatives you've considered**
Other solutions you've thought about.

**Additional context**
Any other context, mockups, or examples.
```

---
## 💬 Community

### Getting Help

- **GitHub Issues** - Bug reports and features
- **GitHub Discussions** - Questions and ideas
- **Email** - your.email@example.com

### Stay Updated

- ⭐ Star the repository
- 👀 Watch for releases
- 📢 Follow announcements

---
## 🎓 Learning Resources

### Go Resources

- [Tour of Go](https://tour.golang.org/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go by Example](https://gobyexample.com/)

### Testing

- [Go Testing](https://golang.org/pkg/testing/)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

### GitHub Actions

- [GitHub Actions Docs](https://docs.github.com/en/actions)
- [Workflow Syntax](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions)

---
## 📄 License

By contributing, you agree that your contributions will be licensed under the MIT License.

---
## 🙏 Thank You!

Every contribution, no matter how small, is appreciated. Thank you for helping make Uniflow better!

---

**Questions?** 
 - Open an issue or start a discussion!.
 - Or contact me on: 
	[![Discord](https://img.shields.io/badge/Discord-%237289DA.svg?logo=discord&logoColor=white)](https://discord.gg/pebble)[![LinkedIn](https://img.shields.io/badge/LinkedIn-%230077B5.svg?logo=linkedin&logoColor=white)](https://linkedin.com/in/oussama-baccara-64552a282)[![Reddit](https://img.shields.io/badge/Reddit-%23FF4500.svg?logo=Reddit&logoColor=white)](https://reddit.com/user/AfraidComposer6150)[![Email](https://img.shields.io/badge/Email-D14836?logo=gmail&logoColor=white)](mailto:oussamabaccara05@gmail.com)

**Ready to contribute?** Check out [good first issues](https://github.com/ignorant05/uniflow/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)!
