FROM golang:1.22-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install any required Go dependencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

# Expose the port specified by the PORT environment variable
EXPOSE 3001

# We're not pre-compiling the binary for development
# The actual run command is specified in docker-compose.yml
CMD ["go", "run", "main.go"]