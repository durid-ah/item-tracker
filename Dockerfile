FROM golang:1.23 AS build-stage
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app
COPY ./dto ./dto
COPY ./helpers ./helpers
COPY ./services ./services

WORKDIR /app/app

RUN CGO_ENABLED=0 GOOS=linux go build -o /item-tracker

FROM alpine AS release
WORKDIR /
COPY --from=build-stage /item-tracker /item-tracker

EXPOSE 8080
USER nonroot:nonroot

CMD [ "/item-tracker" ]