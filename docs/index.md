---
page_title: "Provider: LabPlatform"
description: |-
  The LabPlatform provider manages the DesoLabs training platform as code — users, courses, sessions, connection templates, and student assignments.
---

# LabPlatform Provider

The LabPlatform provider allows you to manage the entire [DesoLabs LabPlatform](https://labplatform.desolabs.it) declaratively using Terraform. You can create and manage users, courses, sessions (classes), connection templates, Git connections, vSphere endpoints, and student assignments.

## Authentication

The provider authenticates with the LabPlatform REST API using an admin username and password.

### Environment Variables (recommended)

```bash
export LABPLATFORM_URL="https://labplatform.example.com"
export LABPLATFORM_USERNAME="admin"
export LABPLATFORM_PASSWORD="your-password"
```

```hcl
provider "labplatform" {}
```

### Inline Configuration

```hcl
provider "labplatform" {
  url      = "https://labplatform.example.com"
  username = "admin"
  password = var.admin_password
}
```

> **Never hardcode credentials.** Use environment variables or `terraform.tfvars` (added to `.gitignore`).

## Example Usage

```hcl
terraform {
  required_providers {
    labplatform = {
      source  = "desotech-it/labplatform"
      version = "~> 0.1"
    }
  }
}

provider "labplatform" {}

resource "labplatform_user" "student" {
  username   = "mario.rossi"
  password   = "SecurePass123!"
  role       = "student"
  first_name = "Mario"
  last_name  = "Rossi"
  email      = "mario@example.com"
}
```

## Schema

### Optional

- `url` (String) — Base URL of the LabPlatform instance. Can also be set with the `LABPLATFORM_URL` environment variable.
- `username` (String) — Admin username. Can also be set with the `LABPLATFORM_USERNAME` environment variable.
- `password` (String, Sensitive) — Admin password. Can also be set with the `LABPLATFORM_PASSWORD` environment variable.
