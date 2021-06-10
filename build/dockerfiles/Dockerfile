# Copyright (c) 2021 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   Red Hat, Inc. - initial API and implementation
#
# https://access.redhat.com/containers/?tab=tags#/registry.access.redhat.com/ubi8-minimal
FROM registry.access.redhat.com/ubi8-minimal:8.4 as builder
RUN microdnf install -y golang && \
  go version

USER root
WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
#COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . .

# Build
RUN export ARCH="$(uname -m)" && if [[ ${ARCH} == "x86_64" ]]; then export ARCH="amd64"; elif [[ ${ARCH} == "aarch64" ]]; then export ARCH="arm64"; fi && \
  CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} GO111MODULE=on go build \
  -a -o _output/bin/header-rewrite-proxy \
  -gcflags all=-trimpath=/ \
  -asmflags all=-trimpath=/ \
  cmd/header-rewrite-proxy/main.go

FROM scratch
COPY --from=builder /workspace/_output/bin/header-rewrite-proxy /

ENTRYPOINT ["/header-rewrite-proxy"]
