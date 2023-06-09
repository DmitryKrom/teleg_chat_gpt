# Base image
FROM golang:1.20
RUN apt-get -y update
RUN apt-get install -y ffmpeg

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Install dependencies
RUN go mod download

# Build the binary
RUN go build -o chatbot .

# Set the entry point
ENTRYPOINT ["/app/chatbot"]
