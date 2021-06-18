FROM golang:1.14.6-alpine3.12 as builder
COPY go.mod go.sum /go/src/pi-backend/
WORKDIR /go/src/go-api
RUN go mod download
COPY . /go/src/go-api
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o build/pi-backend pi-backend

FROM alpine
RUN apk add --no-cache ca-certificates && update-ca-certificates
COPY --from=builder /go/src/pi-backend/build/pi-backend /usr/bin/pi-backend
EXPOSE 8080 8080
ENTRYPOINT ["/usr/bin/pi-backend"]
