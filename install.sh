#!/bin/bash
set -e

# pgmigrate installer
# Usage: curl -fsSL https://raw.githubusercontent.com/matroidbe/pgmigrate/main/install.sh | bash

REPO="github.com/matroidbe/pgmigrate"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
CLONE_DIR=$(mktemp -d)

cleanup() {
    rm -rf "$CLONE_DIR"
}
trap cleanup EXIT

echo "Installing pgmigrate..."

# Check for Go
if ! command -v go &> /dev/null; then
    echo "Error: Go is required but not installed."
    echo "Install Go from https://go.dev/dl/"
    exit 1
fi

# Clone and build
echo "Cloning $REPO..."
git clone --depth 1 "https://$REPO.git" "$CLONE_DIR" 2>/dev/null

echo "Building pgmigrate..."
cd "$CLONE_DIR"

# Get version from git tag or commit
VERSION=$(git describe --tags --always 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "")

go build -ldflags "-X github.com/matroidbe/pgmigrate/internal/cmd.Version=$VERSION -X github.com/matroidbe/pgmigrate/internal/cmd.GitCommit=$COMMIT" -o pgmigrate ./cmd/pgmigrate

# Install
echo "Installing to $INSTALL_DIR..."
mkdir -p "$INSTALL_DIR"
mv pgmigrate "$INSTALL_DIR/"

# Check PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo ""
    echo "Add $INSTALL_DIR to your PATH:"
    echo ""
    if [[ -f "$HOME/.zshrc" ]]; then
        echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.zshrc && source ~/.zshrc"
    else
        echo "  echo 'export PATH=\"$INSTALL_DIR:\$PATH\"' >> ~/.bashrc && source ~/.bashrc"
    fi
    echo ""
fi

echo ""
echo "Successfully installed pgmigrate $VERSION!"
echo ""
echo "Quick start:"
echo "  export DATABASE_URL='postgres://user:pass@localhost/mydb'"
echo "  pgmigrate init        # Create schema.yaml template"
echo "  pgmigrate plan        # Preview changes"
echo "  pgmigrate apply       # Apply changes"
