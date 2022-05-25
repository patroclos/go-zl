FROM golang:1.18-alpine as builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build --buildvcs=false ./cmd/zl
RUN go build --buildvcs=false ./cmd/zlsrv

FROM minlag/mermaid-cli:9.0.3 as prod
RUN npm install
USER root
RUN ln -s /home/mermaidcli/node_modules/.bin/mmdc /usr/local/bin/mmdc
RUN apk add --no-cache ca-certificates bash chromium
RUN ln -s $(which chromium) /usr/bin/chrome
COPY --from=builder /app/zl /app/zlsrv /bin/
COPY --from=builder /app/scripts /filters
RUN mkdir /zettel
ENV ZLPATH=/zettel
ENV GIN_MODE=release
ENV ZLSRV_FILTER__mermaid="/filters/mermaid-filter.sh"
EXPOSE 8000
ENTRYPOINT ["/bin/zlsrv"]
