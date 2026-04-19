# clef

Per-project Claude Code model and effort configuration. Sets `.claude/settings.local.json` in the current directory without touching your global settings.

## Why

Claude Code's `/effort` and `/model` commands write to `~/.claude/settings.json` — global, instant, affects every open session. `clef` writes to the project-local settings file instead, letting you pin different model+effort combinations per project or git worktree.

## Install

**Homebrew (macOS, Linux):**

```bash
brew install sapihav/tap/clef
```

**One-line installer (macOS, Linux):**

```bash
curl -fsSL https://raw.githubusercontent.com/sapihav/clef-cli/main/install.sh | bash
```

Installs to `/usr/local/bin` if writable, otherwise falls back to `$HOME/.local/bin`. Override with `INSTALL_DIR=/custom/path bash install.sh`.

**Manual download:** grab a tarball from [releases](https://github.com/sapihav/clef-cli/releases) and drop the binary on your `$PATH`.

**From source (Go 1.21+):**

```bash
go install github.com/sapihav/clef-cli@latest
```

## Usage

```bash
clef model opus4.7          # set model for this directory
clef effort xhigh           # set effort level
clef show                   # print current local settings
clef reset                  # remove model and effortLevel overrides
clef schema                 # print full command tree as JSON (agent-friendly)
```

`set` is also available for backward compatibility:

```bash
clef set --model sonnet --effort medium
```

### Model values

| Input | Stored as |
|-------|-----------|
| `sonnet` | `sonnet` |
| `opus` | `opus` |
| `haiku` | `haiku` |
| `opus4.7` | `claude-opus-4-7` |
| `opus4.6` | `claude-opus-4-6` |
| `sonnet4.6` | `claude-sonnet-4-6` |
| `haiku4.5` | `claude-haiku-4-5` |
| `claude-*` (any) | passed through as-is |

### Effort values

`low`, `medium`, `high`, `xhigh`

## Global flags

| Flag | Effect |
|------|--------|
| `--json` | Structured JSON to stdout |
| `--json-errors` | Structured JSON errors to stderr |
| `--dry-run` | Print what would be written, no disk write (mutating commands only) |

Examples:

```bash
clef model opus4.7 --json
# {"model":"claude-opus-4-7","file":".claude/settings.local.json"}

clef reset --dry-run
# {"would_remove":["model","effortLevel"],"file":".claude/settings.local.json"}

clef model invalid --json-errors
# stderr: {"error":"invalid model \"invalid\"","code":1}
```

## How it works

Reads and writes `.claude/settings.local.json` in the current working directory. Claude Code's settings hierarchy is:

```
~/.claude/settings.json          (global)
  ↓
.claude/settings.json            (project, committed)
  ↓
.claude/settings.local.json      (local override, gitignored) ← clef writes here
```

`clef set` merges into the existing file — other keys (permissions, hooks, etc.) are preserved.

## Example: per-worktree config

```bash
# Heavy feature branch — use Opus with max thinking
git worktree add ../myproject-feature -b feat/big-thing
cd ../myproject-feature
clef model opus4.7 && clef effort xhigh

# Quick bug fix — Sonnet is fine
git worktree add ../myproject-hotfix -b fix/typo
cd ../myproject-hotfix
clef model sonnet && clef effort medium
```

Each worktree directory maintains its own `.claude/settings.local.json`.

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | User/config error (invalid value, nothing to reset, etc.) |
| `2` | System error (file read/write failure) |

## License

MIT — see [LICENSE](LICENSE).
