FROM golang:1.21

RUN mkdir /app
ADD . /app
WORKDIR /app/cmd/notifications-api

RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /notifications-api

EXPOSE 59191

CMD [ "/notifications-api" ]