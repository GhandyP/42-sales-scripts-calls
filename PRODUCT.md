# Product

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

Go `net/http` with server-rendered HTML templates and minimal JavaScript; Google ADK for Go with Gemini. The user approved the previously proposed stack by asking to implement the plan. Run locally and bind only to loopback.

## Users

A single local operator preparing a sales-call script from customer facts and an approved playbook. The operator's exact job title is not confirmed.

## Product Purpose

Help the operator prepare a useful sales-call script from facts and approved objections supplied for that request, without inventing claims or placing or recording a call. Success means the operator can inspect a structured draft and its playbook evidence before using it.

## Positioning

Ground script claims and objection responses in a playbook supplied by the operator; make evidence and stop conditions visible rather than presenting unsupported text as approved.

## Operating Context

The operator prepares for a future sales conversation in a local web browser. Each request supplies customer facts, an objective, and approved playbook text. Google Gemini receives those inputs through Google ADK only after the operator acknowledges the disclosure and submits the request.

## Capabilities and Constraints

- Generate a structured draft with an opening, questions, objection responses, claims, evidence references, and stop conditions.
- Use an ADK read-only `search_playbook` tool over the playbook supplied for the current request.
- Require an explicit per-request checkbox acknowledging that customer facts and playbook text will be sent to Google Gemini.
- Keep request data transient; do not save customer facts, playbooks, or generated scripts in the MVP.
- Read the Google API key from server-side environment configuration; never expose it to the browser or logs.
- Bind the web server to loopback. No accounts, multiuser support, CRM integration, outbound calling, call recording, or public deployment in the MVP.
- Treat this as a draft for human review, not an approved offer or legal review.

## Evidence on Hand

The repository's README defines the sales-script goal and safety boundaries. No real playbooks, customer data, brand assets, testimonials, or performance evidence have been supplied. Any example playbook or customer facts added for tests or demos must be labeled synthetic.

## Product Principles

1. Do not invent facts, claims, or objections.
2. Show the evidence behind claims and objection responses.
3. Stop or ask for missing information instead of filling gaps with guesses.
4. Keep the human in control; the product drafts only and never calls or records.
5. Be explicit before user-supplied information leaves the local machine for Gemini.
