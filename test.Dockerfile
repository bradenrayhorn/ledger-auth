FROM golang:1.15.7

WORKDIR /go/src/app
COPY . /go/src/app
COPY ./config.test.yml /go/src/app/config.yml
COPY ./wait-for-it.sh /wait-for-it.sh
RUN chmod 555 /wait-for-it.sh

RUN go get -v -t -d ./...
RUN go build -v .

RUN openssl genrsa -out test_rsa 1024
RUN openssl rsa -in test_rsa -pubout > test_rsa.pub

CMD /wait-for-it.sh db:3306 -t 30; go test -v -coverprofile=./reports/coverage.txt -covermode=atomic -coverpkg=./... ./...
