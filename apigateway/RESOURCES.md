# IBM API Connect v10 Resources

## Knowledge

- [IBM Docs: API Connect 10.0.x — Working with the toolkit](https://www.ibm.com/docs/en/api-connect/10.0.x_cd?topic=cli-working-toolkit)
  Official reference for the `apic` CLI: login, command scopes, config. Use for: any CLI command syntax question.
- [IBM Docs: Creating and validating API and Product definitions by using the CLI](https://www.ibm.com/docs/en/api-connect/software/10.0.x_cd?topic=dyaa-creating-validating-api-product-definitions-by-using-cli)
  The canonical create → validate → publish workflow. Use for: lesson exercises on publishing.
- [API Explorer for the v10 Platform APIs](https://apic-api.apiconnect.ibmcloud.com/v10/)
  Interactive reference for the Platform REST API (Provider, Admin, Consumer APIs). Use for: raw REST endpoint paths, payloads, auth.
- [IBM Support: How to consume API Connect Platform — Cloud Management and Provider APIs](https://www.ibm.com/support/pages/how-consume-api-connect-platform-cloud-management-and-provider-apis)
  Walkthrough of token retrieval and calling the Platform API with curl. Use for: CI/CD without the toolkit.
- [IBM Docs: Working with Catalogs](https://www.ibm.com/docs/en/api-connect/10.0.x?topic=apis-working-catalogs)
  Catalog concepts, lifecycle, gateway/portal association. Use for: object-model questions.
- [IBM Cloud API Docs: API Connect Platform — Provider API](https://cloud.ibm.com/apidocs/apiconnect/apic-management-api)
  REST reference with per-endpoint curl examples. Use for: scripting publish/download without `apic`.
- [GitHub: ibm-apiconnect/example-toolkit-scripts](https://github.com/ibm-apiconnect/example-toolkit-scripts)
  Official example shell scripts managing orgs, apps, products, APIs via `apic`. Use for: CI/CD script patterns to steal. Verified working v10 API/product YAMLs live in `bash/` (`findbrancha.yaml`, `brancha_prod.yaml`).
- [IBM: v10 Command Line Intro (example-toolkit-scripts/docs)](https://github.com/ibm-apiconnect/example-toolkit-scripts/blob/master/docs/CommandLine-Intro.md)
  Practitioner-written intro: realms/identity providers, `--format json` + jq patterns, topology extraction. Use for: output-shaping and automation idioms.
- [IBM: v10 REST API First Steps (example-toolkit-scripts/docs)](https://github.com/ibm-apiconnect/example-toolkit-scripts/blob/master/docs/REST-API-FirstSteps.md)
  Worked example: create a client registration, get a token from `/api/token`, make authenticated Provider API calls. Use for: the canonical REST auth recipe (grounds Lesson 3).
- [Medium (IBM Expert Labs): Cloud, pOrg, Catalog, Space, Plan usage recommendation](https://medium.com/ibm-cloud-paks-help-and-guidance-from-ibm-cloud/ibm-api-connect-cloud-porg-catalog-space-plan-usage-recommendation-9c78c397904a)
  How real deployments map orgs/catalogs/spaces to teams and environments. Use for: judgment calls on structure.
- [Cloud Pak Deployment Guides: Promote APIs](https://production-gitops.dev/guides/cp4i/mq/apic/promote-apis/)
  GitOps-style API promotion between catalogs. Use for: CI/CD pipeline design.
- [IBM Docs: GatewayScript policy (execute)](https://www.ibm.com/docs/en/api-connect/10.0.1.x?topic=execute-gatewayscript)
  Reference for the gatewayscript assembly policy. Use for: custom logic in assemblies.
- [IBM Docs: Using context variables in GatewayScript/XSLT (DataPower API Gateway)](https://www.ibm.com/docs/en/api-connect/10.0.x?topic=aplc-using-context-variables-in-gatewayscript-xslt-policies-datapower-api-gateway)
  The authoritative description of the API context tree and `context.get/set`. Use for: anything assembly-related (grounds Lesson 4).
- [IBM Docs: GatewayScript code examples](https://www.ibm.com/docs/en/api-connect/software/10.0.1?topic=gatewayscript-code-examples)
  Official snippets for reading/writing message, headers, variables. Use for: writing gatewayscript policies.
- [IBM example API: utility_1.0.0.yaml (openid repo)](https://raw.githubusercontent.com/ibm-apiconnect/openid/master/utility/utility_1.0.0.yaml)
  Real IBM API with set-variable, switch, and gatewayscript in one assembly. Use for: verified policy YAML syntax.
- [IBM PoT Lab 5: Using GatewayScript in an Assembly](https://ibm-apiconnect.github.io/pot/lab5_logistics_gws.html)
  The invoke-with-output → gatewayscript enrichment pattern. Use for: orchestration examples. (Caveat: v5-era UI/apim.* API — take patterns, not syntax.)
- [Chris Phillips: Debugging an API (APIDevelopment-103)](https://chrisphillips-cminion.github.io/apiconnect/2019/10/01/APIDevelopment-103.html)
  Practitioner guide to debugging assemblies (logs, activity log). Use for: "why does my assembly do that?"
- [GitHub: prichelle/apicv10-UDP](https://github.com/prichelle/apicv10-UDP)
  Best current write-up of user-defined policies on the native API Gateway: catalog-scoped vs global-scoped, config sequences, multistep. Use for: Lesson 6. (Partially fills the UDP gap below.)
- [GitHub: ibm-apiconnect/policy-apigw — user-defined policies](https://github.com/ibm-apiconnect/policy-apigw/blob/master/user-defined-policies/basic/README.md)
  Official examples for packaging reusable user-defined policies on the DataPower API Gateway. Use for: "custom pipeline" reuse.

## Wisdom (Communities)

- [IBM TechXchange Community — API Connect](https://community.ibm.com/community/user/integration/communities/community-home?CommunityKey=2106cca0-a9f9-45c6-9b28-01a28f4ce947)
  IBM's official user community; product managers and Expert Labs people answer. Use for: version-specific gotchas, roadmap questions.
- [Stack Overflow tag: ibm-api-connect](https://stackoverflow.com/questions/tagged/ibm-api-connect)
  Moderate traffic but real practitioner answers. Use for: concrete error messages.

## Gaps

- IBM's docs site (ibm.com/docs) blocks automated fetching — doc links are verified to exist via search but content is grounded via IBM's GitHub repos. If a syntax dispute arises, the user's `apic validate` / `--help` output is the tiebreaker.
- UDP gap partially closed by prichelle/apicv10-UDP and ibm-apiconnect/policy-apigw; still missing: an official end-to-end 10.0.11-era walkthrough for catalog-scoped UDP deployment via `apic`.
