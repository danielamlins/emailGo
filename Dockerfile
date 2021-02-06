############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
RUN apk update && apk add --no-cache bash
# Create appuser.
RUN adduser -D -g '' appuser
WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .
# Fetch dependencies.
# Using go get.
#RUN go get -d -v
# Using go mod.
RUN go mod download
RUN go mod verify
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/api
RUN mv sendgrid.env /go/bin/sendgrid.env
############################
# STEP 2 build a small image
############################
FROM gcr.io/distroless/base-debian10@sha256:abe4b6cd34fed3ade2e89ed1f2ce75ddab023ea0d583206cfa4f960b74572c67
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy our static executable.
COPY --from=builder /go/bin/api /go/bin/api
COPY --from=builder /go/bin/sendgrid.env /go/bin/sendgrid.env
# Copy other CA stuff
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Use an unprivileged user.
USER appuser
# Run the api binary.
ENTRYPOINT ["/go/bin/api"]
