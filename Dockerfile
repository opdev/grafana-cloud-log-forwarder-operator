# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager main.go

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM registry.access.redhat.com/ubi8/ubi-minimal:latest

LABEL \
    com.redhat.component="grafana-cloud-log-forwarder-operator" \
    version="v0.0.1" \
    name="grafana-cloud-log-forwarder-operator" \
    License="Apache-2.0" \
    io.k8s.display-name="grafana-cloud-log-forwarder-operator bundle" \
    io.k8s.description="bundle for the grafana-cloud-log-forwarder-operator" \
    summary="This is the bundle for the grafana-cloud-log-forwarder-operator" \
    maintainer="Grafana <support@grafana.com>" \
    vendor="Grafana Labs" \
    release="v0.0.1" \
    description="Grafana Cloud is a completeobservability stack for metrics, logs, and traces that's tightly integrated with Grafana. Leverage the best open source observability software without the overhead of installing, maintaining, and scaling your observability stack."
    
COPY licenses /licenses

WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
