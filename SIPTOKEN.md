# SipToken modifications

This is a modified version of [New API](https://github.com/QuantumNous/new-api),
licensed under the GNU Affero General Public License v3.0 with the additional
terms in [NOTICE](./NOTICE). SipToken runs it as the console at
https://console.siptoken.ai.

Upstream base: tag `v1.0.0-rc.51`, commit `c9c6da2786b6a0912215127e78e28fbc6650c85f`.

## Changes

| Date       | Files                                   | Change                                                                                                  |
| ---------- | --------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| 2026-11-24 | `i18n/locales/ru.yaml` (new)            | Russian translation of all backend messages: emails, errors returned to API clients, low-quota notices. |
| 2026-11-24 | `i18n/i18n.go`                          | Loads the Russian locale and maps `ru` / `ru-*` (user setting or `Accept-Language`) to it.             |
| 2026-11-24 | `i18n/i18n_test.go`                     | Tests for the Russian language mapping and messages.                                                    |
| 2026-11-24 | `.github/workflows/siptoken-image.yml`  | Builds the Docker image from this source.                                                               |

Nothing else differs from upstream. The web console frontend is unmodified and
keeps the attribution and the link to the original project required by NOTICE.

## Upgrading

1. Check out the new upstream tag.
2. Re-apply the changes above. If upstream added keys to `i18n/locales/en.yaml`,
   add the same keys to `ru.yaml`; `go test ./i18n/` fails until every key is
   translated.
3. Record the new base tag and commit at the top of this file and in
   `UPSTREAM_TAG` in `.github/workflows/siptoken-image.yml`.
4. Push to `main`. The workflow runs the tests and publishes
   `ghcr.io/<owner>/new-api:<upstream tag>-siptoken.<run number>` and `:latest`.
