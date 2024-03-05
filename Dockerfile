FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN apt-get update && apt-get install -y \
    sudo \
    curl \
    && rm -rf /var/lib/apt/lists/*

# Download and install the Ookla Speedtest CLI
RUN curl -s https://packagecloud.io/install/repositories/ookla/speedtest-cli/script.deb.sh | bash
RUN apt-get install -y speedtest

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway-monitor

CMD ["/gateway-monitor"]