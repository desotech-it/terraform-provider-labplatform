---
page_title: "labplatform_class_student Resource - LabPlatform"
description: |-
  Assigns a student to a class with connection templates. Creates the lab and remote connections.
---

# labplatform_class_student

Assigns a student to a class with specific connection templates. When created, the lab environment and remote connections are provisioned automatically. When destroyed, the student is removed and connections are cleaned up.

## Example Usage

```hcl
resource "labplatform_class_student" "mario" {
  class_id = labplatform_class.cka_w13.id
  user_id  = labplatform_user.mario.id
  template_ids = [
    labplatform_connection_template.vnc.id,
    labplatform_connection_template.ssh.id,
  ]
}
```

### Assign multiple students with `for_each`

```hcl
resource "labplatform_class_student" "this" {
  for_each = labplatform_user.students

  class_id     = labplatform_class.cka_w13.id
  user_id      = each.value.id
  template_ids = [
    labplatform_connection_template.vnc.id,
    labplatform_connection_template.ssh.id,
  ]
}
```

## Schema

### Required

- `class_id` (Number) — Class ID.
- `user_id` (Number) — Student user ID.

### Optional

- `template_ids` (List of Number) — Connection template IDs to provision for the student.
