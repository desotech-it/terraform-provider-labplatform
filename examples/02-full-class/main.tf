# =============================================================================
# Esempio 02: Classe completa con tutte le dipendenze
# =============================================================================
# Crea TUTTO il necessario per una classe operativa:
#   1. Connessione Git (per browsare le guide)
#   2. Corso con guida e branch
#   3. Template di connessione (VNC + SSH)
#   4. Trainer
#   5. Studenti (da variabile)
#   6. Classe con 5 giorni
#   7. Assegnazione studenti con connessioni
#
# Uso:
#   cp terraform.tfvars.example terraform.tfvars
#   # Modifica terraform.tfvars con le tue credenziali
#   terraform init
#   terraform apply
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

# --- Variabili ---

variable "github_token" {
  description = "GitHub personal access token per accesso alle guide"
  type        = string
  sensitive   = true
}

variable "trainer_password" {
  description = "Password per il trainer"
  type        = string
  sensitive   = true
  default     = "Trainer2026!"
}

variable "student_password" {
  description = "Password per gli studenti"
  type        = string
  sensitive   = true
  default     = "Student2026!"
}

variable "course_start_date" {
  description = "Data di inizio corso (YYYY-MM-DD)"
  type        = string
  default     = "2026-03-23"
}

variable "students" {
  description = "Lista degli studenti"
  type = list(object({
    username   = string
    first_name = string
    last_name  = string
    email      = string
    company    = optional(string, "")
  }))
  default = [
    {
      username   = "mario.rossi"
      first_name = "Mario"
      last_name  = "Rossi"
      email      = "mario.rossi@example.com"
      company    = "Acme S.r.l."
    },
    {
      username   = "laura.bianchi"
      first_name = "Laura"
      last_name  = "Bianchi"
      email      = "laura.bianchi@example.com"
      company    = "Acme S.r.l."
    },
    {
      username   = "giovanni.verdi"
      first_name = "Giovanni"
      last_name  = "Verdi"
      email      = "giovanni.verdi@example.com"
      company    = "TechCorp"
    },
  ]
}

# --- Locals: calcolo date ---
# Genera automaticamente 5 giorni lavorativi a partire dalla data di inizio

locals {
  # 5 giorni consecutivi (lun-ven)
  course_days = [
    { date = var.course_start_date,                                             start_time = "09:00", end_time = "18:00" },
    { date = timeadd("${var.course_start_date}T00:00:00Z", "24h"),              start_time = "09:00", end_time = "18:00" },
    { date = timeadd("${var.course_start_date}T00:00:00Z", "48h"),              start_time = "09:00", end_time = "18:00" },
    { date = timeadd("${var.course_start_date}T00:00:00Z", "72h"),              start_time = "09:00", end_time = "18:00" },
    { date = timeadd("${var.course_start_date}T00:00:00Z", "96h"),              start_time = "09:00", end_time = "13:00" },
  ]

  # Formatta le date (rimuove la parte orario se presente)
  days = [for d in local.course_days : {
    date       = substr(tostring(d.date), 0, 10)
    start_time = d.start_time
    end_time   = d.end_time
  }]
}

# =============================================================================
# RISORSE
# =============================================================================

# 1. Connessione Git
resource "labplatform_git_connection" "github" {
  name     = "GitHub DesoTech"
  provider = "github"
  org_name = "desotech-it"
  token    = var.github_token
}

# 2. Corso
resource "labplatform_course" "cka" {
  name              = "CKA - Certified Kubernetes Admin"
  description       = "Preparazione alla certificazione CKA"
  guide_repo        = "desotech-it/CKA"
  guide_branch      = "v1.32"
  duration_days     = 5
  git_connection_id = labplatform_git_connection.github.id
}

# 3. Template di connessione
resource "labplatform_connection_template" "vnc" {
  name     = "Linux VNC Desktop"
  protocol = "vnc"
  hostname = "desktop-xfce-1"
  port     = 5901
  password = "user"
}

resource "labplatform_connection_template" "ssh" {
  name     = "SSH Terminal"
  protocol = "ssh"
  hostname = "ssh-server-1"
  port     = 2222
  username = "student"
  password = "student"
}

# 4. Trainer
resource "labplatform_user" "trainer" {
  username   = "trainer.cka"
  password   = var.trainer_password
  role       = "trainer"
  first_name = "Marco"
  last_name  = "Neri"
  email      = "marco.neri@desotech.it"
}

# 5. Studenti
resource "labplatform_user" "students" {
  for_each = { for s in var.students : s.username => s }

  username   = each.value.username
  password   = var.student_password
  role       = "student"
  first_name = each.value.first_name
  last_name  = each.value.last_name
  email      = each.value.email
  company    = each.value.company
}

# 6. Classe
resource "labplatform_class" "this" {
  course_id   = labplatform_course.cka.id
  trainer_ids = [labplatform_user.trainer.id]
  status      = "scheduled"
  notes       = "CKA - Settimana del ${var.course_start_date}"

  dynamic "days" {
    for_each = local.days
    content {
      date       = days.value.date
      start_time = days.value.start_time
      end_time   = days.value.end_time
    }
  }
}

# 7. Assegnazione studenti
resource "labplatform_class_student" "this" {
  for_each = labplatform_user.students

  class_id = labplatform_class.this.id
  user_id    = each.value.id
  template_ids = [
    labplatform_connection_template.vnc.id,
    labplatform_connection_template.ssh.id,
  ]
}

# --- Output ---

output "summary" {
  value = {
    course     = labplatform_course.cka.name
    class_id = labplatform_class.this.id
    trainer    = labplatform_user.trainer.username
    students   = [for s in labplatform_user.students : s.username]
    dates      = [for d in local.days : d.date]
  }
}
