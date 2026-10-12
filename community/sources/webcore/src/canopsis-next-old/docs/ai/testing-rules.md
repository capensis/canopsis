# Testing Rules

## Default Verification

For most migration tasks:

1. Run targeted eslint on changed files.
2. Run `npm run build`.
3. Run focused tests only if they exist and are relevant.

Do not run the dev server unless the user explicitly asks.

## Full Lint

`npm run lint` runs Oxlint and ESLint with auto-fix. During the migration it may fail on unrelated existing files. If that happens:

- Do not fix unrelated files unless asked.
- Run targeted eslint for the files touched by the task.
- Report the unrelated full-lint failures briefly.

## Unit Tests

Current test stack:

- Vitest
- Vue Test Utils 2
- jsdom

Commands:

```bash
npm run test:unit
npx vitest <path>
```

Only add or update tests when the user explicitly asks.

## Existing Test Preservation

When changing functionality that is already covered by tests, do not update, rewrite, delete,
regenerate, or otherwise modify those tests as part of the functional change. Keep the existing
tests untouched until the user explicitly asks to change them, even if the functional change makes
an existing expectation obsolete or causes a test to fail.

Running the relevant tests is still allowed. Report failures caused by the functional change
without adjusting the tests.

## E2E Tests

Current E2E stack:

- Playwright

Command:

```bash
npm run test:e2e
```

Do not run browser/dev-server workflows unless explicitly asked.

## Migration QA Checklist

For migrated UI components:

- Compare behavior against the old source component.
- Verify props, events and slots.
- Verify `@update:<prop>` / `v-model:<prop>` bindings; do not use `v-field` or component `@input`.
- Check Vuetify 4 slot API differences.
- Check labels, validation messages, disabled/read-only state and clear buttons.
- Check layout gaps: vertical 2, horizontal 3.

For form modals:

- Validate submit path.
- Validate required fields.
- Validate click-outside confirmation behavior when form is dirty.
- Validate that no request happens directly from hooks.
- Validate entity form-to-api mapper output.

For stores:

- Check old request route and payload shape.
- Check response extraction (`content`, `data`, `_embedded`, etc.).
- Check hook naming and exposed refs.
- Check list meta/pending/fetching params.

For tables:

- Check pagination.
- Check sorting forward and reverse.
- Check expanded rows.
- Check mass actions.
- Check row hover and borders.
- Check empty state.

## Runtime Error Triage

When the user provides a console stack trace:

1. Identify the top project file in the stack.
2. Inspect the migrated component and old source equivalent.
3. Fix the smallest compatibility mismatch.
4. Run targeted eslint and build.
5. Report changed files and verification.

## Snapshot Testing

Do not add broad snapshots for unstable migration UI unless the user asks. Prefer behavioral tests for helpers, stores and form mappers.
