---
page_title: "labplatform_users Data Source - LabPlatform"
description: |-
  Reads existing users from the platform, optionally filtered by role.
---

# labplatform_users

Reads existing users from the LabPlatform. Can be filtered by role.

## Example Usage

```hcl
data "labplatform_users" "trainers" {
  role = "trainer"
}

# Find a specific trainer by username
locals {
  trainer = [for t in data.labplatform_users.trainers.users : t if t.username == "trainer.cka"][0]
}
```

## Schema

### Optional

- `role` (String) — Filter users by role: `student`, `trainer`, or `admin`. If omitted, returns all users.

### Read-Only

- `users` (List of Object) — List of users. Each user has: `id`, `username`, `email`, `role`, `first_name`, `last_name`.
