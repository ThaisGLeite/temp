# Client Dockerfile
FROM golang:1.20

WORKDIR /app

COPY word-of-wisdom-client.go .

RUN go mod init client && go build -o ./client

CMD [ "./client" ]
