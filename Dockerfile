FROM hammoti/dockergo:basic 
WORKDIR /
RUN go get -u github.com/docker/docker/...
RUN go get -u github.com/hammoti/vertex
RUN go get -u github.com/gorilla/mux
COPY . .
CMD ["go", "run", "controller.go", "dockerfunc.go", "helper.go"]
