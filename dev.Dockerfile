FROM golang:1.17.2

# Set the Current Working Directory inside the container
WORKDIR /stripe-eboekhouden-go

COPY . .

RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]
RUN ["go", "install", "github.com/hooklift/gowsdl/...@latest"]

ENTRYPOINT gowsdl -o ./internal/push/soap/templates/gowsdl/test.go ./internal/push/soap/templates/test.xml && CompileDaemon -log-prefix=false -exclude-dir=.git -build="go build -o stripe-eboekhouden-go ./cmd/api" -command="./stripe-eboekhouden-go"
