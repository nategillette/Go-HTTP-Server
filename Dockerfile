FROM golang:latest AS build

# Copy source
WORKDIR /app
COPY . .

# Get required modules (assumes packages have been added to ./vendor)
RUN go get github.com/jamespearly/loggly
RUN go get github.com/aws/aws-sdk-go/aws/session
RUN go get github.com/aws/aws-sdk-go/service/dynamodb
RUN go get github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute
RUN go get github.com/gorilla/mux

# Build a statically-linked Go binary for Linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -o main .

# New build phase -- create binary-only image
FROM alpine:latest

# Add support for HTTPS
RUN apk update && \
    apk upgrade && \
    apk add ca-certificates

WORKDIR /

# Copy files from previous build container
COPY --from=build /app/main ./

# Add environment variables
#ENV LOGGLY_TOKEN = 
ENV AWS_ACCESS_KEY = AKIA34XNLPJYFUGRBTUO
ENV AWS_SECRET_ACCESS_KEY = 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

# Check results
RUN env && pwd && find .

# Start the application
CMD ["./main"]

# Start the Go app build
#EXPOSE 8080

