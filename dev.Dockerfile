FROM golang:1.17.2

# Set the Current Working Directory inside the container
WORKDIR /stripe-eboekhouden-go

COPY . .

RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]

ENTRYPOINT CompileDaemon -log-prefix=false -exclude-dir=.git -build="go build -o stripe-eboekhouden-go ./cmd/api" -command="./stripe-eboekhouden-go"
