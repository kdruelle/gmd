#!/usr/bin/env sh
set -e

REPO="kdruelle/gmd"
BINARY="gmd"
INSTALL_DIR="/usr/local/bin"

echo "→ Détection OS/ARCH…"
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64 | arm64) ARCH="arm64" ;;
esac

echo "→ Plateforme : $OS-$ARCH"

echo "→ Version : récupération…"
TAG=$(curl -s https://api.github.com/repos/$REPO/releases/latest | grep tag_name | cut -d '"' -f 4)

[ -z "$TAG" ] && echo "Impossible de trouver la version." && exit 1

echo "→ Version : ${TAG}"

ARCHIVE="${BINARY}_${TAG}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$TAG/$ARCHIVE"
CHECKSUM_URL="https://github.com/$REPO/releases/download/$TAG/checksums.txt"

TMP=$(mktemp -d)
cd "$TMP"

echo "→ Téléchargement du binaire…"
curl -sSfLO "$URL"

echo "→ Téléchargement du checksum…"
curl -sSfLO "$CHECKSUM_URL"

echo "→ Vérification SHA256…"
grep "$ARCHIVE" checksums.txt | sha256sum -c -

echo "→ Extraction…"
tar -xzf "$ARCHIVE"

echo "→ Installation…"
sudo mv "$BINARY" "$INSTALL_DIR/$BINARY"
sudo chmod +x "$INSTALL_DIR/$BINARY"

echo "✔ Installé : $BINARY"