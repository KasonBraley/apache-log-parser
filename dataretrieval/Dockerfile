FROM golang:1.17-alpine AS builder

# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . ./

# Build the application.
RUN go build -o /go/bin/dataretrieval ./cmd/app

# build a small image
FROM alpine:3.15
LABEL language="go"

# copy the executable
COPY --from=builder /go/bin/dataretrieval /dataretrieval

# run app
ENTRYPOINT ["./dataretrieval"]
