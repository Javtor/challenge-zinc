FROM golang:latest

# Set the working directory
WORKDIR /indexer

# Clone the repository
RUN git clone https://github.com/Javtor/challenge-zinc.git .

RUN git fetch

RUN git checkout docker

# Download the data.tgz
RUN wget http://www.cs.cmu.edu/~enron/enron_mail_20110402.tgz

# Extract the data
RUN tar -xvzf enron_mail_20110402.tgz

# remove tgz
RUN rm enron_mail_20110402.tgz

# Build the Go file
RUN go build -o main ./indexer/indexer.go

# Run the binary
CMD ["./main", "./enron_mail_20110402"]
