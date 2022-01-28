FROM dependencies AS builder
# Copy the application source code.
COPY ./ /go/src/app/retrievedataservice

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -o /go/bin/retrievedataservice /go/src/app/retrievedataservice/

ENTRYPOINT [ "/go/bin/retrievedataservice" ]

FROM alpine:latest
COPY --from=builder /go/bin/retrievedataservice /bin/retrievedataservice
CMD [ "/bin/retrievedataservice" ]
