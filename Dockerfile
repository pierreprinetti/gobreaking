FROM golang:1.19

COPY ./ /src

RUN cd /src && go build -o /gobreaking .

ENTRYPOINT ["/gobreaking"]
