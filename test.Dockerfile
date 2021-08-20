FROM golang:1.16.4

WORKDIR /go/src/app
COPY . /go/src/app
COPY ./config.test.yml /go/src/app/config.yml
COPY ./wait-for-it.sh /wait-for-it.sh
RUN chmod 555 /wait-for-it.sh

RUN go get -v -t -d ./...
RUN go build -v .

CMD /wait-for-it.sh db:5432 -t 30; go test -v -coverprofile=./reports/coverage.txt -covermode=atomic -coverpkg=./... ./...
