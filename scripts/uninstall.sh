#!/usr/bin/env bash
set -euo pipefail

BINARY="cursor-tool"
INSTALL_DIR="/usr/local/bin"

BIN_PATH="${INSTALL_DIR}/${BINARY}"

if [ ! -f "$BIN_PATH" ]; then
  echo "${BINARY} is not installed at ${BIN_PATH}"
  exit 0
fi

echo "Uninstalling ${BINARY}..."

if [ -w "$INSTALL_DIR" ]; then
  rm -f "$BIN_PATH"
else
  sudo rm -f "$BIN_PATH"
fi

echo "${BINARY} uninstalled successfully."
