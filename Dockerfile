# Final stage only - binary is pre-built by GoReleaser
FROM alpine:3.20
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY zte-cli /usr/local/bin/zte-cli

ENTRYPOINT ["zte-cli"]
CMD ["--help"]
