FROM golang:1.19.3-alpine AS build

WORKDIR /usr/src/app

RUN apk add build-base

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

RUN go build -v -o dist/ ./cmd/scout-api
RUN go build -v -o dist/ ./cmd/scout-scraper

FROM alpine:3.17.0 AS production

WORKDIR /app

COPY --from=build /usr/src/app/dist/scout-api /app
COPY --from=build /usr/src/app/dist/scout-scraper /app

EXPOSE 8080

CMD /app/scout-api