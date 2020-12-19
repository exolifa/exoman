FROM golang:alpine AS builder

# Set necessary environment variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .
WORKDIR /build/main 
# Build the application
RUN go build -o exoman


# Build a small image
FROM alpine:latest
# add TimeZone package 
RUN apk --no-cache add tzdata
# create the /opt/exolifa/bin directory
RUN mkdir -p /opt/exolifa/bin
# copy the compiled code to it
COPY --from=builder /build/main/exoman /opt/exolifa/bin
# set working directory to /opt/exolifa
WORKDIR /opt/exolifa
# create required directories
RUN mkdir -p data templates log
# load parameters and templates 
COPY dockerparams/ data/
COPY templates/ templates/


# Command to run
ENTRYPOINT ["/opt/exolifa/bin/exoman", "/opt/exolifa/data/parameters.json"] 