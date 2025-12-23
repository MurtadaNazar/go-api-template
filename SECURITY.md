# Security Policy

## Supported Versions
We actively maintain and monitor the latest stable release of this repository. Security updates are provided for the latest version only.

## Reporting a Vulnerability
If you discover a security vulnerability, please **do not open a public issue**. Instead, contact us privately via email: **mkm9284@gmail.com**.

Please include:
- A description of the vulnerability
- Steps to reproduce (if applicable)
- Potential impact of the vulnerability
- Your name and contact information (optional but preferred)

We will respond within 48 hours and work with you to resolve the issue responsibly.

## Response Timeline
- **Initial Response**: Within 48 hours
- **Status Update**: Within 1 week
- **Resolution**: We aim to resolve critical vulnerabilities within 30 days

## Responsible Disclosure
We ask that you:
- Do not publicly disclose the vulnerability until we've released a fix
- Do not access other users' data or systems
- Act in good faith and with the spirit of responsible disclosure

## Security Best Practices
- Keep dependencies updated (`go get -u` and review `go.mod`)
- Avoid committing sensitive information (use `.env` and `.gitignore`)
- Use code reviews for all changes
- Report security issues responsibly
- Follow OWASP security guidelines for authentication and data handling
