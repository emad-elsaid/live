FROM golang:1.16
COPY . /app
WORKDIR /app
RUN go build -o live *.go
CMD ./live
