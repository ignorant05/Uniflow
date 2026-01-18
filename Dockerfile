# syntax=docker/dockerfile:1

#######################################
###			 Stage 1: test          ###
#######################################
FROM golang:1.25.4-alpine AS tester

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG BUILDPLATFORM
RUN if [ "$BUILDPLATFORM" = "linux/amd64" ]; then \
	--mount=type=cache,target=/root/.cache/go-build \
	go test -v ./...; \
	fi

#######################################
###			 Stage 2: Build         ###
#######################################
FROM  golang:1.25.4-alpine AS build 

LABEL org.opencontainers.image.title="Uniflow"
LABEL org.opencontainers.image.description="Unified workflow orchestration tool"
LABEL org.opencontainers.image.url="https://github.com/ignorant05/uniflow"
LABEL org.opencontainers.image.source="https://github.com/ignorant05/uniflow"
LABEL org.opencontainers.image.version="0.2.0"
LABEL org.opencontainers.image.licenses="MIT"

ENV CGO_ENABLED=0 \
	GOOS=linux \
	GOARCH=amd64

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download && \
	go mod verify

COPY . . 

ARG TARGETPLATFORM
RUN --mount=type=cache,target=/root/.cache/go-build \
	GOOS=$(echo $TARGETPLATFORM | cut -d'/' -f1) \
	GOARCH=$(echo $TARGETPLATFORM | cut -d'/' -f2) \
	go build -ldflags="-w -s" -trimpath -o uniflow .

#######################################
###			Stage 3: Runtime        ###
#######################################
FROM alpine:latest

RUN apk add --no-cache ca-certificates git

COPY --from=build /app/uniflow /usr/local/bin/uniflow

WORKDIR /workspace

ENTRYPOINT ["uniflow"]
CMD ["--help"]
