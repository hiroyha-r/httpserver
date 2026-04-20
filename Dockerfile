# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /httpserver .

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /httpserver /httpserver
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/httpserver"]
CMD ["--port", "8080"]
