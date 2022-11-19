FROM golang:1.18.0-alpine3.15 AS builder

RUN apk add --no-cache \
  build-base \
  make \
  curl

WORKDIR /workspace
RUN mkdir _out
COPY cmd cmd
COPY imaginer imaginer
COPY logger logger
COPY model model
COPY server server

COPY go.mod go.sum Makefile ./
RUN touch LOCAL_ENV
RUN make clean-compile

FROM scratch
COPY --from=builder /workspace/_out/fmiyrm /usr/local/bin/fmiyrm

ENTRYPOINT ["fmiyrm"]
