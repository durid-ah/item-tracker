FROM node:21.6.2-alpine AS frontend-base

WORKDIR /usr/src/app
################################################################################
# Create a stage for installing production dependecies.
FROM frontend-base AS frontend-deps

# Download dependencies as a separate step to take advantage of Docker's caching.
# Leverage a cache mount to /root/.npm to speed up subsequent builds.
# Leverage bind mounts to package.json and package-lock.json to avoid having to copy them
# into this layer.
RUN --mount=type=bind,source=frontend/package.json,target=package.json \
    --mount=type=bind,source=frontend/package-lock.json,target=package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci --omit=dev

################################################################################
# Create a stage for building the application.
FROM frontend-deps AS frontend-build

# Download additional development dependencies before building, as some projects require
# "devDependencies" to be installed to build. If you don't need this, remove this step.
RUN --mount=type=bind,source=frontend/package.json,target=package.json \
    --mount=type=bind,source=frontend/package-lock.json,target=package-lock.json \
    --mount=type=cache,target=/root/.npm \
    npm ci

COPY ./frontend .

RUN npm run build

FROM golang:1.23-bullseye AS build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app
COPY --from=frontend-build /usr/src/app/public /app/app/public
COPY ./dto ./dto
COPY ./helpers ./helpers
COPY ./services ./services

WORKDIR /app/app

RUN CGO_ENABLED=1 GOOS=linux go build -o /item-tracker -a -ldflags '-linkmode external -extldflags "-static"'

FROM debian:bullseye-slim AS release
WORKDIR /
COPY --from=build-stage /item-tracker /item-tracker

RUN mkdir /var/db

EXPOSE 8080

CMD [ "/item-tracker" ]
