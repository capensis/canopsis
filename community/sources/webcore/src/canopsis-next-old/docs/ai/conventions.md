# Conventions

## Migration Rules

- The source project is `canopsis-next`; always compare with it before implementing migrated behavior.
- Do not invent simplified behavior when the old implementation exists.
- Do not simplify migrated components or helpers as an interim shortcut. Port the old behavior, state, events, slots, validation timing, menu/tooltip behavior, selection handling, and edge cases first; only adapt the parts that are incompatible with Vue 3, Vuetify 4, Vite, Pinia, vee-validate 4, or vue-i18n 11.
- Port all accompanying code used by the migrated source: mixins, hooks, helpers, constants, typedefs, child components, styles and i18n keys. If something is intentionally not ported because it is obsolete or incompatible with the target stack, report that decision and why.
- If exact behavior cannot be preserved in the current target infrastructure, stop and document the gap or ask before replacing it with a reduced implementation.
- Preserve the old component style: Options API stays Options API, Composition API stays Composition API.
- When migrating Vuetify field slots, do not copy old slot names mechanically: old inside-field `append` icons, including help icons, become `append-inner` in Vuetify 4; old outside-field `append-outer` becomes `append`.
- For Vuetify table/item slots whose names contain a dot, use dynamic slot syntax `` #[`item.column-key`] ``; do not write `#item.column-key`, especially when the key has a hyphen (see `docs/ai/coding-standards.md`, dotted `v-slot` names). Do not disable `vue/valid-v-slot` for this.
- Use one explicit spacing mechanism per layout: prefer Vuetify `v-row`/`v-col` spacing props/classes. Do not add old `gap-*` classes or introduce `c-layout`.
- Change only the part that does not fit the new stack.
- Do not refactor unrelated files.
- Do not revert user or other-agent changes.
- Do not touch drag-and-drop/groups sidebar code unless explicitly asked.
- Do not start the dev server unless explicitly asked.
- When finishing any migration stage, immediately update the relevant migration documentation in the same change: `docs/migration-inventory.md` for page/widget/shared-component state, the `Ready For Testing` section at the top of that file when the area becomes runtime-testable, and any page-specific `docs/ai/*-migration-plan.md` checklist when it exists.

## File And Naming Conventions

- Files and directories use `lowercase-with-dashes`.
- Vue components are used as kebab-case in templates.
- Constants use `UPPER_SNAKE_CASE`.
- Variables and functions use `camelCase`.
- Boolean variables should use `is`, `has`, `can` or `should` prefixes.
- Page components should use the `-page` suffix.
- Do not add a Vue `name` option unless recursive rendering requires it.

## Imports

Use the existing import order:

1. External libraries.
2. Config/constants.
3. Services/API/request helpers.
4. Helpers.
5. Hooks/composables.
6. Components.

Combine constants in one import from `@/constants` when practical.

For `.js` helpers, prefer named exports. `.vue` files use `export default`.

## JavaScript Style

- Semicolons are required.
- Follow the configured `arrow-parens` rule: `as-needed`, with parens required for block bodies.
- Prefer `const` and arrow functions where the local style does.
- Use early returns for guard logic.
- Keep comments rare and useful.
- Use JSDoc when old code had JSDoc or when adding shared helpers.
- When migrating helper files from `canopsis-next`, migrate the original JSDoc together with the code, including function comments and all `@typedef` blocks, and keep comments updated if signatures, types or behavior change.
- Do not write `@returns {void}` or `Promise<void>` in JSDoc.

## Vue Style

- In Vue SFC files, keep top-level blocks in this order: `<template>`, then `<script>`, then `<style>`. The `<template>` block must be first whenever it exists.
- Preserve Options API or Composition API from the source component.
- Do not force components into Composition API during migration.
- Do not use `reactive` in the form layer; use `ref`/`computed`.
- Do not destructure `useModals()` or `usePopups()` unless existing code in this project does it intentionally.
- Use `useI18n()` in Composition API and `this.$t` in Options API.
- Use `t`, not `tc`; vue-i18n pluralization is preserved by passing the choice/count as the second argument: `t('common.entity', count, { count })`.
- Replace old `$tc(...)` calls with `$t(...)` using the same choice/count argument.

## UI Conventions

