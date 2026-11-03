# Security

Report vulnerabilities privately to the repository owner. Do not include tokens,
recordings, clone states or another person's personal data in public issues.

My Voice requires an authenticated app session and CSRF token for mutations.
Managed MyAccounts sessions and grants are checked on every sensitive request.
Grant revocation prevents subsequent synthesis and reference access. Shared users
receive use rights; the original reference remains available only to its owner.
Existing generated audio owned by a recipient is not retroactively erased.

Only loopback listeners are permitted. Use Workspace for trusted remote access.
Inputs cannot supply executable paths, arbitrary model URLs or OS commands.
No cloud speech API or telemetry endpoint is configured. Scripts/references and
session tokens are omitted from worker logs. Do not make private data a broad
Drive source or store it under the replaceable program directory.

The preview has native scope/CSRF/queue/path tests, not an independent security
audit. Public voice marketplaces, unrestricted LAN APIs and automatic licence-gate
bypasses are not supported.
