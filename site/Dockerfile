FROM golang:1.27.2-alpine3.24 AS builder

WORKDIR /src

# Download modules first so this layer caches until go.mod/go.sum change
COPY src/go.mod src/go.sum ./
RUN go mod download

COPY src/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/site ./cmd/site

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/site /site

EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/site"]
