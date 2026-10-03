# Teaching Notes

## User profile
- Markus, working developer (German, DATEV context likely). Comfortable with CLIs, YAML, OpenAPI, CI/CD.
- Strong stated preference: **automation-first**. Dislikes the API Manager UI ("cumbersome ugly UI"). Lead every lesson with `apic` CLI / Platform REST API; mention UI only as orientation.
- Has live access to a **v10.x on-prem/OpenShift** environment at work → exercises can be real, but keep them read-only or scoped to a sandbox catalog until we know what he's allowed to touch.

## Environment (confirmed 2026-07-07, see LR-0001)
- Toolkit/platform: **10.0.11.0** (CD stream)
- Gateway type: **DataPower API Gateway (native)** — teach only native policies, skip v5c
- CLI login works against the corporate management endpoint

## Curriculum plan (revise as records accumulate)
1. **0001 — Object model + first CLI contact** (done, verified via LR-0001).
2. **0002 — Publish an OpenAPI file** (done): API yaml + product yaml anatomy, `apic validate`, `products:publish`, `--stage`, verify with `products:list-all`. Based on IBM's official example-toolkit-scripts artifacts.
3. **0003 — Download specs, CLI and raw REST** (done): `products:clone`, `apis:get`/`products:get`, `--debug` endpoint discovery, `/api/token`, curl against the Provider API.
4. **0004 — The assembly ("custom pipeline")** (done): context tree (request/message/api/error), invoke `output:`, set-variable, switch, gatewayscript with `context.*`, `catch`. New reference: assembly-policies.html.
5. 0005 — CI/CD pipeline: scripted publish via REST/toolkit, product replace/supersede, promotion between catalogs.
6. 0006 — User-defined policies on the API Gateway: catalog-scoped vs global-scoped, prichelle/apicv10-UDP + ibm-apiconnect/policy-apigw are the grounding sources.

## Open questions to resolve
- Does he have a sandbox catalog / own provider org where publishing test products is acceptable? (Lesson 2 homework assumes Sandbox is fair game — confirm.)
- For REST auth: does he have the toolkit credentials.json (client_id/secret), or does an admin need to create a registration? (Lesson 3 homework will reveal this.)
- Lesson 2/3 homework not yet reported back — no learning records for publish/download skills yet. Ask before starting Lesson 5.
- Lesson 4 homework asks him to list policies in colleagues' assemblies that the lesson didn't cover → direct input for Lessons 5/6 scoping.
- Teach `context.*` as the idiom; flag `apim.*` as legacy when he encounters it in colleagues' code.

## Style
- Short lessons, tangible win each time. Quizzes with equal-length answers. Cite IBM docs inline.
