FROM kbraley/log-parser-dependencies AS builder
# Copy the application source code.
COPY ./ /go/src/app/logservice

# Build the application.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
go build -o /go/bin/logservice /go/src/app/logservice/

ENTRYPOINT [ "/go/bin/logservice" ]

FROM alpine:latest
COPY --from=builder /go/bin/logservice /bin/logservice
CMD [ "/bin/logservice" ]
