# site

The rnorwood.com portfolio: Go + [templ](https://templ.guide) + [htmx](https://htmx.org) + Tailwind,
server-rendered and shipped as a single distroless binary.

## Quick start

```bash
make dev    # live reload at http://localhost:7331 (templ proxy) — reads content/ from disk
make test   # unit tests, including validation of the real content YAML
```

templ is pinned in `go.mod` (`go tool templ`). The Tailwind standalone CLI is downloaded and
checksum-verified into `bin/` by `make`. No Node toolchain is needed.

| Command | What it does |
|---|---|
| `make dev` | templ watch + Tailwind watch + `go run`, with browser auto-reload |
| `make test` | `go vet` + `go test ./...` |
| `make build` / `make run` | Build `bin/site` (generates templ + CSS first) and run it on :8080 |
| `make build-docker` / `make up` / `make down` | Production image via Compose, on a read-only filesystem |
| `make clean` | Remove `bin/`, generated `*_templ.go` and `site.css` |

## Layout

```
src/
  cmd/site/          main: wiring only (content source, image builder, HTTP server)
  content/           site.yaml + albums/<slug>/album.yaml, embedded into the binary
  internal/gallery/  domain model, Store interface, YAML loader + validation
  internal/images/   URLBuilder: the only code that knows where image bytes live
  internal/server/   routes, handlers, domain -> view-model mapping
  web/components/    templ templates (they render view models only)
  web/static/        CSS (Tailwind input), JS, images, vendored htmx + PhotoSwipe
```

## Content

Each album is `content/albums/<dir>/album.yaml`:

```yaml
title: "Killer Whales Bay"
slug: killer-whales-bay   # optional, defaults to <dir>
date: 2023-07-14
description: "..."
tags: [wildlife, alaska]
featured: true
cover: <sha>              # optional, defaults to the first image
images:                   # display order
  - { sha: <sha>, width: 6000, height: 4000 }
```

Content is validated at startup and by `make test`. A bad album fails with an error naming its file.

**Mock images:** the current content uses fake `sha` values (`mock-...`), and images come from
[picsum.photos](https://picsum.photos) seeded by that value. To switch to real images, add a new
`images.URLBuilder` and change the one line in `cmd/site/main.go`.

## Environment

| Variable | Default | Effect |
|---|---|---|
| `CONTENT_DIR` | unset (embedded content) | Read content YAML from this directory instead. `make dev` sets it. |

## Routes

`/` home · `/albums?tag=&page=` filtered home (always a full page) · `/fragments/albums?tag=&page=`
infinite-scroll batch (htmx) · `/albums/{slug}` album · `/static/*` · `/healthz`
