FROM golang:latest

WORKDIR /app

COPY . .

RUN rm -rf payments.db && go build -o main .

EXPOSE 8080

CMD ["./main"]