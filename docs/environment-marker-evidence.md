# Desktop environment marker evidence

Recorded on 2026-10-16. The WorkBuddy, local Grok Bot and QwenWork evidence below
is based on static inspection of installed macOS distribution artifacts. No agent
was launched and no real environment values, user configuration or session
contents were collected for these tests. Locations refer to the packaged files,
including minified line numbers, not unpublished source repositories. The
synthetic test cases do not establish coverage of every platform or execution mode.

## Product descriptions and official sources

Checked on **2026-10-16**. The descriptions below paraphrase official product
positioning. These pages establish product names and intended uses; they do not
establish environment-variable contracts or prove that a detected process was
started by an agent. The marker evidence and its limits are documented separately.

| Product and result            | Description                                                                                            | Official source                                                                                                                                                                                       |
| ----------------------------- | ------------------------------------------------------------------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Crush (`crush`)               | A terminal coding agent from Charm.                                                                    | [Official source](https://github.com/charmbracelet/crush#readme)                                                                                                                                      |
| Hermes Agent (`hermes-agent`) | An AI agent from Nous Research with persistent memory and reusable skills.                             | [Official source](https://hermes-agent.nousresearch.com/docs)                                                                                                                                         |
| Kilo Code (`kilo-code`)       | An open-source coding agent available in editors and the terminal.                                     | [Official source](https://github.com/Kilo-Org/kilocode#readme)                                                                                                                                        |
| TRAE (`trae`)                 | An AI coding environment with IDE and SOLO modes.                                                      | [Official source](https://www.trae.ai/blog/trae_work_0609)                                                                                                                                            |
| VT Code (`vtcode`)            | A terminal coding agent for planning, implementing and reviewing code.                                 | [Official source](https://github.com/vinhnx/VTCode#readme)                                                                                                                                            |
| WorkBuddy (`workbuddy`)       | Tencent's desktop AI workspace for completing office tasks.                                            | [Official source](https://www.codebuddy.cn/docs/workbuddy/From-Beginner-to-Expert-Guide/Product-Guide)                                                                                                |
| Grok Bot (`grok-bot`)         | AI teammates from Cursor that carry out tasks across apps and websites on a persistent cloud computer. | [Official source](https://cursor.com/docs/grok-bot)                                                                                                                                                   |
| QwenWork (`qwenwork`)         | Alibaba's AI office platform for completing tasks across desktop and web.                              | [Official source](https://qwenwork.cn/docs/product-introduction)                                                                                                                                      |
| Cue (`cue`)                   | Personal AI agents from Manus AI with their own computer, email, phone, and wallet.                    | [Official source](https://cue.im/?lang=zh-CN); [agent API overview](https://open.manus.ai/docs/v2/agents-overview)                                                                                    |
| Muse (`muse`)                 | Meta's personal AI agent that works across connected services on a dedicated cloud computer.           | [Official source](https://research.meta.ai/blog/security-and-safety-for-ai-agents-our-approach-with-muse)                                                                                             |
| Manus (`manus`)               | General-purpose AI agent for research, content creation, and workflow automation.                      | [Official source](https://www.manus.im/about); [product site](https://manus.im/)                                                                                                                      |
| Doubao Work (`doubao-work`)   | An AI work platform for individuals, teams, and enterprises.                                           | [Official source](https://www.feishu.cn/content/article/7677519271848610746); [official English name](https://www.feishu.cn/hc/en-US/articles/428390929435-admin-configure-screen-capture-protection) |
| Zed (`zed`)                   | A code editor designed for fast, collaborative development.                                            | [Official source](https://github.com/zed-industries/zed#readme)                                                                                                                                       |

TRAE's description refers to the coding product with IDE and SOLO modes, not
the separately positioned TRAE Work assistant. Cue's official site identifies
Manus AI in its footer, and the official agent API overview links Cue app agents
to that API. This product attribution does not independently validate Cue's
user-reported runtime marker.

## Zed terminal context

The `ZED_TERM` condition identifies a terminal in the Zed editor, including a
human-operated terminal. It remains the final fallback after agent-specific
rules and does not establish that an agent initiated a command. The product
positioning above therefore describes the editor rather than an agent-only
runtime. The marker is also listed in the
[Hugging Face agent-harness registry](https://github.com/huggingface/huggingface.js/blob/e08240742cd64d8f908b22e52358b84c812e882c/packages/tasks/src/agent-harnesses.ts).

## WorkBuddy

Inspected WorkBuddy 5.8.3, bundle ID `com.tencent.workbuddy.mac`.

- `Contents/Resources/app.asar` SHA-256:
  `a7eff702411704d140d7d8d706fe9e711d9b13e304928481c9eba0af9c06f731`.
- Inside that archive, `main/main-bootstrap.js:77` (`buildDaemonProcessEnv`)
  constructs `WORKBUDDY_APPLICATION_NAME`, `WORKBUDDY_APP_NAME` and
  `WORKBUDDY_IS_PACKAGED`; the daemon spawn at line 75 merges this object into
  the inherited environment. `main/sidecar-manager.js:16`
  (`buildSidecarProcessEnv`) also preserves inherited variables for its spawn.
- `Contents/Resources/app.asar.unpacked/cli/dist/codebuddy-lite-wb.mjs` SHA-256:
  `9eacdcdaaf3d53877ca4eac53836a446770b9a40f5dc6382b43230c73fdc5448`.
  Its `buildChildProcessEnv` at line 1684 supplies the ordinary shell and PTY
  spawn environments. The intervening cleanup at line 1819 removes trace/request
  fields rather than these product-name markers. The sandbox path at line 1134
  carries the target environment in an execution request; native sandbox
  propagation was not dynamically tested.

The rule accepts exact `WorkBuddy` in either complete product-name key. It does
not accept the unprefixed `APPLICATION_NAME` / `APP_NAME`, packaging flags, or
configuration directories alone. Renamed or branded variants may not match.
The public [product guide](https://www.codebuddy.cn/docs/workbuddy/From-Beginner-to-Expert-Guide/Product-Guide)
describes the product, but does not establish these markers as a stable public API.

## Grok Bot

Inspected Grok Bot 0.71.2, bundle ID `com.anysphere.sand`.

- `Contents/Resources/app.asar` SHA-256:
  `055b087a9bacc3be4e5f80c468edf2fd031dd5b928ab2e7da7e01dd10de0e860`.
- `dist/local-exec-daemon/main.cjs:756` (`iZe`) supplies both `SAND_AGENT=1`
  and `CURSOR_AGENT=1` to the local executor. Line 611 (`pI`) passes this
  environment to the Bash initialization at line 602 (`Jae`) and the Zsh
  initialization at line 608 (`Wae`), which merge it into the child environment.

The local branch requires exact `SAND_AGENT=1` and precedes Cursor rules.
`CURSOR_AGENT` alone is shared compatibility context and is not an alternative
Grok Bot marker. The static evidence covers the described local shell path,
not complete cloud or PowerShell behavior. The [public architecture documentation](https://cursor.com/docs/grok-bot/teams#architecture)
distinguishes the cloud service from the local execution daemon. Grok Bot has a
separate result from Grok Build.

A separate user-provided ordinary-shell observation on **2026-10-16** reported
a Grok Bot cloud context with `CURSOR_AGENT` and the exact flags
`SAND_MANIFEST_V2=1` and `SAND_SUPERVISOR_ENABLED=1`, without `SAND_AGENT`.
The cloud branch requires all three conditions together. The infrastructure
flags alone and `CURSOR_AGENT` alone do not match this branch. This observation
was not reproduced by launching the cloud environment for this contribution.
These shared markers do not establish exclusive brand ownership or resolve
nested agents' immediate identity. No tenant, authentication, store, boot or VNC
identifier is used by the rule or included in the fixtures.

## QwenWork

A user supplied the exact product marker
`QODER_WORK_INTEGRATION_PRODUCT=qwenworkcn` on **2026-10-16**. A subsequent
static inspection confirmed it in **QwenWorkCN 1.3.1**, bundle ID
`cn.qwenwork.desktop.mac`. This is a newer artifact than the previously
inspected 1.2.5 distribution.

- `Contents/Resources/app.asar` SHA-256:
  `f654cd3d68406139b326f5d10d3dbade1b85c610d1fee386c2f3dca43fb6596c`.
- Packaged `out/main/main.js` SHA-256:
  `76540a35496c1a0df187081b58c4d640aabc9b0f84e78fdbf7fb9fc36bc6eb3d`.
- At line 688, `buildHostCliDiscoveryEnv` supplies the product marker, using
  the `qwenworkcn` product constant from `chunks/index-CHipc8Ca.js`. At line
  697, `buildQoderWorkIntegrationEnv` passes it through
  `composeQoderAgentSdkRuntimeEnv`; these are the actual internal function names.
- At line 2330, `executeLocalShell` uses the constructed child environment.
  At line 2331, `executeQxcHostShell` rebuilds the trusted discovery environment.
  This establishes the inspected local/host shell paths. Propagation into
  ordinary VM/QXC guest shells was not established, and no agent was launched.

The rule checks for a non-empty `QODER_WORK_INTEGRATION_PRODUCT`, returning
`qwenwork`. The inspected 1.3.1 build sets the fixed value `qwenworkcn`, but the
rule does not limit accepted values to that string. It uses the existing
`env_set` semantics: missing and empty values do not match; all non-empty values,
including case variants and whitespace-only strings, match. Shared use of this
variable by other products has not been established, so this condition is not
proof of exclusive product identity.

The old `QODERWORK_SOURCE_CHAT_ID` and `QODERWORK_IS_CN` conditions remain
excluded. Session IDs, region/site flags, execution mode and service endpoints
do not independently identify the product, even when combined. The
[QwenWork introduction](https://qwenwork.cn/docs/product-introduction) identifies
the product but does not establish this variable as a stable public API.

## Cue

A user-provided Linux desktop-worker runtime observation on **2026-10-16**
identified Cue with `MANUS_LONG_DESKTOP_WORKER=1`. The variable prefix does not
establish a Manus product identity: the rule returns `cue` and accepts only
that exact value.
This observation was not reproduced by launching Cue during this contribution;
no inspected application version or artifact hash is claimed. Other `MANUS_`
path, endpoint or identifier variables are not independent detection conditions.
No real desktop identifier, path or endpoint is included in the evidence or tests.

## Muse

User-provided runtime observations on **2026-10-16** identified
`JARVIS_IS_ASSIGNED=1` in Muse's active-agent context and a non-empty
`JARVIS_TOOL_CALL_ID` during a tool call. Either condition identifies the
reported Muse context. The active-agent value must be exact `1`; the tool-call
ID is checked only for non-emptiness and is never reported by this detector.
Other instance, session, tier, timezone or region fields do not independently
match. These observations were not reproduced by launching Muse for this
contribution, and no inspected application version or artifact hash is claimed.
Fixtures use synthetic tool-call identifiers, never real runtime IDs.

## Manus

A user reported the marker names `MANUS_ADDON_RUNTIME_SPEC_PATH` and
`MANUS_ADDON_RUNTIME_SPECS_REVISION` in an ordinary shell on **2026-10-16**,
without providing their values. The conditions treat them as addon-runtime
descriptors and conservatively require non-empty values, returning `manus`.
This is user-reported shell evidence, not a locally reproduced shell probe or
validation of the previously inspected Manus Studio 2.0.7 artifact. The separate
local static review did not confirm these marker names or their propagation.
No `ENABLED` value is guessed.
Tests use synthetic descriptor values; no real spec path or revision is recorded.

## Doubao Work

A user-provided runtime observation on **2026-10-16** identified
`DOUBAO_OFFICE_AGENT_NAME` as the Doubao Work context marker. The rule uses
this variable's non-emptiness and returns `doubao-work`; it does not compare
the value to a particular agent name. The existing `env_set` behavior treats
missing and empty values as absent and any non-empty value, including whitespace,
as present. This observation was not independently reproduced by launching
Doubao Work for this contribution; no local artifact version, hash or official
stability guarantee is claimed.

Other `DOUBAO_OFFICE_` message, application, citation, proxy, credential, edition,
log and market fields are not detection conditions, separately or together.
Tests use only synthetic metadata and reserved `.invalid` endpoint examples.
No real token, message identifier or log identifier is included in these fixtures.

## Grok Build terminal marker

Public source at commit `2bdd1d6a6369de0e8c68132ea4539e9abd9e14a8`
[documents the dual use of `GROK_AGENT`](https://github.com/xai-org/grok-build/blob/2bdd1d6a6369de0e8c68132ea4539e9abd9e14a8/crates/codegen/xai-grok-tools/src/util/env.rs#L31):
the terminal-child marker is exact `1`, while other values can select an agent.
Both the [ordinary terminal path](https://github.com/xai-org/grok-build/blob/2bdd1d6a6369de0e8c68132ea4539e9abd9e14a8/crates/codegen/xai-grok-tools/src/computer/local/terminal.rs#L3315)
and the [sandbox command path](https://github.com/xai-org/grok-build/blob/2bdd1d6a6369de0e8c68132ea4539e9abd9e14a8/crates/codegen/xai-grok-tools/src/computer/local/terminal.rs#L3391)
apply that marker. This extends the existing `grok` result without treating
arbitrary non-empty selector values as a match.

## Junie launcher marker

Junie requires a non-empty `JUNIE_SHIM_PATH`. The
[launcher at commit `cdbfbcf`](https://github.com/JetBrains/junie/blob/cdbfbcf6a9f17951c94dd2902619dffc37f5881f/templates/junie.shim.sh#L684-L700)
exports its own absolute path in this variable. The launcher also accepts
`JUNIE_DATA` as configurable data-directory context and exports it; detection
intentionally ignores that variable. A missing or empty shim marker does not
match Junie, even when `JUNIE_DATA` is set. Other agent rules retain their
existing order and may still match.

All of these markers can be set or inherited by other processes. They identify
environment context, not an authenticated immediate caller. Existing override
and first-match behavior still apply when several contexts coexist.
