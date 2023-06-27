# Server Dockerfile
FROM golang:1.20

WORKDIR /app

COPY word-of-wisdom-server.go .

RUN go mod init server && go build -o ./server

EXPOSE 8080

CMD [ "./server" ]
