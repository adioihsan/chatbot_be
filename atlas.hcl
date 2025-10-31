# =========================
# Shared data source (GORM -> SQL)
# =========================
data "external_schema" "gorm" {
  program = ["go", "run", "-mod=mod", "./loader"]
}

# =========================
# Reusable migration config
# =========================
variable "MIGRATIONS_DIR" {
  type    = string
  default = "file://migrations"
}

# Common formatting (optional)
format {
  migrate {
    diff = "{{ sql . \"  \" }}"
  }
}

# =========================
# DEV environment (local)
# =========================
env "dev" {
  # Source schema produced by your loader (must be postgres)
  src = data.external_schema.gorm.url

  # Ephemeral scratch DB for planning/diffing
  dev = "docker://postgres/15/dev?search_path=public"

  # The DB you actually APPLY to (local dev DB)
  url = "postgres://dev_user:dev_pass@localhost:5432/dev_db?sslmode=disable&search_path=public"

  migration {
    dir = var.MIGRATIONS_DIR
  }
}

# =========================
# STAGING environment
# =========================
env "staging" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev?search_path=public"

  # Pull from env vars (recommended for secrets)
  url = getenv("STAGING_DATABASE_URL") # e.g. postgres://stg_user:stg_pass@stg-host:5432/stg_db?sslmode=require&search_path=public

  migration {
    dir = var.MIGRATIONS_DIR
  }
}

# =========================
# PROD environment
# =========================
env "prod" {
  src = data.external_schema.gorm.url
  dev = "docker://postgres/15/dev?search_path=public"

  # Always set via secret/env var in CI/CD
  url = getenv("DATABASE_URL") # e.g. postgres://prod_user:xxx@prod-host:5432/prod_db?sslmode=require&search_path=public

  migration {
    dir = var.MIGRATIONS_DIR
    # Safe guardrails for production
    replay = true      # apply in order from the beginning (prevents drift)
    # You can add check flags with CLI (see commands below)
  }
}
