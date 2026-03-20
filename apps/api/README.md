# @awwwkshay/node-ts-api-starter

A lightweight Hono service that exposes `/todos`/`/info`, validates every payload with `zod`, and shares the `todoSchema` contract from `packages/core`.

## Highlights

- **Framework:** [`Hono`](https://hono.dev) with `@hono/node-server` for lean HTTP handling.
- **Schema safety:** `zod` validates the TODO payload, env vars, and response contracts.
- **Shared contracts:** Imports `@awwwkshay/node-ts-core` so the API/web layers share the same `todoSchema`.

## Prerequisites

1. **Node.js** ≥ 24.12.0 (per the root `package.json`).
2. **pnpm** 10.24.0 across the workspace.
3. **Postgres** instance accessible via `DATABASE_URL`.

Install deps from the repo root once:

```bash
pnpm install
```

## Scripts

| Command | Description |
| --- | --- |
| `pnpm --filter @awwwkshay/node-ts-api-starter dev` | Development runner: `rolldown --watch` + `nodemon dist/server.js`, watches files and reloads on change. |
| `pnpm --filter @awwwkshay/node-ts-api-starter build` | Bundle with `rolldown`. |
| `pnpm --filter @awwwkshay/node-ts-api-starter clean` | Remove `dist/`. |
| `pnpm --filter @awwwkshay/node-ts-api-starter format` / `format:check` | Run `oxfmt`. |
| `pnpm --filter @awwwkshay/node-ts-api-starter lint` | Run `oxlint` (type-aware). |
| `pnpm --filter @awwwkshay/node-ts-api-starter watch` | Build with `rolldown --watch`. |
| `pnpm --filter @awwwkshay/node-ts-api-starter run:nodemon` | Restart built bundle (`dist/server.js`) on changes. |
| `pnpm --filter @awwwkshay/node-ts-api-starter start` | Start production node bundle. |
| `pnpm --filter @awwwkshay/node-ts-api-starter auth:generate` | Regenerate `src/db/auth-schema.ts`. |
| `pnpm --filter @awwwkshay/node-ts-api-starter auth:secret` | Print a new `BETTER_AUTH_SECRET`. |
| `pnpm --filter @awwwkshay/node-ts-api-starter db:generate` | Generate the Drizzle schema. |
| `pnpm --filter @awwwkshay/node-ts-api-starter db:migrate` / `db:drop` | Run/drop migrations via `drizzle-kit`. |

## Environment

Values are validated through `apps/api/src/schemas/env.ts`. The app will refuse to boot when validation fails.

```env
NODE_ENV=development      # development | testing | production
PORT=8000
CLIENT_URLS=http://localhost:5173
DATABASE_URL=postgresql://user:pass@host:5432/app
```

- `NODE_ENV` triggers `dotenv` in development and controls logging detail.
- `CLIENT_URLS` is split and supplied to the `@hono/cors` middleware so only the web UI can call this API.
- Use the template under `apps/api/example.env` when running locally.

## Architecture

1. `src/index.ts` wires up the Hono router, env loader, and helper services.
2. `configs/env.ts` loads `.env` (development) or `process.env` (production) and normalizes the payload.
3. `schemas/env.ts` uses `zod` to keep the runtime contract strict, so the server fails fast on misconfiguration.
4. Common `todoSchema` definitions live in `packages/core` so the API, auth, and web layers stay aligned.

## Exposed endpoints

| Method | Path | Description |
| --- | --- | --- |
| `GET /` | `/` | Quick health check returning `Hello Hono`. |
| `GET /info` | `/info` | Echoes the active `NODE_ENV`, `PORT`, and allowed `CLIENT_URLS`. |
| `GET /todos` | `/todos` | Returns a static TODO list validated by `todoSchema`. Replace with your datastore when needed. |

Every route is typed via `Hono<{ Bindings; Variables }>` for strict handler context.

## Development notes

- CORS locks requests to `CLIENT_URLS` and caches preflight responses for 10 minutes (`maxAge: 600`).
- `todoSchema` is shared through `@awwwkshay/node-ts-core`, demonstrating cross-package schemas.
- Logging surfaces the parsed env payload when `NODE_ENV=development` so you can debug faster.

To connect with the UI, run the API dev server and the web service in parallel (`pnpm --filter @awwwkshay/node-ts-web-starter dev`). Adjust `CLIENT_URLS` if the browser runs from another port.

## Deployment

1. Run `pnpm --filter @awwwkshay/node-ts-api-starter build`.
2. Deploy via `NODE_ENV=production node dist/server.js` or your container flavor.
3. Provide the same env vars via your platform or container spec.
4. On container platforms, install only prod deps (`pnpm install --prod`) and point to `dist/server.js`.

## Troubleshooting

- `ZodError` on startup implies missing/wrong env values. Double-check your `.env`.
- CORS failures mean the browser origin doesn’t match `CLIENT_URLS`.

Generated on 2025-12-11
