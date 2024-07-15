FROM golang:1.22-alpine3.19 AS build-stage

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /capstone-project

FROM alpine:3.19 AS build-release-stage

WORKDIR /

#COPY --from=build-stage /app/.env /.env

COPY --from=build-stage /capstone-project /capstone-project

EXPOSE 8080

ENTRYPOINT ["./capstone-project"]
