github-apps-substrate
===

This repository is a reference implementation of GitHub Apps.  
The aim is to enable GitOps in CI/CD.

## Required Perimssions

### Repository permissions

- Contents: Read & Write
- Metadata: Read-only
- Pull requests: Read & Write
- Commit statuses: Read & Write

### Subscribe to events

- Push
- Status

## develop

You will need to specify the following environment variables.

- GITHUB_APP_ID
- GITHUB_WEBHOOK_SECRET
- GITHUB_PRIVATE_KEY
- GITHUB_ENTERPRISE_URL (optional)
- GITHUB_ENTERPRISE_UPLOAD_URL (optional)

## change package name

You can change the package name with the following command.

```fish
⟩ env NAME=example.com/github/apps make rename
```
