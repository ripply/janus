#!/bin/sh

# Latest commit hash
GIT_SHA=`git rev-parse HEAD`

# If working copy has changes, append `-local` to hash
GIT_DIFF=`git diff -s --exit-code || echo "-local"`
GIT_REV=${GIT_SHA}${GIT_DIFF}

go build -ldflags "-X 'github.com/qtumproject/janus/pkg/params.GitSha=${GIT_REV}'" -o $GOPATH/bin $GOPATH/src/github.com/qtumproject/janus/...
