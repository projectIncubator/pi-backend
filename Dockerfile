FROM golang:1.14.6-alpine3.12 as builder
COPY . /go/src/pi-backend/
WORKDIR /go/src/pi-backend/
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/pi-backend

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/pi-backend/build/pi-backend /usr/bin/pi-backend
EXPOSE 8000 8000
ENTRYPOINT ["/usr/bin/pi-backend"]
