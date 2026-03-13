---
page_title: "labplatform_session_student Resource - LabPlatform"
description: |-
  Assigns a student to a session with connection templates. Creates the lab and remote connections.
---

# labplatform_session_student

Assigns a student to a session with specific connection templates. When created, the lab environment and remote connections are provisioned automatically. When destroyed, the student is removed and connections are cleaned up.

## Example Usage

```hcl
resource "labplatform_session_student" "mario" {
  session_id = labplatform_session.cka_w13.id
  user_id    = labplatform_user.mario.id
  template_ids = [
    labplatform_connection_template.vnc.id,
    labplatform_connection_template.ssh.id,
  ]
}
```

### Assign multiple students with `for_each`

```hcl
resource "labplatform_session_student" "this" {
  for_each = labplatform_user.students

  session_id   = labplatform_session.cka_w13.id
  user_id      = each.value.id
  template_ids = [
    labplatform_connection_template.vnc.id,
    labplatform_connection_template.ssh.id,
  ]
}
```

## Schema

### Required

- `session_id` (Number) — Session ID.
- `user_id` (Number) — Student user ID.

### Optional

- `template_ids` (List of Number) — Connection template IDs to provision for the student.
