# Agent Coordination

## Available Agent Threads

Four pinned threads were created from the current migration context:

- `Canopsis Migration: Forms & Validation`
- `Canopsis Migration: Pages & Modals`
- `Canopsis Migration: Stores & API`
- `Canopsis Migration: QA & Regressions`

They are separate chats. They do not have automatic file locks.

## Recommended Ownership

### Forms & Validation

Primary scope:

- `src/components/forms/**`
- `src/hooks/form/**`
- `src/hooks/validator/**`
- `src/mixins/form/**`
- `src/helpers/date/**`

Read-only source:

- Old form fields under `canopsis-next/src/components/forms/**`

### Pages & Modals

Primary scope:

- `src/views/**`
- `src/components/other/**`
- `src/components/modals/**`
- `src/router/**`
- `docs/migration-inventory.md`

### Stores & API

Primary scope:

- `src/stores/**`
- `src/hooks/store/**`
- `src/services/**`
- `src/helpers/entities/**`
- API route constants

### QA & Regressions

Primary scope:

- Read-only by default.
- Runtime stack traces.
- Targeted lint/build/test checks.
- Old/new behavior comparison.

## File Ownership Protocol

Before assigning work, include:

```text
Ownership:
- You may edit:
  - src/path/**
- You may read:
  - src/other-path/**
  - /Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next/src/path/**
- Do not touch:
  - src/protected-path/**

Do not edit outside ownership without asking.
List changed files at the end.
```

## Conflict Avoidance

- One file should have one active owner.
- Shared barrels such as `src/components/forms/fields/index.js` should be edited by one agent at a time.
- Shared constants and i18n files are high-conflict; coordinate before editing.
- If two tasks need the same file, merge them into one task or let the main thread coordinate.

## Suggested Workflow

1. Main thread decides scope and owner.
2. Agent performs bounded work.
3. Agent reports changed files and verification.
4. Main thread integrates or asks for targeted fixes.
5. QA agent can review after integration.

## Example Assignment

```text
Task: migrate pbehavior exception date range field.

Ownership:
- You may edit:
  - src/components/other/pbehavior/**
  - src/components/forms/fields/date-time-picker/**
- You may read:
  - /Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next/src/components/other/pbehavior/**
- Do not touch:
  - src/components/layout/**
  - src/stores/**

Use existing date/time picker fields.
Do not start dev server.
Run targeted eslint and npm run build.
List changed files.
```
