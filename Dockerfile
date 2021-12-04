FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY src/ ./

RUN go build -o /ms-transmission

EXPOSE 8080

CMD [ "/ms-transmission" ]