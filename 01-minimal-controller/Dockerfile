FROM golang:1.26.1 AS builder
WORKDIR /workspace

# Build context must be repository root so go.mod/go.sum are available.
# Example: docker build -f 01-minimal-controller/Dockerfile -t <image> .
# Copy module files first for better build cache usage.
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build the step-01 controller binary.
COPY 01-minimal-controller ./01-minimal-controller
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -o manager ./01-minimal-controller

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager /manager
USER 65532:65532
ENTRYPOINT ["/manager"]
