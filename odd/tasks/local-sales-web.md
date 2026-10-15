Objective: Build the approved personal/local sales-script web app with Go and Google ADK, grounded in user-supplied approved playbooks.
Branch: `feat/local-sales-web` | Delivery strategy: `ask-on-risk` | Test runner: `go test ./...` plus browser verification.

## Specs

1. **S1 — Standalone Go/ADK web app.** The user requested: “ok entonces ahora vamos a construir desde cero este proyecto para que sea una app web con go y google adk y ademas crea el plan con odd y me avisas cuando tengas el plan listo” and later authorized implementation with “mira entonces vamos a implementar ese plan”. Build a standalone Go web app using Google ADK, not the cookbook's private helper.
2. **S2 — Personal/local MVP.** The user selected “MVP personal/local (recomendado)”. Keep the first version single-user and local; no accounts, multiuser data, public deployment, CRM, outbound calls, or recording.
3. **S3 — Grounded sales-script workflow.** Preserve the existing README goal verbatim: “Create a sales call script from supplied customer facts and approved objections without inventing claims or making a live call.” Use supplied customer facts, an objective, and an approved playbook; return a structured script with an opening, questions, objection responses, claims, and stop conditions. Each claim and objection response must cite playbook evidence; invalid or unsupported citations must not be presented as verified.
4. **S4 — ADK playbook handling.** The user delegated this choice with “la idea es usar adk creo que se puede hacer algo bien”. Implement an ADK read-only `search_playbook` function tool over the playbook supplied for that request. Do not add an external playbook, CRM, or telephony integration.
5. **S5 — Explicit Gemini data notice.** The user selected “Aviso y casilla por envío (recomendado)”. Explain that the submitted customer facts and playbook are sent to Google Gemini and require an explicit, server-validated checkbox before each generation. Keep the API key server-side and do not log or persist request content.
6. **S6 — Accepted UI stack and local access.** The accepted plan specifies Go `net/http`, server-rendered HTML templates, minimal JavaScript, and local-only use. Bind to loopback. Prefer transient per-request inputs; do not add persistence unless the user revises scope.
7. **S7 — Verifiable delivery.** Keep deterministic tests independent of live Google credentials; verify behavior through public HTTP/UI interfaces, run applicable Go checks and browser checks, and document setup and limitations.

## Tasks

- **T0** — [S1-S7] inline; in_progress; commit: pending — Record confirmed product truth and the selected Operate-surface direction before code.
- **T1** — [S1,S3,S6,S7] inline; pending; commit: pending — Create the standalone Go module, HTTP/server foundation, sales-script domain contract, and validation tests.
- **T2** — [S3,S4,S7] inline; pending; commit: pending — Implement deterministic playbook evidence chunking/search and the ADK read-only function tool with tests.
- **T3** — [S1,S3,S4,S5,S7] inline; pending; commit: pending — Integrate the Google ADK Gemini agent, structured result parsing, and citation/stop-condition validation with mocked tests.
- **T4** — [S1,S2,S5,S6,S7] inline; pending; commit: pending — Build the local web workflow, explicit consent, server-side validation, and responsive accessible UI with HTTP tests.
- **T5** — [S1-S7] inline; pending; commit: pending — Update run/configuration docs and synthetic examples; run full Go, browser, design-detector, independent security verification, finish review, and final design documentation.

## Log

- **L1 — Original request:** “ok entonces ahora vamos a construir desde cero este proyecto para que sea una app web con go y google adk y ademas crea el plan con odd y me avisas cuando tengas el plan listo”
- **L2 — Scope and authorization:** The user chose “MVP personal/local (recomendado)” and then said “mira entonces vamos a implementar ese plan”.
- **L3 — Playbook direction:** The user said “la idea es usar adk creo que se puede hacer algo bien”; implement a per-request playbook text input with read-only ADK search, without persistent storage.
- **L4 — Data-flow choice:** The user chose “Aviso y casilla por envío (recomendado)”; customer facts and playbook are disclosed as being sent to Google Gemini and require explicit confirmation per request.
- **L5 — Baseline:** The repository was clean on `main`, contained only `README.md`, `main.go`, and `.gitignore`, and had no standalone `go.mod`. `main.go` depends on `github.com/GhandyP/adk-go-enterprise-cookbook/internal/adkhelper`.
- **L6 — Research:** Official `google/adk-go` release `v2.5.0` is available; its tagged `go.mod` requires Go `1.26.6`, while this environment has Go `1.26.4`. Resolve toolchain compatibility before pinning the dependency. Sources: https://github.com/google/adk-go/releases/tag/v2.5.0 and https://raw.githubusercontent.com/google/adk-go/v2.5.0/go.mod.
- **L7 — ODD plan:** Forecast approximately 800 authored changed lines across six work units. Keep commits on `feat/local-sales-web`; use `ask-on-risk` if a future delivery slice needs a PR-size decision. No PR was requested.
- **L8 — Assumptions to preserve:** Inputs are transient; the app binds only to loopback; no customer data or API key is logged. The API key comes from `GOOGLE_API_KEY`. Use an explicit per-request consent checkbox, not a remembered grant.
- **L9 — Product and design preflight:** Wrote `PRODUCT.md`; `impeccable context` found no existing visual implementation and no image generator, so the app surface uses the code-led Operate path. Concept seed `e2ac995d` assigned candidate 5 (theatre promptbook); the assigned direction, top pick, two competitive alternatives, declined challengers, and category standard are shown in the local decision page with five HTML sketches. Awaiting the user's direction choice; no application source code has changed.
- **L10 — External wait:** The decision page was open at http://127.0.0.1:44775/ awaiting the user's selection; the selection has now arrived.
- **L11 — Direction selected:** The user chose the assigned theatre-promptbook direction with `{"optionId":"assigned","steer":"","buildPath":"code","buildPathFlipped":false}`. Continue code-led. T0 remains in progress until the direction contract is committed.
- **L12 — Surface contract:** Wrote and read back `.impeccable/surfaces/internal-web-templates-index-html.md`; verified all six required direction-contract blocks and seed key `e2ac995d`. T0 is ready for its documentation work-unit commit; no application source code has changed.
- **L13 — Repository hygiene:** Added root ignore rules for generated local Impeccable decision mocks, question logs, review captures, and build state. `git check-ignore` confirms those paths are ignored while `.impeccable/surfaces/internal-web-templates-index-html.md` remains eligible for versioning.
