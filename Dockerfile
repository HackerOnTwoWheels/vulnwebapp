FROM golang:1.18

RUN mkdir $(go env GOPATH)/src/vulnwepapp/
WORKDIR $(go env GOPATH)/src/vulnwepapp/
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN cp vulndb.db  $(go env GOPATH)/bin/vulndb.db
RUN go build -v -o $(go env GOPATH)/bin/vulnwebapp ./...

CMD ["vulnwebapp"] 
