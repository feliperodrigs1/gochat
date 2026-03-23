FROM golang:1.26.1

WORKDIR /app

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD ["air"]