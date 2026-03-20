# @awwwkshay/node-ts-core

Shared schemas, types, and helpers consumed by `apps/api`, `apps/web`, and any additional workspace packages.

## Key exports

- `todoSchema`: A `zod` schema describing each TODO item (`id`, `title`, `completed`).
- Re-exported types/helpers inside `src/index.ts` so downstream packages can import everything from a single entry point.
- Build outputs (CJS + ESM + `.d.ts`) under `dist/` so consuming apps stay aligned without duplicating code.

## Installation & usage

Install workspace dependencies from the repo root:

```bash
pnpm install
```

Then import the schema anywhere in the workspace:

```ts
import { todoSchema } from "@awwwkshay/node-ts-core";

const parsed = todoSchema.safeParse({
  id: "1",
  title: "Share schemas",
  completed: false,
});
```

The package exposes both ESM (`dist/index.mjs`) and CJS (`dist/index.cjs`) bundles via the `exports` field so consumers can choose their preferred loader.

## Development workflow

| Script | Description |
| --- | --- |
| `pnpm --filter @awwwkshay/node-ts-core dev` | Run `tsdown --watch` for incremental builds. |
| `pnpm --filter @awwwkshay/node-ts-core build` | Run `tsdown` once to emit CJS/ESM and TypeScript declarations. |
| `pnpm --filter @awwwkshay/node-ts-core type-check` | Run `tsc --noEmit`. |
| `pnpm --filter @awwwkshay/node-ts-core clean` | Remove `dist/`. |
| `pnpm --filter @awwwkshay/node-ts-core format` / `format:check` | Run `oxfmt`. |
| `pnpm --filter @awwwkshay/node-ts-core lint` | Run `oxlint` (type-aware). |

After rebuilding, commit both `src/` and the generated `dist/` files so the other packages immediately consume the updated artifacts.

## Project structure

```bash
packages/core/
├── src/index.ts        # Re-export shared schemas/helpers
├── src/todo.ts         # TODO schema definition
├── package.json        # Scripts + build config
├── tsdown.config.ts    # tsdown build options
└── tsconfig.json       # Compiler settings
```

`tsdown` already outputs CJS (`dist/index.cjs`), ESM (`dist/index.mjs`), and typings (`dist/index.d.mts`). Adjust `tsdown.config.ts` only if you need extra entry points or custom bundling logic.

## Testing contracts

Treat schema changes as potential breaking changes. Update dependent applications (especially `apps/api` and `apps/web`) after modifying exports and rerun their build/lint scripts to ensure everything still type checks.

## Publishing

This package is workspace-local. No registry publishing is required—keep the version in sync manually if you ever extract it for npm.

Generated on 2025-12-11
