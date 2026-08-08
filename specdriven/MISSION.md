# Mission: Spec-Driven Development with GitHub Spec Kit

## Why
Markus has a real work project — a Spring Boot backend and an Angular frontend in **two separate repos**, communicating via OpenAPI-generated REST APIs — and wants to build features on it spec-first using GitHub Spec Kit with GitHub Copilot (CLI and VS Code). The goal is shipping real features through the SDD workflow, not just understanding the theory.

## Success looks like
- Can run the full Spec Kit loop (constitution → specify → clarify → plan → tasks → implement) on a real feature in the work project, driving Copilot in both VS Code and the Copilot CLI
- Has a working answer for **where specs live** in a two-repo (backend + frontend) architecture, and can defend it
- Can set up Spec Kit and Copilot CLI on Windows + Git Bash **behind the corporate proxy** (SSL inspection, CA bundles) without help
- Writes specs that produce correct OpenAPI contracts, so generated Spring Boot and Angular API code stays in sync

## Constraints
- Windows machine, Git Bash as the shell (not PowerShell-first)
- Enterprise environment: corporate proxy with SSL inspection is a standing constraint — every tool-setup lesson must address it
- Tooling is GitHub Copilot (VS Code + Copilot CLI); other agents (Claude, Gemini) are out of scope

## Out of scope
- Other SDD tools (Kiro, OpenSpec, Tessl) — mention only for contrast
- Copilot autocomplete basics (already fluent)
- Spring Boot / Angular fundamentals (already the day job)
