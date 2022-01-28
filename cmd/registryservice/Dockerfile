FROM dependencies AS builder
# Copy the application source code.
COPY ./ /go/src/app/registryservice

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -o /go/bin/registryservice /go/src/app/registryservice/

ENTRYPOINT [ "/go/bin/registryservice" ]

FROM alpine:latest
COPY --from=builder /go/bin/registryservice /bin/registryservice
CMD [ "/bin/registryservice" ]
