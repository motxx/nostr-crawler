FROM golang:1.22.2-alpine3.19 as builder

ENV ROOT=/go/src/app
WORKDIR ${ROOT}

RUN apk update && apk add git
COPY go.mod go.sum ./
RUN go mod download

COPY . ${ROOT}
RUN CGO_ENABLED=0 GOOS=linux go build -o $ROOT/nostr-crawler


FROM scratch as prod

ENV ROOT=/go/src/app
WORKDIR ${ROOT}
COPY --from=builder ${ROOT}/nostr-crawler ${ROOT}
COPY --from=builder ${ROOT}/run-crawler.sh ${ROOT}

CMD ["./run-crawler.sh"]
