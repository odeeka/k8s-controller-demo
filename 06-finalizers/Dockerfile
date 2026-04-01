FROM golang:1.26.1 AS builder
WORKDIR /workspace

# Build context must be repository root so go.mod/go.sum are available.
# Example: docker build -f 06-finalizers/Dockerfile -t <image> .
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build the step-06 controller binary.
COPY 06-finalizers ./06-finalizers
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -a -o manager ./06-finalizers

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/manager /manager
USER 65532:65532
ENTRYPOINT ["/manager"]
