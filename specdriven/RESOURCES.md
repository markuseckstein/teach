# Spec-Driven Development / GitHub Spec Kit Resources

## Knowledge

- [Docs: Spec Kit Documentation (official)](https://github.github.com/spec-kit/)
  The primary source. Covers the Spec → Plan → Tasks → Implement workflow, all `/speckit.*` commands, integrations, extensions. Last updated May 2026. Use for: anything about how Spec Kit itself works.
- [Docs: Spec Kit Quick Start](https://github.github.com/spec-kit/quickstart.html)
  Prerequisites (git, uvx/pipx), install commands, the 8 slash commands in order, `--script sh` vs `--script ps` selection. Use for: installation and first-run steps.
- [Repo: github/spec-kit](https://github.com/github/spec-kit)
  Source, issues, and discussions. The discussions are where multi-repo/monorepo guidance actually lives. Use for: edge cases the docs don't cover yet.
- [Blog: "Spec-driven development with AI" — The GitHub Blog](https://github.blog/ai-and-ml/generative-ai/spec-driven-development-with-ai-get-started-with-a-new-open-source-toolkit/)
  The announcement post; best articulation of *why* SDD exists and the philosophy behind it. Use for: motivation, explaining SDD to colleagues.
- [Blog: "Diving Into Spec-Driven Development With GitHub Spec Kit" — Microsoft for Developers](https://developer.microsoft.com/blog/spec-driven-development-spec-kit)
  Practitioner walkthrough from Microsoft. Use for: a second worked example of the full loop.
- [Docs: Configuring network settings for GitHub Copilot — GitHub Docs](https://docs.github.com/en/copilot/managing-copilot/configure-personal-settings/configuring-network-settings-for-github-copilot)
  Official proxy/certificate guidance: `HTTP_PROXY`/`HTTPS_PROXY`, `NODE_EXTRA_CA_CERTS`, OS trust store behavior. Use for: all corporate-proxy setup.
- [Docs: Network settings concepts for Copilot — GitHub Docs](https://docs.github.com/en/copilot/concepts/network-settings)
  Domains to allowlist (github.com, api.github.com, *.githubusercontent.com …). Use for: firewall allowlist requests to IT.
- [Docs: TLS certificates — uv (Astral)](https://docs.astral.sh/uv/concepts/authentication/certificates/)
  uv's trust model: bundled Mozilla roots by default, `UV_SYSTEM_CERTS` for the OS store, `SSL_CERT_FILE` override. Use for: any `uvx`/`specify` failure behind the proxy.
- [Repo: github/copilot-cli README](https://github.com/github/copilot-cli)
  Authoritative install doc: PowerShell 6+ Windows prerequisite, winget/npm install paths, enterprise-policy disable note. Use for: Copilot CLI setup and "is it even enabled for my org" questions.
- [Discussion: Multi-Repository Setup — spec-kit #1743](https://github.com/github/spec-kit/discussions/1743)
  Maintainer guidance on separate specs/backend/frontend repos; the constitution-overlap heuristic; `--no-git` flag. Use for: the "where do specs live" decision.
- [Discussion: Spec-Kit artifacts in a monorepo — spec-kit #769](https://github.com/github/spec-kit/discussions/769)
  Community patterns for shared spec folders. Use for: contrast with the multi-repo approach.

- [Video: "The ONLY guide you'll need for GitHub Spec Kit" — Den Delimarsky (YouTube, 40 min)](https://www.youtube.com/watch?v=a9eR1xsfvHg)
  By the Spec Kit maintainer at GitHub. Full workflow walkthrough, constitution to implement. Use for: seeing the whole loop driven by someone who built the tool.
- [Training: "Get Started with Spec-Driven Development and GitHub Spec Kit" — Microsoft Learn](https://learn.microsoft.com/en-us/training/modules/spec-driven-development-github-spec-kit-greenfield-intro/)
  Official free module, greenfield scenario. Use for: structured practice with checks.
- [Training: "Implement Spec-Driven Development using the GitHub Spec Kit" — Microsoft Learn](https://learn.microsoft.com/en-us/training/modules/spec-driven-development-github-spec-kit-enterprise-developers/)
  Follow-up module aimed at enterprise developers, VS Code + Copilot. Use for: enterprise-standards angle on the same workflow.
- [Templates: github/spec-kit/templates](https://github.com/github/spec-kit/tree/main/templates)
  The actual constitution/spec/plan/tasks templates the commands fill in. Use for: knowing what "complete" looks like for each artifact.

## Wisdom (Communities)

- [github/spec-kit Discussions](https://github.com/github/spec-kit/discussions)
  The de-facto community; maintainers answer directly. Use for: architecture questions (like spec placement), workflow critique.
- [GitHub Community Discussions — Copilot](https://github.com/orgs/community/discussions/categories/copilot)
  Where enterprise proxy/cert issues get triaged (e.g. discussions #29127, #35544). Use for: proxy problems that survive the official docs.

## Gaps

- No high-trust written guide yet found for **Spec Kit + OpenAPI contract-first** across separate Spring Boot / Angular repos — will likely need to synthesize from #1743 + OpenAPI Generator docs.
- Copilot CLI behavior specifically under **Git Bash on Windows** (vs PowerShell) is thinly documented; verify hands-on in an install lesson.
