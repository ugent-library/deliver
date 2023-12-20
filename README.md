# deliver

## Prerequisites

To run the application you will need:

* A PostgreSQL database
* An OpenID Connect endpoint
* An S3 compatible object store

## Configuration

Deliver is configured with these environment variables:

* `DELIVER_ENV`
* `DELIVER_HOST`
* `DELIVER_PORT`
* `DELIVER_ADMINS`
* `DELIVER_REPO_CONN`
* `DELIVER_STORAGE_BACKEND`
* `DELIVER_STORAGE_CONN`
* `DELIVER_OIDC_URL`
* `DELIVER_OIDC_ID`
* `DELIVER_OIDC_SECRET`
* `DELIVER_OIDC_REDIRECT_URL`
* `DELIVER_COOKIE_SECRET`
* `DELIVER_MAX_FILE_SIZE`
* `DELIVER_TIMEZONE`

## Docker Setup

You will need Docker and git to get started

* `cp docker-compose.example.yml docker-compose.yml`
* configure the oidc provider
* `docker compose up`

## Local development setup with live reload

For local development you will also need:

* Go >= 1.20
* A recent version of node.js
* A recent version of Postgres
* An S3 compatible object store

Initial setup:

```sh
cp .env.example .env
cp reflex.conf.example reflex.conf
make install-dev
```

To run the development server:

```sh
make dev
```
