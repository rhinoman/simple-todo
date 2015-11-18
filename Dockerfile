FROM fedora:20
MAINTAINER noone

RUN  yum -y update && yum clean all
RUN  yum -y install couchdb && yum clean all

RUN  sed -e 's/^bind_address = .*$/bind_address = 0.0.0.0/' -i /etc/couchdb/default.ini

# Install Go
RUN yum -y groupinstall "Development Tools"; yum clean all
RUN mkdir /goroot && curl https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz | tar xzvf - -C /goroot --strip-components=1

# Setup environment
ENV GOROOT /goroot
ENV GOPATH /go
ENV PATH $PATH:$GOROOT/bin:$GOPATH/bin


ADD . /go/src/github.com/rhinoman/simple-todo

RUN go get github.com/rhinoman/couchdb-go
RUN go get github.com/twinj/uuid
RUN go install github.com/rhinoman/simple-todo

COPY ./start.sh /
EXPOSE 5984 8085

CMD ["/bin/sh", "/start.sh"]

