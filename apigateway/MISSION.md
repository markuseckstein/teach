# Mission: IBM API Connect (v10.x) for API Publishers

## Why
Markus must publish and manage his team's APIs (OpenAPI YAML) on his company's
IBM API Connect v10.x installation (on-prem / OpenShift). He wants to do this
through the `apic` toolkit CLI and the Platform REST API — scriptable and
CI/CD-friendly — treating the API Manager UI as a fallback, not the workflow.

## Success looks like
- Publish an `openapi.yml` to a catalog from the command line, end to end.
- Download/export existing API and Product definitions from a catalog via CLI or REST.
- Read and modify an API's assembly ("custom pipeline") — policies in `x-ibm-configuration`.
- Explain the object hierarchy (provider org → catalog → space → product → plan; consumer org → app → subscription) and where an OpenAPI file fits into it.
- Wire publishing into a CI/CD pipeline using the Platform REST API or `apic` in scripts.

## Constraints
- Automation-first: CLI and REST before UI, always.
- Has a live work environment (v10.x on OpenShift) to practice against — real exercises are possible, but must be safe (no destructive operations against shared catalogs without care).
- Working developer: lessons must be short and immediately applicable.

## Out of scope
- Installing/operating the platform itself (OpenShift operators, gateway sizing, HA).
- The Developer Portal customization side.
- Legacy v5 / v2018 compatibility topics.
