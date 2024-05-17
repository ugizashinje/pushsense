FROM golang:1.21-bookworm
WORKDIR /app

COPY ./main /app/main
RUN mkdir /app/config


CMD ["./main"]




