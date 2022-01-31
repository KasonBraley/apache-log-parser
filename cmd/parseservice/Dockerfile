FROM kbraley/log-parser-dependencies AS builder
# Copy the application source code.
COPY ./ /go/src/app/parseservice

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -o /go/bin/parseservice /go/src/app/parseservice/

ENTRYPOINT [ "/go/bin/parseservice" ]

FROM alpine:latest
COPY --from=builder /go/bin/parseservice /bin/parseservice
CMD [ "/bin/parseservice" ]
