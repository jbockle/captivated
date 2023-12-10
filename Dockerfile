# example:
# docker run --env ASB_CONNECTION_STRING=$env:CAPTIVATED_ACS --env BLOB_CONNECTION_STRING=$env:CAPTIVATED_BCS -p 8080:8080 captivated:0

# Build the application from source
FROM golang:1.21.5-bullseye AS build-stage

COPY ./ /app
RUN ls -la /app/*

WORKDIR /app

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /captivated

# Run the tests in the container
# FROM build-stage AS run-test-stage
# RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /captivated /captivated

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/captivated"]