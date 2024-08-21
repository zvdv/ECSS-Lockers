# TODOs

- [ ] Script to send locker expiry emails, a Go binary will run in as cron job
      in the Docker image. - For Zoe
- [ ] Cache static asset (middleware?)
- [ ] Auth token invalidating
- [ ] CSFR middleware
- [ ] Admin routes
  - CSFR middleware will be applied to all routes (other than auth of course)
  - [ ] Middleware admin token checker
  - [ ] Auth: `PUT /auth/admin/`
    - Sends user name and password via url form to authenticate
    - Redirects to `/admin/dash/` on success with cookies set
  - [ ] Dash: `GET /admin/dash/`
    - Query database for all registration, dump data onto an html table
      - Each cell has a form button to `DELETE /admin/api/registration`
  - [ ] Remove registration: `DELETE /admin/api/registration`
    - Form data `locker` should be sent, containing the locker ID (ie, `ELW 120`).
  - [ ] Export: `GET /admin/api/export`
    - Empty body
    - Exports the current `registration` table into a csv file, then self-email

## Getting started with Emily

- go over how templating works - this is simple, she can read the docs but for
  the sake of hitting the ground running. Docs: https://pkg.go.dev/text/template.
  Note that we are using `html/template` instead, but the syntax is similar
- go over how auth works, and what would be added to the auth flow with CSFR token.
- go over general architect? the code is self-explanatory but if she wants
- go over generating env variables:
  - app password for gmail
  - database url and token with Turso
  - generating key
