---
page_title: "labplatform_session Resource - LabPlatform"
description: |-
  Manages a session (class) — a scheduled instance of a course with dates, times, and trainers.
---

# labplatform_session

Manages a session (class) in the LabPlatform. A session is a scheduled instance of a course with specific dates, times, and assigned trainers.

## Example Usage

```hcl
resource "labplatform_session" "cka_w13" {
  course_id   = labplatform_course.cka.id
  trainer_ids = [labplatform_user.trainer.id]
  status      = "scheduled"
  notes       = "CKA - Week 13"

  days = [
    { date = "2026-03-23", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-24", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-25", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-26", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-27", start_time = "09:00", end_time = "13:00" },
  ]
}
```

## Schema

### Required

- `course_id` (Number) — Reference to the course.
- `days` (List of Object) — List of session days. See [days](#days) below.

### Optional

- `trainer_ids` (List of Number) — List of trainer user IDs.
- `status` (String) — Session status: `scheduled` (default), `active`, `completed`, `cancelled`.
- `notes` (String) — Free-text notes.

### Read-Only

- `id` (Number) — Session ID.

### days

Each day object has the following attributes:

| Attribute | Type | Required | Description |
|-----------|------|----------|-------------|
| `date` | String | Yes | Date in `YYYY-MM-DD` format. |
| `start_time` | String | No | Start time. Default: `09:00`. |
| `end_time` | String | No | End time. Default: `18:00`. |

## Import

```bash
terraform import labplatform_session.this 12
```
