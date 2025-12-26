# 🚀 Uniflow

A CLI tool for managing and triggering workflows across multiple CI/CD platforms (GitHub Actions, Jenkins, GitLab CI, CircleCI).

![Go Version](https://img.shields.io/badge/Go-1.24.4+-00ADD8?style=flat&logo=go) ![License](https://img.shields.io/badge/license-MIT-blue.svg) ![Status](https://img.shields.io/badge/status-active-success.svg)

---
## ✨ Features

- 🔧 **Multi-Platform Support**: GitHub Actions, Jenkins, GitLab CI, CircleCI
- 🔐 **Secure Configuration**: Environment variables and OS keyring integration
- 📦 **Profile Management**: Separate configs for dev, staging, production
- 🔄 **Real-Time Log Streaming**: Follow workflow execution with colored output
- ✅ **Config Validation**: Catch errors before running workflows
- 🎯 **Simple CLI**: Easy-to-use commands with helpful error messages

---
## 🎬 Quick Demo

```bash
# Initialize configuration
uniflow init

# List available workflows
uniflow workflows

# Trigger a workflow
uniflow trigger deploy.yml --input environment=prod

# Stream logs in real-time
uniflow logs deploy.yml --follow 

# Check workflow status
uniflow status deploy.yml
```

---
## 📦 Installation

### Quick Install (Recommended)

```bash
# Clone the repository
git clone https://github.com/ignorant05/uniflow.git
cd uniflow

# Build and install
make build
sudo mv uniflow /usr/local/bin/

# Verify installation
uniflow --version
```

### From Source

```bash
git clone https://github.com/ignorant05/uniflow.git
cd uniflow
make install
```

### Using installation script (NOT RECOMMENDED)
```bash 
# Clone the repository
git clone https://github.com/ignorant05/uniflow.git
cd uniflow

# Run installation script
./install.sh
```

> #NOTE: Please **verify** the [install.sh script](https://github.com/ignorant05/Uniflow/blob/main/install.sh) before proceeding with this installation                          method.

### Prerequisites

- Go 1.24.4 or higher
- Git
- GitHub personal access token (for GitHub Actions)

---
## 🚀 Quick Start

### 1. Initialize Configuration

```bash
uniflow init
```

This creates `~/.uniflow/config.yaml` with default settings.

### 2. Set Your API Tokens

```bash
export GITHUB_TOKEN="ghp_your_token_here" (only github for now is supported)
```

Or add to your shell profile (`~/.bashrc`, `~/.zshrc`):

```bash
echo 'export GITHUB_TOKEN="ghp_your_token_here"' >> ~/.bashrc (or `~/.zshrc`)
source ~/.bashrc (or `~/.zshrc`)
```

### 3. Configure Your Repository

```bash
uniflow config set profiles.default.github.default_repo "ignorant05/Uniflow"
```

### 4. Verify Configuration

```bash
# Listing all configuration for a profile (default)
uniflow config list

# This is needed to validate config
uniflow config validate
```

### 5. Trigger Your First Workflow

```bash
uniflow trigger deploy.yml --input environment=dev
```

---
## 📖 Documentation

- **Installation Guide** - Detailed setup instructions
- **Configuration Guide** - Complete config reference
- **Commands Reference** - All commands with examples
- **GitHub Actions Guide** - GitHub-specific setup
- **Contributing Guide** - How to contribute

---
## 🎯 Common Use Cases

### Monitor Deployments

```bash
# Trigger deployment
uniflow trigger deploy.yml --input environment=production

# Follow logs in real-time
uniflow logs deploy.yml --follow 
```

### Check Recent Runs

```bash
# Show status of all workflows
uniflow status

# Show detailed status for specific workflow
uniflow status deploy.yml --limit 10 --verbose
```

### Debug Failed Runs

```bash
# Find failed run
uniflow status deploy.yml

# View logs of specific run
uniflow logs --run-id 123456 --tail 100
```

### Multi-Environment Deployments

```bash
# Deploy to staging
uniflow trigger deploy.yml --profile staging --input env=staging

# Deploy to production
uniflow trigger deploy.yml --profile prod --input env=production
```

---
## 🔧 Configuration

Basic configuration file (`~/.uniflow/config.yaml`):

```yaml
default_platform: github
version: "1.0"

profiles:
  default:
    github:
      token: ${GITHUB_TOKEN}
      default_repo: owner/repo
      base_url: https://api.github.com
    
    # jenkins isn't supported yet 
    jenkins:
      url: https://jenkins.company.com
      username: admin
      token: ${JENKINS_TOKEN}
```

See Configuration Guide for complete reference.

---
## 📝 Available Commands

| Command     | Description              | Example                            |
| ----------- | ------------------------ | ---------------------------------- |
| `init`      | Initialize configuration | `uniflow init`                     |
| `config`    | Manage configuration     | `uniflow config list`              |
| `workflows` | List available workflows | `uniflow workflows`                |
| `trigger`   | Trigger a workflow       | `uniflow trigger deploy.yml`       |
| `status`    | Check workflow status    | `uniflow status deploy.yml`        |
| `logs`      | View workflow logs       | `uniflow logs deploy.yml --follow` |

See [Commands Reference](https://github.com/ignorant05/Uniflow/blob/main/doc/commands.md) for detailed commands documentation.

---
## 🎨 Features Showcase

### Real-Time Log Streaming

```bash
uniflow logs deploy.yml --follow 
```

- ✅ Color-coded output (errors in red, success in green)
- ⏰ Timestamps for each line
- 🔄 Live updates every 3 seconds
- ⚡ Graceful Ctrl+C handling

### Multi-Profile Support

```bash
# Development environment
uniflow trigger deploy.yml --profile dev

# Production environment
uniflow trigger deploy.yml --profile prod
```

### Status Monitoring

```bash
uniflow status deploy.yml
```

Shows:

- ✅ Run number and status
- 🔄 Success/failure conclusion
- ⏰ Triggered time
- 🔗 Direct link to run

---
## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](https://github.com/ignorant05/Uniflow/blob/main/doc/github-actions.md) for details.

### Quick Start for Contributors

```bash
# Fork and clone
git clone https://github.com/ignorant05/uniflow.git
cd uniflow

# Install dependencies
make install

# Run tests (don't forget to export your access tokens)
make test

# Build
make build
```

---
## 🐛 Troubleshooting

### Common Issues

**Error: "config file not found"**

```bash
# Solution: Initialize first
uniflow init
```

**Error: "GitHub token is required"** 

```bash
# Solution: Set environment variable
export GITHUB_TOKEN="ghp_your_token"
```

**Error: "workflow not found"**

```bash
# Solution: Check available workflows
uniflow workflows
```

See Installation Guide for more troubleshooting.

---
## 📊 Project Status

**Current Version:** 1.0.0

**Supported Platforms:**

- ✅ GitHub Actions (Full support)
- 🚧 Jenkins (Coming soon)
- 🚧 GitLab CI (Coming soon)
- 🚧 CircleCI (Coming soon)

---
## 📄 License

This project is licensed under the [MIT License](https://github.com/ignorant05/Uniflow/blob/main/LICENSE).

---
## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- Configuration powered by [Viper](https://github.com/spf13/viper)
- GitHub API via [go-github](https://github.com/google/go-github)
- Colored output using [fatih/color](https://github.com/fatih/color)

---
## 📧 Contact & Support

- **GitHub Issues**: [Report bugs](https://github.com/ignorant05/uniflow/issues)
- **Discussions**: [Ask questions](https://github.com/ignorant05/uniflow/discussions)
- **Email**: [oussamabaccara05@gmail.com](mailto:oussamabaccara05@gmail.com)
- **Discord**: [pebble](https://discord.gg/pebble)

---
## ⭐ Star History

If you find this project useful, please consider giving it a star!

---
## 💻 Contribution

As for contributions, see the [contribution guidelines](https://github.com/ignorant05/Uniflow/blob/main/doc/github-actions.md) 

---
Made by [ignorant05](https://github.com/ignorant05)
