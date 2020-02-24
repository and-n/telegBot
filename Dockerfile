FROM golang:1.13.8
#FROM golang:latest

WORKDIR /telegBot

COPY go.mod starter.go ./
COPY src ./src

#ENV GOPATH=./telegBot

RUN go mod download

#COPY . .

RUN go build -o main .

#EXPOSE 8080

CMD ["./main"]

