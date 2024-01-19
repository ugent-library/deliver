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

## Dev Containers

This project supports [Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers). Following these steps
will auto setup a containerized development environment for this project. In VS Code, you will be able to start a terminal
that logs into a Docker container. This will allow you to write and interact with the code inside a self-contained sandbox.

**Installing the Dev Containers extension**

1. Open VS Code.
2. Go to the [Dev Containers](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers) extension page.
3. Click the `install` button to install the extension in VS Code.

**Open in Dev Containers**

1. Open the project directory in VS Code.
2. Click on the green "Open a remote window" button in the lower left window corner.
3. Choose "reopen in container" from the popup menu.
4. The green button should now read "Dev Container: App name" when successfully opened.
5. Open a new terminal in VS Code from the `Terminal` menu link.

You are now logged into the dev container and ready to develop code, write code, push to git or execute commands.

**Run the project**

1. Open a new terminal in VS Code from the `Terminal` menu link.
2. Execute this command `reflex -d none -c reflex.docker.conf`.
3. Once the application has started, VS Code will show a popup with a link that opens the project in your browser.