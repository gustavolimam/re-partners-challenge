# Use Golang image
FROM golang:1.24

# Set the working directory inside the container
WORKDIR /app

# Copy the Go modules files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Expose the port used by your application
EXPOSE 3000

# Build the Go app
RUN go build -o main .

# Command to run the app
CMD ["./main"]