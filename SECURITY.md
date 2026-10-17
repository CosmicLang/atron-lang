# Security Policy

## Supported versions

Atron has not had a release yet. Milestone **M0** provides only the `atron`
command with `version` and `help`; it cannot parse or execute Atron source, so
the attack surface is limited to command-line argument handling.

| Version | Supported |
|---|---|
| `main` branch (pre-release, M0) | Yes — fixes land on `main` |
| Any tagged release | None yet; first release planned for milestone M10 (0.1.0) |

Once 0.1.0 is published, this table will list the supported minor versions. The
project intends to support the latest minor release and the one before it.

## Reporting a vulnerability

**Do not open a public issue, pull request or discussion for a suspected
vulnerability.** A public report removes the maintainer's ability to prepare and
ship a fix before the details are widely known.

Atron uses **coordinated disclosure through GitHub Private Vulnerability
Reporting (PVR)**. PVR is the project's only accepted private vulnerability
channel.

### Activation status

**PVR is enabled for `CosmicLang/atron-lang`.** The repository API reports
`{"enabled": true}` for
`repos/CosmicLang/atron-lang/private-vulnerability-reporting`, and the endpoint
returns HTTP 200. The channel is operational: a report filed through it creates
a private draft security advisory visible only to the reporter and the
maintainers.

### Submitting a report

The **Security** tab of the repository shows a **Report a vulnerability** button
that opens the GitHub Security Advisory draft form. Reports are visible only to
you and the repository maintainers until an advisory is published.

1. Open `https://github.com/CosmicLang/atron-lang/security/advisories/new`, or
   use the **Report a vulnerability** button on the repository's **Security**
   tab.
2. Give the advisory a concise title and fill in the affected ecosystem and
   package (`golang`, `github.com/CosmicLang/atron-lang`).
3. Describe the issue using the fields listed below.
4. Submit. GitHub creates a private draft advisory and notifies the maintainers.

Please include:

- the affected version (`atron version` output) and platform,
- the exact input that triggers the problem (source file or command line),
- the observed and expected behaviour,
- an assessment of impact, and whether the issue is already public,
- any suggested fix, if you have one.

## What to expect

The following are targets, not guarantees, and they apply from the moment a
report arrives through PVR:

| Stage | Target |
|---|---|
| Acknowledgement of the report | within 3 business days |
| Initial assessment and severity triage | within 10 business days |
| Fix or mitigation for high and critical issues | within 30 days |
| Public disclosure | coordinated with the reporter after a fix ships |

Reporters are credited in the release notes unless they ask not to be.

## Scope

In scope once the language exists:

- the compiler front end (lexer, parser, resolver) on hostile input,
- the bytecode compiler and virtual machine, including stack and bounds safety,
- the module loader, including path traversal outside the project root,
- the standard library,
- the embedding API in `pkg/atron`,
- memory exhaustion or unbounded recursion reachable from ordinary source.

Out of scope:

- vulnerabilities in third-party dependencies that are already public and have
  no Atron-specific exploit path (please report those upstream),
- issues that require a modified binary or a hostile local environment,
- denial of service that requires a deliberately unbounded program, unless the
  specification defines a resource limit that Atron fails to enforce,
- missing hardening that has no demonstrated impact.

## Security posture of the design

These design decisions exist partly for security and are enforced in review:

- No mutable global state, so an embedder cannot be affected by another
  interpreter instance in the same process.
- The VM is expected to check stack and operand bounds on every operation; the
  bytecode is treated as untrusted input once a serializer exists (M9+).
- The module loader must never escape the resolved project root.
- `gosec` and `govulncheck` run in CI on every push and pull request.

## Disclosure policy

The project follows coordinated disclosure. A fix is prepared and released
before details are published. Advisories are published through GitHub Security
Advisories so that dependents can track them automatically, and are requested as
CVE identifiers where the impact warrants it.

The reporter is consulted on the disclosure date. The project will not publish
an advisory while the reporter has an outstanding, reasonable objection, except
where the issue is already being exploited in the wild.

### Separate channels

Vulnerability reports and code-of-conduct reports are handled through different
channels and are not interchangeable. The code-of-conduct moderation address in
`CODE_OF_CONDUCT.md` is an interpersonal-conduct channel; it is deliberately not
published here as a security contact. Do not send vulnerability details to it.
