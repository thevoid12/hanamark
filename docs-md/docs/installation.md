---
created_on: 2024-12-26
---

# Installation

Hanamark is distributed as a single binary with no dependencies.

## Download Binary (Recommended)

Download the pre-built binary for your platform from the [GitHub Releases](https://github.com/thevoid12/hanamark/releases) page.

| Platform | Download |
| --- | --- |
| macOS (Intel) | `hanamark_Darwin_x86_64.tar.gz` |
| macOS (Apple Silicon) | `hanamark_Darwin_arm64.tar.gz` |
| Linux (64-bit) | `hanamark_Linux_x86_64.tar.gz` |
| Linux (ARM 64-bit) | `hanamark_Linux_arm64.tar.gz` |
| Windows (64-bit) | `hanamark_Windows_x86_64.zip` |
| Windows (ARM 64-bit) | `hanamark_Windows_arm64.zip` |

download the binary and use gunzip to extract the tarball and then get the executable

```bash
# Example: Download and make executable (macOS/Linux)
chmod +x hanamark-darwin-arm64
mv hanamark-darwin-arm64 hanamark

# Run
./hanamark init
```

### Setup

```bash
# Download and make executable (macOS/Linux)
chmod +x hanamark-darwin-arm64
mv hanamark-darwin-arm64 hanamark

# Verify installation
./hanamark -version
```

> **No dependencies required** - the binary is self-contained.

## Build from Source

Only needed if you want to build from source or contribute to development.

### Prerequisites

- Go 1.20 or higher

### Build Steps

```bash
# Clone repository
git clone https://github.com/thevoid12/hanamark.git
cd hanamark

# Build binary
go build -o hanamark .
```

Or using Make:

```bash
make build
```

## Verify Installation

```bash
./hanamark -version
```

This should display the current version of Hanamark.
