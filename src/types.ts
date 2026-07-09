import spec from '../agents.json';

// The list of known agents (their stable keys and emitted names) is the
// source of truth in agents.json. `KnownAgents` and `KnownAgentNames` below
// are both derived from it, so adding an agent to the JSON updates both.
type Agent = (typeof spec.agents)[number];

/**
 * Union of every canonical agent name declared in agents.json. Custom names
 * supplied via the AI_AGENT variable are not part of this union.
 */
export type KnownAgentNames = Agent['name'];

/**
 * Map of UPPER_SNAKE key -> canonical agent name, e.g.
 * `KNOWN_AGENTS.CURSOR_CLI === 'cursor-cli'`. Derived from agents.json.
 */
export type KnownAgents = {
  [A in Agent as A['key']]: A['name'];
};

export interface KnownAgentDetails {
  name: KnownAgentNames;
}

export type AgentResult =
  | {
      isAgent: true;
      agent: KnownAgentDetails;
    }
  | {
      isAgent: false;
      agent: undefined;
    };

// Structural types for the conditions loaded from agents.json. These mirror
// agents.schema.json; the JSON is the source of truth for the actual logic.
interface EnvSetCondition {
  type: 'env_set';
  name: string;
}
interface EnvValueCondition {
  type: 'env_value';
  name: string;
  value: string;
}
interface EnvMatchesCondition {
  type: 'env_matches';
  name: string;
  pattern: string;
}
interface FileExistsCondition {
  type: 'file_exists';
  path: string;
}
interface NoTtyCondition {
  type: 'no_tty';
}
interface AnyOfCondition {
  type: 'anyOf';
  conditions: Condition[];
}
interface AllOfCondition {
  type: 'allOf';
  conditions: Condition[];
}
export type Condition =
  | EnvSetCondition
  | EnvValueCondition
  | EnvMatchesCondition
  | FileExistsCondition
  | NoTtyCondition
  | AnyOfCondition
  | AllOfCondition;

export interface AgentSpec {
  key: string;
  name: string;
  match: Condition;
}
