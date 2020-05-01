# Using ubuntu because the GoLang docker image is unusually big
FROM ubuntu:18.04

# Setting up some stuff
RUN apt-get update
RUN apt-get install -y libssl1.0.0 wget apt-utils lsb-release curl

# Installing GoLang 1.14.1
ENV GOLANG_VERSION 1.14.1

RUN curl -sSL https://storage.googleapis.com/golang/go$GOLANG_VERSION.linux-amd64.tar.gz | tar -v -C /usr/local -xz
ENV PATH /usr/local/go/bin:$PATH

RUN mkdir -p /go/src /go/bin && chmod -R 777 /go
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

# Set current working directory inside the container
WORKDIR /app

# Download all dependencies. 
COPY go.mod go.sum ./
RUN go mod download

# Copying source code
COPY . .

# Build
RUN go build -o main .

EXPOSE 8080

CMD ["./main"]