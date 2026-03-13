# =============================================================================
# Esempio 03: Più classi nella stessa settimana con risorse condivise
# =============================================================================
# Scenario: due corsi diversi nella stessa settimana, con trainer diversi
# ma che condividono gli stessi template di connessione.
#
# Usa data sources per leggere risorse già esistenti (corsi, template, trainer).
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

variable "student_password" {
  type      = string
  sensitive = true
  default   = "Student2026!"
}

# --- Data sources: risorse già esistenti ---

# Legge tutti i corsi dalla piattaforma
data "labplatform_courses" "all" {}

# Legge tutti i template di connessione
data "labplatform_templates" "all" {}

# Legge tutti i trainer
data "labplatform_users" "trainers" {
  role = "trainer"
}

# --- Locals: trova risorse per nome ---

locals {
  # Trova corso CKA dall'elenco
  cka_course = [for c in data.labplatform_courses.all.courses : c if c.name == "CKA - Certified Kubernetes Admin"][0]

  # Trova corso Docker dall'elenco
  docker_course = [for c in data.labplatform_courses.all.courses : c if c.name == "DOC101 - Docker Fundamentals"][0]

  # Trova template VNC e SSH
  vnc_template = [for t in data.labplatform_templates.all.templates : t if t.protocol == "vnc"][0]
  ssh_template = [for t in data.labplatform_templates.all.templates : t if t.protocol == "ssh"][0]

  # Trova trainer per username
  trainer_marco = [for t in data.labplatform_users.trainers.users : t if t.username == "trainer.cka"][0]
  trainer_luca  = [for t in data.labplatform_users.trainers.users : t if t.username == "trainer02"][0]
}

# --- Studenti per il corso CKA ---

resource "labplatform_user" "cka_students" {
  for_each = {
    "anna.ferrari"   = { first_name = "Anna",   last_name = "Ferrari",  email = "anna@example.com" }
    "paolo.romano"   = { first_name = "Paolo",  last_name = "Romano",   email = "paolo@example.com" }
    "chiara.colombo" = { first_name = "Chiara", last_name = "Colombo",  email = "chiara@example.com" }
  }

  username   = each.key
  password   = var.student_password
  role       = "student"
  first_name = each.value.first_name
  last_name  = each.value.last_name
  email      = each.value.email
}

# --- Studenti per il corso Docker ---

resource "labplatform_user" "docker_students" {
  for_each = {
    "luca.moretti"    = { first_name = "Luca",    last_name = "Moretti",   email = "luca@example.com" }
    "elena.ricci"     = { first_name = "Elena",   last_name = "Ricci",     email = "elena@example.com" }
    "alessio.galli"   = { first_name = "Alessio", last_name = "Galli",     email = "alessio@example.com" }
    "sara.martinelli" = { first_name = "Sara",    last_name = "Martinelli", email = "sara@example.com" }
  }

  username   = each.key
  password   = var.student_password
  role       = "student"
  first_name = each.value.first_name
  last_name  = each.value.last_name
  email      = each.value.email
}

# --- Classe CKA: lunedì-venerdì ---

resource "labplatform_session" "cka" {
  course_id   = local.cka_course.id
  trainer_ids = [local.trainer_marco.id]
  status      = "scheduled"
  notes       = "CKA - Settimana 13"

  days = [
    { date = "2026-03-23", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-24", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-25", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-26", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-27", start_time = "09:00", end_time = "13:00" },
  ]
}

resource "labplatform_session_student" "cka" {
  for_each = labplatform_user.cka_students

  session_id   = labplatform_session.cka.id
  user_id      = each.value.id
  template_ids = [local.vnc_template.id, local.ssh_template.id]
}

# --- Classe Docker: stessa settimana, trainer diverso ---

resource "labplatform_session" "docker" {
  course_id   = local.docker_course.id
  trainer_ids = [local.trainer_luca.id]
  status      = "scheduled"
  notes       = "Docker Fundamentals - Settimana 13"

  days = [
    { date = "2026-03-23", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-24", start_time = "09:00", end_time = "18:00" },
    { date = "2026-03-25", start_time = "09:00", end_time = "13:00" },
  ]
}

resource "labplatform_session_student" "docker" {
  for_each = labplatform_user.docker_students

  session_id   = labplatform_session.docker.id
  user_id      = each.value.id
  template_ids = [local.vnc_template.id, local.ssh_template.id]
}

# --- Output ---

output "classes" {
  value = {
    cka = {
      session_id = labplatform_session.cka.id
      students   = [for s in labplatform_user.cka_students : s.username]
    }
    docker = {
      session_id = labplatform_session.docker.id
      students   = [for s in labplatform_user.docker_students : s.username]
    }
  }
}
