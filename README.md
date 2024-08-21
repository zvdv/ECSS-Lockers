# Lockers

Built using the braindead stack of Golang, HTMX, and TailwindCSS

## Getting started

Tailwind server and dependencies:

```sh
npm i # install dependencies
npm run tw # starts the tailwind "compiler"
```

### Database migration

```sh
go run ./cmd/migration.go
```

### Starts the app

```sh
go run ./cmd/app.go
```

Note: for auth cookie to work, go on your browser `http://127.0.0.1:8080`

### Environment variables

```txt
EMAIL_HOST_ADDRESS=
EMAIL_HOST_PASSWORD=
CIPHER_KEY=
DOMAIN=
DATABASE_URL=
DATABASE_AUTH_TOKEN=
```
