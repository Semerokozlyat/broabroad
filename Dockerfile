# syntax=docker/dockerfile:1

# Build the application
FROM golang:1.19 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
ADD . /app
# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /broabroad /app/cmd/broabroad/*

# Run tests
FROM build-stage AS run-tests-stage
RUN go test -v ./...

# Deploy app
FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /broabroad /broabroad
EXPOSE 8080
USER nonroot:nonroot
# Run the binary
ENTRYPOINT ["/broabroad"]