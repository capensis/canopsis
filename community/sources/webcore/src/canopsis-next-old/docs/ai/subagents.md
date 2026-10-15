# Migration Subagents

This file is the persistent registry for specialized subagents used during the Canopsis Vue 3 migration.

Old source is always read-only:

```text
/Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next
```

Current target project:

```text
/Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next-vue-3
```

## Registry

| Layer | Nickname | Agent ID |
| --- | --- | --- |
| Stores/API | Bohr | `019f8951-7408-7cc3-9a49-0170b997003a` |
| Forms | Poincare | `019f8951-751b-7b31-8921-7ae6facd4c67` |
| Tables/Actions | Mendel | `019f8951-7561-7d33-9835-1c38824f0467` |
| Vuetify/Layout | Linnaeus | `019f8951-75bd-7590-aae2-e8edda195fd8` |
| Page Migration | Descartes | `019f8951-7626-7d63-b621-540cd2648875` |
| QA/Regression | Avicenna | `019f8951-767e-7650-96c8-6bb30a494f85` |

## Coordination Rules

- Assign work by layer, not by convenience.
- Do not let two agents edit the same file set at the same time.
- Each agent must read relevant `docs/ai/*.md` before implementation work.
- Each agent must compare with the old source before changing migrated behavior.
- Agents must not edit `../canopsis-next`.
- Agents must list changed files in their final response.
- For code changes, run targeted ESLint on touched files; run `npm run build` for app-code changes.
- Shared high-conflict files such as constants, i18n, router records, barrels, and common components require explicit ownership.

## Stores/API: Bohr

Owns store and backend integration work.

Primary areas when assigned:

- `src/stores/**`
- `src/hooks/store/**`
- `src/services/**`
- `src/helpers/entities/**`
- API route constants, usually `src/config.js`

Checklist:

- Preserve old API payloads and response extraction exactly.
- Keep request flow: component -> hook -> Pinia store -> `src/services/request`.
- Hooks wrap stores and must not import `request` or API routes.
- Components consume hooks, not Pinia stores or request directly.
- Preserve `items`, `meta`, `pending`, `fetchingParams`, previous params, and `WithoutStore` methods where old behavior requires them.
- Preserve bulk action routes and payloads.

Do not touch without explicit assignment:

- Components, pages, modals, forms, layout, router, i18n.
- Drag-and-drop/groups sidebar code.

## Forms: Poincare

Owns form architecture and validation work.

Primary areas when assigned:

- `src/components/forms/**`
- `src/hooks/form/**`
- `src/hooks/validator/**`
- `src/mixins/form/**`
- `src/helpers/date/**`

Checklist:

- Form roots use `useForm()`.
- Subforms use `useFormContext()`, `useFormModels()`, or `useFormModel(path)`.
- Leaf field components receive model props and emit `update:*`.
- Avoid immediate validation during initial form hydration.
- Do not use `v-field` or component `@input`.
- Prefer `c-*` field wrappers and `c-form-block`.
- Use `t` / `$t`, not `tc` / `$tc`; preserve plural choice/count arguments.

Do not touch without explicit assignment:

- Store/API layers, pages, modals, router.
- Shared constants/i18n/barrels unless coordinated.

## Tables/Actions: Mendel

Owns table behavior, selection, mass actions, and table action panels.

Primary areas when assigned:

- `src/components/common/table/c-advanced-data-table.vue`
- `src/components/common/table/c-table-mass-actions-panel.vue`
- `src/components/common/actions-panel/**`
- Table usage in assigned `src/components/other/**`
- Table option/query glue in assigned `src/views/**`

Checklist:

- Preserve old table behavior and visual density.
- `select-all` enables selection; disabled rows must not be selectable.
- Keep entity-specific mass actions in feature components unless explicitly assigned otherwise.
- Mass action slots use `{ selected, selectedKeys, count, clearSelected }`.
- Preserve pagination, sorting, search, expand-row, and keep-selected behavior.
- Teleports use `TELEPORTS_NAMES.massActionsPanel` with modal scoping when needed.

Do not touch without explicit assignment:

- Stores/API, forms/validation, router, i18n, shared constants.
- Groups sidebar drag-and-drop code.

## Vuetify/Layout: Linnaeus

Owns visual parity and Vuetify adaptation.

Primary areas when assigned:

- Assigned page/component layout files.
- `src/components/layout/**` only when specifically assigned.
- Shared/common layout wrappers only when specifically assigned.
- Modal, tab, card, tooltip, overflow, spacing, and density fixes in assigned scopes.

Checklist:

- Preserve old visual behavior, density, spacing, and interactions.
- Use `c-layout`.
- Use vertical field gap `2` and horizontal field gap `3`.
- Prefer `c-*` wrappers over raw Vuetify fields.
- Use Vuetify 4 activator slots with `{ props }`.
- Keep modal bodies scrollable without breaking headers/actions.
- Watch overflow, text overlap, tooltips, and top-bar sizing.

Do not touch without explicit assignment:

- Stores/API, form architecture internals, router, i18n, shared constants.
- Drag-and-drop/groups sidebar code.

## Page Migration: Descartes

Owns end-to-end page migration orchestration.

Primary areas when assigned:

- `src/views/**`
- `src/components/other/**`
- `src/components/modals/**`
- `src/router/**`
- `docs/migration-inventory.md`
- Page-specific migration plans under `docs/ai/*-migration-plan.md`

Checklist:

- Inventory old route/page, child components, modals, stores, permissions, i18n, helpers.
- Preserve old behavior and component style.
- Use `c-page` / `c-page-header` where the old page did.
- Map missing stores/forms/tables/helpers before coding.
- Preserve permissions, refresh flows, polling/socket flows, query/options restoration, and tab lazy behavior.
- Coordinate with layer agents instead of patching unrelated layers casually.

Do not touch without explicit assignment:

- Stores/API internals, form internals, shared constants/i18n.
- Drag-and-drop/groups sidebar code.

## QA/Regression: Avicenna

Owns read-only checks and regression review unless explicitly assigned edits.

Primary inspection areas:

- `src/views/**`
- `src/components/other/**`
- `src/components/modals/**`
- `src/components/forms/**`
- `src/stores/modules/**`
- `src/hooks/store/modules/**`
- `src/helpers/entities/**`
- `src/i18n/messages/**`
- `src/router/**`

Checklist:

- Run targeted ESLint first, then `npm run build` for app-code changes.
- Compare migrated behavior against old source.
- Watch console for Vue/Vuetify/i18n errors.
- Verify store payloads and response extraction.
- Verify table pagination, sorting, expansion, selection, mass actions, and empty states.
- Verify form validation timing, submit paths, dirty close behavior, and API mappers.
- Verify modal sizing/overflow, tab switching, permissions, and visual parity.

Do not touch without explicit assignment:

- Any target files by default.
- Old source.
- Shared high-conflict files unless ownership is explicitly granted.
