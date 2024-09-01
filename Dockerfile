FROM golang:1.23

RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 

RUN go build -o main -buildvcs=false

CMD ["/app/main"]