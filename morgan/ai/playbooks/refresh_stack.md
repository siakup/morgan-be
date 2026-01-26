# Playbook: Refresh Stack

## Purpose
Update the `docs/STACK.md` and `ai/rules/` to reflect the current reality of the codebase.
Use this when the codebase has drifted from documentation or after major refactors.

## Triggers
- User request "Refresh Stack".
- Periodic maintenance.
- After a large merge that changes architecture or dependencies.

## Inputs
- Current codebase.

## Preconditions
- Access to full workspace.

## Rules to Load
- [ai/rules/00_principles.md](../rules/00_principles.md)

## Templates to Use
- None.

## Execution Phases

### Phase 0 — Environment & Context Check
1. Read `go.mod` for dependency updates.
2. Read `Dockerfile` for OS/Env updates.
3. Scan `cmd/` and `module/` for structural changes.

### Phase 1 — Analysis & Decomposition
1. Compare `docs/STACK.md` against findings.
2. Identify new patterns that should be codified into rules.
3. Identify obsolete rules.

### Phase 2 — Implementation
1. **Update STACK.md**:
   - Update Go version, libraries, and versions.
   - Update "Key Directories" if changed.
2. **Update Rules**:
   - If a new pattern is dominant (e.g., everyone uses `slog` instead of `zerolog`), update the rule.
   - Create new rules for new standard components.

### Phase 3 — Validation
1. **Self-Review**: Does the documentation match the code?
2. **Checklist**: Execute `ai/checklists/definition_of_done.md` (Documentation section).

## Outputs
- Updated `docs/STACK.md`
- Updated `ai/rules/*.md`

## Failure Modes
- **Hallucination**: Documenting libraries that are not actually used.
- **Over-prescription**: Creating rules for one-off hacks. Only document *patterns*.
