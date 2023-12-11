# Use the official Go image as the base image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Download dependencies and build the application
RUN go mod tidy
RUN go build -o server .

# Expose the port the application runs on
EXPOSE 8001

# Command to run the application
CMD ["./server"]
