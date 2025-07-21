FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS app_build_stage

ARG TARGETOS
ARG TARGETARCH
ARG TARGETPLATFORM
ARG BUILDPLATFORM

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o book-store-service ./cmd/main.go

RUN chmod +x book-store-service

FROM --platform=$TARGETPLATFORM alpine:latest AS runner
WORKDIR /

COPY --from=app_build_stage /app/book-store-service /book-store-service

COPY --from=app_build_stage /etc/passwd /etc/passwd
USER 1000

EXPOSE 8080

HEALTHCHECK CMD ["./book-store-service", "--health"] || exit 1

ENTRYPOINT ["/book-store-service"]