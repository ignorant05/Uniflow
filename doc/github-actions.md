
# GitHub Actions Guide

Complete guide for using Uniflow with GitHub Actions.

---
## 🎯 Overview

Uniflow provides seamless integration with GitHub Actions, allowing you to:
- Trigger workflows remotely from the command line
- Monitor workflow execution in real-time
- Stream logs with colored output
- Check workflow status and history

---
## 📋 Prerequisites

- GitHub account
- Repository with GitHub Actions workflows
- Personal access token
- Uniflow installed and configured

---
## 🔑 Setup GitHub Token

### Step 1: Generate Personal Access Token

1. Go to https://github.com/settings/tokens/new
2. Token name: "Uniflow CLI" (or whatever)
3. Expiration: 90 days (or as needed)
4. Select scopes:
   - ✅ `repo` - Full control of repositories
   - ✅ `workflow` - Update GitHub Action workflows
5. Click "Generate token"
6. **Copy the token** (starts with `ghp_`)

### Step 2: Set Environment Variable

```bash
# Add to ~/.bashrc or ~/.zshrc
export GITHUB_TOKEN="ghp_your_token_here"

# Reload shell
source ~/.bashrc
```

### Step 3: Configure Uniflow

```bash
# Set default repository
uniflow config set profiles.default.github.default_repo "owner/repo"

# Verify configuration
uniflow config list
uniflow config validate
```

---
## 📝 Creating Triggerable Workflows

### Required: `workflow_dispatch` Trigger

Workflows must have `workflow_dispatch` to be triggered remotely:

```yaml
# .github/workflows/deploy.yml
name: Deploy Application

on:
  workflow_dispatch:  # THIS IS REQUIRED!
    inputs:
      environment:
        description: 'Deployment environment'
        required: true
        type: choice
        options:
          - dev
          - staging
          - production
      
      version:
        description: 'Version to deploy'
        required: false
        type: string
      
      dry_run:
        description: 'Run in dry-run mode'
        required: false
        type: boolean
        default: false

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      
      - name: Deploy
        run: |
          echo "Deploying to ${{ inputs.environment }}"
          echo "Version: ${{ inputs.version }}"
          echo "Dry run: ${{ inputs.dry_run }}"
```

### Input Types

#### String Input

```yaml
inputs:
  version:
    description: 'Version number'
    required: false
    type: string
    default: 'latest'
```

#### Choice Input

```yaml
inputs:
  environment:
    description: 'Environment'
    required: true
    type: choice
    options:
      - dev
      - staging
      - production
```

#### Boolean Input

```yaml
inputs:
  debug:
    description: 'Enable debug mode'
    required: false
    type: boolean
    default: false
```

### Complete Example

```yaml
name: Full CI/CD Pipeline

on:
  workflow_dispatch:
    inputs:
      environment:
        type: choice
        options: [dev, staging, prod]
        default: dev
      version:
        type: string
        required: false
      run_tests:
        type: boolean
        default: true
      notify_slack:
        type: boolean
        default: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build
        run: npm run build
  
  test:
    needs: build
    if: ${{ inputs.run_tests }}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Test
        run: npm test
  
  deploy:
    needs: [build, test]
    if: always()
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to ${{ inputs.environment }}
        run: |
          echo "Deploying version ${{ inputs.version || 'latest' }}"
          echo "Environment: ${{ inputs.environment }}"
  
  notify:
    needs: deploy
    if: ${{ inputs.notify_slack }}
    runs-on: ubuntu-latest
    steps:
      - name: Notify Slack
        run: echo "Sending Slack notification"
```

---
## 🚀 Using Uniflow with GitHub Actions

### List Available Workflows

```bash
# See all workflows in repository
uniflow workflows

# Output:
# 1. Deploy Application
#    File: .github/workflows/deploy.yml
#    State: active
```

### Trigger Workflows

#### Basic Trigger

```bash
uniflow trigger deploy.yml
```

#### With Branch

```bash
uniflow trigger deploy.yml --branch develop
```

#### With Inputs

```bash
# Single input
uniflow trigger deploy.yml --input environment=production

# Multiple inputs
uniflow trigger deploy.yml \
  --input environment=staging \
  --input version=v2.1.0 \
  --input dry_run=false
```

#### Complete Example

```bash
uniflow trigger deploy.yml \
  --branch main \
  --input environment=production \
  --input version=v3.0.0 \
  --input run_tests=true \
  --input notify_slack=true \
  --verbose
```

### Monitor Workflow Status

```bash
# Check latest run
uniflow status deploy.yml

# Check recent runs
uniflow status deploy.yml --limit 10

# Show all runs
uniflow status deploy.yml --all
```

### View Logs

#### Basic Logs

```bash
# Latest run
uniflow logs deploy.yml

# Specific run
uniflow logs --run-id 123456
```

#### Real-Time Streaming

```bash
# Follow logs as they happen
uniflow logs deploy.yml --follow

# Last 100 lines only
uniflow logs deploy.yml --follow --tail 100
```

#### Job-Specific Logs

```bash
# View specific job
uniflow logs deploy.yml --job "build"

# Follow specific job
uniflow logs deploy.yml --job "deploy" --follow
```

---
## 🎨 Workflow Examples

