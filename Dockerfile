FROM golang:1.22.4 as builder

ENV GOOS linux
ENV CGO_ENABLED 0

WORKDIR /go/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:3.16 as production

WORKDIR /go/src/app

RUN apk add --no-cache ca-certificates

COPY --from=builder /go/src/app/main /go/src/app/
COPY --from=builder /go/src/app/.env /go/src/app/

# ENV GOMAXPROCS 4
# ENV GOGC 50
# ENV GOTRACEBACK all

ENV PORT 8080

EXPOSE $PORT

CMD ["./main"]
