# APP Dockerfile
# For more control, you can copy and build manually
FROM golang:alpine

LABEL Name=userservice Version=0.0.1

RUN mkdir -p /go/src \
    && mkdir -p /go/bin \
    && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/gitlab.com/salismt/microservice-pattern-user-service

ADD . $GOPATH/src/gitlab.com/salismt/microservice-pattern-user-service
WORKDIR $GOPATH/src/gitlab.com/salismt/microservice-pattern-user-service

RUN go build -o main .
CMD ["/go/src/gitlab.com/salismt/microservice-pattern-user-service/main"]

EXPOSE 3000

