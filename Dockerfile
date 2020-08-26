FROM golang:latest

COPY . $GOPATH/src/app/

WORKDIR $GOPATH/src/app/

# create .kube directory
RUN mkdir -p /root/.kube

# install dep
#RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

#RUN dep init && \
#    dep ensure --add github.com/gorilla/mux

RUN go build -o main .

EXPOSE 80

CMD ["./main"]
