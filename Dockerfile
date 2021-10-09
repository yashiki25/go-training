FROM golang:1.17-buster

RUN mkdir /app
WORKDIR /app
ADD . /app

CMD ["go", "run", "main.go"]

RUN GO111MODULE=off go get -u github.com/oxequa/realize
EXPOSE 8080
CMD ["realize", "start", "--run"]