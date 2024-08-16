# Lockers

## Getting started

Tailwind server and dependencies:

```sh
npm run i # install dependencies
npm run tw # starts a tailwind compiler
```

Starts the app

```sh
go run ./cmd/app.go
```

Note: for auth cookie to work, go on your browser `http://127.0.0.1:8080`

### Environment variables

```txt
GMAIL_USER=<ecss email>
GMAIL_PASSWORD=<ecss email password/app password>
ORIGIN=<hosting domain>
CIPHER_KEY=<32 char key>
```
