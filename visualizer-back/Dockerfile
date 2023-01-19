FROM golang:latest

# Set the working directory
WORKDIR /visualizer-back

# Clone the repository
RUN git clone https://github.com/Javtor/challenge-zinc.git .

RUN git fetch

RUN git checkout docker

# Go to the appropiate folder
WORKDIR ./visualizer-back

RUN go mod tidy

# Build the Go file
RUN go build -o main .

EXPOSE 3000

# Run the binary
CMD ["./main", "./enron_mail_20110402"]
