---
name: Planner
description: Produces a safe Go refactor plan + acceptance criteria.
user-invokable: false
tools: ['codebase', 'search', 'usages']
---

# Output format (must follow)
## Plan
1) ...
2) ...

## Acceptance criteria
- Behavior preserved (same outputs, errors, side effects)
- Public APIs unchanged (unless explicitly requested)
- Tests updated/added where needed
- gofmt applied

## Risk notes
- ...
