FROM  golang as builder

WORKDIR /go/src/go-postgres-app/goapp/goapi

ENV GOPATH /go
ENV GO111MODULE on

COPY go.mod . 
COPY go.sum . 
RUN go mod download
COPY . . 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build 

FROM scratch 

COPY --from=builder /go/src/go-postgres-app/goapp/goapi /goapp/

EXPOSE 4444 
ENTRYPOINT [ "/goapp/goapi" ]