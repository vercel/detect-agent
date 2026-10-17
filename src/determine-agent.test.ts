import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { determineAgent, KNOWN_AGENTS } from './index';
import { access } from 'node:fs/promises';
import testCases from '../testcases.json';

vi.mock('node:fs/promises', () => ({ access: vi.fn() }));

vi.setConfig({ testTimeout: 6 * 60 * 1000 });

type TestCase = (typeof testCases)[number];

const ALL_AGENT_ENVS = [
  'AGENT_CONTEXT_OUT',
  'AI_AGENT',
  'ANTIGRAVITY_AGENT',
  'ANTIGRAVITY_CLI_ALIAS',
  'APPLICATION_NAME',
  'APP_NAME',
  'AUGMENT_AGENT',
  'CLAUDECODE',
  'CLAUDE_CODE',
  'CLAUDE_CODE_IS_COWORK',
  'CLINE_ACTIVE',
  'CODEBUDDY_CONFIG_DIR',
  'CODEX_CI',
  'CODEX_SANDBOX',
  'CODEX_SANDBOX_NETWORK_DISABLED',
  'CODEX_THREAD_ID',
  'CONFIG_DIR',
  'COPILOT_ALLOW_ALL',
  'COPILOT_GITHUB_TOKEN',
  'COPILOT_MODEL',
  'CRUSH',
  'CURSOR_AGENT',
  'CURSOR_EXTENSION_HOST_ROLE',
  'CURSOR_TRACE_ID',
  'DOUBAO_OFFICE_AGENT_NAME',
  'DOUBAO_OFFICE_ANSWER_MESSAGE_ID',
  'DOUBAO_OFFICE_APP_ID',
  'DOUBAO_OFFICE_CLI_CITATION',
  'DOUBAO_OFFICE_CLI_FORWARD_PROXY',
  'DOUBAO_OFFICE_CLI_TOKEN',
  'DOUBAO_OFFICE_EDITION',
  'DOUBAO_OFFICE_FORWARD_PROXY',
  'DOUBAO_OFFICE_LOG_ID',
  'DOUBAO_OFFICE_MARKET',
  'DOUBAO_OFFICE_PLATFORM_APP_ID',
  'GEMINI_CLI',
  'GOOSE_PROVIDER',
  'GOOSE_TERMINAL',
  'GROK_AGENT',
  'GROK_PLUGIN_DATA',
  'GROK_PLUGIN_ROOT',
  'HATCHLING_ID',
  'HERMES_AGENT',
  'HERMES_SESSION_ID',
  'IS_PACKAGED',
  'JARVIS_IS_ASSIGNED',
  'JARVIS_TOOL_CALL_ID',
  'JUNIE_DATA',
  'JUNIE_SHIM_PATH',
  'KILOCODE_FEATURE',
  'KIMI_CODE_HOME',
  'KIMI_PLUGIN_ROOT',
  'MANUS_ADDON_PREVIEW_PUBLIC_ORIGIN',
  'MANUS_ADDON_REGISTRY_ENABLED',
  'MANUS_ADDON_RUNTIME_PUBLIC_ORIGIN',
  'MANUS_ADDON_RUNTIME_SPECS_REVISION',
  'MANUS_ADDON_RUNTIME_SPEC_PATH',
  'MANUS_ADDON_SERVER_BRANCH',
  'MANUS_CONFIG_DIR',
  'MANUS_GROUP_ENABLED',
  'MANUS_HOME',
  'MANUS_LONG_DESKTOP_WORKER',
  'OPENCLAW_SHELL',
  'OPENCODE',
  'OPENCODE_CLIENT',
  'QODERCN_SERVER_ENDPOINT',
  'QODERWORK_EXECUTION_ENV',
  'QODERWORK_IS_CN',
  'QODERWORK_SOURCE_CHAT_ID',
  'QODER_SITE',
  'QODER_WORK_INTEGRATION_PRODUCT',
  'QWENWORK_BASE_URL',
  'REGION',
  'REPL_ID',
  'SAND_AGENT',
  'SAND_MANIFEST_V2',
  'SAND_SUPERVISOR_ENABLED',
  'SESSION_ID',
  'TERM_PROGRAM',
  'TIER',
  'TIMEZONE',
  'TRAE_AI_SHELL_ID',
  'VTCODE',
  'WORKBUDDY_APPLICATION_NAME',
  'WORKBUDDY_APP_NAME',
  'WORKBUDDY_CONFIG_DIR',
  'WORKBUDDY_IS_PACKAGED',
  'ZED_TERM',
];

describe('determineAgent', () => {
  const originalIsTTY = process.stdout.isTTY;
  const files = new Map<string, 'present' | 'EACCES'>();

  beforeEach(() => {
    files.clear();
    vi.mocked(access).mockReset();
    vi.mocked(access).mockImplementation(async path => {
      const state = files.get(String(path));
      if (state !== 'present') {
        const code = state ?? 'ENOENT';
        throw Object.assign(new Error(code), { code });
      }
    });
    vi.unstubAllEnvs();
    for (const v of ALL_AGENT_ENVS) {
      vi.stubEnv(v, '');
    }
    vi.stubEnv('PATH', '/usr/bin');
  });

  afterEach(() => {
    vi.unstubAllEnvs();
    process.stdout.isTTY = originalIsTTY;
    vi.mocked(access).mockReset();
  });

  for (const tc of testCases as TestCase[]) {
    it(tc.name, async () => {
      for (const [key, value] of Object.entries(tc.env)) {
        vi.stubEnv(key, value);
      }

      if (tc.tty !== undefined && tc.tty !== null) {
        process.stdout.isTTY = tc.tty;
      }

      if (tc.files && tc.files.length > 0) {
        for (const path of tc.files) {
          files.set(path, 'present');
        }
      }

      const result = await determineAgent();

      const expectedName =
        tc.expectedName ??
        (tc.expectedAgentKey
          ? (KNOWN_AGENTS as Record<string, string>)[tc.expectedAgentKey]
          : undefined);

      if (tc.expectedIsAgent) {
        expect(result).toEqual({
          isAgent: true,
          agent: { name: expectedName },
        });
      } else {
        expect(result).toEqual({ isAgent: false });
      }
    });
  }

  it('handles file system errors gracefully for devin detection', async () => {
    files.set('/opt/.devin', 'EACCES');
    const result = await determineAgent();
    expect(result).toEqual({ isAgent: false });
  });
});
