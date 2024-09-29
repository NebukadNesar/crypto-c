# Step 1: Use the official Go image as the base image
FROM golang:1.18-alpine AS build

# Step 2: Set the working directory inside the container
WORKDIR /app

# Step 3: Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Step 4: Copy all the Go source files into the container
COPY . .

# Step 5: Build the Go app (it will automatically find all source files in the directory)
RUN go build -o main .

# Step 6: Use a minimal base image for running the app
FROM alpine:latest

# Step 7: Set the working directory inside the container
WORKDIR /app

# Step 8: Copy the binary from the build stage
COPY --from=build /app/main .

# Step 9: Expose the port the app runs on
EXPOSE 8080

# Step 10: Define the command to run the app
CMD ["./main"]
