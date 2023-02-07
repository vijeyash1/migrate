FROM golang:latest AS builder
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o dbmigrate ./main.go
FROM alpine:latest AS production
COPY --from=builder /app .
CMD ["./dbmigrate"]