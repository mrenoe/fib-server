FROM golang:alpine

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor
ENV APP_USER app

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .


WORKDIR /dist

RUN cp /build/main .

EXPOSE 8080

CMD ["/dist/main"]