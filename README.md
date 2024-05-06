# Lockers

This is a website for the registration and management of lockers. It was built for the University of Victoria's Engineering and Computer Science Students' Society, so it may not meet your needs exactly.

## Running

The recommended way to run this software is via the Dockerfile included with the repo. These are the required environment variables for it to work:

- `JWT_SECRET`: a randomly chosen value for signing tokens
- `ADMIN_PASSWORD`: the password for the admin panel. Enter the on the login page instead of the email to access the admin panel.
- `DATABASE_URL`: the connection url for a Planetscale database
- `GMAIL_USER`: the email for a gmail account
- `GMAIL_PASSWORD` an application password for a gmail account. This is specially generated for the application, not your regular login.
- `ORIGIN`: the expected origin URL. This is required to prevent CSRF attacks. You may disable this in the SvelteKit config file.
- `TZ`: optional. This is a TZ identifier. Set this if you don't want your partially hydrated content to be rendered in UTC.

You should also expose port 3000, as this serves the HTTP connection.

## Development

You should set all of the environment variables above in a `.env` file for use during development.

Run these commands to get setup:

```console
$ npm install
$ npm run db:codegen
```

To run:

```console
$ npm run dev
```

To build:

```console
$ npm run build
```

You may use any of the scripts in the `db` folder to access the database. All of them require the `DATABASE_URL` environment variable to be set.

Note that this project uses SvelteKit's NodeJS adapter. This means that shared singletons like database connections are acceptable for performance reasons, since it won't be running in a serverless environment. However, any state that MUST persist between connections should be put in the database.

If `git push` hangs, you may need to refer to [this solution](https://stackoverflow.com/a/68711337). I'm assuming this occurs because there are large files in the `vendor` directory.

## Deployment

The app is currently hosted with Google Cloud Platform on a micro instance of Google Compute
Engine (GCE) and with Docker.

### Build image

Build the image and push it to docker registry, since our local machines are a lot faster at
this than having to build it on the tiny GCE instance.

```sh
npm run docker:deploy
```

- Note that currently it's being pushed onto my registry at `hn275/ecss-locker`. To whom might
  be dealing with deployment in the future, you will need to modify the
  [docker-compose.yml](./docker-compose.yml) file accordingly.

To apply the changes, ssh into the GCE server and changed into the directory <insert dir here>
then restart the docker container

```sh
docker compose down && docker compose up -d
```

### Database migration

The [schema.sql](./db/schema.sql) file now contains the code queries to create the table and to seed
all the lockers.

TODO: write this procedure down

## Design

This website was designed for maximum usability on mobile devices with unstable or slow internet. Excluding admin functionality, the entire website has a mobile-first layout and works without JavaScript. JavaScript is used to progressively enhance the experience with features such as loading states and a client-side router.