### Example 1: Simple Deployment

```yaml
# .github/workflows/deploy.yml
name: Deploy

on:
  workflow_dispatch:
    inputs:
      environment:
        type: choice
        options: [dev, prod]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Deploy
        run: ./deploy.sh ${{ inputs.environment }}
```

**Trigger:**
```bash
uniflow trigger deploy.yml --input environment=prod
```

### Example 2: Docker Build & Push

```yaml
# .github/workflows/docker.yml
name: Build Docker Image

on:
  workflow_dispatch:
    inputs:
      tag:
        type: string
        required: true
      push:
        type: boolean
        default: false

jobs:
  docker:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Build
        run: docker build -t myapp:${{ inputs.tag }} .
      
      - name: Push
        if: ${{ inputs.push }}
        run: docker push myapp:${{ inputs.tag }}
```

**Trigger:**
```bash
uniflow trigger docker.yml \
  --input tag=v1.2.3 \
  --input push=true
```

### Example 3: Multi-Stage Pipeline

```yaml
# .github/workflows/pipeline.yml
name: CI/CD Pipeline

on:
  workflow_dispatch:
    inputs:
      stage:
        type: choice
        options: [all, build, test, deploy]
        default: all

jobs:
  build:
    if: ${{ inputs.stage == 'all' || inputs.stage == 'build' }}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: npm run build
  
  test:
    needs: build
    if: ${{ inputs.stage == 'all' || inputs.stage == 'test' }}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: npm test
  
  deploy:
    needs: test
    if: ${{ inputs.stage == 'all' || inputs.stage == 'deploy' }}
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - run: ./deploy.sh
```

**Trigger:**
```bash
# Run all stages
uniflow trigger pipeline.yml --input stage=all

# Run only tests
uniflow trigger pipeline.yml --input stage=test
```

---
## 🔧 Advanced Usage

### Using Different Profiles

```bash
# Development
uniflow trigger deploy.yml --profile dev --input env=dev

# Production
uniflow trigger deploy.yml --profile prod --input env=prod
```

### Automated Monitoring Script

```bash
#!/bin/bash
# monitor-deployment.sh

WORKFLOW="deploy.yml"

# Trigger workflow
echo "Triggering deployment..."
uniflow trigger $WORKFLOW --input environment=production

# Wait for start
sleep 5

# Stream logs
echo "Streaming logs..."
uniflow logs $WORKFLOW --follow 

# Check final status
echo "Checking final status..."
uniflow status $WORKFLOW --limit 1
```

### CI/CD Integration

```bash
#!/bin/bash
# ci-deploy.sh

set -e

# Trigger deployment
uniflow trigger deploy.yml \
  --input environment=production \
  --input version=${CI_COMMIT_TAG}

# Follow logs and capture output
uniflow logs deploy.yml --follow --no-color > deployment.log

# Check if successful
if uniflow status deploy.yml | grep -q "✅ Success"; then
  echo "Deployment successful!"
  exit 0
else
  echo "Deployment failed!"
  cat deployment.log
  exit 1
fi
```

---
## 🐛 Troubleshooting

### Issue: "workflow does not have 'workflow_dispatch' trigger"

**Solution:** Add `workflow_dispatch:` to your workflow file:

```yaml
on:
  workflow_dispatch:  # Add this
  push:
    branches: [main]
```

### Issue: "workflow not found"

**Check available workflows:**
```bash
uniflow workflows
```

**Verify file exists:**
```bash
# Check GitHub
https://github.com/owner/repo/tree/main/.github/workflows
```

### Issue: "403 rate limit exceeded"

**Solution:** Wait or use multiple tokens:

```bash
# Check rate limit
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/rate_limit

# Use different token
export GITHUB_TOKEN="different_token"
```

### Issue: Logs not showing

**Common causes:**
- Workflow hasn't started yet
- Logs expired (90 days)
- Permissions issue

**Debug:**
```bash
# Check workflow status
uniflow status deploy.yml

# Try specific run ID
uniflow logs --run-id 123456

# Check with verbose
uniflow logs deploy.yml --verbose
```

---
## 📚 Best Practices

### 1. Organize Workflows

```
.github/
  workflows/
    deploy-dev.yml
    deploy-staging.yml
    deploy-prod.yml
    test.yml
    build.yml
```

### 2. Use Descriptive Names

```yaml
name: Deploy to Production  # ✅ Clear
# vs
name: Deploy  # ❌ Vague
```

### 3. Document Inputs

```yaml
inputs:
  version:
    description: 'Semantic version (e.g., v1.2.3)'  # ✅ Helpful
    required: true
```

### 4. Set Sensible Defaults

```yaml
inputs:
  environment:
    default: 'dev'  # ✅ Safe default
```

### 5. Add Validation

```yaml
jobs:
  validate:
    runs-on: ubuntu-latest
    steps:
      - name: Validate version format
        run: |
          if [[ ! "${{ inputs.version }}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            echo "Invalid version format"
            exit 1
          fi
```

---
## 🔗 Related Documentation

- [Commands Reference](commands.md)
- [GitHub Actions Documentation](https://docs.github.com/en/actions)

---

**Ready to use Uniflow with GitHub Actions!** 🚀
