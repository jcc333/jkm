# Dockerfile
FROM golang:1.23-alpine

# Install SSH server
RUN apk add --no-cache openssh bash

# Set up SSH
RUN ssh-keygen -A && \
    echo "root:password" | chpasswd && \
    adduser -D -s /bin/bash app && \
    echo "app:password" | chpasswd

# Set up SSH config
RUN echo "PermitRootLogin yes" >> /etc/ssh/sshd_config && \
    echo "PasswordAuthentication yes" >> /etc/ssh/sshd_config

# Create app directory
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/jkm ./cmd/jkm/main.go

RUN echo '#!/bin/bash\n\
    clear\n\
    exec /usr/local/bin/jkm' > /home/app/.bashrc && \
    chown app:app /home/app/.bashrc

COPY README.md /home/app/README.md
COPY welcome.md /etc/motd
COPY .env.biz /home/app/.env

EXPOSE 22

# Run SSH server instead of directly running the app
CMD ["/usr/sbin/sshd", "-D"]
