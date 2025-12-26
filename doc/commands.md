# Commands Reference

Complete reference for all Uniflow commands.

## 📋 Command Overview

| Command     | Description              | Aliases |
| ----------- | ------------------------ | ------- |
| `init`      | Initialize configuration | `i`     |
| `config`    | Manage configuration     | `c`     |
| `workflows` | List available workflows | `w`     |
| `trigger`   | Trigger a workflow       | `t`     |
| `status`    | Check workflow status    | `s`     |
| `logs`      | View workflow logs       | `l`     |

## 🎯 Global Flags

Available for all commands:

| Flag        | Short | Description           | Default   |
| ----------- | ----- | --------------------- | --------- |
| `--verbose` | `-v`  | Enable verbose output | `false`   |
| `--profile` | `-p`  | Config profile to use | `default` |
| `--help`    | `-h`  | Show help             | -         |
| `--version` | -     | Show version          | -         |

---

## `init` Command

Initialize Uniflow configuration.

### Usage

```bash
uniflow init 
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--force` / `-f` | Overwrite existing config | `false` |

### Examples

```bash
# Basic initialization
uniflow init

# Force overwrite existing config
uniflow init --force

# With verbose output
uniflow init --verbose
```

### What It Does

1. Creates `~/.uniflow/` directory
2. Generates `config.yaml` with defaults
3. Creates `logs/` subdirectory
4. Sets proper file permissions

### Output

```
❯ Initializing uniflow configuration...
✓ Created config directory: /home/user/.uniflow
✓ Created config file: /home/user/.uniflow/config.yaml
✓ Created logs directory: /home/user/.uniflow/logs

✓ Initialization complete!

</> Info:  Next steps:
   1. Set your API tokens as environment variables:
	export GITHUB_TOKEN=your_token_here
	export JENKINS_TOKEN=your_token_here

   2. Or update the config file directly:
   /home/user/.uniflow/config.yaml

   3. Verify your configuration:
   uniflow config list

   4. Start triggering workflows:
   uniflow trigger my-workflow
```

---

## `config` Command

Manage configuration settings.

### Subcommands

- `list` - Show current configuration
- `set` - Update configuration value
- `get` - Get specific configuration value
- `validate` - Validate configuration

### `config list`

Show current configuration.

#### Usage

```bash
uniflow config list [flags]
```

#### Flags

| Flag             | Short | Description             | Default   |
| ---------------- | ----- | ----------------------- | --------- |
| `--profile`      | `-p`  | Profile to display      | `default` |
| `--show-secrets` | `-s`  | Show sensitive values   | `false`   |
| `--force`        | `-f`  | Force show long secrets | `false`   |

#### Examples

```bash
# List default profile
uniflow config list

# List production profile
uniflow config list --profile prod

# Show with secrets unmasked
uniflow config list --show-secrets

# Verbose mode
uniflow config list --verbose
```

### `config set`

Update a configuration value.

#### Usage

```bash
uniflow config set <key> <value>
```

#### Key Format

```
profiles.<profile>.<platform>.<field>
```

#### Examples

```bash
# Set default repository
uniflow config set profiles.default.github.default_repo "owner/repo"

# Set Jenkins URL
uniflow config set profiles.default.jenkins.url "https://jenkins.local"

# Set Jenkins username
uniflow config set profiles.default.jenkins.username "admin"

# Change default platform
uniflow config set default_platform github
```

### `config get`

Get a specific configuration value.

#### Usage

```bash
uniflow config get <key>
```

#### Examples

```bash
# Get default platform
uniflow config get default_platform

# Get GitHub repo
uniflow config get profiles.default.github.default_repo

# Get Jenkins URL
uniflow config get profiles.default.jenkins.url
```

### `config validate`

Validate configuration file.

#### Usage

```bash
uniflow config validate
```

#### Examples

```bash
# Validate current config
uniflow config validate

# Validate with verbose output
uniflow config validate --verbose
```

#### Output (Success)

```
❯ Validating configuration...
✅ Configuration is valid!
```

#### Output (Errors)

```
✓ Validating configuration...
<?> Error: Configuration validation failed:

  1. profiles.default.github.token: token is required
  2. profiles.default.jenkins.url: must be a valid URL

Please fix these issues.
```

---

## `workflows` Command

List available workflows in the repository.

### Usage

```bash
uniflow workflows [flags]
```

### Aliases

- `w`
- `workflows`

### Flags

