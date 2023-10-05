FROM golang:1.21.1-alpine3.18

WORKDIR /app

COPY . .

RUN mkdir bin

RUN go mod download && go build -o ./bin/safechildhood ./cmd/safechildhood/main.go

EXPOSE 8080

CMD [ "./bin/safechildhood" ]


