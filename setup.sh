#!/bin/bash

# Function to check if a command exists
command_exists() {
  command -v "$1" >/dev/null 2>&1
}

# Function to install Go
install_go() {
  echo "Go is not installed. Installing Go..."
  GO_VERSION="1.22"
  OS="$(uname | tr '[:upper:]' '[:lower:]')"
  ARCH="$(uname -m)"
  
  if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
  elif [ "$ARCH" = "aarch64" ]; then
    ARCH="arm64"
  fi
  
  GO_TAR="go$GO_VERSION.$OS-$ARCH.tar.gz"
  GO_URL="https://golang.org/dl/$GO_TAR"
  
  wget $GO_URL
  sudo tar -C /usr/local -xzf $GO_TAR
  rm $GO_TAR
  
  export PATH=$PATH:/usr/local/go/bin
  echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.profile
  source ~/.profile
}

# Check if Go is installed
if command_exists go; then
  echo "Go is already installed."
else
  install_go
fi

# Set up project
echo "Setting up the project..."
cd "$(dirname "$0")"
go mod tidy

# Build the project
echo "Building the project..."
go build -o ./bin/cmd ./cmd

echo "Setup complete. You can now run the project with './bin/cmd start'"