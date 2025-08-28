# syntax=docker/dockerfile:1

# 1) Build the frontend (Vite) so //go:embed can include dist/*
FROM node:20-alpine AS ui
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
# Copy only files needed for the UI build
COPY index.html ./
COPY main.css ./
COPY vite.config.ts ./
COPY tsconfig.json ./
COPY postcss.config.js ./
COPY tailwind.config.js ./
COPY view/ ./view/
RUN npm run build:ui

# 2) Build the Go binary
FROM golang:1.23.1-alpine AS build
WORKDIR /src
RUN apk add --no-cache ca-certificates tzdata
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Bring in the prebuilt UI assets for go:embed
COPY --from=ui /src/dist ./dist
# Build a static binary for use with scratch
ENV CGO_ENABLED=0 GOOS=linux
RUN go build -ldflags="-s -w" -o /out/random ./...

# 3) Minimal runtime image
FROM scratch
# (Optional but nice) SSL certs & tzdata for HTTPS/timezones
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
# Our app
COPY --from=build /out/random /random
# Document the typical web port; Railway will still inject $PORT
EXPOSE 8080
ENV PORT=8080
ENTRYPOINT ["/random"]
