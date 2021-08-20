#!/bin/bash -e
cd "$(dirname $0)"
PATH=$HOME/go/bin:$PATH
unset GOPATH
export GO111MODULE=on

if ! type -p goveralls; then
  echo go get github.com/mattn/goveralls
  go get github.com/mattn/goveralls
fi

echo logrotate...
go test -v -covermode=count -coverprofile=logrotate.out .
go tool cover -func=logrotate.out
[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=logrotate.out -service=travis-ci -repotoken $COVERALLS_TOKEN

echo safe...
go test -v -covermode=count -coverprofile=safe.out ./safe
go tool cover -func=safe.out
[ -z "$COVERALLS_TOKEN" ] || goveralls -coverprofile=safe.out -service=travis-ci -repotoken $COVERALLS_TOKEN

# check that non Linux builds can succeed
GOOS=windows go build ./...
