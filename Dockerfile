FROM golang:1.19-alpine
WORKDIR /app/

COPY . .

RUN go mod download
RUN go build -o ./borsch-runner-service ./cmd/main.go

FROM docker:20.10.17-alpine3.16
WORKDIR /app/
COPY --from=0 /app/borsch-runner-service ./

ENTRYPOINT ./borsch-runner-service
