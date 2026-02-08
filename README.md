# Copilot Multi-Agent Refactor Guide

This repo demonstrates a simple, repeatable workflow for refactoring Go code with a set of specialized Copilot agents. The agents are defined under [.github/agents](.github/agents) and are designed to be orchestrated by a Conductor prompt.

## Agents

- `planner`: Produces a refactor plan, acceptance criteria, and risk notes.
- `refactorer`: Applies the plan with minimal, behavior-preserving changes.
- `critic`: Reviews the changes and blocks unless the result is approved.
- `tester`: Runs tests and reports failures or coverage gaps.
- `conductor`: Orchestrates the pipeline across all agents until approved.

## Workflow (Conductor-Driven)

Use the Conductor to run a gated pipeline:

1) Planner: draft plan + acceptance criteria + risks.
2) Refactorer: implement the plan.
3) Critic: review; if changes required, loop back to Refactorer.
4) Tester: run tests + report results; if failures, loop back.
5) Finish only when Critic approves and tests pass.

This sequence keeps changes safe, reviewable, and test-backed while encouraging small, idiomatic refactors.

## Conductor Prompt

Use this prompt when starting a refactor session with the Conductor agent:

```text
We are refactoring internal/billing/invoice.go. Constraints:
- No behavior changes.
- Keep diffs small.
- Must pass tests.

Orchestration:
1) Use Planner as a subagent to produce plan + acceptance criteria + risk notes.
2) Use Refactorer as a subagent to implement.
3) Use Critic as a subagent to review; block unless APPROVED.
4) Use Tester as a subagent to run tests / report failures.
5) Loop fix -> review -> test until APPROVED.
Return only final result + changelog.
```
