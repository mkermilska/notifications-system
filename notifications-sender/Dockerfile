FROM golang:1.21

RUN mkdir /app
ADD . /app
WORKDIR /app/cmd/notifications-sender

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /notifications-sender

CMD [ "/notifications-sender" ]