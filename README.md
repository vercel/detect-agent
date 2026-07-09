# detect-agent

A lightweight utility for detecting if code is being executed by an AI agent or automated development environment.

## Installation

```bash
npm install detect-agent
```

## Usage

```typescript
import { determineAgent } from 'detect-agent';

const { isAgent, agent } = await determineAgent();

if (isAgent) {
  console.log(`Running in ${agent.name} environment`);
  // Adapt behavior for AI agent context
}
```

## Agent Definitions

The agent detection rules live in [`agents.json`](https://raw.githubusercontent.com/vercel/detect-agent/main/agents.json), validated against [`agents.schema.json`](https://raw.githubusercontent.com/vercel/detect-agent/main/agents.schema.json).

> [!NOTE]
> You can use these definitions two ways: install the [npm package](https://www.npmjs.com/package/detect-agent) for a ready-to-use API, **or** consume the JSON file directly from GitHub if you're not in a JavaScript environment or want to build your own detection logic:
>
> ```
> https://raw.githubusercontent.com/vercel/detect-agent/main/agents.json
> ```

## Supported Agents

This package can detect the following AI agents and development environments:

- **Cursor** (Anysphere)
- **Claude Code** (Anthropic)
- **Claude Cowork** (Anthropic)
- **Devin** (Cognition Labs)
- **Gemini CLI** (Google)
- **Codex** (OpenAI)
- **Antigravity** (Google DeepMind)
- **Augment**
- **OpenCode**
- **Goose** (Block)
- **Junie** (JetBrains)
- **Kiro** (AWS)
- **Pi**
- **GitHub Copilot**
- **Replit**
- **Custom agents** (via the `AI_AGENT` environment variable)

See [`agents.json`](https://raw.githubusercontent.com/vercel/detect-agent/main/agents.json) for the exact environment variables and conditions used to detect each one.

## The AI_AGENT Standard

We're promoting `AI_AGENT` as a universal environment variable standard for AI development tools. This allows any tool or library to easily detect when it's running in an AI-driven environment.

### For AI Tool Developers

Set the `AI_AGENT` environment variable to identify your tool:

```bash
export AI_AGENT="your-tool-name"
# or
AI_AGENT="your-tool-name" your-command
```

### Recommended Naming Convention

- Use lowercase with hyphens for multi-word names
- Include version information if needed, separated by an `@` symbol
- Examples: `claude-code`, `cursor-cli`, `devin@1`, `custom-agent@2.0`

## Use Cases

### Adaptive Behavior

```typescript
import { determineAgent } from 'detect-agent';

async function setupEnvironment() {
  const { isAgent, agent } = await determineAgent();

  if (isAgent) {
    // Running in AI environment - adjust behavior
    process.env.LOG_LEVEL = 'verbose';
    console.log(`🤖 Detected AI agent: ${agent.name}`);
  }
}
```

### Telemetry and Analytics

```typescript
import { determineAgent } from 'detect-agent';

async function trackUsage(event: string) {
  const result = await determineAgent();

  analytics.track(event, {
    agent: result.isAgent ? result.agent.name : 'human',
    timestamp: Date.now(),
  });
}
```

### Feature Toggles

```typescript
import { determineAgent } from 'detect-agent';

async function shouldEnableFeature(feature: string) {
  const result = await determineAgent();

  // Enable experimental features for AI agents
  if (result.isAgent && feature === 'experimental-ai-mode') {
    return true;
  }

  return false;
}
```

## Contributing

We welcome contributions!

### Adding New Agent Support

To add support for a new AI agent:

1. Add the agent's detection rule to [`agents.json`](./agents.json), following the structure defined in [`agents.schema.json`](./agents.schema.json). Agents are evaluated in array order — the first match wins, so place more specific rules before more general ones.
2. Add test cases in `src/determine-agent.test.ts`.
3. Add the agent to the **Supported Agents** list above.

## Links

- [GitHub Repository](https://github.com/vercel/detect-agent)
- [npm Package](https://www.npmjs.com/package/detect-agent)
- [Vercel Documentation](https://vercel.com/docs)
