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
ENV APT_KEY_DONT_WARN_ON_DANGEROUS_USAGE=DEB_CONTROL_STRICT_LEVEL=no

# Update package lists with retry logic
RUN set -e && \
    for i in 1 2 3; do \
        echo "[INFO] Attempting apt-get update (attempt $i/3)..." && \
        apt-get update && break || \
        if [ $i -eq 3 ]; then echo "[ERROR] Failed to update apt after 3 attempts"; exit 1; fi && \
        sleep 10; \
    done && \
    echo "[INFO] apt-get update successful"

# Upgrade and install dependencies
RUN echo "[INFO] Installing dependencies..." && \
    apt-get upgrade -y --allow-unauthenticated && \
    apt-get install -y --no-install-recommends --allow-unauthenticated \
        git \
        wget \
        ca-certificates \
        python3 \
        python-is-python3 \
        calibre \
        ffmpeg \
        libsm6 \
        libxext6 && \
    echo "[INFO] Dependencies installed successfully" && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Verify Calibre installation
RUN echo "[INFO] Verifying Calibre..." && \
    which ebook-convert && \
    ebook-convert --version && \
    echo "[INFO] Calibre verification successful"

# Create working directory
WORKDIR /app

# Create files directory and symlink (bot expects /files/)
RUN echo "[INFO] Creating file directories..." && \
    mkdir -p /app/files && \
    chmod 777 /app/files && \
    ln -s /app/files /files && \
    echo "[INFO] File directories created"

# Copy go.mod and go.sum first for better caching
COPY go.mod go.sum ./
RUN echo "[INFO] Downloading Go dependencies..." && \
    go mod download && \
    echo "[INFO] Go dependencies downloaded"

# Copy source code
COPY . .

# Build the application
RUN echo "[INFO] Building application..." && \
    go build -o send-to-kindle-telegram-bot . && \
    echo "[INFO] Build completed successfully"

# Make binary executable
RUN chmod +x ./send-to-kindle-telegram-bot && \
    echo "[INFO] Dockerfile build complete"

# Run the bot
CMD ["./send-to-kindle-telegram-bot"]
