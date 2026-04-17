# syntax=docker/dockerfile:1.7

# ---------- Frontend build ----------
FROM node:25-alpine AS frontend
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ---------- Backend build ----------
FROM golang:1.26-alpine AS backend
WORKDIR /src
# Go modules first for layer caching.
COPY backend/go.mod backend/go.sum ./backend/
WORKDIR /src/backend
RUN go mod download

# Source + frontend dist for embed.
WORKDIR /src
COPY backend/ ./backend/
COPY --from=frontend /app/dist ./backend/web/dist

WORKDIR /src/backend
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux go build -tags prod \
    -ldflags="-s -w -X main.version=${VERSION}" \
    -o /out/monkey-planner ./cmd/monkey-planner

# ---------- Runtime ----------
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=backend /out/monkey-planner /monkey-planner

VOLUME ["/data"]
EXPOSE 8080

ENV MP_ADDR=":8080" \
    MP_DSN="sqlite:///data/monkey.db"

USER nonroot:nonroot
ENTRYPOINT ["/monkey-planner"]