| Flag        | Short | Description           | Default   |
| ----------- | ----- | --------------------- | --------- |
| `--profile` | `-p`  | Config profile to use | `default` |

### Examples

```bash
# List workflows
uniflow workflows

# Use alias
uniflow wf

# Use specific profile
uniflow workflows --profile prod

# Verbose mode
uniflow workflows --verbose
```

### Output

```
❯ Listing available workflows...

✓ Found 3 workflow(s):

1. Deploy Application
   File: .github/workflows/deploy.yml
   State: active

2. Run Tests
   File: .github/workflows/test.yml
   State: active

3. Build Docker
   File: .github/workflows/build.yml
   State: active

💡 Trigger a workflow with:
   uniflow trigger deploy.yml
```

---

## `trigger` Command

Trigger a workflow execution.

### Usage

```bash
uniflow trigger <workflow> [flags]
```

### Aliases

- `t`

### Arguments

| Argument   | Description       | Required |
| ---------- | ----------------- | -------- |
| `workflow` | Workflow filename | ✅ Yes    |

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--branch` | `-b` | Branch to run on | `main` |
| `--input` | `-i` | Workflow inputs (key=value) | - |
| `--profile` | `-p` | Config profile to use | `default` |
| `--platform` | - | Platform to use | `github` |

### Examples

```bash
# Basic trigger
uniflow trigger deploy.yml

# Use alias
uniflow t deploy.yml

# Trigger on specific branch
uniflow trigger deploy.yml --branch develop

# With workflow inputs
uniflow trigger deploy.yml --input environment=prod --input version=v1.0

# Multiple inputs
uniflow trigger deploy.yml \
  --input environment=staging \
  --input version=v2.1.0 \
  --input dry_run=false

# Use specific profile
uniflow trigger deploy.yml --profile prod

# Verbose mode
uniflow trigger deploy.yml --verbose
```

### Output

```
❯ Triggering workflow: deploy.yml
</> Info: Testing connection...
✓ Successfully authenticated as: ignorant05
✓ Testing connection passed...
</> Info: ignorant05/my-repo
   Workflow triggered successfully!
   Repository: ignorant05/my-repo
   Workflow: deploy.yml
   Branch: main
   View at: https://github.com/ignorant05/my-repo/actions

   Inputs:
   environment: prod
   version: v1.0
```

---

## `status` Command

Check workflow run status.

### Usage

```bash
uniflow status [workflow] [flags]
```

### Aliases

- `s`

### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `workflow` | Workflow filename (optional) | ❌ No |

### Flags

| Flag        | Short | Description            | Default   |
| ----------- | ----- | ---------------------- | --------- |
| `--all`     | `-a`  | Show all runs          | `false`   |
| `--limit`   | `-l`  | Number of runs to show | `5`       |
| `--profile` | `-p`  | Config profile to use  | `default` |

### Examples

```bash
# Show all workflows (latest run each)
uniflow status

# Show specific workflow
uniflow status deploy.yml

# Show more runs
uniflow status deploy.yml --limit 10

# Show all available runs
uniflow status deploy.yml --all

# Use alias
uniflow s deploy.yml

# Verbose mode
uniflow status deploy.yml --verbose
```

### Output (All Workflows)

```
❯ Checking status of all workflows...

   Deploy Application
   File: deploy.yml
────────────────────────────────────────────────────────────────────────────────
  Run #47
    Status:     Completed
    Conclusion: Success
    Branch:     main
    Triggered:  2 hours ago

   Run Tests
   File: test.yml
────────────────────────────────────────────────────────────────────────────────
  Run #152
    Status:     In Progress
    Conclusion: Pending
    Branch:     develop
    Triggered:  5 minutes ago

✓ Displayed status for 2 workflow(s)
```

### Output (Specific Workflow)

```
❯ Checking status of workflow: deploy.yml

   Workflow: Deploy Application
   File: deploy.yml
────────────────────────────────────────────────────────────────────────────────

  Recent Runs (showing 5 of 47):

  Run #47
    Status:     Completed
    Conclusion: Success
    Branch:     main
    Triggered:  2 hours ago

  Run #46
    Status:     Completed
    Conclusion: Failure
    Branch:     develop
    Triggered:  1 day ago

  Run #45
    Status:     Completed
    Conclusion: Success
    Branch:     main
    Triggered:  2 days ago
