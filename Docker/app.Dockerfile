# syntax=docker/dockerfile:1
FROM golang:latest

# Creates an app directory to hold your app’s source code
WORKDIR /app

# Copies everything from your root directory into /app
COPY . .

# Installs Go dependencies
RUN go mod download

# Builds your app with optional configuration
RUN go build ./cmd/api

# Tells Docker which network port your container listens on
EXPOSE 8080

ENV DB_HOST="db"
ENV DB_USER="postgres"
ENV DB_PASS="postgres"
ENV DB_PORT="5432"
ENV WEB_PORT="8080"
ENV DB_NAME="library"
# Specifies the executable command that runs when the container starts
CMD /app/api