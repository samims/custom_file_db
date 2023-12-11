# Start from a Golang base image with the desired version
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the local package files to the container's workspace
COPY . .

# Download dependencies if any are required
# RUN go mod download

# Build the Go application
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
