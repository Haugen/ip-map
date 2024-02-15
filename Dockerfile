FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

RUN go mod download && go mod verify

COPY . ./

CMD ["air"]
