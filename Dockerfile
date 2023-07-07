FROM golang:1.20.2-alpine as builder
WORKDIR /opt
COPY . .
RUN go mod tidy; \
    go fmt ./...; \
    go mod download;
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./node-wizard .

FROM alpine:3.17
WORKDIR /app
COPY --from=builder /opt/node-wizard /app/node-wizard
CMD [ "./node-wizard" ]