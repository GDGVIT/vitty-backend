FROM golang:1.19

WORKDIR /usr/src/app

COPY ./vitty-backend-gofiber/go.mod ./vitty-backend-gofiber/go.sum ./

RUN go mod download && go mod verify

COPY ./vitty-backend-gofiber .

RUN go build -v -o bin/vitty-backend-gofiber .

CMD ["./bin/vitty-backend-gofiber"]