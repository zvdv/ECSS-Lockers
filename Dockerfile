FROM node:20.11-alpine AS build

WORKDIR /app
COPY package.json .
COPY package-lock.json .
COPY vendor vendor
RUN npm ci
COPY . .
RUN npm run build
RUN npm run node:build

FROM node:20.11-alpine AS deploy

WORKDIR /app

RUN chown node:node ./
USER node
EXPOSE 3000
COPY --from=build /app/package-lock.json .
COPY --from=build /app/package.json .
COPY --from=build /app/vendor vendor

ARG NODE_ENV=production
ENV NODE_ENV $NODE_ENV

RUN npm ci --omit dev

COPY --from=build /app/build build

CMD ["node", "build"]
