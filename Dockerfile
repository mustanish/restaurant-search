# Get golang's alpines image as the base image
FROM golang:alpine

# Set working directory inside container
WORKDIR /usr/src/search

# Copy .mod and .sum files to the working directory inside container
COPY go.mod go.sum ./

# Run go mod to download all dependancies
RUN go mod download

# Copy entire soucre code from host to the working directory inside container
COPY  ./ ./

# Get compiledaemon to rebuild and restart the application
RUN go get github.com/githubnemo/CompileDaemon

# Expose port 5000 to the host machine
EXPOSE 5000

# Configure compiledaemon to rebuild and restart the application
CMD CompileDaemon --build="go build search.go" --command=./search



