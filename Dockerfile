FROM golang:1.16
ENV GOOS=linux \
    CGO_ENABLED=0 \
    GOARCH=amd64

WORKDIR /usr/src/app
COPY . .
RUN go mod tidy
RUN go build -o main .
ENTRYPOINT ["./main"]
#ENTRYPOINT ["sleep","9999"]: