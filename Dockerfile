FROM golang:alpine as builder

RUN apk update && apk add gcc libc-dev

WORKDIR /source
ADD . /source/.
ARG VERSION

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s -X main.version='${VERSION}'" -o grade

FROM scratch

COPY --from=builder /source/grade /

ENTRYPOINT [ "/grade" ]
CMD ["--help"]
