# Builder stage
FROM golang:1.20.4-buster as builder

# Set working directory
WORKDIR /go/src/gitlab.com/tokend/course-certificates/ccp

# Install necessary packages
RUN apt-get update && apt-get install -y gcc g++ make git
RUN apt-get install -y imagemagick libmagickwand-dev

# Copy project files
COPY . .

# Build the application
RUN export CGO_ENABLED=1 && export GO111MODULE=off && export GOOS=linux && go build -o /usr/local/bin/ccp

# Service stage
FROM alpine:3.14.6 as service

# Copy the built binary from the builder stage
COPY --from=builder /usr/local/bin/ccp /usr/local/bin/ccp

# Set the entrypoint command
ENTRYPOINT ["ccp"]

# Install necessary runtime dependencies
RUN apk add --no-cache ca-certificates
