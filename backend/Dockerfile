# Use the official Go image as the base image
FROM golang:1.23

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files to the workspace
COPY go.mod go.sum ./

# Download all the dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code to the workspace
COPY . .

# Build the Go application
RUN go build -o main .

# Set the environment variable for the API key (this is optional, you can also set it directly in your code)
# ENV CHUNGYAK_INFO_API_KEY  

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
