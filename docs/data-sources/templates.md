---
page_title: "labplatform_templates Data Source - LabPlatform"
description: |-
  Reads all existing connection templates from the platform.
---

# labplatform_templates

Reads all existing connection templates from the LabPlatform.

## Example Usage

```hcl
data "labplatform_templates" "all" {}

# Find templates by protocol
locals {
  vnc_template = [for t in data.labplatform_templates.all.templates : t if t.protocol == "vnc"][0]
  ssh_template = [for t in data.labplatform_templates.all.templates : t if t.protocol == "ssh"][0]
}
```

## Schema

### Read-Only

- `templates` (List of Object) — List of connection templates. Each template has: `id`, `name`, `protocol`, `hostname`, `port`.
