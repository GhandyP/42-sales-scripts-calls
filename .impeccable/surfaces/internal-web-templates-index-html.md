---
version: 1
slug: "internal-web-templates-index-html"
primary_target: "internal/web/templates/index.html"
related_targets: []
---

# Call preparation surface

## Scope

Operate-mode local web app for one operator preparing a sales-call draft. The task is to enter customer facts, an objective, and approved playbook text; submit only after acknowledging the Google Gemini data transfer; then inspect the structured draft and its evidence. Proof is the exact supplied playbook passage behind each claim or objection response. The app never calls or records. Inputs remain transient. The first viewport covers the complete form, draft state, and evidence relationship; responsive layouts preserve the workflow on narrow screens. The interface language follows the repository's English README. No persistent library, user accounts, CRM, multiuser mode, or public deployment is in scope.

## Direction contract

THESIS: Prepare a conversation as a marked rehearsal script: every claim or objection cue leads to the exact approved passage behind it.

OWN-WORLD: Translate a stage manager's promptbook into a working screen, not a theater costume. Use a full palette: ink blue `#152B78` for the shell, clean paper `#F6F7F2` for the ground, cue coral `#F44F3E` for critical actions/stops, and acid lime `#D8F05B` for the active evidence cue. Use a distinctive condensed face for display headings, a readable workhorse sans for form/body text, and monospaced type only for cue IDs. Standard web controls remain standard. Cues are labels, never color-only marks.

STORY: The operator supplies allowed facts and approved material, sees the external-data disclosure, confirms it per request, and receives a draft whose evidence can be inspected without losing the surrounding script.

FIRST VIEWPORT: At desktop width, a compact ink-blue header sits above three aligned work areas: a narrow left form with the three required inputs and consent, a broad central draft with a real empty state before generation, and a right evidence rail. On a generated draft, selecting a cue highlights its source excerpt and every section that cites it. On mobile, stack form, draft, and evidence in that order; keep consent beside the action it gates.

FORM: The user chose the assigned candidate, the theatre promptbook, from concept seed `e2ac995d` (assigned index 5). This is a code-led Operate surface; no visual comp is owed.

FINISH: unreviewed and undocumented is unfinished; this build ends with the finish review, the verdict, DESIGN.md, and every shipping raster carrying its provenance
