FROM golang:1.22-alpine

# Set the working directory to /app
WORKDIR /app

# Copy the go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and install any required Go dependencies
RUN go mod download

# Copy the entire source code to the working directory
COPY . .

#TODO fix the caching part
# ENV GOCACHE=/root/.cache/go-build
# RUN --mount=type=cache,target="/root/.cache/go-build" go build -o app
# COPY --from=builder /app/app .

# Build the Go application
RUN go build -o main .

# Expose the port specified by the PORT environment variable
EXPOSE 3001

# Set the entry point of the container to the executable
CMD ["./main"]