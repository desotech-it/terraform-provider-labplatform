# Terraform Provider — DesoLabs LabPlatform

Provider Terraform per gestire l'intera piattaforma LabPlatform come codice. Permette di creare e gestire utenti, corsi, classi, template di connessione, connessioni Git, endpoint vSphere e assegnazione studenti in modo dichiarativo e riproducibile.

## Indice

- [Installazione](#installazione)
- [Configurazione](#configurazione)
- [Quick Start](#quick-start)
- [Risorse](#risorse)
- [Data Sources](#data-sources)
- [Esempi](#esempi)
- [Workflow del tecnico](#workflow-del-tecnico)
- [Pubblicazione su Terraform Registry](#pubblicazione-su-terraform-registry)
- [Troubleshooting](#troubleshooting)

---

## Installazione

### Da Terraform Registry (consigliato)

Una volta pubblicato il provider sul [Terraform Registry](https://registry.terraform.io), basta dichiararlo nel file `.tf`:

```hcl
terraform {
  required_providers {
    labplatform = {
      source  = "desotech-it/labplatform"
      version = "~> 0.1"
    }
  }
}
```

Poi eseguire:

```bash
terraform init
```

Terraform scarica automaticamente il provider per la tua piattaforma (Linux, macOS, Windows).

### Build locale (solo per sviluppo)

Se stai sviluppando il provider o non è ancora pubblicato sul registry:

```bash
cd terraform/
go mod tidy
make install
```

Questo compila il binario e lo copia in `~/.terraform.d/plugins/`.

---

## Configurazione

Il provider si autentica con le API REST di LabPlatform usando username e password di un account admin.

### Variabili d'ambiente (consigliato)

```bash
export LABPLATFORM_URL="https://labplatform.desolabs.it"
export LABPLATFORM_USERNAME="admin"
export LABPLATFORM_PASSWORD="la-tua-password"
```

```hcl
provider "labplatform" {}
```

### Configurazione inline (solo per test locali)

```hcl
provider "labplatform" {
  url      = "https://labplatform.desolabs.it"
  username = "admin"
  password = var.admin_password   # MAI hardcodare, usare variabili
}
```

> **Non committare mai credenziali.** Usa env var o `terraform.tfvars` (aggiunto al `.gitignore`).

| Parametro | Env var | Descrizione |
|-----------|---------|-------------|
| `url` | `LABPLATFORM_URL` | URL base della piattaforma |
| `username` | `LABPLATFORM_USERNAME` | Username admin |
| `password` | `LABPLATFORM_PASSWORD` | Password admin (sensitive) |

---

## Quick Start

Crea un file `main.tf`:

```hcl
terraform {
  required_providers {
    labplatform = {
      source  = "desotech-it/labplatform"
      version = "~> 0.1"
    }
  }
}

provider "labplatform" {}

# Crea uno studente
resource "labplatform_user" "mario" {
  username   = "mario.rossi"
  password   = "Student2026!"
  role       = "student"
  first_name = "Mario"
  last_name  = "Rossi"
  email      = "mario@example.com"
}
```

```bash
export LABPLATFORM_URL="https://labplatform.desolabs.it"
export LABPLATFORM_USERNAME="admin"
export LABPLATFORM_PASSWORD="xxx"

terraform init
terraform plan     # anteprima
terraform apply    # applica
terraform destroy  # cancella
```

---

## Risorse

### `labplatform_user`

Utente della piattaforma (studente, trainer o admin).

```hcl
resource "labplatform_user" "studente" {
  username   = "mario.rossi"
  password   = var.student_password   # write-only
  role       = "student"              # student | trainer | admin
  first_name = "Mario"
  last_name  = "Rossi"
  email      = "mario@example.com"
  company    = "Acme S.r.l."
  phone      = "+39 333 1234567"
  language   = "it"                   # default: it
}
```

| Campo | Tipo | Richiesto | Note |
|-------|------|:---------:|------|
| `username` | string | si | Univoco |
| `password` | string | si | Write-only, non letto dall'API |
| `role` | string | si | `student`, `trainer`, `admin` |
| `email` | string | | |
| `first_name` | string | | |
| `last_name` | string | | |
| `company` | string | | |
| `phone` | string | | |
| `language` | string | | Default: `it` |

---

### `labplatform_course`

Corso con guida e durata.

```hcl
resource "labplatform_course" "cka" {
  name              = "CKA - Certified Kubernetes Admin"
  description       = "Preparazione certificazione CKA"
  guide_repo        = "desotech-it/CKA"
  guide_branch      = "v1.32"
  duration_days     = 5
  git_connection_id = labplatform_git_connection.github.id
}
```

| Campo | Tipo | Richiesto | Note |
|-------|------|:---------:|------|
| `name` | string | si | |
| `description` | string | | |
| `guide_repo` | string | | Formato `org/repo` |
| `guide_branch` | string | | Default: `main` |
| `duration_days` | number | | Default: 5 |
| `git_connection_id` | number | | Riferimento a `labplatform_git_connection` |

---

### `labplatform_connection_template`

Template di connessione ai laboratori. Definisce come gli studenti si collegano.

```hcl
# VNC
resource "labplatform_connection_template" "vnc" {
  name     = "Linux Desktop"
  protocol = "vnc"
  hostname = "desktop-xfce-1"
  port     = 5901
  password = "user"
}

# SSH
resource "labplatform_connection_template" "ssh" {
  name     = "Terminale SSH"
  protocol = "ssh"
  hostname = "ssh-server-1"
  port     = 2222
  username = "student"
  password = "student"
}

# RDP
resource "labplatform_connection_template" "rdp" {
  name     = "Windows Desktop"
  protocol = "rdp"
  hostname = "rdp-desktop-1"
  port     = 3389
}

# vSphere VM
resource "labplatform_connection_template" "vsphere" {
  name                = "Windows Server VM"
  protocol            = "vsphere"
  hostname            = "/Datacenter/vm/Labs/win-01"
  vsphere_endpoint_id = labplatform_vsphere_endpoint.vcenter.id
}
```

| Campo | Tipo | Richiesto | Note |
|-------|------|:---------:|------|
| `name` | string | si | |
| `protocol` | string | si | `vnc`, `rdp`, `ssh`, `vsphere` |
| `hostname` | string | | IP/hostname o path VM (vsphere) |
| `port` | number | | 5901 (vnc), 3389 (rdp), 2222 (ssh) |
| `username` | string | | |
| `password` | string | | Sensitive |
| `parameters` | string | | JSON con parametri extra per la connessione |
| `vsphere_endpoint_id` | number | | Richiesto per protocol `vsphere` |
| `course_id` | number | | Associa al corso (opzionale) |

---

### `labplatform_class`

Classe — un'istanza schedulata di un corso con date, orari e trainer.

```hcl
resource "labplatform_class" "cka_w13" {
  course_id   = labplatform_course.cka.id
  trainer_ids = [labplatform_user.trainer.id]
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
```

| Campo | Tipo | Richiesto | Note |
|-------|------|:---------:|------|
| `course_id` | number | si | Riferimento al corso |
| `trainer_ids` | list(number) | | Lista ID trainer |
| `status` | string | | `scheduled` (default), `active`, `completed`, `cancelled` |
| `notes` | string | | Note libere |
| `days` | list(object) | si | Vedi sotto |

**Oggetto `days`:**

| Campo | Tipo | Richiesto | Note |
|-------|------|:---------:|------|
| `date` | string | si | `YYYY-MM-DD` |
| `start_time` | string | | Default: `09:00` |
| `end_time` | string | | Default: `18:00` |

---

### `labplatform_class_student`

Assegna uno studente a una classe con i template di connessione. Crea il lab e le connessioni remote.

```hcl
resource "labplatform_class_student" "mario" {
  class_id = labplatform_class.cka_w13.id
  user_id    = labplatform_user.mario.id
  template_ids = [
    labplatform_connection_template.vnc.id,
    labplatform_connection_template.ssh.id,
  ]
}
```

| Campo | Tipo | Richiesto | Note |
|-------|------|:---------:|------|
| `class_id` | number | si | |
| `user_id` | number | si | |
| `template_ids` | list(number) | | Template da creare per lo studente |

Al `destroy`, lo studente viene rimosso dalla classe e le connessioni remote cancellate.

---

### `labplatform_git_connection`

Connessione a GitHub o Gitea per browsing repository e branch.

```hcl
resource "labplatform_git_connection" "github" {
  name     = "GitHub DesoTech"
  provider_name = "github"
  org_name = "desotech-it"
  token    = var.github_token
}
```

---

### `labplatform_vsphere_endpoint`

Endpoint vCenter per accesso console VM.

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

---

## Data Sources

Leggono risorse esistenti dalla piattaforma (sola lettura).

```hcl
# Tutti i trainer
data "labplatform_users" "trainers" {
  role = "trainer"
}

# Tutti i corsi
data "labplatform_courses" "all" {}

# Tutti i template di connessione
data "labplatform_templates" "all" {}

# Uso: trovare un corso per nome
locals {
  cka = [for c in data.labplatform_courses.all.courses : c if c.name == "CKA"][0]
}
```

---

## Esempi

La cartella `examples/` contiene 4 scenari progressivi:

| Cartella | Descrizione |
|----------|-------------|
| `01-basic-course/` | Corso base con template VNC e SSH |
| `02-full-class/` | Classe completa: corso, trainer, studenti, classe, assegnazioni |
| `03-multiple-classes/` | Due classi nella stessa settimana con data sources |
| `04-vsphere-lab/` | Laboratorio con VM vSphere |

Per eseguire un esempio:

```bash
cd examples/02-full-class/
cp terraform.tfvars.example terraform.tfvars
# Modifica terraform.tfvars con le tue credenziali

terraform init
terraform plan
terraform apply
```

---

## Workflow del tecnico

### Preparare un laboratorio per la settimana

```
1. Crea cartella:  mkdir labs/cka-2026-w13 && cd labs/cka-2026-w13
2. Copia esempio:  cp ../../terraform/examples/02-full-class/* .
3. Modifica:       vim terraform.tfvars  (studenti, date, credenziali)
4. Applica:        terraform init && terraform apply
```

### Aggiungere studenti dopo

Aggiungi studenti alla lista in `terraform.tfvars`:

```bash
terraform apply    # crea solo i nuovi studenti e le assegnazioni
```

### Fine corso

```bash
# Opzione 1: cambia stato nel .tf → status = "completed"
terraform apply

# Opzione 2: cancella tutto
terraform destroy
```

### Riusare per un'altra settimana

```bash
cp -r labs/cka-2026-w13 labs/cka-2026-w14
cd labs/cka-2026-w14
# Cambia date e studenti in terraform.tfvars
terraform init && terraform apply
```

### Struttura consigliata

```
labs/
├── cka-2026-w13/
│   ├── main.tf
│   └── terraform.tfvars      # gitignored
├── docker-2026-w13/
│   ├── main.tf
│   └── terraform.tfvars
└── .gitignore                 # *.tfvars, .terraform/, *.tfstate*
```

---

## Pubblicazione su Terraform Registry

Per evitare il build locale e permettere `terraform init` diretto, il provider va pubblicato sul **Terraform Registry**.

### Prerequisiti

1. **Repository GitHub dedicata** (nome: `terraform-provider-labplatform`)
2. **Chiave GPG** per firmare i rilasci
3. **Account su registry.terraform.io** collegato a GitHub

### Passi

#### 1. Creare la repo dedicata

Il Terraform Registry richiede che il provider sia in una repo con nome `terraform-provider-labplatform`. Copia la cartella `terraform/` in una nuova repo:

```bash
# Crea nuova repo su GitHub: desotech-it/terraform-provider-labplatform
gh repo create desotech-it/terraform-provider-labplatform --public

# Copia i file
mkdir /tmp/tf-provider && cd /tmp/tf-provider
cp -r ~/desolabs-labplatform/terraform/* .
cp -r ~/desolabs-labplatform/terraform/.goreleaser.yml .
cp -r ~/desolabs-labplatform/terraform/.github .
git init && git add -A
git commit -m "Initial commit: Terraform provider for LabPlatform"
git remote add origin git@github.com:desotech-it/terraform-provider-labplatform.git
git push -u origin main
```

#### 2. Generare chiave GPG

```bash
gpg --full-generate-key
# Tipo: RSA, 4096 bit, nessuna scadenza
# Nome: DesoTech Release Signing Key

# Esporta la chiave pubblica (servirà per il registry)
gpg --armor --export "DesoTech Release Signing Key"

# Esporta la chiave privata (servirà come secret GitHub)
gpg --armor --export-secret-keys "DesoTech Release Signing Key"
```

#### 3. Configurare i secrets GitHub

Nella repo `terraform-provider-labplatform`, vai in **Settings > Secrets and variables > Actions** e aggiungi:

| Secret | Valore |
|--------|--------|
| `GPG_PRIVATE_KEY` | Output di `gpg --armor --export-secret-keys` |
| `GPG_PASSPHRASE` | Passphrase della chiave GPG |

#### 4. Registrare su Terraform Registry

1. Vai su https://registry.terraform.io/sign-in
2. Accedi con il tuo account GitHub
3. Vai su **Publish > Provider**
4. Seleziona la repo `desotech-it/terraform-provider-labplatform`
5. Carica la **chiave GPG pubblica**
6. Conferma

#### 5. Creare un rilascio

```bash
cd terraform-provider-labplatform/
git tag v0.1.0
git push origin v0.1.0
```

La GitHub Action (`.github/workflows/release.yml`) si attiva automaticamente:
- Compila il provider per Linux, macOS, Windows (amd64 + arm64)
- Firma i checksum con GPG
- Crea la release su GitHub

Il Terraform Registry rileva la nuova release e la rende disponibile.

#### 6. Usare il provider pubblicato

```hcl
terraform {
  required_providers {
    labplatform = {
      source  = "desotech-it/labplatform"
      version = "~> 0.1"
    }
  }
}
```

```bash
terraform init    # scarica automaticamente il provider
```

### Alternativa: OpenTofu Registry

Il processo è identico ma su https://github.com/opentofu/registry. Basta aggiungere la repo al registry OpenTofu e funziona sia con `terraform` che con `tofu`.

---

## Troubleshooting

| Errore | Causa | Soluzione |
|--------|-------|-----------|
| `Authentication Failed` | Credenziali errate | Verificare env var |
| `API error (HTTP 409)` | Risorsa duplicata | Importarla o cambiare nome |
| `API error (HTTP 404)` | Risorsa non trovata | `terraform state rm` |
| `API error (HTTP 403)` | Permessi insufficienti | Usare account admin |

### Importare risorse esistenti

Se hai risorse create dalla UI e vuoi gestirle con Terraform:

```bash
terraform import labplatform_user.mario 42
terraform import labplatform_course.cka 5
terraform import labplatform_class.this 12
```

### Debug

```bash
TF_LOG=DEBUG terraform plan
```
