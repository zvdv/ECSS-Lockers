FROM golang:1.23.0-alpine3.20 AS gobuild
WORKDIR /locker
COPY go.mod go.sum ./
COPY /internal ./internal
COPY /cmd/app/ ./cmd/app/
RUN go mod tidy
RUN go build -o app ./cmd/app/main.go

FROM node:22-alpine3.19 AS nodebuild
WORKDIR /locker
COPY package.json package-lock.json index.css tailwind.config.js ./
COPY /internal/router ./internal/router
COPY /templates ./templates
RUN npm install
RUN npm run tw:buildonly

FROM alpine:3.20
WORKDIR /locker
COPY /templates ./templates
COPY /assets ./assets
COPY --from=gobuild /locker/app ./
COPY --from=nodebuild /locker/assets/css/index.css ./assets/css/index.css
RUN apk add --no-cache tzdata
ENV TZ=Canada/Pacific
EXPOSE 8080
CMD [ "./app" ]