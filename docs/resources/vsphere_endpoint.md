---
page_title: "labplatform_vsphere_endpoint Resource - LabPlatform"
description: |-
  Manages a vSphere (vCenter) endpoint for VM console access.
---

# labplatform_vsphere_endpoint

Manages a vSphere vCenter endpoint for VM console access in lab environments.

## Example Usage

```hcl
resource "labplatform_vsphere_endpoint" "vcenter" {
  name       = "vCenter Lab"
  url        = "https://vcenter.example.com"
  username   = "administrator@vsphere.local"
  password   = var.vcenter_password
  datacenter = "Datacenter-Lab"
  insecure   = true
}
```

## Schema

### Required

- `name` (String) — Display name.
- `url` (String) — vCenter URL.
- `username` (String) — vCenter username.
- `password` (String, Sensitive) — vCenter password. Write-only.
- `datacenter` (String) — Datacenter name.

### Optional

- `insecure` (Boolean) — Skip TLS verification. Default: `false`.

### Read-Only

- `id` (Number) — Endpoint ID.
