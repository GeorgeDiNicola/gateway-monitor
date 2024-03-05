# syntax=docker/dockerfile:1

FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

USER root

RUN apt-get update && apt-get install -y sudo

RUN sudo apt-get install curl && \
    curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | sudo bash && \
    sudo apt-get install speedtest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway-monitor

CMD ["/gateway-monitor"]