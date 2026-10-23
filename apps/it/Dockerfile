FROM node:24.22.1-alpine3.24@sha256:06d7adf3e0ef61a13f1d28682726366f97a15aa1888c934d3ecec86c33918b53 AS builder
WORKDIR /workspace
COPY package.json package-lock.json ./
COPY apps ./apps
COPY packages ./packages
RUN npm ci --no-audit --no-fund
RUN npm run build -w @trickreport/it

FROM node:24.22.1-alpine3.24@sha256:06d7adf3e0ef61a13f1d28682726366f97a15aa1888c934d3ecec86c33918b53 AS runtime
WORKDIR /workspace
ENV NODE_ENV=production HOSTNAME=0.0.0.0 PORT=3000 NEXT_TELEMETRY_DISABLED=1
COPY --from=builder --chown=node:node /workspace/apps/it/.next/standalone ./
COPY --from=builder --chown=node:node /workspace/apps/it/.next/static ./apps/it/.next/static
COPY --from=builder --chown=node:node /workspace/apps/it/public ./apps/it/public
USER node
EXPOSE 3000
HEALTHCHECK --interval=10s --timeout=5s --retries=5 --start-period=20s CMD wget -qO- http://127.0.0.1:3000/login || exit 1
CMD ["node", "apps/it/server.js"]
