# 42 — Sales scripts and calls

## Goal

Create a sales call script from supplied customer facts and approved objections without inventing claims or making a live call.

This repository is being rebuilt as a local, single-user web app that drafts evidence-grounded sales-call scripts from supplied customer facts and approved playbook text. Drafts require a human review; the app does not place or record calls.

## Current foundation

The standalone Go module includes the sales-script domain contract, structural validation, a loopback HTTP server, and a request-scoped playbook search tool. The current server binds only to `127.0.0.1:8080`; application routes and Gemini generation are not implemented yet.

Use Go 1.26.6 or newer. The module pins Google ADK v2.5.0; the Go command can acquire the required toolchain when automatic toolchain selection is enabled.

Run the current server foundation with:

```bash
go run .
```

The Gemini agent and its API-key setup will be documented when the generation workflow is in place. Do not put API keys in source files or commit them.

## Playbook retrieval

`internal/playbook` builds an immutable index from the approved playbook supplied for one request. It returns exact source excerpts of at most 1,000 runes with stable IDs such as `playbook-001`. Search ranks distinct matching terms, keeps source order for ties, and returns three passages by default (maximum five).

The ADK `search_playbook` tool accepts only a query and optional result limit; it cannot replace the playbook or access another request's index. Matching is lexical and case-insensitive over Unicode letters and digits. It does not normalize Unicode or establish that a passage semantically supports a generated claim. Generation and a separate support-validation gate remain future work.

## Tests

No live Google credentials are needed:

```bash
go test ./internal/playbook -count=1
go test ./...
```
