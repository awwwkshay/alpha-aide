# Postgres helper

This folder hosts `create-multiple-postgresql-databases.sh`, a helper script that runs inside the official Postgres Docker image to provision multiple databases/users in one go. Tilt mounts this script into the container under `/docker-entrypoint-initdb.d`, so every time Postgres boots (and detects a fresh data directory) it can create the extra databases defined via `POSTGRES_MULTIPLE_DATABASES`.

## Usage

### With Tilt or Docker Compose

1. Ensure the container mounts this directory to `/docker-entrypoint-initdb.d`.
2. Set environment variables before starting Postgres:

```yaml
POSTGRES_MULTIPLE_DATABASES: "app,auth"
POSTGRES_USER: postgres
POSTGRES_PASSWORD: securepassword
```

Postgres runs the script once when the data directory is empty, and it will create a user/database pair for every comma-separated name listed in `POSTGRES_MULTIPLE_DATABASES`.

### Custom image

If you build your own Postgres image, copy this script into it or mount the directory at runtime. The script expects `POSTGRES_USER` to exist (it uses the official entrypoint defaults) and requires `POSTGRES_MULTIPLE_DATABASES` to be set to a comma-separated list of names.

### Non-standard names

Wrap names with quotes if they contain hyphens, uppercase letters, or other special characters. Example:

```yaml
POSTGRES_MULTIPLE_DATABASES: "test-db-1","test-db-2"
```

## Scripts

`create-multiple-postgresql-databases.sh` looks for `POSTGRES_MULTIPLE_DATABASES`, splits the value on commas, and for each entry it:

1. Creates a user named after the database.
2. Creates a matching database.
3. Grants all privileges on that database to the user.

Environment variables such as `POSTGRES_DB`, `POSTGRES_USER`, and `POSTGRES_PASSWORD` still respect the official Postgres image defaults.

Generated on 2025-12-11
