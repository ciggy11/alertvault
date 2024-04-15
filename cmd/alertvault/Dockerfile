FROM golang:1.22

ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /alertvault

COPY       alertvault-amd64  .
COPY       config.yaml .
RUN        chmod +x /alertvault/alertvault-amd64

EXPOSE     8080
ENTRYPOINT ["./alertvault-amd64"]
CMD ["-config", "config.yaml"]


