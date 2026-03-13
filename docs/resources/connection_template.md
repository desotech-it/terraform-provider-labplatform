---
page_title: "labplatform_connection_template Resource - LabPlatform"
description: |-
  Manages a connection template (VNC, RDP, SSH, or vSphere) that defines how students connect to lab machines.
---

# labplatform_connection_template

Manages a connection template in the LabPlatform. Templates define how students connect to lab machines via VNC, RDP, SSH, or vSphere console.

## Example Usage

### VNC Desktop

```hcl
resource "labplatform_connection_template" "vnc" {
  name     = "Linux Desktop"
  protocol = "vnc"
  hostname = "desktop-xfce-1"
  port     = 5901
  password = "user"
}
```

### SSH Terminal

```hcl
resource "labplatform_connection_template" "ssh" {
  name     = "SSH Terminal"
  protocol = "ssh"
  hostname = "ssh-server-1"
  port     = 2222
  username = "student"
  password = "student"
}
```

### RDP Desktop

```hcl
resource "labplatform_connection_template" "rdp" {
  name     = "Windows Desktop"
  protocol = "rdp"
  hostname = "rdp-desktop-1"
  port     = 3389
}
```

### vSphere VM

```hcl
resource "labplatform_connection_template" "vsphere" {
  name                = "Windows Server VM"
  protocol            = "vsphere"
  hostname            = "/Datacenter/vm/Labs/win-01"
  vsphere_endpoint_id = labplatform_vsphere_endpoint.vcenter.id
}
```

## Schema

### Required

- `name` (String) — Template name.
- `protocol` (String) — Connection protocol: `vnc`, `rdp`, `ssh`, or `vsphere`.

### Optional

- `hostname` (String) — Target hostname, IP, or VM path (for vSphere).
- `port` (Number) — Target port (e.g. 5901 for VNC, 3389 for RDP, 22 for SSH).
- `username` (String) — Connection username.
- `password` (String, Sensitive) — Connection password. Write-only.
- `parameters` (String) — Additional connection parameters as JSON string.
- `vsphere_endpoint_id` (Number) — Required for `vsphere` protocol.
- `course_id` (Number) — Associate this template with a specific course.
- `guest_id` (String) — vSphere guest OS ID.

### Read-Only

- `id` (Number) — Template ID.
