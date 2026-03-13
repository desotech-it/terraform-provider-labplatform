# =============================================================================
# Esempio 01: Corso base con template di connessione
# =============================================================================
# Crea un corso e i relativi template di connessione.
# Questo è il primo passo: prima si creano i "mattoni" (corso + template),
# poi si usano negli esempi successivi per creare classi e assegnare studenti.
#
# Uso:
#   terraform init
#   terraform plan
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

provider "labplatform" {
  # Credenziali via env:
  #   LABPLATFORM_URL, LABPLATFORM_USERNAME, LABPLATFORM_PASSWORD
}

# --- Corso ---

resource "labplatform_course" "docker" {
  name          = "DOC101 - Docker Fundamentals"
  description   = "Corso introduttivo su Docker e container"
  guide_repo    = "desotech-it/DOC101"
  guide_branch  = "v2.1"
  duration_days = 3
}

# --- Template di connessione ---
# I template definiscono COME gli studenti si connettono ai lab.
# Ogni studente assegnato riceverà una copia di queste connessioni.

resource "labplatform_connection_template" "vnc_desktop" {
  name     = "Linux Desktop (VNC)"
  protocol = "vnc"
  hostname = "desktop-xfce-1"   # Nome del service Kubernetes
  port     = 5901
  password = "user"
}

resource "labplatform_connection_template" "ssh_terminal" {
  name     = "Terminale SSH"
  protocol = "ssh"
  hostname = "ssh-server-1"
  port     = 2222
  username = "student"
  password = "student"
}

# --- Output ---

output "course_id" {
  description = "ID del corso creato"
  value       = labplatform_course.docker.id
}

output "template_ids" {
  description = "ID dei template di connessione"
  value = {
    vnc = labplatform_connection_template.vnc_desktop.id
    ssh = labplatform_connection_template.ssh_terminal.id
  }
}
