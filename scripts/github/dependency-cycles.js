'use strict';

const fs = require('fs');

/**
 * Extrahiert Issue-IDs aus einem Markdown-Textblock mit Dependency-Hinweisen.
 * Unterstützt Heading-Block "## Dependencies", Zeilen "Depends on:" sowie Tasklisten.
 * @param {string} text
 * @returns {Set<number>} Menge der referenzierten Issue-IDs
 */
function extractDependencyIdsFromText(text) {
  const ids = new Set();
  if (!text) {
    return ids;
  }

  const blockRegex = /(##+\s+Dependencies[\s\S]*?)(?=\n## |$)/i;
  let block = '';
  const match = text.match(blockRegex);
  if (match) {
    block = match[1];
  } else {
    const fallbackLines = text
      .split(/\n/)
      .filter((line) => /^Depends on:/i.test(line.trim()));
    if (fallbackLines.length) {
      block = fallbackLines.join('\n');
    }
  }

  if (!block) {
    return ids;
  }

  const directRefs = block.match(/#[0-9]+/g) || [];
  for (const ref of directRefs) {
    const value = parseInt(ref.substring(1), 10);
    if (!Number.isNaN(value)) {
      ids.add(value);
    }
  }

  for (const taskMatch of block.matchAll(/- \[[ xX]\]\s*#([0-9]+)/g)) {
    const value = parseInt(taskMatch[1], 10);
    if (!Number.isNaN(value)) {
      ids.add(value);
    }
  }

  return ids;
}

/**
 * Lädt Issues eines Repos (pull requests werden gefiltert).
 * @param {import('@actions/github').GitHub} github
 * @param {string} owner
 * @param {string} repo
 * @param {{state?: string, perPage?: number}} [options]
 * @returns {Promise<Array<object>>}
 */
async function fetchAllIssues(github, owner, repo, options = {}) {
  const { state = 'open', perPage = 100 } = options;
  const issues = [];
  let page = 1;
  while (true) {
    const { data } = await github.rest.issues.listForRepo({
      owner,
      repo,
      state,
      per_page: perPage,
      page,
    });
    if (!data.length) {
      break;
    }
    for (const item of data) {
      if (!item.pull_request) {
        issues.push(item);
      }
    }
    page += 1;
  }
  return issues;
}

/**
 * Baut einen gerichteten Graphen aus Issues -> Dependencies.
 * @param {Array<object>} issues
 * @param {(body: string) => Set<number>} extractor
 * @returns {Map<number, Set<number>>}
 */
function buildDependencyGraph(issues, extractor = extractDependencyIdsFromText) {
  const graph = new Map();
  for (const issue of issues) {
    const deps = extractor(issue.body || '');
    if (deps.size) {
      graph.set(issue.number, deps);
    }
  }
  return graph;
}

/**
 * Normalisiert einen Zyklus (Rotation zum kleinsten Knotenwert) und erzeugt einen Key.
 * @param {Array<number>} segment
 * @returns {{key: string, nodes: Array<number>} | null}
 */
function normalizeCycle(segment) {
  if (!segment || segment.length < 2) {
    return null;
  }
  const rotated = segment.slice();
  const numeric = rotated.map((value) => Number(value));
  const min = Math.min(...numeric);
  if (!Number.isFinite(min)) {
    return null;
  }
  while (Number(rotated[0]) !== min && rotated.length) {
    rotated.push(rotated.shift());
  }
  const key = rotated.join('->');
  return { key, nodes: rotated };
}

/**
 * Findet alle Zyklen des Abhängigkeitsgraphen per DFS.
 * @param {Map<number, Set<number>>} graph
 * @returns {Array<{key: string, nodes: Array<number>}>}
 */
function findCycles(graph) {
  const cycles = [];
  const seen = new Set();
  const visited = new Set();
  const stack = new Set();
  const path = [];

  function recordCycle(segment) {
    const normalized = normalizeCycle(segment);
    if (!normalized) {
      return;
    }
    if (seen.has(normalized.key)) {
      return;
    }
    seen.add(normalized.key);
    cycles.push(normalized);
  }

  function dfs(node) {
    if (stack.has(node)) {
      const idx = path.indexOf(node);
      if (idx !== -1) {
        recordCycle(path.slice(idx).concat(node));
      }
      return;
    }
    if (visited.has(node)) {
      return;
    }
    visited.add(node);
    stack.add(node);
    path.push(node);
    const neighbors = graph.get(node) || new Set();
    for (const dep of neighbors) {
      dfs(dep);
    }
    stack.delete(node);
    path.pop();
  }

  for (const node of graph.keys()) {
    dfs(node);
  }

  return cycles;
}

/**
 * Erzeugt ein Mermaid-Diagramm aus Zyklen (max. limit Zyklen für Übersicht).
 * @param {Array<{nodes: Array<number>}>} cycles
 * @param {number} [limit]
 * @returns {string}
 */
function createMermaidDiagram(cycles, limit = 20) {
  if (!cycles.length) {
    return 'graph LR\n%% no cycles detected\n';
  }
  const edges = new Set();
  let mermaid = 'graph LR\n';
  for (const { nodes } of cycles.slice(0, limit)) {
    for (let i = 0; i < nodes.length - 1; i += 1) {
      const edgeKey = `${nodes[i]}-->${nodes[i + 1]}`;
      if (!edges.has(edgeKey)) {
        mermaid += `  ${nodes[i]} --> ${nodes[i + 1]}\n`;
        edges.add(edgeKey);
      }
    }
  }
  return mermaid;
}

/**
 * Schreibt JSON- und Mermaid-Artefakte für gefundene Zyklen.
 * @param {Array<{key: string, nodes: Array<number>}>} cycles
 * @param {{prefix?: string, limit?: number}} [options]
 * @returns {{jsonPath: string, mermaidPath: string}}
 */
function writeCycleArtifacts(cycles, options = {}) {
  const { prefix = 'cycles', limit = 20 } = options;
  const jsonPath = `${prefix}.json`;
  const mermaidPath = `${prefix}.mmd`;
  const mermaid = createMermaidDiagram(cycles, limit);
  fs.writeFileSync(jsonPath, JSON.stringify(cycles, null, 2));
  fs.writeFileSync(mermaidPath, mermaid);
  return { jsonPath, mermaidPath };
}

/**
 * Convenience-Helfer: lädt Issues, baut Graph und liefert Zyklen plus Mermaid.
 * @param {import('@actions/github').GitHub} github
 * @param {string} owner
 * @param {string} repo
 * @param {{state?: string, limit?: number}} [options]
 * @returns {Promise<{cycles: Array<{key: string, nodes: Array<number>}>, issues: Array<object>, graph: Map<number, Set<number>>, mermaid: string}>}
 */
async function collectDependencyCycles(github, owner, repo, options = {}) {
  const { state = 'open', limit = 20 } = options;
  const issues = await fetchAllIssues(github, owner, repo, { state });
  const graph = buildDependencyGraph(issues);
  const cycles = findCycles(graph);
  const mermaid = createMermaidDiagram(cycles, limit);
  return { cycles, issues, graph, mermaid };
}

module.exports = {
  extractDependencyIdsFromText,
  fetchAllIssues,
  buildDependencyGraph,
  findCycles,
  createMermaidDiagram,
  writeCycleArtifacts,
  collectDependencyCycles,
};
