FROM golang:1.23.3

ARG importPath
ARG pkg

WORKDIR /go/src/${importPath}

COPY . .
