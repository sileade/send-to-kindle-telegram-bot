# Optimized Dockerfile for send-to-kindle-telegram-bot
FROM amd64/ubuntu:24.04

# Copy Go from official image
COPY --from=golang:1.22-bullseye /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"

# Set timezone
ENV TZ=Europe/Minsk
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Install dependencies in single layer to reduce image size
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends \
        git \
        wget \
        python3 \
        python-is-python3 \
        xz-utils \
        xdg-utils \
        ffmpeg \
        libsm6 \
        libxext6 \
        ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install Calibre and add to PATH
RUN wget -nv -O- https://download.calibre-ebook.com/linux-installer.sh | sh /dev/stdin && \
    ln -s /opt/calibre/ebook-convert /usr/local/bin/ebook-convert

# Verify Calibre installation
RUN ebook-convert --version

# Create working directory
WORKDIR /app

# Create files directory with proper permissions
RUN mkdir -p /app/files && chmod 777 /app/files

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o send-to-kindle-telegram-bot .

# Make binary executable
RUN chmod +x ./send-to-kindle-telegram-bot

# Set environment for Calibre
ENV PATH="/opt/calibre:${PATH}"

# Run the bot
CMD ["./send-to-kindle-telegram-bot"]