FROM golang:1.18-alpine as builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build --buildvcs=false ./cmd/zl
RUN go build --buildvcs=false ./cmd/zlsrv

FROM alpine:3.15 as prod
RUN apk add --no-cache ca-certificates
COPY --from=builder /app/zl /app/zlsrv /bin/
RUN mkdir /zettel
ENV ZLPATH=/zettel
ENV GIN_MODE=release
EXPOSE 8000
ENTRYPOINT ["/bin/zlsrv"]
