# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| latest  | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in cursor-tool, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

### How to Report

1. Go to the [Security tab](https://github.com/WarnetBes/cursor-tool/security/advisories/new) of this repository
2. Click "Report a vulnerability"
3. Fill in the details of the vulnerability

Alternatively, you can contact the maintainer directly through GitHub.

### What to Include

- Description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

### Response Timeline

- **Acknowledgement**: Within 48 hours
- **Status Update**: Within 7 days
- **Fix Release**: As soon as possible, typically within 30 days

## Security Considerations

cursor-tool requires elevated permissions on some platforms to:
- Read and write `storage.json` in application data directories
- Write to Windows Registry (Windows only)

Always download releases from the [official GitHub releases page](https://github.com/WarnetBes/cursor-tool/releases) and verify checksums.
