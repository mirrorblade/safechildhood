FROM golang:1.21.1-alpine3.18 as builder

WORKDIR /builder

COPY . .

RUN apk update
RUN apk add make

RUN make build


FROM alpine:3.18

WORKDIR /app

COPY --from=builder /builder/bin ./bin
COPY --from=builder /builder/configs ./configs
COPY --from=builder /builder/resources ./resources
COPY --from=builder /builder/.env ./.env
COPY --from=builder /builder/key.json ./key.json

CMD ["./bin/safechildhood"]



