# Stage 1 : builder stage


FROM golang:latest AS builder

LABEL maintainer="Vrashabh Sontakke <vrashabhsontakke@outlook.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


# Stage 2 : runtime stage


FROM alpine:latest AS runtime

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary from builder stage to runtime stage
COPY --from=builder /app/main .

# Expose port 8080 to the outside world
EXPOSE 8080

# run executable binary
CMD ["./main"]