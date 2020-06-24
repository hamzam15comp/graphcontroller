FROM gocv:latest 
RUN go get -u gocv.io/x/gocv
RUN go get -u github.com/hamzam15comp/vertex
COPY . .
EXPOSE 7000
CMD ./startup.sh

