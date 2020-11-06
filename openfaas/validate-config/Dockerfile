# Stage 0: Watchdog
FROM openfaas/of-watchdog:0.8.1 as watchdog

# Stage 1: Build
FROM golang:1.15.3-alpine3.12 as build

# Setup the go working directory and copy source
ENV CGO_ENABLED=0
RUN mkdir -p /go/src/app
WORKDIR /go/src/app
COPY . .

# Make native image executable
RUN chmod +x /go/src/app/target/sourcehawk

# Add user overrides to the root go.mod, which is the only place "replace" can be used
RUN cat function/GO_REPLACE.txt >> ./go.mod || exit 0

# Run gofmt
RUN test -z "$(gofmt -l $(find . -type f -name '*.go'))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }

# Configure GO Environment
ARG GO111MODULE="off"
ARG GOPROXY=""
ARG GOFLAGS=""

# Run the go tests for the function
WORKDIR /go/src/app/function
RUN go test ./... -cover

# Compile the go handler into a native image
WORKDIR /go/src/app
RUN go build --ldflags "-s -w" -a -installsuffix cgo -o handler .

# Stage 2: Final Image
FROM alpine:3.12.1

# Setup non-root app user
RUN addgroup -S app && adduser -S -g app --uid 1001 app \
    && mkdir -p /home/app \
    && chown app /home/app

# Initialize working directory and switch to app user
WORKDIR /home/app
USER app

# Copy required executables into image
COPY --from=build --chown=app /go/src/app/handler .
COPY --from=watchdog --chown=app /fwatchdog /usr/bin/fwatchdog
COPY --from=build --chown=app /go/src/app/target/sourcehawk /usr/bin/sourcehawk

# Setup the watchfog environment and run command
ENV fprocess="./handler" mode="http" upstream_url="http://127.0.0.1:8082"
CMD ["/usr/bin/fwatchdog"]