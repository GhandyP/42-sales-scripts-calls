# 42 — Sales scripts and calls

## Goal

Create a sales call script from supplied customer facts and approved objections without inventing claims or making a live call.

This repository is being rebuilt as a local, single-user web app that drafts evidence-grounded sales-call scripts from supplied customer facts and approved playbook text. Drafts require a human review; the app does not place or record calls.

## Current foundation

The standalone Go module, domain contract, validation, and loopback HTTP server are being established before the playbook and Gemini workflow. The current server binds only to `127.0.0.1:8080`; application routes are added in a later implementation task.

Run the current server foundation with:

```bash
go run .
```

The ADK integration and its API-key setup will be documented when the generation workflow is in place. Do not put API keys in source files or commit them.
