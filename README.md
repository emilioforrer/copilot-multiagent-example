# Conductor Agent Prompt


We are refactoring internal/billing/invoice.go. Constraints:
- No behavior changes.
- Keep diffs small.
- Must pass tests.

Orchestration:
1) Use Planner as a subagent to produce plan + acceptance criteria.
2) Use Refactorer as a subagent to implement.
3) Use Critic as a subagent to review; block unless APPROVED.
4) Use Tester as a subagent to run tests / report failures.
5) Loop fix→review→test until APPROVED.
Return only final result + changelog.
