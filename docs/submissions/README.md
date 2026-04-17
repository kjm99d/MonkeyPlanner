# External Registry Submissions

Ready-to-paste PR bodies and listing entries for every directory we want
MonkeyPlanner to appear in. Each file is dedicated to one upstream repo;
open the linked PR form, paste the snippet, and you are done.

| Target | File | Priority | Expected impact |
|--------|------|----------|-----------------|
| `modelcontextprotocol/servers` | [`mcp-official-servers.md`](./mcp-official-servers.md) | **P0** | Single largest source of MCP-aware discovery |
| `awesome-mcp-servers` (punkpeye/appcypher) | [`awesome-mcp-servers.md`](./awesome-mcp-servers.md) | P0 | GitHub Trending regular |
| `mcp.so` | [`mcp-so.md`](./mcp-so.md) | P1 | SEO — indexed by search engines |
| `awesome-claude-code` | [`awesome-claude-code.md`](./awesome-claude-code.md) | P1 | Target persona directly |
| `awesome-ai-agents` | [`awesome-ai-agents.md`](./awesome-ai-agents.md) | P1 | Broader agent ecosystem |
| `awesome-self-hosted` | [`awesome-self-hosted.md`](./awesome-self-hosted.md) | P2 | Local-first / privacy audience |
| `awesome-go` | [`awesome-go.md`](./awesome-go.md) | P2 | Go contributors |

## Order of operations

1. Submit `modelcontextprotocol/servers` first — it has the largest audience
   and their review turnaround is fast.
2. While that PR is open, submit `awesome-mcp-servers` and `mcp.so` in
   parallel — different maintainers, no conflict.
3. After the first listing lands and stars start arriving, submit the rest.
   Fresh stars help prove the project is active, which some awesome-list
   maintainers explicitly look for.

## Shared boilerplate

Every submission leads with the same one-liner so the value prop is
consistent across directories:

> **MonkeyPlanner — Local-first task memory for your AI coding agents.
> Approve with a click; Claude Code, Cursor, and any MCP client do the
> rest.**
