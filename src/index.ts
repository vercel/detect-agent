import spec from '../agents.json';
import { evaluateCondition } from './evaluate-condition';
import type {
  AgentResult,
  AgentSpec,
  KnownAgentNames,
  KnownAgents,
} from './types';

export type {
  AgentResult,
  KnownAgentDetails,
  KnownAgentNames,
  KnownAgents,
} from './types';

export const KNOWN_AGENTS = Object.fromEntries(
  spec.agents.map(({ key, name }) => [key, name])
) as KnownAgents;

const agents = spec.agents as unknown as AgentSpec[];
const aiAgentVar = spec.aiAgentVar;

/**
 * Resolve the AI_AGENT variable, which takes highest priority. Its trimmed
 * value is emitted verbatim as the agent name. Returns undefined when the
 * variable is unset or empty after trimming.
 */
function resolveAiAgentStandard(): KnownAgentNames | undefined {
  const raw = process.env[aiAgentVar];
  if (!raw) {
    return undefined;
  }
  const value = raw.trim();
  if (!value) {
    return undefined;
  }

  return value as KnownAgentNames;
}

export async function determineAgent(): Promise<AgentResult> {
  const aiAgentStandard = resolveAiAgentStandard();
  if (aiAgentStandard) {
    return { isAgent: true, agent: { name: aiAgentStandard } };
  }

  for (const agent of agents) {
    if (await evaluateCondition(agent.match)) {
      return { isAgent: true, agent: { name: agent.name as KnownAgentNames } };
    }
  }

  return { isAgent: false, agent: undefined };
}
