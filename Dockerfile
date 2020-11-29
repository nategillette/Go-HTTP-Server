FROM golang
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go get github.com/jamespearly/loggly
RUN go get github.com/gorilla/mux
RUN go build -o main .
CMD ["./main"]
EXPOSE 8080