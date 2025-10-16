FROM golang:1.25-alpine3.22 AS build

LABEL maintainer="meetmorrowsolonmars@gmail.com"

ARG APP_PATH

WORKDIR /app

COPY . /app/

RUN go mod download

RUN go build -o ./bin/app $APP_PATH

FROM ubuntu:22.04 AS app

LABEL maintainer="meetmorrowsolonmars@gmail.com"

WORKDIR /app

EXPOSE 80/tcp

EXPOSE 84/tcp

RUN groupadd -r appuser && useradd -r -s /bin/false -g appuser appuser

COPY --from=build /app/bin /app

RUN chown -R appuser:appuser /app

RUN apt-get update && apt-get install -y ca-certificates

USER appuser

ENTRYPOINT [ "./app" ]
