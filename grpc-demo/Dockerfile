FROM golang:alpine AS development

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

CMD ["sh", "-c", "air --build.cmd 'go build -o /go/bin/grpcdemo ./server' --build.bin '/go/bin/grpcdemo'"]

# Use scratch instead of alpine for a minimal final image
FROM scratch AS production

COPY --from=development /go/bin/grpcdemo /go/bin/grpcdemo

EXPOSE 8080

ENTRYPOINT ["/go/bin/grpcdemo"]
