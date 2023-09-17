# Use an official Go runtime as a parent image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /books-authors

# Copy the local package files to the container's workspace
COPY . .

# Build the Go application inside the container
RUN go build -o books-authors

# Expose the port your application listens on
EXPOSE 8080

# Run your application
CMD ["./books-authors"]
