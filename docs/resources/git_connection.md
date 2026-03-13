---
page_title: "labplatform_git_connection Resource - LabPlatform"
description: |-
  Manages a Git connection (GitHub or Gitea) for repository and branch browsing.
---

# labplatform_git_connection

Manages a Git connection for browsing repositories and branches when configuring courses.

## Example Usage

```hcl
resource "labplatform_git_connection" "github" {
  name          = "GitHub DesoTech"
  provider_name = "github"
  org_name      = "desotech-it"
  token         = var.github_token
}
```

## Schema

### Required

- `name` (String) — Display name.
- `provider_name` (String) — Git provider type: `github` or `gitea`.
- `org_name` (String) — Organization or username to list repos from.
- `token` (String, Sensitive) — API access token. Write-only.

### Optional

- `base_url` (String) — API base URL (for Gitea). Not needed for GitHub.

### Read-Only

- `id` (Number) — Connection ID.
