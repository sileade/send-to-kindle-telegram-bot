# Optimized Dockerfile for send-to-kindle-telegram-bot
FROM amd64/ubuntu:24.04

# Copy Go from official image
COPY --from=golang:1.22-bullseye /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

# Set timezone
ENV TZ=Europe/Minsk
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Set non-interactive mode for apt
ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies including Calibre from Ubuntu repos
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
        git \
        wget \
        python3 \
        python-is-python3 \
        calibre \
        ffmpeg \
        libsm6 \
        libxext6 \
        ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Verify Calibre installation
RUN which ebook-convert && ebook-convert --version

# Create working directory
WORKDIR /app

# Create files directory and symlink (bot expects /files/)
RUN mkdir -p /app/files && \
    chmod 777 /app/files && \
    ln -s /app/files /files

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o send-to-kindle-telegram-bot .

# Make binary executable
RUN chmod +x ./send-to-kindle-telegram-bot

# Run the bot
CMD ["./send-to-kindle-telegram-bot"]