FROM golang:1.19.5-bullseye

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build
RUN go build -o /razorblade cmd/shortener/main.go

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# Run
CMD [ "/razorblade" ]
