FROM golang:1.20 as build

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o testapp

EXPOSE 8080

CMD ["/app/testapp"]