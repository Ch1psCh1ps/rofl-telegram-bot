FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o geniemap .
RUN chmod +x geniemap

CMD ./geniemap