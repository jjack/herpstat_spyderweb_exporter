# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

# directory for herpstat_exporter's source code
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -o herpstat_exporter

FROM scratch
COPY --from=builder /app/herpstat_exporter /herpstat_exporter
EXPOSE 10010
ENTRYPOINT [ "./herpstat_exporter" ]
