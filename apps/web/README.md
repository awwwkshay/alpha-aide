# @awwwkshay/node-ts-web-starter

React 19 client that pairs Vite, Nitro, and TanStack Router/Query with shared `todoSchema` contracts from `packages/core`.

## Prerequisites

1. Node.js ≥ 24.12.0 and pnpm 10.24.0 (workspace-wide requirements).
2. Copy `apps/web/example.env`, remove the comments, and set `NITRO_API_BASE_URL`, `VITE_API_BASE_URL`, `NITRO_PORT`, and `VITE_APP_NAME` appropriately.
3. Install dependencies from the workspace root:

```bash
pnpm install
```

## Scripts

| Command | Description |
| --- | --- |
| `pnpm --filter @awwwkshay/node-ts-web-starter dev` | Vite dev server with Nitro dev mode and TanStack Router devtools. |
| `pnpm --filter @awwwkshay/node-ts-web-starter build` | Nitro + Vite production build (`dist/`). |
| `pnpm --filter @awwwkshay/node-ts-web-starter preview` | Preview the production build via `vite preview`. |
| `pnpm --filter @awwwkshay/node-ts-web-starter test` | Run `vitest run` for unit/integration tests. |
| `pnpm --filter @awwwkshay/node-ts-web-starter start` | Start the compiled Nitro SSR server (`node ./dist/server/index.mjs`). |
| `pnpm --filter @awwwkshay/node-ts-web-starter format` / `format:check` | Run `oxfmt`. |
| `pnpm --filter @awwwkshay/node-ts-web-starter lint` | Run `oxlint` (type-aware). |
| `pnpm --filter @awwwkshay/node-ts-web-starter clean` | Remove `dist/`. |

## Environment

```bash
NITRO_PORT=3000
NITRO_API_BASE_URL=http://localhost:8000
VITE_API_BASE_URL=http://localhost:8000
VITE_APP_NAME=Node TS Starter
```

- `NITRO_API_BASE_URL` points to the API server for SSR fetches; `VITE_API_BASE_URL` mirrors this for browser-only calls.
- `NITRO_PORT` controls where the SSR server listens during `pnpm start` or `nitro preview`.
- Any env referenced above can be overridden via `apps/web/example.env` or while running inside Tilt.

## Architecture

- `src/router.tsx` wires TanStack Router with the generated `routeTree` and shared query client for SSR-friendly fetching.
- `src/routes/` (e.g., `__root.tsx`, `index.tsx`) define the layout and hero experience.
- `src/integrations/tanstack-query/` exposes the DevTools-enabled TanStack Query client used by loaders.
- Styling is handled by Tailwind 4 plus `src/styles.css`, with Lucide icons and Motion for gradients.
- The shared `todoSchema` lives in `packages/core` so the UI validates API responses with the same contract.

## Development flow

1. Start the API: `pnpm --filter @awwwkshay/node-ts-api-starter dev`.
2. Start the web app: `pnpm --filter @awwwkshay/node-ts-web-starter dev`.
3. Visit the Vite URL (usually `http://localhost:5173`). Changes to `src/routes/` regenerate `routeTree.gen.ts` via the plugin configured in `vite.config.ts`.

## Testing & quality

- `vitest`: run `pnpm --filter @awwwkshay/node-ts-web-starter test` to execute unit/integration suites.
- `oxfmt` / `oxlint`: formatting and lint commands keep this package consistent with the API/auth services.
- Tailwind + Motion ensure the UI stays responsive and expressive without relying on default styles.

## Deployment

1. Build the SSR bundle: `pnpm --filter @awwwkshay/node-ts-web-starter build`.
2. Set `NITRO_API_BASE_URL` to your production API and run `pnpm --filter @awwwkshay/node-ts-web-starter start` or `node ./dist/server/index.mjs`.
3. Container deployments should copy `dist/`, install production deps, expose `NITRO_PORT`, and run the same start command.

Generated on 2025-12-11
