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
| 2026-11-24 | `.github/dependabot.yml` (removed)      | No dependency update pull requests here; dependencies follow upstream releases.                         |
| 2026-11-24 | `web/src/lib/api-base-address.ts` (new) | Picks the address API clients should use: the first API info entry if the admin set one, otherwise the server address. |
| 2026-11-24 | `web/src/features/keys/components/dialogs/cc-switch-dialog.tsx`, `web/src/features/keys/components/data-table-row-actions.tsx`, `web/src/features/chat/hooks/use-chat-presets.ts` | CC Switch import, chat app imports (Cherry Studio and others) and "Copy Connection Info" use that address instead of the server address. SipToken's server address is the console (needed for email links), while clients should call the API domain. |
| 2026-11-24 | `web/src/lib/api-base-address.ts`, `web/src/features/keys/components/dialogs/cc-switch-dialog.tsx` | The CC Switch import names the provider after the site (System Name, e.g. "SipToken") instead of "My Claude" / "My Codex" / "My Gemini". Falls back to upstream's names while the System Name is unset or still "New API". |
| 2026-11-24 | `web/src/features/chat/lib/chat-links.ts` | The Cherry Studio import creates a provider named after the site (ID `siptoken`, name "SipToken") instead of filling Cherry's built-in "New API" preset. Importing into the preset replaced the address of any other New API site the user had set up there, while that site's key stayed first, so requests failed with "invalid token". Falls back to upstream's `new-api` while the System Name is unset, still "New API", or has no letters or digits. |

Nothing else differs from upstream: compare with the first commit of this
repository, which is the unmodified upstream source. The web console frontend
keeps the attribution and the link to the original project required by NOTICE.

## Upgrading

1. Check out the new upstream tag.
2. Re-apply the changes above. If upstream added keys to `i18n/locales/en.yaml`,
   add the same keys to `ru.yaml`; `go test ./i18n/` fails until every key is
   translated. Search the frontend for other places that read `server_address`
   for client imports, and route new ones through `getApiBaseAddress`.
3. Record the new base tag and commit at the top of this file and in
   `UPSTREAM_TAG` in `.github/workflows/siptoken-image.yml`.
4. Push to `main`. The workflow runs the tests and publishes
   `ghcr.io/<owner>/new-api:<upstream tag>-siptoken.<run number>` and `:latest`.
