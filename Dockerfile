# Build the manager binary
FROM golang:1.10.3 as builder

# Copy in the go src
WORKDIR /go/src/github.com/jecho/ksops-test
COPY pkg/    pkg/
COPY cmd/    cmd/
COPY vendor/ vendor/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager github.com/jecho/ksops-test/cmd/manager

# Copy the controller-manager into a thin image
FROM alpine:3.8
RUN apk add gnupg && apk add git
WORKDIR /root/
RUN git clone https://github.com/mozilla/sops.git
WORKDIR /root/sops
RUN gpg --import pgp/sops_functional_tests_key.asc
#   RUN gpg --import pgp/some_mount_path
WORKDIR /root/
COPY --from=builder /go/src/github.com/jecho/ksops-test/manager .
ENTRYPOINT ["./manager"]