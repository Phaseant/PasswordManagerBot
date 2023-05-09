FROM golang:alpine AS builder
WORKDIR /go/src/app
COPY . .
RUN go mod download
RUN go build -o botserver  cmd/main.go

FROM alpine
WORKDIR /app
COPY --from=builder /go/src/app /app/
CMD [ "./botserver" ]


