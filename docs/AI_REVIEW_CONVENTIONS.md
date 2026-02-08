
# AI Review Conventions (Go)

This document defines **mandatory, enforceable rules** for AI-assisted refactors
in this repository.

These rules are primarily enforced by the **Critic agent** in the Copilot
multi-agent pipeline, but they also serve as human-readable guardrails.

The goal is simple:  
👉 **No refactor is approved without execution evidence.**

---

## Scope

These conventions apply to:
- Any AI-assisted refactor
- Any non-trivial code restructuring
- Any change submitted for review by an AI Critic agent

They apply regardless of perceived refactor size.

---

## Mandatory Gates (Hard Rules)

### 1. Unit tests must pass

No refactor may be approved unless **all unit tests pass**:

```bash
go test ./...
````

#### Enforcement

* If test results are missing → ❌ **CHANGES REQUIRED**
* If tests fail → ❌ **CHANGES REQUIRED**
* Code inspection alone is **never sufficient**

If tests cannot be run due to environment limitations, the review must explicitly state:

> ❌ CHANGES REQUIRED — tests not executed

---

### 2. Race detector required for concurrency code

If a refactor:

* touches concurrency code, OR
* modifies code that already uses concurrency

then approval **requires**:

```bash
go test -race ./...
```

#### Enforcement

* If concurrency is involved and the race detector was not run → ❌ **CHANGES REQUIRED**
* Exceptions must be explicitly justified (e.g. unsupported platform)

---

### 3. No “looks good” approvals

The following are **not valid approval reasons**:

* “Looks good”
* “Seems safe”
* “No obvious issues”
* “Small refactor”

Approval must be based on:

* Tests passing
* Behavior preservation verified
* Risks explicitly assessed

---

## Definition of Concurrency Code

Code is considered **concurrency-related** if it includes or modifies any of:

* `go func()`
* Channels
* `sync.Mutex`, `sync.RWMutex`
* `sync.WaitGroup`
* `sync/atomic`
* Goroutine lifecycle management
* `context.Context` used to control goroutines

If in doubt → treat the change as concurrency-related.

---

## Reviewer (Critic Agent) Responsibilities

The Critic agent **must**:

* Enforce all mandatory gates in this document
* Block approval if any gate is unmet
* Explicitly cite which gate failed
* Never approve based on confidence or inspection alone
* Require execution evidence (test output or summary)

The Critic agent **may not**:

* Assume correctness
* Skip gates for “small” changes
* Approve without test results

---

## Tester Agent as Source of Truth

The Tester agent is responsible for producing **authoritative execution evidence**.

Expected Tester output includes:

* Commands executed
* PASS / FAIL status
* Key output lines (especially failures)

The Critic agent must rely on this evidence when issuing a verdict.

---

## Philosophy

These conventions exist to prevent:

* Silent regressions
* False confidence
* AI “politeness bias”
* Review fatigue

They intentionally make AI reviews **stricter than human reviews**.

> If it didn’t run, it didn’t happen.

---

## Summary (TL;DR)

* ❌ No `go test ./...` → no approval
* ❌ Concurrency + no `-race` → no approval
* ❌ “Looks good” → invalid
* ✅ Evidence-based approval only

This keeps the pipeline honest and the codebase safe.



