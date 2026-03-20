# -------------------------------------------------------------------
# 1. Base setup (pnpm + workspace files)
# -------------------------------------------------------------------
FROM node:24-slim AS setup

WORKDIR /app
RUN npm install -g pnpm

COPY pnpm-workspace.yaml ./
COPY package.json pnpm-lock.yaml ./


# -------------------------------------------------------------------
# 2. Build packages (for production)
# -------------------------------------------------------------------
FROM setup AS packages_builder
WORKDIR /app

COPY packages ./packages/
RUN pnpm install --frozen-lockfile
RUN pnpm -r --filter "./packages/*" build


# -------------------------------------------------------------------
# 3. Build API (for production)
# -------------------------------------------------------------------
FROM setup AS api_builder
WORKDIR /app

# Bring built packages
COPY --from=packages_builder /app/packages ./packages

COPY apps/api/package.json ./apps/api/
RUN pnpm install --frozen-lockfile

COPY apps/api ./apps/api/
RUN pnpm -r --filter ./apps/api build


# -------------------------------------------------------------------
# 4. Development Stage (Live dev: pnpm dev)
# -------------------------------------------------------------------
FROM setup AS development

ENV NODE_ENV=development
ENV PORT=8000
ENV CLIENT_URLS=http://localhost:3000
ENV DATABASE_URL="postgresql://postgres:Akshay@123@host.docker.internal:5432/app"

WORKDIR /app

# Copy EVERYTHING needed for dev mode
COPY . .

# Install all deps (root + packages + apps)
RUN pnpm install --frozen-lockfile

EXPOSE ${PORT} 9229

WORKDIR /app/apps/api
CMD ["pnpm","dev"]


# -------------------------------------------------------------------
# 5. Production Stage (small + built only)
# -------------------------------------------------------------------
FROM node:24-slim AS production

ENV NODE_ENV=production
ENV PORT=8000
ENV CLIENT_URLS=http://localhost:3000
ENV DATABASE_URL="postgresql://postgres:Akshay@123@host.docker.internal:5432/app"

WORKDIR /app

COPY pnpm-workspace.yaml ./
COPY package.json pnpm-lock.yaml ./

# Copy only the built artifacts (no src)
COPY --from=packages_builder /app/packages ./packages
COPY --from=api_builder /app/apps/api ./apps/api

# Install production-only deps
# Install production dependencies
RUN npm install -g pnpm \
    && pnpm install --frozen-lockfile --prod --ignore-scripts

EXPOSE ${PORT}

WORKDIR /app/apps/api
CMD ["node", "dist/server.js"]
