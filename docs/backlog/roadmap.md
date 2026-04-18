# clef-cli improvements roadmap

Current commands: `set --model --effort`, `show`, `reset`

---

## M1 — Dedicated model/effort subcommands + version support

**Goal:** natural UX for the two most common operations; fix model version rejection.

### Changes

**`clef model <value>`** — new command, sets `model` key
- Positional arg, no flags needed
- Replaces the need for `clef set --model`

**`clef effort <value>`** — new command, sets `effortLevel` key
- Positional arg
- Replaces the need for `clef set --effort`

**Model validation overhaul** (fixes `clef set --model opus4.6` bug):

| Input | Stored as | Notes |
|---|---|---|
| `sonnet` | `sonnet` | existing alias, keep |
| `opus` | `opus` | existing alias, keep |
| `haiku` | `haiku` | existing alias, keep |
| `opus4.7` | `claude-opus-4-7` | shorthand alias |
| `opus4.6` | `claude-opus-4-6` | shorthand alias |
| `sonnet4.6` | `claude-sonnet-4-6` | shorthand alias |
| `haiku4.5` | `claude-haiku-4-5` | shorthand alias |
| `claude-opus-4-7` | `claude-opus-4-7` | full ID passthrough |
| `claude-sonnet-4-6` | `claude-sonnet-4-6` | full ID passthrough |
| `claude-haiku-4-5` | `claude-haiku-4-5` | full ID passthrough |

Accept any `claude-*` string as passthrough (future-proofs new releases without code changes). Keep the short family names as convenience aliases.

**`set` stays** — backward compat, delegates to the same logic.

**Effort validation** — no change needed (low/medium/high/xhigh are not versioned).

### Estimate
~150 LoC: `cmd/model.go`, `cmd/effort.go`, updated validation in `cmd/set.go`.

---

## M2 — AI agent friendliness

**Goal:** make clef scriptable and agent-discoverable, following posthog-cli baseline.

### Changes

**`clef schema`** — emit full command tree as JSON
```json
{
  "cli": "clef",
  "version": "0.3.0",
  "description": "...",
  "commands": [
    {
      "name": "model",
      "path": "model",
      "description": "Set Claude model for this project",
      "arguments": [{"name": "value", "required": true}],
      "options": [],
      "output": {"type": "object", "fields": {"model": "string", "file": "string"}}
    },
    ...
  ]
}
```

**`--json` global flag** — structured stdout on every command
- `clef model opus4.7 --json` → `{"model":"claude-opus-4-7","file":".claude/settings.local.json"}`
- `clef show --json` → already JSON, no change in content, just honor the flag
- `clef reset --json` → `{"removed":["model","effortLevel"],"file":"..."}`

**`--json-errors` global flag** — structured errors to stderr
- `clef model invalid --json-errors` → `{"error":"invalid model \"invalid\"","code":1}`

**Exit codes** — align to workspace convention:
- `0` success (unchanged)
- `1` user/config error (unchanged)
- `2` system error (read/write failure) — currently also `1`, split these

**`--dry-run` on mutating commands** (`model`, `effort`, `set`, `reset`)
- Prints the JSON that would be written, no disk write

### Estimate
~200 LoC: `cmd/schema.go`, global flag wiring in `cmd/root.go`, output helpers.

---

## Out of scope (YAGNI)
- Interactive TUI for model selection
- Remote settings sync
- Shell completion (cobra has built-in, enable if asked)
