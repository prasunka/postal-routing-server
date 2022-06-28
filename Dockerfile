FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build  -o server ./cmd/web/

EXPOSE 4000

CMD ["./server"]