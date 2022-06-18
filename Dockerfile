FROM golang:1.17-alpine

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /app/discount

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the Go app
RUN go run main.go serve .


# This container exposes port 8080 to the outside world
EXPOSE 1337

# Run the binary program produced by `go install`
CMD ["serve"]