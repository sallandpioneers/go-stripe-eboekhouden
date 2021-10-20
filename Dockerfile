FROM golang:1.17.2 as build_base

# Set the Current Working Directory inside the container
WORKDIR /stripe-eboekhouden-go

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

FROM build_base as builder

WORKDIR /stripe-eboekhouden-go

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/stripe-eboekhouden-go ./cmd/api

######## Start a new stage from scratch #######
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /

COPY --from=builder /go/bin/stripe-eboekhouden-go .

EXPOSE 8080

CMD ["./stripe-eboekhouden-go"] 