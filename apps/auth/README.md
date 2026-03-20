# Auth service

Hono service wired with Better Auth and Drizzle so authentication, session, and JWT tooling run via a shared Postgres database.

## Prerequisites

1. Node.js ≥ 24.12.0 and pnpm 10.24.0 (workspace requirements).
2. Postgres database reachable via `DATABASE_URL`.
3. Copy `apps/auth/example.env`, remove the leading `#`, and populate values such as `BETTER_AUTH_SECRET`, `CLIENT_URLS`, and `BETTER_AUTH_URL` before running locally.

Install workspace deps once from the repo root:

```bash
pnpm install
```

## Scripts

| Command | Description |
| --- | --- |
| `pnpm --filter @awwwkshay/node-ts-auth-starter dev` | Runs `rolldown --watch` in parallel with `nodemon dist/server.js` for a live loop. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter build` | Bundles via `rolldown`. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter clean` | Removes `dist/`. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter format` / `format:check` | Run `oxfmt`. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter lint` | Run `oxlint` (type-aware). |
| `pnpm --filter @awwwkshay/node-ts-auth-starter watch` | `rolldown --watch`. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter run:nodemon` | Restart the built server (`dist/server.js`) when it changes. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter start` | Run the production bundle with `NODE_ENV=production node dist/server.js`. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter auth:generate` | Regenerate `src/db/auth-schema.ts` via the Better Auth CLI. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter auth:secret` | Print a new `BETTER_AUTH_SECRET` for cookies/tokens. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter db:generate` | Run `drizzle-kit generate`. |
| `pnpm --filter @awwwkshay/node-ts-auth-starter db:migrate` / `db:drop` | Execute or drop migrations via `drizzle-kit`. |

## Environment variables

| Name | Purpose | Notes |
| --- | --- | --- |
| `NODE_ENV` | Controls build/runtime mode | Defaults to `production`. |
| `PORT` | Port for the Hono server | Defaults to `8000`. |
| `CLIENT_URLS` | Allowed CORS origins | Split on commas inside `app.ts`. |
| `BETTER_AUTH_SECRET` | Cookie/jwt signing key | Minimum 18 characters; rotate with `auth:secret`. |
| `BETTER_AUTH_URL` | Public base URL | Used by Better Auth for cookies, docs, and redirect URIs. |
| `DATABASE_URL` | Postgres connection string | Consumed by Drizzle + Better Auth adapter in `src/lib/auth.ts`. |

The app validates this payload via `zod`, so it will refuse to start if a key is missing or malformed.

## Architecture highlights

- `apps/auth/src/app.ts` wires CORS, Better Auth plugins, and forwards every request to `better-auth` via `auth.handler`, giving you login/logout, session, user, verification, and docs endpoints out of the box.
- `apps/auth/src/lib/auth.ts` configures the Drizzle adapter so a single Postgres database drives both ORM models and Better Auth tables.
- `apps/auth/src/server.ts` bootstraps `@hono/node-server`, logs the configuration, and listens for `SIGINT`/`SIGTERM` for graceful shutdown.

## Database & auth tooling

- `auth:generate` keeps `src/db/auth-schema.ts` aligned with Better Auth’s schema definitions.
- `auth:secret` prints a new `BETTER_AUTH_SECRET` for signing cookies and tokens.
- `db:generate`, `db:migrate`, and `db:drop` use `drizzle-kit` to generate models and run migrations.

## Deployment tips

1. Ensure `NODE_ENV=production`, `CLIENT_URLS` includes your web origin, and your database migrations are applied.
2. Run `pnpm --filter @awwwkshay/node-ts-auth-starter build` before starting with `pnpm --filter @awwwkshay/node-ts-auth-starter start`.
3. The `/docs` UI, JWT endpoints, and cookie signing all rely on `BETTER_AUTH_SECRET`; rotate it carefully.

Generated on 2025-12-11
