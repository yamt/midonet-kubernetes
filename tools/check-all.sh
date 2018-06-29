#! /bin/sh

set -x
set -e

cd $(dirname $0)/..

./tools/check-copyrights.sh

dep ensure
go list ./...|grep -v github.com/midonet/midonet-kubernetes/pkg/client|xargs go vet
GOOS=linux GOARCH=arm64 go build ./...
GOOS=darwin GOARCH=amd64 go build ./...
