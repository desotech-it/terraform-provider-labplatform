---
page_title: "labplatform_course Resource - LabPlatform"
description: |-
  Manages a course with guide repository and duration.
---

# labplatform_course

Manages a course in the LabPlatform. Courses define training content with an optional guide repository and duration.

## Example Usage

```hcl
resource "labplatform_course" "cka" {
  name              = "CKA - Certified Kubernetes Admin"
  description       = "CKA certification preparation"
  guide_repo        = "desotech-it/CKA"
  guide_branch      = "v1.32"
  duration_days     = 5
  git_connection_id = labplatform_git_connection.github.id
}
```

## Schema

### Required

- `name` (String) — Course name.

### Optional

- `description` (String) — Course description.
- `guide_repo` (String) — Guide repository in `org/repo` format.
- `guide_branch` (String) — Guide branch. Default: `main`.
- `duration_days` (Number) — Course duration in days. Default: `5`.
- `git_connection_id` (Number) — Reference to a `labplatform_git_connection` for repository browsing.

### Read-Only

- `id` (Number) — Course ID.

## Import

```bash
terraform import labplatform_course.cka 5
```
