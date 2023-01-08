FROM golang:1.19-alpine as build

WORKDIR /app
COPY server server
COPY internal internal
# Fetch dependencies
COPY go.mod ./
RUN go mod download

# Build
COPY . ./
RUN CGO_ENABLED=0 go build ./server/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/main /app/ticker-beats-admin

EXPOSE 31034
ENTRYPOINT ["/app/ticker-beats-admin"]
