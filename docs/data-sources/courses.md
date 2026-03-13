---
page_title: "labplatform_courses Data Source - LabPlatform"
description: |-
  Reads all existing courses from the platform.
---

# labplatform_courses

Reads all existing courses from the LabPlatform.

## Example Usage

```hcl
data "labplatform_courses" "all" {}

# Find a course by name
locals {
  cka = [for c in data.labplatform_courses.all.courses : c if c.name == "CKA"][0]
}
```

## Schema

### Read-Only

- `courses` (List of Object) — List of courses. Each course has: `id`, `name`, `description`, `guide_repo`, `guide_branch`, `duration_days`.
