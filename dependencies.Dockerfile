FROM golang:1.17
# Add the module files and download dependencies.
COPY ./go.mod /go/src/app/go.mod
COPY ./go.sum /go/src/app/go.sum
WORKDIR /go/src/app
RUN go mod download
# Add the shared packages.
COPY ./logger /go/src/app/logger
COPY ./parse /go/src/app/parse
COPY ./registry /go/src/app/registry
COPY ./retrieve /go/src/app/retrieve
COPY ./service /go/src/app/service
