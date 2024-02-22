FROM golang:1.12

# Set destination for COPY
WORKDIR /app

# Download Go modules
#COPY go.mod go.sum ./
#RUN go mod download
RUN go mod init github.com/the-hidden-eye/Reverse-Proxy-Cache-Redis
#RUN go get github.com/go-redis/redis/v8
RUN go get github.com/go-redis/redis
#RUN go install https://github.com/go-redis/redis@latest
# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /dockercacheproxy

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
## https://docs.docker.com/reference/dockerfile/#expose
#EXPOSE 8080

# Run
#CMD ["/dockercacheproxy"]
CMD sh -c "export ; /dockercacheproxy 2>&1 |grep -e Init -e Ready -e Requested -e Creating -e panic"