- For conditional `v-alert` / `c-alert`, use `v-if` on the component; do not use `:value`, `:model-value`, or `v-model` (see `docs/ai/coding-standards.md`, `v-alert` visibility).
- Use Vuetify `v-row`/`v-col` for layout in forms and feature components. Do not introduce `c-layout`; when an existing `c-layout` is encountered in code being edited, rewrite it to the equivalent `v-row`/`v-col` structure while preserving alignment, spacing and responsive behavior.
- Preserve vertical and horizontal field spacing through the corresponding `v-row`/`v-col` spacing props/classes.
- Use icons in icon buttons when an icon exists.
- Button and tab text should be uppercase.
- Always use the global `c-tooltip` wrapper in application templates. Do not use raw `v-tooltip` or restore `c-simple-tooltip`; see `docs/ai/coding-standards.md` (`Use c-tooltip`).
- Tooltips and menus should use fade transitions by default.
- Global popups must be rendered above every app layer: keep `the-popups` teleported to `body` and do not lower its top-level z-index below modal/sidebar/overlay layers.
- Tables should use Vuetify 2-like row borders and hover.
- Avoid card nesting and decorative page sections.
- Preserve old visual density unless asked to redesign.

## Forms

- Use `useForm()` from `src/hooks/form` for form roots.
- Components that call `useForm()`, `useFormContext()`, `useFormModels()` or `useFormModel()` must be Composition API components; do not add `setup()` with form hooks to an otherwise Options API component.
- Subforms (`*-form.vue`) use `useFormContext()`, `useFormModels()` or `useFormModel()`. They read and write through the parent form context; do not pass `form` by props or emit `update:form`.
- Field components under `src/components/**/fields/` (leaf/composite inputs, not `*-form` subforms) receive model props and emit `update:*`; do not use `useForm()`, `useFormContext()`, `useFormModels()` or `useFormModel()` inside them.
- Widget settings sidebars use `useWidgetSettings()` instead of copying submit, dirty-close and user-preference logic locally; the hook owns `useForm()` and the shared `widget-settings` wrapper reuses that provided context instead of being validated through a component ref.
- Use wrapper components such as `c-text-field`, `c-select`, `c-combobox`, `c-radio-group`, `text-editor-field`.
- Do not use raw `v-text-field`, `v-select`, `v-combobox`, etc. in app templates when a `c-*` wrapper exists; see `docs/ai/coding-standards.md` (`Use c-* Field Wrappers`).
- Leaf wrappers may use `useField()` for validation state.
- Values should be passed with `v-model`, `v-model:<prop>`, or explicit `:prop` + `@update:<prop>`.
- Treat `name` as the validation field name only. Do not build field paths from `props.name` in `setup()` or add computed/helper aliases for them; concatenate child names directly in the template, for example `:name="`${name}.value`"`. For array items, append the item's stable `key`, not its index.
- Do not use `v-field` or component `@input`; see `docs/ai/coding-standards.md` (`Replace v-field And @input`).
- Avoid direct `setFieldValue`; vee-validate should not be the form value store.

## Stores And Hooks

- Components use hooks from `src/hooks/store/modules`.
- Hooks wrap Pinia stores and **must not** call `request` or API routes directly.
- Pinia stores under `src/stores/modules` own all server request/API logic.
- Expose domain-specific method names from hooks: `createUser`, `removeRole`, `fetchBroadcastMessagesList`.
- Preserve old API payloads exactly when porting stores.
- Port the store module before the hook; port the hook before any component that needs the entity.

Layer:

```text
component -> hooks/store/modules -> stores/modules -> services/request
```

## i18n

- Messages live in `src/i18n/messages/en` and `src/i18n/messages/fr`.
- Keep key structures aligned.
- Do not rename old message keys unless the feature migration requires it.
- Controlled HTML in translations is allowed only where the component intentionally renders sanitized/controlled HTML.
- The new project does not expose `tc`. For plural messages, use `t(key, choice, params)` / `$t(key, choice, params)` and keep the old plural message format.

## Agent Coordination

When delegating work to an agent, include:

```text
Scope:
- work only in:
  - src/path/**
- read-only:
  - /Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next/src/path/**
- do not touch:
  - src/other-path/**

Do not edit files outside scope without asking.
List every changed file in the final response.
```

Do not let two agents edit the same file at the same time.
