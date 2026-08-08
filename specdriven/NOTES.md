# Working Notes

## User preferences
- Include good videos (YouTube/Vimeo) in lessons when relevant content exists (stated 2026-07-06). Vet the presenter — prefer maintainers/official channels.
- Practice style: **mix per topic** — real-repo hands-on steps for workflow/tooling topics, quizzes for concepts and terminology (stated 2026-07-06).
- Shell context for all commands: **Git Bash on Windows**. Give `sh` variants, not PowerShell, unless contrasting.
- Every tool-setup lesson must include the corporate-proxy angle (SSL inspection, `NODE_EXTRA_CA_CERTS`, allowlists) — it's a standing constraint, not an aside.

## Open threads
- User asked in session 1: "where do the specs go — backend or frontend repo?" Answered in chat (grounded in spec-kit discussion #1743); candidate for a dedicated lesson + reference doc once they've done the basic loop once.
- Terminology hazard to reinforce early: **"spec" (SDD feature spec, Markdown) ≠ "spec" (OpenAPI spec, API contract)**. Their stack uses both meanings daily.

## Teaching sequence sketch (revise as records accumulate)
1. ✅ 0001 — The SDD loop & Spec Kit's artifacts (concepts + quiz)
2. ✅ 0002 — Toolchain behind the proxy: uv/specify, git schannel, Copilot CLI (hands-on; verify user actually completed the steps before building on them)
3. ✅ 0003 — Constitution for the two-repo project (hands-on; enforceability test + "hypothetical PR" feedback loop)
4. ✅ 0004 — specify → clarify → plan on the bookmark practice feature (hands-on; contract born in contracts/)
5. ✅ 0005 — Where do specs live: constitution heuristic, umbrella workspace + --no-git (feature.json verified in scripts/bash/common.sh), specs-repo graduation path. Closes the session-1 question.
6. ✅ 0006 — tasks → analyze → implement (hands-on). Key teaching: tests-optional-by-default trap vs constitution Article III; analyze read-only + constitution conflicts CRITICAL; phase-wise implement; contract-sync tasks into both repos' generators.
7. Review session: retrieval practice across lessons 1–6 (spaced recall — user has had no verified practice yet; check hands-on status first!)
8. Copilot CLI vs VS Code agent mode: when to drive from which
9. Apply to the real work project (graduation)

⚠ Lessons 2–4 are hands-on and unverified — no evidence yet the user ran them. Before lesson 5, ask what actually happened at their machine (proxy errors, generated artifacts) and write learning records from that.
