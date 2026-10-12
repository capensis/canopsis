# Decision Log

## 2026-06-24: Rewrite Strategy

The project is being rewritten from scratch on the current stack instead of doing a direct in-place Vue 2 to Vue 3 migration.

Consequences:

- New code lives in `/Users/verokhov/Workspace/canopsis-pro-vue-3/canopsis`.
- Old behavior is ported from `canopsis-next`.
- Migration can introduce new wrappers and Pinia stores, but old feature behavior should remain intact.

## 2026-06-24: JavaScript During Active Migration

The active migration phase prefers JavaScript and avoids TypeScript churn.

Consequences:

- JS files are acceptable and expected.
- Files with no TypeScript-specific syntax can be converted from TS to JS.
- ESLint must work for JS and TS.

## 2026-06-24: Preserve Component Style

Components should not be rewritten just because they are old.

Consequences:

- Old Options API components stay Options API when possible.
- Old Composition API components stay Composition API.
- Rewrites are allowed only where the old library contract no longer works.

## 2026-06-24: New Form Layer

`src/hooks/form` separates form values from vee-validate.

Consequences:

- `form.values` is the only source of form data.
- vee-validate handles validation state, errors, meta and submit only.
- Subforms use `form.model(path)`.
- Field wrappers use `useField()` only for validation and error rendering.
- Do not synchronize every value into vee-validate with `setFieldValue`.

## 2026-06-24: Store Access Through Hooks

Components should not call Pinia stores directly.

Consequences:

- Store modules own request logic.
- Hooks in `src/hooks/store/modules` expose store refs and methods.
- Hooks do not make server requests directly.

## 2026-06-26: Enforce Store Request Layer

Hooks under `src/hooks/store/modules` must not call `request` directly.

Consequences:

- Every entity with API access needs a Pinia module in `src/stores/modules`.
- Hooks only wrap store refs and methods.
- Alarm tag, meta alarm rule, webhook token rule, external data table and pbehavior type hooks were aligned to this rule.

## 2026-06-24: Layout Rules

Use `c-layout` as the flex wrapper.

Consequences:

- Replace migrated `d-flex flex-column` patterns with `c-layout`.
- Vertical field gaps are `2`.
- Horizontal field gaps are `3`.

## 2026-06-24: Date And Time Picker Migration

Date/time picker fields were ported from old Vuetify 2 components to Vuetify 4.

Consequences:

- `time-picker-field` uses separate combobox search state.
- Invalid time input filters the list but does not update the model.
- `date-time-picker-field` works with `Date` values.
- Broadcast message form now uses `date-time-picker-field`.

## 2026-06-24: c-patterns-field Migration

`c-patterns-field` is being ported without rewriting its old component tree.

Consequences:

- Old validator APIs must be replaced with vee-validate 4 backed logic.
- Old `$tc` usage is replaced with `t`.
- Component logic stays close to the source implementation.
- Consumers are integrated gradually.

## 2026-07-22: vue-i18n Pluralization API

The Vue 3 project uses the current vue-i18n API and does not expose `tc`.

Consequences:

- Replace `tc(key, choice, params)` with `t(key, choice, params)`.
- Replace `$tc(key, choice, params)` with `$t(key, choice, params)`.
- Preserve the choice/count argument so plural messages keep working.
- Do not add compatibility wrappers that reintroduce `tc`.

## 2026-06-24: Drag And Drop Caution

Groups sidebar drag-and-drop had regressions during migration attempts.

Consequences:

- Do not touch DnD/groups sidebar code unless the user explicitly asks.
- If work resumes there, isolate it and verify manually.

## 2026-08-13: Dotted Table Slot Names

Vuetify data-table item/header slots are named `item.<key>` and `header.<key>`. Vue treats `#item.key` as a slot modifier, and hyphenated keys fail `vue/valid-v-slot`.

Consequences:

- Use `` #[`item.row-icons`] `` / `` #[`header.${key}`] `` for dotted slot names.
- Do not write `#item.row-icons`.
- Do not disable `vue/valid-v-slot` to silence that error.
- `#item.actions` and other unhyphenated Vuetify shorthands may stay as-is.

## 2026-06-24: Agent Threads

Four pinned agent threads were created for specialized follow-up work.

Agents:

- Forms & Validation
- Pages & Modals
- Stores & API
- QA & Regressions

Consequences:

- Agents do not automatically know file ownership.
- Give each agent a scope.
- Do not let multiple agents edit the same file at the same time.
