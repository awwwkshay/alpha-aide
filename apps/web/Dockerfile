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
# 3. Build Web (for production)
# -------------------------------------------------------------------
FROM setup AS web_builder
WORKDIR /app

# Bring built packages
COPY --from=packages_builder /app/packages ./packages

COPY apps/web/package.json ./apps/web/
RUN pnpm install --frozen-lockfile

COPY apps/web ./apps/web/
RUN pnpm -r --filter ./apps/web build


# -------------------------------------------------------------------
# 4. Development Stage (pnpm dev)
# -------------------------------------------------------------------
FROM setup AS development

ENV NODE_ENV=development
ENV NITRO_PORT=3000
ENV NITRO_API_BASE_URL=http://api:8000
ENV VITE_API_BASE_URL=http://api:8000
ENV VITE_APP_NAME="Node TS Starter"

WORKDIR /app

# Copy everything (source for rebuilds)
COPY . .
# Install pnpm
RUN npm install -g pnpm

# Install deps
RUN pnpm install --frozen-lockfile

# Expose ports for Nitro and debugger
EXPOSE ${NITRO_PORT} 9229

# Start Nitro + Vite dev server
WORKDIR /app/apps/web
CMD ["pnpm", "dev"]


# -------------------------------------------------------------------
# 5. Production Stage (small + built only)
# -------------------------------------------------------------------
FROM node:24-slim AS production

ENV NODE_ENV=production
ENV NITRO_PORT=3000
ENV NITRO_API_BASE_URL=http://localhost:8000
ENV VITE_API_BASE_URL=http://localhost:8000
ENV VITE_APP_NAME="Node TS Starter"

WORKDIR /app

COPY pnpm-workspace.yaml ./
COPY package.json pnpm-lock.yaml ./

# Copy only the built artifacts (no src)
COPY --from=packages_builder /app/packages ./packages
COPY --from=web_builder /app/apps/web ./apps/web

# Install production dependencies
RUN npm install -g pnpm \
    && pnpm install --frozen-lockfile --prod --ignore-scripts

EXPOSE ${NITRO_PORT}

WORKDIR /app/apps/web
CMD ["node", "dist/server/index.mjs"]
