FROM golang:1.13.8

WORKDIR /telegBot

COPY go.mod starter.go ./
COPY src ./src

ARG key
ENV API_KEY $key

RUN go mod download

#COPY . .

RUN go build -o main .

#EXPOSE 8080

CMD ["./main"]

