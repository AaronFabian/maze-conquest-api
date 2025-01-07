FROM golang:latest

COPY . /usr/src/app

WORKDIR /usr/src/app

RUN go build -o main .

EXPOSE 8000

CMD ["./main"]

# docker build -t maze-conquest-api .
# docker run --name maze-conquest-api --publish 8000:8000 maze-conquest-api