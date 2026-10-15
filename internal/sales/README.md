# Sales script domain

This package defines the request and response contract for a draft generated from customer facts, a call objective, and approved playbook text. Claims and objection responses refer to evidence passage IDs so the web layer can display the exact source text.

`SalesScript.ValidateStructure` checks required sections and passage-ID references only. It does not establish that the cited passage semantically supports a claim; the generation path must pass a separate support-verification gate before presenting a draft as grounded.

Request data is validated but not persisted. The model caller remains responsible for obtaining playbook passages through the read-only search tool.
