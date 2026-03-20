# Node TypeScript Starter

A pnpm-powered monorepo that boots a React 19 + TanStack front-end, a Hono-based API, a Better Auth session service, shared TypeScript contracts, and Kubernetes-ready infrastructure so you can experiment with full-stack patterns from a single workspace.

## Workspace layout

- `apps/` – The main applications: `api`, `auth`, and `web`. Each has its own `package.json`, source code, and build scripts.
  - `apps/api/` – Hono API server with `rolldown` bundling and `nodemon` for hot reloads.
  - `apps/auth/` – Better Auth session management service using Drizzle ORM for database interactions.
  - `apps/web/` – React 19 front-end with Vite and Nitro for server-side rendering and API proxying.
  - `apps/postgres/` – Helper scripts for managing Postgres databases, including multi-database creation for API and Auth services.
- `packages/` – Shared TypeScript code, currently just `core` for shared Zod schemas and types, but open to utilities, UI components, or anything else that can be consumed across apps.
  - `packages/core/` – Exports shared TypeScript types and Zod schemas, bundled with `tsdown` for easy consumption by both API and web clients.
- `infrastructure/` – Kubernetes & Docker manifests to bring the stack up locally or deploy it to a cluster.
  - `infrastructure/k8s/` - Kubernetes manifests for API, Auth, Web, and Postgres services, plus any shared resources like ConfigMaps or Secrets.
  - `infrastructure/docker/` - Docker compose files and helper scripts for running the stack without Tilt or Kubernetes, if desired.
- [`Tiltfile`](./Tiltfile) – The Tilt configuration that defines how to build and run each service in a Kubernetes-like environment, including watching for changes and rebuilding as needed.
- [`turbo.json`](./turbo.json) – Turborepo configuration for orchestrating builds, tests, and other scripts across the monorepo with caching and parallel execution.

## Prerequisites

1. **Node.js** ≥ v24 (`package.json` `engines`).
2. **pnpm** 10.24.0 (root `packageManager`).
3. **Docker** when you need Postgres containers or Tilt’s cluster emulation.
4. **Tilt** (`pnpm dev:k8s`) for the Kubernetes-style dev workflow.

## Getting started

### Running the project locally for development

```bash
# Clone the repo and install dependencies:
$ git clone git@github.com:awwwkshay/node-ts-starter.git
# Make sure to use pnpm 10.24.0 as specified in package.json:
$ npm install -g pnpm@10.24.0
# Install all dependencies across the monorepo:
$ pnpm install

# Login to whatever docker registries you plan to push to (e.g., GitHub Container Registry, Docker Hub):
$ docker login ghcr.io 
$ docker login docker.io

# Add these lines to .zshrc if on mac
export GITHUB_TOKEN=<your-github-token>
export GITHUB_USERNAME=<your-github-username>
export GITHUB_EMAIL=<your-github-email>

# Load the .zshrc changes:
$ source ~/.zshrc

# Verify the variables are set:
$ echo $GITHUB_TOKEN $GITHUB_USERNAME $GITHUB_EMAIL

# For development run using tilt:
$ tilt up
```

## Environment templates

| App | Template |
| --- | --- |
| API | [`apps/api/example.env`](./apps/api/example.env) |
| Auth | [`apps/auth/example.env`](./apps/auth/example.env) |
| Web | [`apps/web/example.env`](./apps/web/example.env) |

Copy the template, remove the leading `#` comments, and update the values to match your local credentials or hosted secrets.

## Contributing

1. Make changes in the package that owns the feature (`apps/api`, `apps/auth`, `apps/web`, or `packages/core`).
2. Rebuild the touched package (`pnpm --filter <package> build`).
3. Run `pnpm --filter <package> format` / `lint` to stay consistent.
4. Update `packages/core/dist` whenever shared exports change so the apps consume the latest bundle.

Generated on 2026-02-14
