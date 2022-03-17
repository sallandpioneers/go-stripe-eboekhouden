FROM golang:1.17.8

# Set the Current Working Directory inside the container
WORKDIR /go-stripe-eboekhouden

COPY . .

RUN ["go", "install", "github.com/hooklift/gowsdl/...@master"]
RUN ["go", "install", "github.com/githubnemo/CompileDaemon@latest"]
RUN ["go", "install", "github.com/kyleconroy/sqlc/cmd/sqlc@latest"]

ENTRYPOINT cd ./internal/push/soap && gowsdl -o eboekhouden.go -p generated ./templates/eboekhouden.wsdl && cd ../../../ && cd ./internal/storage/mysql/queries && sqlc generate && cd ../../../../ && CompileDaemon -log-prefix=false -exclude-dir=.git -build="go build -o go-stripe-eboekhouden ." -command="./go-stripe-eboekhouden serve"
