# Security Policy

## Supported Versions

We actively support the following versions with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 0.x.x   | :white_check_mark: |

## Reporting a Vulnerability

We take security vulnerabilities seriously. If you discover a security vulnerability, please follow these steps:

### 1. **Do Not** Create a Public Issue

Please do not report security vulnerabilities through public GitHub issues, discussions, or pull requests.

### 2. Report Privately

Send an email to **security@boldminds.tech** with the following information:

- **Subject**: Security Vulnerability in bold-minds/cond
- **Description**: Detailed description of the vulnerability
- **Steps to Reproduce**: Clear steps to reproduce the issue
- **Impact**: Potential impact and severity assessment
- **Suggested Fix**: If you have ideas for a fix (optional)

### 3. Response Timeline

- **Initial Response**: Within 48 hours
- **Status Update**: Within 7 days
- **Resolution**: Varies based on complexity, typically within 30 days

### 4. Disclosure Process

1. We will acknowledge receipt of your vulnerability report
2. We will investigate and validate the vulnerability
3. We will develop and test a fix
4. We will coordinate disclosure timing with you
5. We will release a security update
6. We will publicly acknowledge your responsible disclosure (if desired)

## Security Considerations

`cond` is a pure-computation library with a very small attack surface:

- **No network I/O.** `cond` does not make network calls.
- **No file I/O.** `cond` does not read or write files.
- **No reflection.** `cond` uses concrete type switches only.
- **No external dependencies.** `cond` is pure Go stdlib.
- **No mutation.** `cond` never modifies input values.

### Panic Contract

`If[T]` will **intentionally panic** when a branch argument is neither a `T` nor a `func() T`. This is by design: a mistyped branch is a programming error, and silently returning the zero value would hide bugs. The panic message identifies which branch failed its type assertion (`cond: type assertion failed for trueVal` / `cond: type assertion failed for falseVal`).

Callers wrapping untrusted input should never pass that input directly as a `trueVal`/`falseVal` of type `any` — a hostile caller could exploit the panic as a DoS vector. Validate types at the boundary and only call `If[T]` with known-typed values.

### Known Limitations

- `IsEmpty` cannot introspect arbitrary user-defined types without reflection, and by design does not use reflection. Types not explicitly handled return `false` (treated as present). If you need presence checks on custom structs, perform them yourself before calling `IsEmpty`.

## Security Updates

Security updates will be:

- Released as patch versions (e.g., 0.1.1)
- Documented in the CHANGELOG.md
- Announced through GitHub releases
- Tagged with security labels

## Acknowledgments

We appreciate responsible disclosure and will acknowledge security researchers who help improve the security of this project.

Thank you for helping keep our project and users safe!
