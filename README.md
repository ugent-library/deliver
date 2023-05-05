# deliver

## Development

### Configuration

Deliver is configured with these environment variables:

* `DELIVER_PRODUCTION`
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
* `DELIVER_BANNER`

If a `.env` file is present in the project root, it's environment variables will be loaded.

To get started: 

```sh
cp .env.example .env
```

### Live reload

This project uses [reflex](https://github.com/cespare/reflex) to reload the app
server and recompile assets after changes.

```sh
go install github.com/cespare/reflex@latest
cp reflex.conf.example reflex.conf
reflex -d none -c reflex.conf
```