# Environment established: v10.0.11.0 toolkit, native DataPower API Gateway

Markus completed Lesson 1's hands-on and reported `apic version` = **10.0.11.0** and that his work catalogs use the **DataPower API Gateway** (native, not v5-compatible). He can log in and list catalogs from the CLI, so the object model + login flow is demonstrated, not just covered.

**Implications:** All assembly/policy lessons should target the native API Gateway policy set (e.g. `invoke` version 2.x, gatewayscript for the API Gateway) and ignore v5c-only material. Doc links can pin to the 10.0.x CD stream. Exercises may assume a working login against the corporate management endpoint.
