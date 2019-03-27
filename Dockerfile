FROM golang:1.11

WORKDIR $GOPATH/src/github.com/tenzi-muraru/paymentApp
COPY . .

# Download all the dependencies
RUN go get -d -v ./...
# Install the package
RUN go install -v ./...

EXPOSE 8080
# Run the executable
CMD ["paymentApp"]