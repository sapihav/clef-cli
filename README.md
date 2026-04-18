# clef

Per-project Claude Code model and effort configuration. Sets `.claude/settings.local.json` in the current directory without touching your global settings.

## Why

Claude Code's `/effort` and `/model` commands write to `~/.claude/settings.json` — global, instant, affects every open session. `clef` writes to the project-local settings file instead, letting you pin different model+effort combinations per project or git worktree.

## Install

```bash
go install github.com/sapihav/clef@latest
```

Requires Go 1.21+. Binary installs to `$GOPATH/bin` (usually `~/go/bin`).

## Usage

```bash
clef set --model opus --effort xhigh   # pin model and effort for this directory
clef set --effort medium               # change effort only
clef set --model sonnet                # change model only
clef show                              # print current local settings
clef reset                             # remove model and effortLevel overrides
```

Valid values:
- `--model`: `sonnet`, `opus`, `haiku`
- `--effort`: `low`, `medium`, `high`, `xhigh`

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
clef set --model opus --effort xhigh

# Quick bug fix — Sonnet is fine
git worktree add ../myproject-hotfix -b fix/typo
cd ../myproject-hotfix
clef set --model sonnet --effort medium
```

Each worktree directory maintains its own `.claude/settings.local.json`.

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | Success |
| `1` | User/config error (invalid flag value, nothing to reset, etc.) |

## License

MIT — see [LICENSE](LICENSE).
