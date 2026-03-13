---
page_title: "labplatform_user Resource - LabPlatform"
description: |-
  Manages a user (student, trainer, or admin) in the LabPlatform.
---

# labplatform_user

Manages a user in the LabPlatform. Users can be students, trainers, or admins.

## Example Usage

```hcl
resource "labplatform_user" "student" {
  username   = "mario.rossi"
  password   = var.student_password
  role       = "student"
  first_name = "Mario"
  last_name  = "Rossi"
  email      = "mario@example.com"
  company    = "Acme S.r.l."
  phone      = "+39 333 1234567"
  language   = "it"
}
```

### Multiple students with `for_each`

```hcl
resource "labplatform_user" "students" {
  for_each = {
    "anna.ferrari"   = { first_name = "Anna",  last_name = "Ferrari" }
    "paolo.romano"   = { first_name = "Paolo", last_name = "Romano" }
  }

  username   = each.key
  password   = var.student_password
  role       = "student"
  first_name = each.value.first_name
  last_name  = each.value.last_name
}
```

## Schema

### Required

- `username` (String) — Unique username.
- `password` (String, Sensitive) — Password. Write-only, not returned by the API.
- `role` (String) — User role: `student`, `trainer`, or `admin`.

### Optional

- `email` (String) — Email address.
- `first_name` (String) — First name.
- `last_name` (String) — Last name.
- `company` (String) — Company name.
- `phone` (String) — Phone number.
- `language` (String) — Language code. Default: `it`.

### Read-Only

- `id` (Number) — User ID.

## Import

Users can be imported using their ID:

```bash
terraform import labplatform_user.student 42
```
