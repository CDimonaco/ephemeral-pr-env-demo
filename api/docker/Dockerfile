FROM golang:1.21.3 as builder

WORKDIR /app

RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz

COPY ./api /app

RUN task build-server
RUN task build-seeder

RUN chmod +x build/server && chmod +x build/seeder

FROM alpine:3

LABEL org.opencontainers.image.source="https://github.com/cdimonaco/ephimeral-pr-env-demo"
LABEL org.opencontainers.image.authors="Carmine Di Monaco <carmine.dimonaco@gmail.com>"
LABEL org.opencontainers.image.title="Ephimeral PR env demo - API"

WORKDIR /app

COPY --from=builder /app/build/server /app
COPY --from=builder /app/build/seeder /app
COPY --from=builder /app/migrate /app

COPY ./api/docker/entrypoint.sh /app/entrypoint.sh
COPY ./api/internal/persistence/seed/test.sql /app/seed.sql
COPY ./api/internal/persistence/migrations/ /app/migrations

RUN chmod +x /app/entrypoint.sh

entrypoint ["/app/entrypoint.sh"]
