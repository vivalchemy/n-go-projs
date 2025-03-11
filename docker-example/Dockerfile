FROM golang:alpine AS base

WORKDIR /app

COPY go.* .

RUN go mod download

COPY . .

RUN go build -o /go/bin/app .

EXPOSE 8080

ENTRYPOINT ["/go/bin/app"]

FROM scratch AS prod

COPY --from=base /go/bin/app /app

ENTRYPOINT [ "/app" ]
