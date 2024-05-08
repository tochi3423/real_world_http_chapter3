FROM golang:1.16

WORKDIR /app 

COPY go.* ./ 
RUN go mod download 

COPY . . 

RUN CGO_ENABLED=0 go build -o /go/bin/server 

EXPOSE 18888 

CMD ["/go/bin/server"]