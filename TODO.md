# TODOs

- [ ] Script to send locker expiry emails, a Go binary will run in as cron job
      in the Docker image. - For Zoe
- [ ] Cache static asset (middleware?)
- [x] Auth token invalidating
- [x] CSFR middleware
- [x] Admin routes
  - CSFR middleware will be applied to all routes (other than auth of course)
  - [x] Middleware admin token checker
  - [x] Auth: `PUT /auth/admin`
    - Sends user name and password via url form to authenticate
    - Redirects to `/admin` on success with cookies set
  - [x] Dash: `GET /admin/dash`
    - Query database for all registration, dump data onto an html table
      - Each cell has a form button to `DELETE /admin/api/registration`
  - [x] Remove registration: `DELETE /admin/registration`
    - Form data `locker` should be sent, containing the locker ID (ie, `ELW 120`).
  - [x] Export: `GET /admin/api/export`
    - Exports the current `registration` table into a csv file, then self-email
- [ ] Option to deregister locker/switch registered locker (optional)
- [ ] Make low numbered lockers searchable without leading 0s (low priority)