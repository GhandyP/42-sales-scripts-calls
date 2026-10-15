# 42 — Sales scripts and calls

## Goal

Create a sales call script from supplied customer facts and approved objections without inventing claims or making a live call.

## Run the current prototype

From this case directory:

```bash
export GOOGLE_API_KEY="your-key"
go run .
```

`internal/adkhelper.Run` uses ADK Go v2.2.0, Gemini, and `gemini-flash-latest`; set `ADK_MODEL` to override the model.

Sample prompt:

```text
Create an opening and objection responses from this approved playbook.
```

## What exists now

- `main.go (exists)` sends customer facts and playbook text to one focused Gemini ADK agent.
- The result is free-form text; call stages, claims, objection IDs, questions, and prohibited statements are not validated.
- This is intentionally text-only: it does not retrieve an objection library, record a call, or dial a customer.

## Missing elements

- **Contract:** Define `CallBrief{customerFacts, objective, playbookVersion}` -> `SalesScript{opening, questions[], objections[], claims[], stopConditions[]}`.
- **ADK layer:** Add a script schema and call-prep session; use `search_playbook` as a read-only function tool and keep telephony tools out of the draft path.
- **Domain boundary:** Add approved-claim/objection retrieval, CRM context, call-recording, and telephony adapters.
- **Fixtures and validation:** Cover objections, prohibited claims, missing facts, competitor mentions, regulated products, and script length.
- **Safety:** Require sales/legal approval for claims and consent for recording; no outbound call or offer can happen automatically.

## Target structure

```text
42-sales-scripts-calls/
├── main.go (exists)
├── README.md (exists)
├── schema.go (next)
├── script.go (next)
├── adapters/playbook.go (next)
├── testdata/playbook.json (next)
└── main_test.go (next)
```

## Implementation order

1. Define call stages, claim, objection, question, and stop-condition schemas.
2. Add approved playbook retrieval and a customer-scoped session.
3. Validate forbidden claims, script length, and missing facts.
4. Test regulated/competitor objections and escalation triggers.
5. Add recording/telephony adapters only after consent and human approval design.

## Done when

- [ ] Each script claim and objection response cites approved playbook evidence.
- [ ] Missing facts and prohibited topics produce stop conditions.
- [ ] Fixtures cover regulated and adversarial conversations.
- [ ] Calling and recording require consent and human-controlled launch.

## Next improvement

Make approved playbook citations mandatory in `SalesScript` before adding call tooling.
