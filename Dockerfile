FROM golang:1.20-alpine

EXPOSE 33333
WORKDIR /app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN mkdir -p /usr/local/bin/
RUN go build -v -o /usr/local/bin/ ./...

ENV WAIT_VERSION 2.7.2
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

CMD ["cmd"]
#EXPOSE 33333/tcp

