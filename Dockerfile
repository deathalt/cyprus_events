FROM golang:1.21.1 as build
WORKDIR /app
COPY . .
RUN go build -o event_exporter main.go

FROM ubuntu:latest
RUN apt update && apt install -y ca-certificates && apt clean -y
COPY --from=build /app/event_exporter /usr/bin
ENTRYPOINT [ "/usr/bin/event_exporter" ]
