# =============================================================================
# Esempio 04: Laboratorio con VM vSphere
# =============================================================================
# Crea un laboratorio che usa VM su vSphere (vCenter) per i desktop
# invece dei container Kubernetes. Utile per corsi che richiedono
# Windows o ambienti specifici.
# =============================================================================

terraform {
  required_providers {
    labplatform = {
      source  = "desotech-it/labplatform"
      version = "~> 0.1"
    }
  }
}

provider "labplatform" {}

variable "vcenter_password" {
  description = "Password vCenter"
  type        = string
  sensitive   = true
}

variable "student_password" {
  type      = string
  sensitive = true
  default   = "Student2026!"
}

# --- Endpoint vSphere ---

resource "labplatform_vsphere_endpoint" "lab" {
  name       = "vCenter Lab"
  url        = "https://vcenter.desotech.io"
  username   = "administrator@vsphere.local"
  password   = var.vcenter_password
  datacenter = "Datacenter-Lab"
  insecure   = true
}

# --- Corso Windows ---

resource "labplatform_course" "windows" {
  name          = "WIN201 - Windows Server Administration"
  description   = "Amministrazione Windows Server 2022"
  duration_days = 3
}

# --- Template: VM Windows via vSphere ---

resource "labplatform_connection_template" "win_vm" {
  name                = "Windows Server Lab VM"
  protocol            = "vsphere"
  hostname            = "/Datacenter-Lab/vm/Labs/win-server-template"
  vsphere_endpoint_id = labplatform_vsphere_endpoint.lab.id
}

# --- Template: desktop Linux VNC per strumenti ---

resource "labplatform_connection_template" "linux_tools" {
  name     = "Linux Tools Desktop"
  protocol = "vnc"
  hostname = "desktop-xfce-2"
  port     = 5901
  password = "user"
}

# --- Trainer ---

resource "labplatform_user" "trainer_win" {
  username   = "trainer.windows"
  password   = "Trainer2026!"
  role       = "trainer"
  first_name = "Andrea"
  last_name  = "Verdi"
  email      = "andrea.verdi@desotech.it"
}

# --- Studenti ---

resource "labplatform_user" "students" {
  for_each = {
    "student.win01" = { first_name = "Marco",    last_name = "Bianchi" }
    "student.win02" = { first_name = "Giulia",   last_name = "Neri" }
    "student.win03" = { first_name = "Federico", last_name = "Costa" }
  }

  username   = each.key
  password   = var.student_password
  role       = "student"
  first_name = each.value.first_name
  last_name  = each.value.last_name
}

# --- Classe ---

resource "labplatform_class" "windows" {
  course_id   = labplatform_course.windows.id
  trainer_ids = [labplatform_user.trainer_win.id]
  status      = "scheduled"
  notes       = "Windows Server - Lab con VM vSphere"

  days = [
    { date = "2026-04-06", start_time = "09:00", end_time = "18:00" },
    { date = "2026-04-07", start_time = "09:00", end_time = "18:00" },
    { date = "2026-04-08", start_time = "09:00", end_time = "13:00" },
  ]
}

# --- Assegnazioni: ogni studente riceve VM Windows + desktop Linux ---

resource "labplatform_class_student" "this" {
  for_each = labplatform_user.students

  class_id = labplatform_class.windows.id
  user_id    = each.value.id
  template_ids = [
    labplatform_connection_template.win_vm.id,
    labplatform_connection_template.linux_tools.id,
  ]
}
