FROM golang:1.23
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./app ./app
COPY ./dto ./dto
COPY ./helpers ./helpers
COPY ./services ./services

WORKDIR /app/app

RUN CGO_ENABLED=0 GOOS=linux go build -o /item-tracker

EXPOSE 8080

CMD [ "/item-tracker" ]