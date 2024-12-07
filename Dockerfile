# Builder
FROM golang:1.23.3 AS builder

RUN useradd -u 1000 -g 65534 pismo

ARG importPath
ARG pkg

WORKDIR /go/src/${importPath}

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
     -o app $pkg

# Runner
FROM scratch

ARG importPath

COPY --from=builder /etc/passwd /etc/passwd

USER pismo

COPY --from=builder /go/src/${importPath}/app app

EXPOSE 8080

CMD [ "/app" ]
