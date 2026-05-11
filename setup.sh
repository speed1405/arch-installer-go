#!/usr/bin/env bash
# setup.sh – Prepare an Arch Linux live environment to build and run arch-installer-go.
# Must be run as root on a live Arch ISO.
set -euo pipefail

# ── Colour helpers ────────────────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BOLD='\033[1m'
NC='\033[0m'

info()  { echo -e "${GREEN}[INFO]${NC}  $*"; }
warn()  { echo -e "${YELLOW}[WARN]${NC}  $*"; }
error() { echo -e "${RED}[ERROR]${NC} $*" >&2; exit 1; }

# ── 1. Root check ─────────────────────────────────────────────────────────────
if [[ $EUID -ne 0 ]]; then
    error "This script must be run as root (use: sudo $0)"
fi
info "Running as root – OK"

# ── 2. Internet connectivity ──────────────────────────────────────────────────
info "Checking internet connectivity…"
if ! curl -s --max-time 5 https://archlinux.org > /dev/null 2>&1; then
    error "No internet connection detected. Connect to the network and try again."
fi
info "Internet connectivity – OK"

# ── 3. Sync system clock ──────────────────────────────────────────────────────
info "Enabling NTP time synchronisation…"
timedatectl set-ntp true
info "System clock sync – OK"

# ── 4. Refresh pacman keyring ─────────────────────────────────────────────────
info "Refreshing pacman keyring (needed for package verification)…"
pacman-key --init
pacman-key --populate archlinux
pacman -Sy --noconfirm archlinux-keyring
info "Keyring refresh – OK"

# ── 5. Required system tools ──────────────────────────────────────────────────
REQUIRED_TOOLS=(lsblk fdisk pacstrap arch-chroot genfstab)
MISSING=()

for tool in "${REQUIRED_TOOLS[@]}"; do
    if ! command -v "$tool" > /dev/null 2>&1; then
        MISSING+=("$tool")
    fi
done

if [[ ${#MISSING[@]} -gt 0 ]]; then
    warn "Missing tools: ${MISSING[*]}"
    info "Installing arch-install-scripts to provide missing tools…"
    pacman -S --noconfirm arch-install-scripts
fi

# Final verification
for tool in "${REQUIRED_TOOLS[@]}"; do
    command -v "$tool" > /dev/null 2>&1 \
        || error "Required tool '$tool' still not found after install attempt."
done
info "All required system tools present – OK"

# ── 6. Go toolchain ───────────────────────────────────────────────────────────
# Minimum required: Go 1.22 (major=1, minor≥22)
REQUIRED_GO_MAJOR=1
REQUIRED_GO_MINOR=22

# Returns 0 (true) if the installed Go version is below the requirement.
go_too_old() {
    local ver
    ver=$(go version | grep -oP '(?<=go)\d+\.\d+' | head -1)
    local major minor
    major=$(echo "$ver" | cut -d. -f1)
    minor=$(echo "$ver" | cut -d. -f2)
    if [[ $major -lt $REQUIRED_GO_MAJOR ]]; then
        return 0
    elif [[ $major -eq $REQUIRED_GO_MAJOR && $minor -lt $REQUIRED_GO_MINOR ]]; then
        return 0
    fi
    return 1
}

need_go=false
if ! command -v go > /dev/null 2>&1; then
    warn "Go not found – will install via pacman"
    need_go=true
elif go_too_old; then
    CURRENT_GO=$(go version | grep -oP '(?<=go)\d+\.\d+' | head -1)
    warn "Go $CURRENT_GO found, but >= ${REQUIRED_GO_MAJOR}.${REQUIRED_GO_MINOR} is required – will upgrade"
    need_go=true
else
    CURRENT_GO=$(go version | grep -oP '(?<=go)\d+\.\d+\.\d+' | head -1)
    info "Go $CURRENT_GO found – OK"
fi

if $need_go; then
    info "Installing Go via pacman…"
    pacman -S --noconfirm go
    info "Go $(go version | grep -oP '(?<=go)\d+\.\d+\.\d+' | head -1) installed – OK"
fi

# ── 7. Build the installer binary ────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY="${SCRIPT_DIR}/arch-installer"

info "Downloading Go module dependencies…"
(cd "$SCRIPT_DIR" && go mod download)

info "Building arch-installer…"
(cd "$SCRIPT_DIR" && go build -o "$BINARY" .)

if [[ ! -x "$BINARY" ]]; then
    error "Build succeeded but binary not found at $BINARY"
fi
info "Build complete – binary at ${BOLD}${BINARY}${NC}"

# ── Done ──────────────────────────────────────────────────────────────────────
echo ""
echo -e "${GREEN}${BOLD}Setup complete!${NC}"
echo -e "Run the installer with: ${BOLD}sudo ${BINARY}${NC}"
