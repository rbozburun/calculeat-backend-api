# Start from golang base image
FROM golang:1.21

# Specify that we now need to execute any commands in this directory.
WORKDIR /go/src/github.com/calculeat/main_rest_api

# Copy everything from this project into the filesystem of the container.
COPY . .

# Get packages
RUN go get -v

# Compile the binary exe for our app.
RUN go build -o main .


# Start the application.
CMD ["./main"]







