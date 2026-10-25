FROM golang:1.27.2-alpine3.24 AS builder

WORKDIR /src

# Download modules first so this layer caches until go.mod/go.sum change
COPY src/go.mod src/go.sum ./
RUN go mod download

# Tailwind standalone CLI (musl build for Alpine), pinned and checksum-verified.
# Keep in sync with TAILWIND_VERSION / the musl checksum in the Makefile's release.
ARG TAILWIND_VERSION=v4.3.3
ARG TAILWIND_SHA256=a04d34ceacc8f52cbe8920ad846cdeb61d3d0021dba32db0d1f77c9d9fad7a6c

RUN apk add --no-cache libstdc++ libgcc \
 && wget -qO /usr/local/bin/tailwindcss \
      https://github.com/tailwindlabs/tailwindcss/releases/download/${TAILWIND_VERSION}/tailwindcss-linux-x64-musl \
 && echo "${TAILWIND_SHA256}  /usr/local/bin/tailwindcss" | sha256sum -c - \
 && chmod +x /usr/local/bin/tailwindcss



COPY src/ ./

# *_templ.go is gitignored, so generate it here with the templ tool pinned in go.mod
RUN go tool templ generate

# site.css is gitignored, so build it here (after templ, so all templates exist)
RUN tailwindcss -i web/static/css/input.css -o web/static/css/site.css --minify


RUN CGO_ENABLED=0 GOOS=linux go build -o /out/site ./cmd/site

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/site /site

EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/site"]
