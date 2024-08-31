# Lockers

Built using the braindead stack of Golang, SQLite, HTMX, and TailwindCSS. Fully SSR, etc etc.

## Getting started

### Dependencies

Tailwind build tools are the only thing of JS ecosystem that is used in this project.

To install dependencies:

```sh
npm i
```

To start a Tailwind _"compiler"_:

```sh
npm run tw
```

### Database migration

```sh
go run ./cmd/migration
```

### Starts the app

```sh
go run ./cmd/app
```

Note: for auth cookie to work, go on your browser `http://127.0.0.1:8080`

### Environment variables

- `EMAIL_HOST_ADDRESS`: ECSS associated email (gmail) for sending locker-related email from
- `EMAIL_HOST_PASSWORD`: ECSS gmail's App password (if using Gmail, which is likely...)
- `SUPPORT_EMAIL`: Email (any type) for questions to be directed to
- `CIPHER_KEY`: base64 encoding for a cipher key, run `go run ./cmd/keygen` to generate one.
- `DOMAIN`: Hosting domain
- `DATABASE_URL`: Turso database url
- `DATABASE_AUTH_TOKEN`: Tursor database auth token
- `ADMIN_USERNAME`: admin username
- `ADMIN_PASSWORD`: admin password