```

---
## `logs` Command

View workflow execution logs.

### Usage

```bash
uniflow logs [workflow] [flags]
```

### Aliases

- `l`

### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `workflow` | Workflow filename (optional) | ❌ No |

### Flags

| Flag         | Short | Description              | Default   |
| ------------ | ----- | ------------------------ | --------- |
| `--run-id`   | -     | Specific run ID          | `0`       |
| `--job`      | `-j`  | Specific job name        | `""`      |
| `--follow`   | `-f`  | Follow logs in real-time | `false`   |
| `--tail`     | `-t`  | Show last N lines        | `0` (all) |
| `--no-color` | -     | Disable colored output   | `false`   |
| `--platform` | -     | Platform to use          | `github`  |
| `--profile`  | `-p`  | Config profile to use    | `default` |

### Examples

```bash
# Show logs for latest run
uniflow logs deploy.yml

# Show logs for specific run
uniflow logs --run-id 123456

# Follow logs in real-time
uniflow logs deploy.yml --follow

# Use alias with follow
uniflow l deploy.yml -f

# Show last 50 lines
uniflow logs deploy.yml --tail 50

# Combine options
uniflow logs deploy.yml --follow --tail 100

# Show specific job logs
uniflow logs deploy.yml --job "build"

# No colors (for piping)
uniflow logs deploy.yml --no-color > logs.txt

# Verbose mode
uniflow logs deploy.yml --verbose
```

### Output (Basic)

```
❯❯❯ Fetching logs for workflow: deploy.yml

📋 Workflow Run
────────────────────────────────────────────────────────────────────────────────
Name:       Deploy Application
Run #:      47
Status:     Completed
Branch:     main
Commit:     a1b2c3d
Actor:      ignorant05
────────────────────────────────────────────────────────────────────────────────

   Job: build
   Status: Completed
   Result: Success
────────────────────────────────────────────────────────────────────────────────
Run actions/checkout@v4
Syncing repository: owner/repo
Successfully checked out main branch
Run npm install
Dependencies installed successfully
Run npm run build
Warning: Large bundle size detected
Build completed successfully
```

### Output (Follow Mode)

```
❯ Fetching logs for workflow: deploy.yml

Workflow Run
────────────────────────────────────────────────────────────────────────────────
Name:       Deploy Application
Run #:      48
Status:     In Progress
Branch:     main
Commit:     b2c3d4e
Actor:      ignorant05
────────────────────────────────────────────────────────────────────────────────

Following logs (press Ctrl+C to stop)...

   Job: build
   Status: In Progress
────────────────────────────────────────────────────────────────────────────────
10:30:01 Run actions/checkout@v4
10:30:02 Syncing repository
10:30:03 Checkout complete
10:30:05 Run npm install
... (logs continue updating) ...

^C
Received interrupt signal...

Log streaming stopped
```

### Output (List Jobs)

```
Jobs for Run #47:

1. build
   Status: Completed
   Conclusion: Success

2. test
   Status: Completed
   Conclusion: Success

3. deploy
   Status: Completed
   Conclusion: Success

   To view logs for a specific job:
   uniflow logs --run-id 47 --job <job-name>
```

---
## 🎨 Output Colors

### Status Indicators

-  Queued
-  In Progress
-  Completed
-  Waiting

### Conclusion Indicators

-  Success (Green)
-  Failure (Red)
-  Cancelled (Yellow)
-  Skipped
-  Timed Out

### Log Levels

- **Red**: Errors, failures, fatal
- **Yellow**: Warnings
- **Green**: Success, passed, completed
- **Gray**: Debug information
- **White**: Normal info

---
## 💡 Common Workflows

### Deploy to Production

```bash
# 1. Trigger deployment
uniflow trigger deploy.yml --input environment=production

# 2. Follow logs
uniflow logs deploy.yml --follow 

# 3. Check final status
uniflow status deploy.yml
```

### Debug Failed Run

```bash
# 1. Find failed run
uniflow status deploy.yml

# 2. List jobs to find which failed
uniflow logs --run-id 123456 --jobs

# 3. View logs of failed job
uniflow logs --run-id 123456 --job "test" --tail 100
```

### Monitor Active Run

```bash
# Stream logs in real-time
uniflow logs deploy.yml --follow --verbose
```

### Save Logs to File

```bash
# Without colors
uniflow logs deploy.yml --no-color > deployment.log

# With timestamps
uniflow logs deploy.yml --no-color > deployment.log
```

---
## 🔗 Related Documentation

- [GitHub Actions Guide](github-actions.md)

---

**Next:** [GitHub Actions Guide](github-actions.md) →
