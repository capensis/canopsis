# Canopsis Next Migration Inventory

Source project: `/Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next`

This file tracks migration progress for top-level route pages and widget types from `canopsis-next`.

Status legend:

- `[x]` migrated
- `[~]` partially migrated
- `[ ]` not migrated

## Immediate Tasks

- `[ ]` Remove `normalizr` from the old store migration path before porting larger entity modules. Prefer plain normalized-by-id state, Pinia stores and explicit mapper helpers instead of carrying over the old Vuex/normalizr shape.
  - Do not port the old `SCHEMA_EMBEDDED_KEY` / `entities/removeFromStore` mechanism. Embedded metadata was only used internally by the normalizr plugin for mark-and-sweep and parent cleanup.
  - Single pbehavior delete reaches `entities/removeFromStore` from `pbehavior-actions`, but the affected tables refresh themselves. Revisit only when migrating `AlarmsList`, `Context` and `ServiceWeather` pbehavior tabs; if needed, add explicit local refresh/parent update there instead of restoring global embedded cleanup.
- `[x]` Port remaining old entity store modules to Pinia facades.
  - Simple list/detail/bulk modules are wrapped with `createCRUDStore` where the old Vuex module was only CRUD around an API route.
  - Request-only or widget-scoped modules stay as explicit Pinia stores: `metrics` keeps per-widget buckets for aggregated/group/vector/availability data and meta, `template-testing` also covers old `copy/vars`, `snmp` is split into `snmp-rule` and `snmp-mib`, and `info` remains covered by `app-info` plus `infos`.
  - Empty old namespace folders `auth-token-job`, `instruction-job` and `webhook-job` do not require Pinia modules.
  - Consumer runtime parity still belongs to the page/widget migration where those hooks become used.

## Cross-Cutting Migration Rules

- `[x]` Target dependency baseline for UI migration is Vuetify `^4.1.2` from `package.json`. Do not describe the target as Vuetify 3; adapt old Vuetify 2 markup against the actual Vuetify 4 API and the compatibility wrappers already present in this project.
- `[x]` Replace `v-field` and `@input` on Vue components with explicit Vue 3 bindings:
  - remove `v-field="..."`;
  - use `:prop` + `@update:prop`, or `v-model:prop` when a form model field already exists;
  - replace `model: { prop: 'X', event: 'input' }` with `emits: ['update:X']`;
  - replace `emit('input', ...)` with `emit('update:X', ...)`.
  - use `useModelField()` / `useArrayModelField()` for object and array subforms.
  - details: `docs/ai/coding-standards.md` (`Replace v-field And @input`).
- `[x]` Use `c-*` Vuetify wrappers instead of raw `v-*` form controls in app templates (`c-select` instead of `v-select`, etc.); details: `docs/ai/coding-standards.md` (`Use c-* Field Wrappers`).
- `[x]` Removed `c-select-field` compat wrapper; bind `c-select` directly with `:model-value` / `@update:model-value` and `item-title` / `item-value`.
- `[x]` Merge project-component listener objects into a single `v-bind` object; rename `update:<prop>` listeners to `onUpdate:<prop>` instead of using paired `v-bind` + `v-on`; details: `docs/ai/coding-standards.md` (`v-bind` And Listener Objects For Project Components`).

## Shared Components

- `[x]` `c-fab-btn`: migrated to `src/components/common/buttons/c-fab-btn.vue`.
- `[x]` `c-refresh-btn`: migrated to `src/components/common/buttons/c-refresh-btn.vue` because `c-fab-btn` depends on it.
- `[x]` Modal base dependencies for maintenance mode: `c-modal-wrapper`, `c-form-block`, `c-form-block-row`, `text-editor-field`, `c-color-picker-field`.
- `[x]` Common modals: `confirmation`, `create-maintenance`.
- `[x]` `c-duration-field`: full old duration value/unit field is migrated in the compat field barrel, with `long`, `medium`, `sameWidth`, custom `units`, clearable unit handling and pluralized unit labels adapted to Vue I18n `t`.
- `[~]` Text editor stack: `text-editor`, `text-editor-field`, `text-editor-with-template`, `c-widget-template-field`, variables menu/list, upload hook, template-vars hook, image/images viewer modals and sidebar text-editor fields are migrated. Remaining: connect full vee-validate form validation inside editor modals when related pages/widgets start using complex rules.
- `[x]` `c-draggable-list-field`: migrated with Vue 3 `vuedraggable` API; emits `update:value` / `update:model-value` matching its `value` / `modelValue` props.
- `[~]` `c-advanced-data-table` migration plan:
  - `[x]` Inventory old dependencies and current target gaps.
  - `[x]` Leaf table placeholder: `c-empty-data-table-columns`.
  - `[~]` Table-focused completion plan, excluding CRUD page migration:
    - `[~]` Freeze supported Vue 3 API/compat surface: props, events and slots for `headers.text/value`, `headers.title/key`, `options`, `sortBy/sortDesc`, `item.<field>`, `header.*`, `toolbar`, `expand` and `mass-actions`.
    - `[~]` Verify and harden core `v-data-table-server` behavior: server-side pagination, sorting, loading/no-data/footer states and whether local item slicing is still needed. Local item slicing is removed; Vuetify 4 sort updates are converted back to canonical `sortBy/sortDesc` before emitting `update:options`.
    - `[ ]` Verify and harden selection behavior: disabled rows, select all, selected objects/keys, clear on options/data changes and mass-actions slot contract.
    - `[ ]` Verify and harden expand rows: expandable predicate, transition cleanup, colspan with select/expand columns and expanded state reset.
    - `[~]` Verify slot compatibility: `item`, `item.<field>`, `header`, `header.<field>`, `body.prepend`, `progress`, `toolbar`, `expand`, `mass-actions`. `item.<field>` now accepts both legacy `header.value` and Vuetify 4 `header.key`.
    - `[ ]` Verify table-owned search integration: simple search and advanced-search submit/reset/loading event flow. Deep advanced-search editors remain tracked under the search stack.
    - `[ ]` Verify pagination components: empty totals, last page after deletions, page reset when `itemsPerPage` grows and no-options mode.
    - `[ ]` Keep `c-table-mass-actions-panel` outside the core table scope. It is a consumer of the `mass-actions` slot with domain actions, not an internal dependency of `c-advanced-data-table`.
    - `[ ]` Runtime checklist against existing table consumers: users, roles, broadcast messages, template-testing, dynamic-info, event-filter, idle-rules, scenarios.
  - `[~]` Direct helper components used by `c-advanced-data-table`:
    - `[x]` Foundation helpers: `c-layout` and `c-icon` are migrated enough for table usage.
    - `[~]` Simple search: `c-search` and `c-search-field` are migrated; table now passes `options.search` into `c-search`, local search value syncs with external reset, submit/clear emit `update:modelValue`, and combobox history remove/pin controls are restored. Remaining table-specific pass is runtime focus/keyboard verification inside the toolbar.
    - `[~]` Advanced search shell: `c-advanced-search` and `c-advanced-search-field` are migrated; remaining table-specific pass is submit/reset/loading event flow.
    - `[~]` Advanced search internals: `advanced-search-rules`, `advanced-search-rule`, chips, duration/range chips and history button are migrated but need parity checks for keyboard flow, lazy values, date/duration editors and disabled/pending states.
    - `[~]` Pagination helpers: `c-table-pagination` and `c-pagination` are migrated; remaining table-specific pass is edge cases for zero totals, page recalculation after `itemsPerPage` changes and out-of-range pages after deletion.
    - `[x]` `c-empty-data-table-columns` exists as a leaf placeholder; it is not part of the runtime path for ordinary `c-advanced-data-table` lists.
  - `[~]` Pagination base: `c-pagination`, `c-table-pagination`, `c-items-per-page-field`, `helpers/pagination`.
  - `[~]` Search base: `c-search`, `c-search-field`, advanced-search public wrappers, rule/chip/history components and query helpers are migrated; remaining work is deep old parity for date/duration-specific editors and store-backed attribute value fetching.
  - `[~]` Vuetify table compatibility: initial mapping for old Vuetify 2 `headers.text/value`, `server-items-length`, `options`, `sortDesc`, `show-select`, `show-expand`, expanded slot and selectable API onto Vuetify 4 `v-data-table-server` is in place; verify against first real consumer.
  - `[~]` Main component: `components/common/table/c-advanced-data-table.vue` is ported without changing Options API; remaining work is runtime verification and consumer-specific slot fixes.
  - `[~]` First real consumer: `users-list` is switched to `c-advanced-data-table`; remaining users-table work is `c-table-mass-actions-panel` and bulk user actions UI.
- `[~]` `c-patterns-field` migration strategy:
  - `[x]` Scope rule: old component tree is ported without rewriting the main logic; Options API components stay Options API and Composition API hooks stay hooks.
  - `[x]` Phase 1, shared model/helpers: `helpers/entities/pattern/form.js`, `helpers/entities/filter/form.js`, `helpers/entities/pattern/fields/form.js` and required constants are available.
  - `[x]` Phase 2, Pinia stores/hooks: migrated `pattern`, `pbehavior-patterns`, `patterns-fields`, `pattern-entities-optimize`, entity-category and pbehavior dictionary facades. Hooks wrap stores; components do not call request directly.
  - `[x]` Phase 3, base rule editor leaves: migrated `pattern-attribute-field`, `pattern-operator-field`, `pattern-rule-field-date-value`, `pattern-rule-field`, `pattern-group-field`, `pattern-groups-field`, `pattern-editor-field`, `pattern-advanced-editor-field`, `c-pattern-field`, `c-filter-field`.
  - `[x]` Phase 4, domain pattern fields: migrated `c-alarm-patterns-field`, `c-entity-patterns-field`, `c-pbehavior-patterns-field`, `c-event-filter-patterns-field`, `c-service-weather-patterns-field`; added lightweight compatibility fields needed by their dynamic `valueField` configs.
  - `[~]` Phase 5, counters and preview modals: migrated `pattern-counters`, `pattern-count-message`, `pattern-count-alarms-modal`, `pattern-count-entities-modal`. Remaining: real `alarms-list`, `entities-list`, `entities-comparison` modal bodies still depend on porting alarm/entity list stacks.
  - `[x]` Phase 6, optimization/suggestions: migrated `pattern-optimization`, `pattern-progress`, `pattern-try-optimization`, `pattern-field-suggestions-wrapper`, `pattern-suggestions`, `pattern-suggestions-list`, `pattern-suggestions-header`, `pattern-suggestion-tab`, `pattern-suggestion-content`, `pattern-fix-with-ai-btn`; API calls go through pattern-entities-optimize store/hook.
  - `[x]` Phase 7, `c-patterns-field` wrapper: migrated `c-patterns-field.vue`; old `v-field` usage is replaced with explicit `:value` / `@update:value` bindings.
  - `[ ]` Phase 8, consumers: integrate first in read-only scenario action pattern tab, then in forms that already depend on patterns: scenario modal, event-filter, dynamic-info, meta-alarm-rule, declare-ticket-rule, pbehavior/remediation/storage settings as they are migrated.
  - `[~]` Verification checklist: build and lint pass. Remaining runtime checks: read-only rendering, add/remove groups/rules, date and duration values, infos fields, required/someRequired validation, counters, see alarms/entities modals, optimization cancel/retry/apply/reject, and existing query payload parity.
- `[~]` Remaining common button components from old `components/common/buttons`: `c-expand-btn` is migrated for scenarios expand cards, `c-btn-with-error` is migrated for patterns. Pending: `c-action-fab-btn`, `c-fab-expand-btn`, `c-download-btn`, `c-copy-btn`, `c-speed-dial`, `c-db-export-btn`, `c-grid-btns`, density toggle dependencies.

## Shared Hooks And Mixins

- `[x]` Query hooks: `hooks/query/local-query.js`, `hooks/query/options.js`, `hooks/query/shared.js`.
- `[x]` Query pending hook: `hooks/query/pending.js`; extended with `initialPending`.
- `[x]` Query mixins: `mixins/query/query.js`, `mixins/query/interval.js`.
- `[x]` Query helper dependencies: `helpers/query.js`, `helpers/entities/shared/query.js`, `helpers/collection.js`, `helpers/date/date-intervals.js`.
- `[~]` Groups side bar stack: `groups-side-bar`, nested group/view/playlists/draggable components, navigation editing mode, view/group/playlist/active-view Pinia stores and related mixins are migrated. Remaining: wire create/edit/duplicate group/view modals after those modals are migrated.
- `[ ]` Remaining query-related entity helpers from old `helpers/entities/*/query.js`: migrate when corresponding pages/widgets need them.

## Pages

Pages are taken from `src/router.js` in the source project.

| Status | Route name                                     | Path                                  | Source component                                             | Target component                                                  | Notes                                                                                                                                                                                                                                                                                                                  |
| ------ | ---------------------------------------------- | ------------------------------------- | ------------------------------------------------------------ | ----------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `[~]`  | `login`                                        | `/login`                              | `src/views/login.vue`                                        | `src/views/login-page.vue`                                        | Main login page and form are migrated; messages will be moved later.                                                                                                                                                                                                                                                   |
| `[~]`  | `home`                                         | `/`                                   | `src/views/home.vue`                                         | `src/views/home-page.vue`                                         | Stub/basic page exists.                                                                                                                                                                                                                                                                                                |
| `[ ]`  | `view`                                         | `/view/:id`                           | `src/views/view.vue`                                         | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `view-kiosk`                                   | `/kiosk-views/:id/:tabId`             | `src/views/view-kiosk.vue`                                   | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `alarms`                                       | `/alarms/:id`                         | `src/views/alarm.vue`                                        | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[~]`  | `admin-rights`                                 | `/admin/rights`                       | `src/views/admin/permissions.vue`                            | `src/views/admin/rights-page.vue`                                 | Main rights table migrated with lightweight checkbox; final visual parity pass pending.                                                                                                                                                                                                                                |
| `[~]`  | `admin-users`                                  | `/admin/users`                        | `src/views/admin/users.vue`                                  | `src/views/admin/users-page.vue`                                  | Main users tab/table and create/edit user form modal are migrated; shared tokens pending.                                                                                                                                                                                                                              |
| `[~]`  | `admin-roles`                                  | `/admin/roles`                        | `src/views/admin/roles.vue`                                  | `src/views/admin/roles-page.vue`                                  | Main list and create/edit/duplicate form are migrated. Remaining: role mass actions and deeper permissions editor parity.                                                                                                                                                                                              |
| `[~]`  | `admin-broadcast-messages`                     | `/admin/broadcast-messages`           | `src/views/admin/broadcast-messages.vue`                     | `src/views/admin/broadcast-messages-page.vue`                     | Main list, active messages and create/edit form are migrated; runtime parity pass pending.                                                                                                                                                                                                                             |
| `[x]`  | `admin-playlists`                              | `/admin/playlists`                    | `src/views/admin/playlists.vue`                              | `src/views/admin/playlists-page.vue`                              | Main list, refresh, copy/play actions, delete/bulk actions, create/edit/duplicate playlist modal, manage tabs modal and expand panel are migrated.                                                                                                                                                                     |
| `[~]`  | `admin-planning-administration`                | `/admin/planning`                     | `src/views/admin/planning.vue`                               | `src/views/admin/planning-page.vue`                               | Tabs, lists, expand panels, create buttons and pbehavior type/reason/exception create/import modal forms are migrated. Remaining: edit/duplicate action wiring and runtime parity pass.                                                                                                                                |
| `[ ]`  | `admin-healthcheck`                            | `/admin/healthcheck`                  | `src/views/admin/healthcheck.vue`                            | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `admin-kpi`                                    | `/admin/kpi`                          | `src/views/admin/kpi.vue`                                    | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `admin-events-records`                         | `/admin/events-records`               | `src/views/admin/events-records.vue`                         | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[~]`  | `admin-anomaly-monitored-connectors`           | `/admin/anomaly-monitored-connectors` | `src/views/admin/anomaly-monitored-connectors.vue`           | `src/views/admin/anomaly-monitored-connectors-page.vue`           | Main page, list, search/pagination, single create/edit/delete modal and mass enable/disable/remove actions are migrated. Remaining: history modal parity for the healthcheck connector tile and runtime visual parity pass.                                                                                            |
| `[~]`  | `admin-template-testing`                       | `/admin/template-testing`             | `src/views/admin/template-testing.vue`                       | `src/views/admin/template-testing-page.vue`                       | Main data/tests tabs, lists, create/edit data modal and edit/remove test modals are migrated; test-variables sidebar stack remains pending for other forms.                                                                                                                                                            |
| `[~]`  | `admin-jobs-management`                        | `/admin/jobs`                         | `src/views/admin/jobs.vue`                                   | `src/views/admin/jobs-page.vue`                                   | Main jobs tabs/table, filters, expand panel, play/pause/repeat actions and routing are migrated. Remaining: edit modal/form after declare-ticket check-status field layer is ported.                                                                                                                                   |
| `[~]`  | `admin-custom-objects-external-auth-tokens`    | `/admin/external-auth-tokens`         | `src/views/admin/custom-objects/external-auth-tokens.vue`    | `src/views/admin/custom-objects/external-auth-tokens-page.vue`    | Main list, expand panel, create/edit/delete modals, request form fields, token variables, route and webhook-token-rule CRUD are migrated. Remaining: old template testing tab parity after shared form-general-patterns/test-variables stack is ported, plus final runtime parity pass.                                |
| `[ ]`  | `admin-custom-objects-external-data-tables`    | `/admin/external-data-tables`         | `src/views/admin/custom-objects/external-data-tables.vue`    | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[~]`  | `admin-custom-objects-entity-infos-properties` | `/admin/entity-infos-properties`      | `src/views/admin/custom-objects/entity-infos-properties.vue` | `src/views/admin/custom-objects/entity-infos-properties-page.vue` | Main page, list, type filter, create/edit/duplicate form modal, delete and bulk delete actions are migrated. Remaining: runtime visual parity pass.                                                                                                                                                                    |
| `[~]`  | `admin-custom-objects-icons`                   | `/admin/icons`                        | `src/views/admin/custom-objects/icons.vue`                   | `src/views/admin/custom-objects/icons-page.vue`                   | Main page, icons list, custom icon preview/registration, create/edit upload modal and delete action are migrated. Remaining: runtime visual parity pass.                                                                                                                                                               |
| `[~]`  | `admin-custom-objects-maps`                    | `/admin/maps`                         | `src/views/admin/custom-objects/maps.vue`                    | `src/views/admin/custom-objects/maps-page.vue`                    | Main page, CRUD table, search/pagination, single delete, mass delete wiring, expand panel shell, geo/Mermaid/tree/flowchart expand previews, create/edit/duplicate map modals and map form/editor stack are migrated. Remaining: runtime visual parity pass and new-validator parity for nested editor required rules. |
| `[ ]`  | `admin-custom-objects-tags`                    | `/admin/tags`                         | `src/views/admin/custom-objects/tags.vue`                    | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `admin-custom-objects-llms`                    | `/admin/llms`                         | `src/views/admin/custom-objects/llms.vue`                    | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `admin-settings-views-import-export`           | `/admin/views-import-export`          | `src/views/admin/settings/views-import-export.vue`           | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `admin-settings-notifications`                 | `/admin/notifications`                | `src/views/admin/settings/notifications-settings.vue`        | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `admin-settings-comment-templates`             | `/admin/comment-templates`            | `src/views/admin/settings/comment-templates.vue`             | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `admin-settings-widget-templates`              | `/admin/widget-templates`             | `src/views/admin/settings/widget-templates.vue`              | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `admin-settings-storage-settings`              | `/admin/storage-settings`             | `src/views/admin/settings/storage-settings.vue`              | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `admin-settings-state-settings`                | `/admin/state-settings`               | `src/views/admin/settings/state-settings.vue`                | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[~]`  | `exploitation-pbehaviors`                      | `/exploitation/pbehaviors`            | `src/views/exploitation/pbehaviors.vue`                      | `src/views/exploitation/pbehaviors-page.vue`                      | Main page, list table, expand tabs, advanced search, delete/export, bulk actions, alarm filtering and pbehavior planning calendar modal are migrated. Create/edit/duplicate list actions and recurrence calendar modal remain pending.                                                                                 |
| `[~]`  | `exploitation-event-filters`                   | `/exploitation/event-filters`         | `src/views/exploitation/event-filters.vue`                   | `src/views/exploitation/event-filters-page.vue`                   | Main list, expand content, failures tab, remove/export, bulk actions and store hooks are migrated; create/edit/duplicate modal form remains pending and actions are hidden until it is adapted.                                                                                                                        |
| `[ ]`  | `exploitation-snmp-rules`                      | `/exploitation/snmp-rules`            | `src/views/exploitation/snmp-rules.vue`                      | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[~]`  | `exploitation-dynamic-infos`                   | `/exploitation/dynamic-infos`         | `src/views/exploitation/dynamic-infos.vue`                   | `src/views/exploitation/dynamic-infos-page.vue`                   | Page shell, dynamic infos list with advanced search, templates list, stores and create/edit modal are migrated.                                                                                                                                                                                                        |
| `[~]`  | `exploitation-meta-alarm-rules`                | `/exploitation/meta-alarm-rule`       | `src/views/exploitation/meta-alarm-rules.vue`                | `src/views/exploitation/meta-alarm-rules-page.vue`                | Main list, read-only expand patterns, route, helper and bulk/delete actions are migrated. Create/edit/duplicate modal form remains pending for the forms pass.                                                                                                                                                         |
| `[~]`  | `exploitation-scenarios`                       | `/exploitation/scenarios`             | `src/views/exploitation/scenarios.vue`                       | `src/views/exploitation/scenarios-page.vue`                       | Main list/read-only expand content and scenario store are migrated; create/edit modal is not migrated yet.                                                                                                                                                                                                             |
| `[~]`  | `exploitation-idle-rules`                      | `/exploitation/idle-rules`            | `src/views/exploitation/idle-rules.vue`                      | `src/views/exploitation/idle-rules-page.vue`                      | Main list, expand, delete/bulk actions and create/edit/duplicate form are migrated; pbehavior action parameters depend on future pbehavior stack migration.                                                                                                                                                            |
| `[~]`  | `exploitation-flapping-rules`                  | `/exploitation/flapping-rules`        | `src/views/exploitation/flapping-rules.vue`                  | `src/views/exploitation/flapping-rules-page.vue`                  | Main list, read-only expand, delete/bulk actions and create/edit/duplicate form are migrated; AI sidebar parity pending.                                                                                                                                                                                               |
| `[~]`  | `exploitation-resolve-rules`                   | `/exploitation/resolve-rules`         | `src/views/exploitation/resolve-rules.vue`                   | `src/views/exploitation/resolve-rules-page.vue`                   | Main list, read-only expand, delete/bulk actions and create/edit/duplicate form are migrated; AI sidebar parity pending.                                                                                                                                                                                               |
| `[ ]`  | `exploitation-declare-ticket-rules`            | `/exploitation/declare-ticket-rules`  | `src/views/exploitation/declare-ticket-rules.vue`            | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[~]`  | `exploitation-link-rules`                      | `/exploitation/link-rules`            | `src/views/exploitation/link-rules.vue`                      | `src/views/exploitation/link-rules-page.vue`                      | Page, list, mass actions, DB export, create/edit/duplicate modal, link rule form and Monaco-based advanced editor are migrated. AI chat remains a deferred stub per `docs/ai/architecture.md`.                                                                                                                         |
| `[ ]`  | `exploitation-remediation`                     | `/exploitation/remediation`           | `src/views/exploitation/remediation.vue`                     | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `playlist`                                     | `/playlist/:id`                       | `src/views/playlist.vue`                                     | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `profile-patterns`                             | `/filters`                            | `src/views/profile/patterns.vue`                             | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `profile-themes`                               | `/themes`                             | `src/views/profile/themes.vue`                               | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[ ]`  | `notifications`                                | `/notifications/:tabId`               | `src/views/notification/notifications.vue`                   | -                                                                 |                                                                                                                                                                                                                                                                                                                        |
| `[x]`  | `error`                                        | `/error`                              | `src/views/error.vue`                                        | `src/views/error-page.vue`                                        | Migrated with polling and redirect logic.                                                                                                                                                                                                                                                                              |

## Page Blocks

Each page is tracked by functional blocks. Add concrete component/file names inside the block as we start migrating that page.

### `login`

- `[~]` Main content: `src/views/login.vue` -> `src/views/login-page.vue`; login card, auth variants and footer are started.
- `[ ]` Modal windows: none found directly on route page; verify nested login components later.
- `[ ]` Sidebars/drawers: none expected.
- `[~]` State/API/socket: app-info polling and auth flow are started.
- `[~]` Routing/access: `/login`, `hideHeader`, anonymous access.
- `[ ]` Messages/i18n: user asked to move messages later.

### `home`

- `[~]` Main content: `src/views/home.vue` -> `src/views/home-page.vue`; target is currently a stub.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[~]` Routing/access: `/` route exists.
- `[ ]` Messages/i18n: not migrated.

### `view`

- `[ ]` Main content: `src/views/view.vue`.
- `[ ]` Tabs: migrate view tabs loading, active tab selection, tab permissions and tab navigation.
- `[ ]` Grid system: migrate layout shell, breakpoint handling, drag/resize rules, edit mode and widget grid parameters.
- `[ ]` Grid system content: migrate widget rendering inside grid cells, widget wrapper, widget edit controls, widget loading/error states and widget component resolution.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/view/:id`, view permission by route id.
- `[ ]` Messages/i18n: not migrated.

### `view-kiosk`

- `[ ]` Main content: `src/views/view-kiosk.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/kiosk-views/:id/:tabId`, simple navigation, view permission by route id.
- `[ ]` Messages/i18n: not migrated.

### `alarms`

- `[ ]` Main content: `src/views/alarm.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: alarm fetch, widget fetch and alarm widget integration not migrated.
- `[ ]` Routing/access: `/alarms/:id`, technical view permission.
- `[ ]` Messages/i18n: not migrated.

### `admin-rights`

- `[~]` Main content: `src/views/admin/permissions.vue` -> `src/views/admin/rights-page.vue`; tabs, search, expand/collapse, roles matrix and sticky submit/cancel are migrated.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: permission and role Pinia stores/hooks are migrated for rights page; broader role CRUD page reuse remains pending.
- `[x]` Routing/access: `/admin/rights`, permission management access.
- `[~]` Messages/i18n: existing permission messages are wired; final copy parity verification pending.

### `admin-users`

- `[~]` Main content: `src/views/admin/users.vue`; users tab and users table are migrated.
- `[~]` Modal windows: confirmation modal and create/edit user form modal are migrated; shared token actions are not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: user CRUD store integration is migrated for list/remove/update hooks and create/edit modal submit action; shared token flow is pending.
- `[x]` Routing/access: `/admin/users`, user management access.
- `[ ]` Messages/i18n: not migrated.

### `admin-roles`

- `[~]` Main content: `src/views/admin/roles.vue` -> `src/views/admin/roles-page.vue`; roles table with search, pagination, edit, duplicate and remove actions is migrated. Remaining: role mass actions parity.
- `[~]` Modal windows: `create-role` modal and role form are migrated; remaining: deeper permissions editor parity if needed.
- `[ ]` Sidebars/drawers: not migrated.
- `[x]` State/API/socket: role Pinia store/hook is reused and extended with previous-params fetching.
- `[x]` Routing/access: `/admin/roles`, role management access.
- `[ ]` Messages/i18n: not migrated.

### `admin-broadcast-messages`

- `[~]` Main content: `src/views/admin/broadcast-messages.vue` -> `src/views/admin/broadcast-messages-page.vue`; table, search, pagination, status and actions are migrated.
- `[~]` Modal windows: `create-broadcast-message` modal and form are migrated with new `useForm`; remaining runtime parity check for date-time picker behavior and views tree selection.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: broadcast message Pinia store/hook, active broadcast messages, mark-as-read and maintenance edit integration are migrated; runtime socket verification pending.
- `[x]` Routing/access: `/admin/broadcast-messages`, broadcast message access.
- `[~]` Messages/i18n: existing broadcast message and modal messages are wired; final copy parity verification pending.

### `admin-playlists`

- `[x]` Main content: `src/views/admin/playlists.vue` -> `src/views/admin/playlists-page.vue`; route is wired to a real page and the playlists table is migrated with search, pagination, play link, copy link, create/edit/duplicate, delete, bulk actions and read-only expand panel.
- `[x]` Modal windows: `create-playlist` and `manage-playlist-tabs` are migrated with `useForm`; selected tabs are stored as tab objects in form state and converted to `_id` list for API submit.
- `[x]` Sidebars/drawers: no playlist sidebars/drawers are required for this page.
- `[x]` State/API/socket: playlist Pinia CRUD store/hook is reused; current user refresh runs after single delete and create/update.
- `[x]` Routing/access: `/admin/playlists`, playlist management access.
- `[x]` Messages/i18n: existing playlist/common modal messages are wired.

### `admin-planning-administration`

- `[~]` Main content: `src/views/admin/planning.vue` -> `src/views/admin/planning-page.vue`; route is wired to a real page with pbehavior type/reason/exception tabs, migrated list tables, read-only expand panels and create/import FAB actions. Type expand shows description/type tabs with icon/color, reason expand shows description, exception expand shows description/periods tabs with formatted dates and type metadata.
- `[~]` Modal windows: `create-pbehavior-type`, `create-pbehavior-reason`, `create-pbehavior-exception` and `import-pbehavior-exception` are migrated and registered with `useForm`. The forms keep the old `c-form-block` layout, modal date picker variant, primary date header, adjacent months, and no premature validation on open. Remaining: wire table row edit/duplicate actions where needed and do a runtime parity pass for exception date editing.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: pbehavior type/reason/exception Pinia CRUD stores/hooks are reused for list, delete and bulk hide/unhide/remove.
- `[x]` Routing/access: `/admin/planning`, grouped planning permission and per-tab read permissions.
- `[~]` Messages/i18n: existing pbehavior/common/modal messages are wired; final copy parity pending.

### `admin-healthcheck`

- `[ ]` Main content: `src/views/admin/healthcheck.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: top-bar healthcheck chips are migrated; admin page is not.
- `[ ]` Routing/access: `/admin/healthcheck`, healthcheck `can` access.
- `[ ]` Messages/i18n: not migrated.

### `admin-kpi`

- `[~]` Main content: `src/views/admin/kpi-page.vue` is wired; Graphs tab has SLI, alarm metrics and rating metrics fetch/render/export paths; KPI filters, rating settings and collection settings tabs are migrated; runtime parity checks are still pending.
- `[~]` Modal windows: `create-kpi-filter` is migrated; rating/collection related modals are not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: filter/rating/metrics stores already exist; page healthcheck socket and KPI graph metric fetch/export paths are wired.
- `[x]` Routing/access: `/admin/kpi`, KPI `can` access.
- `[~]` Messages/i18n: KPI messages already exist; runtime parity still pending.

### `admin-events-records`

- Working plan: `docs/ai/admin-events-records-migration-plan.md`.
- `[x]` Main content: `src/views/admin/events-records-page.vue`.
- `[x]` Modal windows: `events-record`, `apply-event-filter`.
- `[ ]` Sidebars/drawers: not migrated.
- `[x]` State/API/socket: events-record Pinia store/hook, current status polling, recording/resending/export flows.
- `[x]` Routing/access: `/admin/events-records`, events record `can` access.
- `[x]` Messages/i18n: base events-record/modal/page-header messages already present.

### `admin-anomaly-monitored-connectors`

- `[~]` Main content: `src/views/admin/anomaly-monitored-connectors.vue` -> `src/views/admin/anomaly-monitored-connectors-page.vue`; list, search, pagination, refresh, create/edit/delete entry points and mass enable/disable/remove actions are migrated.
- `[~]` Modal windows: `create-anomaly-monitored-connector` modal and form are migrated. Remaining: `anomaly-monitored-connector-history` modal parity for healthcheck connector tile.
- `[ ]` Sidebars/drawers: not migrated.
- `[x]` State/API/socket: Pinia store/hook, CRUD, enabled patch, states/history fetches and bulk actions are migrated.
- `[x]` Routing/access: `/admin/anomaly-monitored-connectors`, anomaly monitored connector access.
- `[x]` Messages/i18n: base page-header, modal and anomaly monitored connector messages already present.

### `admin-template-testing`

- `[~]` Main content: `src/views/admin/template-testing.vue` -> `src/views/admin/template-testing-page.vue`; data and tests tabs with search, pagination, type filters and row actions are migrated.
- `[~]` Modal windows: `create-template-testing-data` and `create-template-testing-test` are migrated with `useForm`; remaining runtime parity pass for JSON field parse/reset behavior.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: template testing Pinia store/hook reused for data/tests CRUD and list fetching.
- `[x]` Routing/access: `/admin/template-testing`, template testing `can` access.
- `[~]` Messages/i18n: existing template testing and modal messages are wired; final copy parity verification pending.

### `admin-jobs-management`

- `[~]` Main content: `src/views/admin/jobs.vue`; migrated to `src/views/admin/jobs-page.vue` with tabs, table, filters, row expansion and play/pause/repeat actions.
- `[ ]` Modal windows: edit modal/form is deferred until the shared declare-ticket check-status field and form mappers are ported.
- `[x]` Sidebars/drawers: none in the migrated jobs page.
- `[x]` State/API/socket: existing ticket status job Pinia store/hook and bulk play/pause endpoints are reused.
- `[x]` Routing/access: `/admin/jobs`, job management `can` access.
- `[~]` Messages/i18n: existing jobs/modal messages are reused; final runtime copy parity pass pending after modal/form migration.

### `admin-custom-objects-external-auth-tokens`

- `[~]` Main content: `src/views/admin/custom-objects/external-auth-tokens.vue`; migrated to `src/views/admin/custom-objects/external-auth-tokens-page.vue` with list, search, pagination, expand panel and linked-rule delete protection.
- `[~]` Modal windows: create/edit external auth token and confirmation phrase delete modals are migrated.
- `[x]` Form fields: request URL, timeout, repeat request, skip verify, auth, headers, payload, token, token variables switch and token expiration are migrated. Supporting shared fields now include `c-duration-field`, `c-retry-field`, `c-payload-text-field`, `c-payload-textarea-field`, request pair fields and payload variable menus.
- `[ ]` Remaining form parity: port the old `c-form-general-patterns-tabs` / test-variables stack and restore the template testing tab behavior in `external-auth-token-form.vue`.
- `[x]` Sidebars/drawers: none in the migrated page.
- `[x]` State/API/socket: existing webhook-token-rule Pinia store/hook and template-vars endpoint are reused.
- `[x]` Routing/access: `/admin/external-auth-tokens`, external auth token access.
- `[~]` Messages/i18n: existing external auth token and modal messages are reused; final runtime copy parity pass pending.
- `[ ]` Final runtime parity pass: verify headers add/remove/edit, payload variables menu, token allow variables mode, repeat request, clearable duration, linked-rules delete tooltip and failed-request expand panel.

### `admin-custom-objects-external-data-tables`

- `[~]` Main content: `src/views/admin/custom-objects/external-data-tables.vue`; migrated to `src/views/admin/custom-objects/external-data-tables-page.vue` with list, pagination, linked-rule delete protection, removed-from-config warning and records expand panel.
- `[~]` Modal windows: create/edit table, create/edit/duplicate record and import records modals are migrated; runtime parity check pending.
- `[ ]` Sidebars/drawers: not migrated.
- `[x]` State/API/socket: external data table Pinia store now uses CRUD helper for table CRUD and exposes import, export and record actions through migrated hooks.
- `[x]` Routing/access: `/admin/external-data-tables`, external data table access.
- `[~]` Messages/i18n: existing external data and modal messages are reused; final English/French parity pass pending.
- `[ ]` Final runtime parity pass: verify table create/edit/delete, records expand load, column type/tag edits, record create/edit/duplicate/delete/bulk delete, import upload/preview/complete and schema export.

### `admin-custom-objects-entity-infos-properties`

- `[~]` Main content: `src/views/admin/custom-objects/entity-infos-properties.vue` -> `src/views/admin/custom-objects/entity-infos-properties-page.vue`; list, type filter, search, pagination and row actions are migrated.
- `[x]` Modal windows: `create-entity-info-property` modal and form are migrated; `c-entity-info-property-key-field` and type field are migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[x]` State/API/socket: entity info property Pinia store/hook already exists and is wired to page CRUD and bulk remove.
- `[x]` Routing/access: `/admin/entity-infos-properties`, entity info property access.
- `[~]` Messages/i18n: existing old keys are reused; final copy parity verification pending.

### `admin-custom-objects-icons`

- `[~]` Main content: `src/views/admin/custom-objects/icons.vue` -> `src/views/admin/custom-objects/icons-page.vue`; icons list, pagination, refresh, custom SVG preview and actions are migrated.
- `[x]` Modal windows: `create-icon` modal, icon form and SVG file selector are migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: material icon font/list generator and custom icon registry are migrated; page uses icon Pinia store hook for multipart CRUD. Remaining: socket refresh parity if needed.
- `[x]` Routing/access: `/admin/icons`, icon access.
- `[~]` Messages/i18n: existing old keys are reused; final copy parity verification pending.

### `admin-custom-objects-maps`

- `[~]` Main content: `src/views/admin/custom-objects/maps.vue` -> `src/views/admin/custom-objects/maps-page.vue`; table/search/pagination/actions, expand panel shell, geo/Mermaid/tree-of-dependencies/flowchart expand previews, common flowchart stack and cluster graph editor stack are migrated.
- `[~]` Modal windows: `create-map`, `create-geo-map`, `create-flowchart-map`, `create-mermaid-map`, `create-tree-of-dependencies-map` are migrated and registered with the Vue 3 modal plugin. Their forms/editors now use explicit `v-model:*`/`update:*` events instead of the old `v-field` directive. Remaining: runtime visual parity pass and reintroducing nested editor required-rule validation through the new validator APIs.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: existing Pinia CRUD map store/hook is used, including bulk remove and `fetchMapStateWithoutStore`; Mermaid, geo, tree and flowchart preview helpers, popup support and form/editor-specific map helpers are migrated.
- `[x]` Routing/access: `/admin/maps`, map access.
- `[~]` Messages/i18n: existing old map/geomap/modal keys are reused; final parity verification pending.

### `admin-custom-objects-tags`

- `[ ]` Main content: `src/views/admin/custom-objects/tags.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/admin/tags`, tag access.
- `[ ]` Messages/i18n: not migrated.

### `admin-custom-objects-llms`

- `[ ]` Main content: `src/views/admin/custom-objects/llms.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/admin/llms`, LLM access.
- `[ ]` Messages/i18n: not migrated.

### `admin-settings-views-import-export`

- `[ ]` Main content: `src/views/admin/settings/views-import-export.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/admin/views-import-export`, view import/export `can` access.
- `[ ]` Messages/i18n: not migrated.

### `admin-settings-notifications`

- `[ ]` Main content: `src/views/admin/settings/notifications-settings.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/admin/notifications`, notification settings `can` access.
- `[ ]` Messages/i18n: not migrated.

### `admin-settings-comment-templates`

- `[ ]` Main content: `src/views/admin/settings/comment-templates.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/admin/comment-templates`, comment template access.
- `[ ]` Messages/i18n: not migrated.

### `admin-settings-widget-templates`

- `[ ]` Main content: `src/views/admin/settings/widget-templates.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/admin/widget-templates`, widget template access.
- `[ ]` Messages/i18n: not migrated.

### `admin-settings-storage-settings`

- `[ ]` Main content: `src/views/admin/settings/storage-settings.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/admin/storage-settings`, storage settings access.
- `[ ]` Messages/i18n: not migrated.

### `admin-settings-state-settings`

- `[ ]` Main content: `src/views/admin/settings/state-settings.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/admin/state-settings`, state setting access.
- `[ ]` Messages/i18n: not migrated.

### `exploitation-pbehaviors`

- `[~]` Main content: `src/views/exploitation/pbehaviors.vue` -> `src/views/exploitation/pbehaviors-page.vue`; list table, advanced search, pagination, expand tabs (patterns, entities, comments, recurrence), delete/export action, bulk enable/disable/remove, and alarm filtering button are migrated.
- `[~]` Modal windows: `pbehaviorPlanning` and recurrent changes confirmation are migrated and registered. `create-pbehavior`, recurrence rule editor modal and recurrence calendar modal are not migrated yet; list edit/duplicate actions stay disabled until the remaining forms/calendar pass.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: pbehavior store/hooks are reused for list fetch, delete, DB export, bulk actions and alarm filtering socket room. Do not port old `entities/removeFromStore` for single pbehavior delete; use API delete plus list refresh.
- `[x]` Routing/access: `/exploitation/pbehaviors`, pbehavior access.
- `[x]` Messages/i18n: reused existing pbehavior messages.

### `exploitation-event-filters`

- `[~]` Main content: `src/views/exploitation/event-filters.vue` -> `src/views/exploitation/event-filters-page.vue`; list, search, pagination, remove/export actions, DB export button, expand tabs, summary, patterns, external data and failures tab are migrated. Create/edit/duplicate actions are hidden until the modal form is adapted.
- `[ ]` Modal windows: `create-event-filter` source files are copied but not registered yet; adapt form to the new `useForm` layer, wrappers and validation before enabling create/edit/duplicate. `event-filter-general-form`, `event-filter-enrichment-form` and `event-filter-change-entity-form` use `useFormModels` (no `:form` / `update:form` bridging).
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: event-filter Pinia store/hook, failures API, bulk actions, DB export facade and supporting external-data/pbehavior lookup facades are migrated for the page.
- `[x]` Routing/access: `/exploitation/event-filters`, event filter access.
- `[~]` Messages/i18n: existing event-filter messages are wired; final copy parity verification pending.

### `exploitation-snmp-rules`

- `[ ]` Main content: `src/views/exploitation/snmp-rules.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/exploitation/snmp-rules`, SNMP rule access.
- `[ ]` Messages/i18n: not migrated.

### `exploitation-dynamic-infos`

- `[~]` Main content: `src/views/exploitation/dynamic-infos.vue` -> `src/views/exploitation/dynamic-infos-page.vue`; tabs, dynamic infos list with `c-advanced-data-table`/advanced-search, expand item and templates list are migrated.
- `[~]` Modal windows: `create-dynamic-info` and `create-dynamic-info-template` are migrated for create/edit with `vee-validate` form; `add-dynamic-info-infos-from-template` is migrated for selecting names from templates; remaining work is final visual/behavior parity pass.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: dynamic info Pinia CRUD store/hook and dynamic info templates hook are migrated; bulk/mass action UI is pending.
- `[x]` Routing/access: `/exploitation/dynamic-infos`, dynamic info access.
- `[~]` Messages/i18n: existing `dynamic-info` and dynamic info template modal/form messages are wired; remaining work is final copy parity verification.

### `exploitation-meta-alarm-rules`

- `[~]` Main content: `src/views/exploitation/meta-alarm-rules.vue` -> `src/views/exploitation/meta-alarm-rules-page.vue`; list, search, pagination, read-only pattern expand panel, DB export and delete/bulk actions are migrated. Create/edit/duplicate actions are hidden until the form/modal layer is ported.
- `[ ]` Modal windows: `src/components/modals/meta-alarm-rule/create-meta-alarm-rule.vue` and related form stack are not migrated yet.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: Pinia CRUD store/hook and meta alarm rule form helper are migrated; AI chat/socket support remains deferred.
- `[x]` Routing/access: `/exploitation/meta-alarm-rule`, meta alarm rule access.
- `[x]` Messages/i18n: existing meta alarm rule messages are wired.

### `exploitation-scenarios`

- `[~]` Main content: `src/views/exploitation/scenarios.vue` -> `src/views/exploitation/scenarios-page.vue`; table, search, pagination, row expand and action general tab are migrated. Remaining: full pattern tab rendering depends on porting `c-patterns-field`; row db export button depends on porting `c-db-export-btn` and db-export store.
- `[ ]` Modal windows: `create-scenario` intentionally not migrated yet; page create/edit/duplicate actions are hidden until that modal is ported. Confirmation modal is reused for remove.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: scenario Pinia store and hook are migrated for list/remove/bulk enable/bulk disable/bulk remove/test execution API wrappers; test execution UI is not migrated.
- `[x]` Routing/access: `/exploitation/scenarios`, scenario access.
- `[~]` Messages/i18n: existing `scenario` and `modals.createScenario` messages are wired; modal/form messages will be verified with modal migration.

### `exploitation-idle-rules`

- `[~]` Main content: `src/views/exploitation/idle-rules.vue` -> `src/views/exploitation/idle-rules-page.vue`; list, search, pagination, summary/patterns expand, delete/bulk actions and create/edit/duplicate entry points are migrated.
- `[~]` Modal windows: `create-idle-rule` form is migrated; pbehavior action parameters still require pbehavior stack migration.
- `[ ]` Sidebars/drawers: not migrated.
- `[x]` State/API/socket: idle rules Pinia store/hook are migrated for list/create/update/delete/bulk actions.
- `[x]` Routing/access: `/exploitation/idle-rules`, idle rules access.
- `[~]` Messages/i18n: existing `idle-rules` messages are wired; final copy parity verification pending.

### `exploitation-flapping-rules`

- `[~]` Main content: `src/views/exploitation/flapping-rules.vue` -> `src/views/exploitation/flapping-rules-page.vue`; shared alarm-status rules list, search, pagination, expand, delete/bulk actions and create/edit/duplicate entry points are migrated.
- `[~]` Modal windows: `create-alarm-status-rule` modal and form are migrated with general/patterns tabs; remaining parity: old AI chat sidebar.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: flapping rules Pinia store/hook are migrated for list/delete/bulk actions.
- `[x]` Routing/access: `/exploitation/flapping-rules`, flapping rules access.
- `[~]` Messages/i18n: existing `flapping-rule` and page header messages are wired; final copy parity verification pending.

### `exploitation-resolve-rules`

- `[~]` Main content: `src/views/exploitation/resolve-rules.vue` -> `src/views/exploitation/resolve-rules-page.vue`; shared alarm-status rules list, search, pagination, expand, delete/bulk actions and create/edit/duplicate entry points are migrated.
- `[~]` Modal windows: `create-alarm-status-rule` modal and form are migrated with general/patterns tabs; remaining parity: old AI chat sidebar.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: resolve rules Pinia store/hook are migrated for list/delete/bulk actions.
- `[x]` Routing/access: `/exploitation/resolve-rules`, resolve rules access.
- `[~]` Messages/i18n: existing `resolve-rule` and page header messages are wired; final copy parity verification pending.

### `exploitation-declare-ticket-rules`

- `[ ]` Main content: `src/views/exploitation/declare-ticket-rules.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/exploitation/declare-ticket-rules`, declare ticket rule access.
- `[ ]` Messages/i18n: not migrated.

### `exploitation-link-rules`

- `[x]` Main content: `src/views/exploitation/link-rules-page.vue`, `src/components/other/link-rule/link-rules-list.vue`.
- `[x]` Modal windows: `src/components/modals/link-rule/create-link-rule.vue`.
- `[~]` Sidebars/drawers: AI chat sidebar is intentionally stubbed until the deferred AI migration pass.
- `[x]` State/API/socket: Pinia store/hook already available in `src/stores/modules/link-rule.js` and `src/hooks/store/modules/link-rule.js`.
- `[x]` Routing/access: `/exploitation/link-rules`, link rule access.
- `[x]` Messages/i18n: link rule and modal messages already present.

### `exploitation-remediation`

- `[ ]` Main content: `src/views/exploitation/remediation.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/exploitation/remediation`, grouped remediation permission.
- `[ ]` Messages/i18n: not migrated.

### `playlist`

- `[ ]` Main content: `src/views/playlist.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/playlist/:id`, playlist read permission by route id.
- `[ ]` Messages/i18n: not migrated.

### `profile-patterns`

- `[ ]` Main content: `src/views/profile/patterns.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[ ]` State/API/socket: not migrated.
- `[ ]` Routing/access: `/filters`, authenticated access.
- `[ ]` Messages/i18n: not migrated.

### `profile-themes`

- `[ ]` Main content: `src/views/profile/themes.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: base Vuetify theme runtime was migrated; profile theme page is not.
- `[ ]` Routing/access: `/themes`, profile theme permission.
- `[ ]` Messages/i18n: not migrated.

### `notifications`

- `[ ]` Main content: `src/views/notification/notifications.vue`.
- `[ ]` Modal windows: not migrated.
- `[ ]` Sidebars/drawers: not migrated.
- `[~]` State/API/socket: top-bar notification menu shell exists; notifications page is not.
- `[ ]` Routing/access: `/notifications/:tabId`, authenticated access.
- `[ ]` Messages/i18n: not migrated.

### `error`

- `[x]` Main content: `src/views/error.vue` -> `src/views/error-page.vue`.
- `[x]` Modal windows: none expected.
- `[x]` Sidebars/drawers: none expected.
- `[x]` State/API/socket: app-info polling and redirect logic migrated.
- `[x]` Routing/access: `/error`, `hideHeader`, query props.
- `[ ]` Messages/i18n: not fully audited.

## Router

Router rules are taken from `src/router.js`, `src/constants/common.js`, `src/config.js` and `src/helpers/router.js` in the source project.

### Base Setup

- `[x]` Router mode/history: source uses `ROUTER_MODE = 'history'`; target uses `createWebHistory`.
- `[x]` Routes table: all source routes are registered; not-yet-migrated pages use `src/views/placeholder-page.vue`.
- `[x]` Route constants: `ROUTES_NAMES`, `ROUTES` and `OLD_ROUTES` are available in target constants.
- `[x]` Async route components: real pages and placeholders are lazy-loaded.
- `[x]` `replaceAsync`: source compatibility helper is available as `router.replaceAsync`.

### Auth Redirect Rules

- `[x]` `access_token` query rule: if route query contains `access_token`, call auth token application logic, then replace the route with the same target and query without `access_token`.
- `[x]` Anonymous to protected redirect: if user is not logged in and route requires login, redirect to `login` with `redirect: to.fullPath` and pass through `errorMessage`.
- `[x]` Logged user to anonymous page redirect: if user is logged in and route has `requiresLogin: false`, redirect to `home`.
- `[x]` Public/neutral route support: preserve source behavior where routes without `requiresLogin` can be visited both anonymously and authenticated.

### Access Rules

- `[x]` App-info access guard: `checkAppInfoAccessForRoute(to)` checks edition-based page availability through `USER_PERMISSIONS_TO_PAGES_RULES`; on failure shows `common.notFound` and redirects to `home`.
- `[x]` User access guard: `checkUserAccessForRoute(to)` checks `requiresPermission` with default action `read`; supports permission id functions and grouped permission arrays.
- `[~]` Permission waiting: current implementation fetches current user when needed; `info/popupTimeout` waiting is not required yet because popup defaults already exist.
- `[x]` Permission algorithm: target uses `checkUserAnyAccess`, grouped permissions and CRUD action behavior.
- `[x]` Guard failure redirect: source redirects failed `beforeResolve` checks to `home`.

### Post-navigation Rules

- `[~]` Entity sweep: after route change hook exists; actual `entities/sweep` call will be wired when entity store is migrated.
- `[x]` Initial error popup: on router ready, if initial route has `errorMessage`, show persistent error popup.
- `[ ]` Scroll/restoration behavior: not explicit in source router; verify if any page-level logic depends on old Vue Router behavior.

## Legacy Redirects

Redirects are taken from the bottom of `src/router.js`.

| Status | Old path                         | Target route                                   | Notes               |
| ------ | -------------------------------- | ---------------------------------------------- | ------------------- |
| `[x]`  | `/admin/remediation`             | `exploitation-remediation`                     |                     |
| `[x]`  | `/admin/external-data-tables`    | `admin-custom-objects-external-data-tables`    |                     |
| `[x]`  | `/admin/entity-infos-properties` | `admin-custom-objects-entity-infos-properties` |                     |
| `[x]`  | `*`                              | `home`                                         | Catch-all redirect. |

## Widgets

Widgets are taken from `WIDGET_TYPES` and `COMPONENTS_BY_WIDGET_TYPES` in `src/constants/widget.js`.

| Status | Widget type         | Source component name        | Source file                                                                 | Icon                     | Edition / scope              | Notes |
| ------ | ------------------- | ---------------------------- | --------------------------------------------------------------------------- | ------------------------ | ---------------------------- | ----- |
| `[ ]`  | `AlarmsList`        | `alarms-list-widget`         | `src/components/widgets/alarm/alarms-list.vue`                              | `view_list`              | Core, top-level              |       |
| `[ ]`  | `Context`           | `entities-list-widget`       | `src/components/widgets/context/entities-list.vue`                          | `view_list`              | Core, top-level              |       |
| `[ ]`  | `ServiceWeather`    | `service-weather-widget`     | `src/components/widgets/service-weather/service-weather.vue`                | `view_module`            | Core, top-level              |       |
| `[ ]`  | `StatsCalendar`     | `stats-calendar-widget`      | `src/components/widgets/stats/calendar/stats-calendar.vue`                  | `calendar_today`         | Pro, top-level               |       |
| `[ ]`  | `Text`              | `text-widget`                | `src/components/widgets/text/text.vue`                                      | `view_headline`          | Core, top-level              |       |
| `[ ]`  | `Counter`           | `counter-widget`             | `src/components/widgets/counter/counter.vue`                                | `numbers`                | Core, top-level              |       |
| `[ ]`  | `Junit`             | `testing-weather-widget`     | `src/components/widgets/testing-weather/testing-weather.vue`                | `view_module`            | Pro, top-level               |       |
| `[ ]`  | `Map`               | `map-widget`                 | `src/components/widgets/map/map.vue`                                        | `location_on`            | Pro, top-level               |       |
| `[ ]`  | `BarChart`          | `bar-chart-widget`           | `src/components/widgets/chart/bar-chart-widget.vue`                         | `stacked_bar_chart`      | KPI/chart child widget       |       |
| `[ ]`  | `LineChart`         | `line-chart-widget`          | `src/components/widgets/chart/line-chart-widget.vue`                        | `show_chart`             | KPI/chart child widget       |       |
| `[ ]`  | `PieChart`          | `pie-chart-widget`           | `src/components/widgets/chart/pie-chart-widget.vue`                         | `donut_small`            | KPI/chart child widget       |       |
| `[ ]`  | `Numbers`           | `numbers-widget`             | `src/components/widgets/chart/numbers-widget.vue`                           | `functions`              | KPI/chart child widget       |       |
| `[ ]`  | `UserStatistics`    | `user-statistics-widget`     | `src/components/widgets/statistics/user-statistics-widget.vue`              | `perm_identity`          | Pro, statistics child widget |       |
| `[ ]`  | `AlarmStatistics`   | `alarm-statistics-widget`    | `src/components/widgets/statistics/alarm-statistics-widget.vue`             | `notification_important` | Pro, statistics child widget |       |
| `[ ]`  | `Availability`      | `availability-widget`        | `src/components/widgets/availability/availability-widget.vue`               | `published_with_changes` | Pro, top-level               |       |
| `[ ]`  | `ExternalDataTable` | `external-data-table-widget` | `src/components/widgets/external-data-table/external-data-table-widget.vue` | `database`               | Core, top-level              |       |

## Widget Blocks

Each widget is tracked by main rendering logic and settings/sidebar logic. Source setting sidebars are mapped from `SIDE_BARS_BY_WIDGET_TYPES` in `src/constants/widget.js`.

### `AlarmsList`

- `[ ]` Main content: `src/components/widgets/alarm/alarms-list.vue`.
- `[ ]` Settings: `alarm-settings`; migrate widget parameters forms, columns, filters, actions, export, expand panel settings and remediation-related settings.
- `[ ]` Store/API/socket: alarm list queries, periodic refresh, alarm socket rooms, user preferences and export. When migrating the alarm expand pbehavior tab, do not rely on old embedded entity cleanup; `pbehaviors-simple-list` should refresh its own list after pbehavior delete.
- `[ ]` Messages/i18n: not migrated.

### `Context`

- `[ ]` Main content: `src/components/widgets/context/entities-list.vue`.
- `[ ]` Settings: `context-settings`; migrate columns, filters, category/action settings and export settings.
- `[ ]` Store/API/socket: entity list queries, advanced search, user preferences and export. When migrating the entity expand pbehavior tab, use local pbehavior list refresh after delete instead of porting old `entities/removeFromStore`.
- `[ ]` Messages/i18n: not migrated.

### `ServiceWeather`

- `[ ]` Main content: `src/components/widgets/service-weather/service-weather.vue`.
- `[ ]` Settings: `service-weather-settings`; migrate service display, counters, pause/event settings and linked alarm list settings.
- `[ ]` Store/API/socket: service weather queries, counters and actions. `service-entity` uses `entity.pbehaviors?.length` to show the pbehavior tab; handle any stale-parent UX explicitly during widget migration, not by restoring normalizr embedded cleanup.
- `[ ]` Messages/i18n: not migrated.

### `StatsCalendar`

- `[ ]` Main content: `src/components/widgets/stats/calendar/stats-calendar.vue`.
- `[ ]` Settings: `stats-calendar-settings`; migrate calendar parameters, criticity levels, time range and linked alarm list settings.
- `[ ]` Store/API/socket: stats calendar queries and drilldown alarm list handoff.
- `[ ]` Messages/i18n: not migrated.

### `Text`

- `[ ]` Main content: `src/components/widgets/text/text.vue`.
- `[ ]` Settings: `text-settings`; migrate template editor and display options.
- `[ ]` Store/API/socket: runtime template variables and refresh behavior.
- `[ ]` Messages/i18n: not migrated.

### `Counter`

- `[ ]` Main content: `src/components/widgets/counter/counter.vue`.
- `[ ]` Settings: `counter-settings`; migrate levels, margins, columns, height, block template and linked alarm list settings.
- `[ ]` Store/API/socket: counter queries and alarm list handoff.
- `[ ]` Messages/i18n: not migrated.

### `Junit`

- `[ ]` Main content: `src/components/widgets/testing-weather/testing-weather.vue`.
- `[ ]` Settings: `testing-weather-settings`; migrate test suite parameters and display options.
- `[ ]` Store/API/socket: testing weather queries and statistics.
- `[ ]` Messages/i18n: not migrated.

### `Map`

- `[ ]` Main content: `src/components/widgets/map/map.vue`.
- `[ ]` Settings: `map-settings`; migrate map source, layers, breadcrumbs, marker options and filters.
- `[ ]` Store/API/socket: map queries and navigation state.
- `[ ]` Messages/i18n: not migrated.

### `BarChart`

- `[ ]` Main content: `src/components/widgets/chart/bar-chart-widget.vue`.
- `[ ]` Settings: `bar-chart-settings`; migrate chart metrics, sampling, time range, filters and display options.
- `[ ]` Store/API/socket: KPI metric queries.
- `[ ]` Messages/i18n: not migrated.

### `LineChart`

- `[ ]` Main content: `src/components/widgets/chart/line-chart-widget.vue`.
- `[ ]` Settings: `line-chart-settings`; migrate chart metrics, sampling, time range, filters and display options.
- `[ ]` Store/API/socket: KPI metric queries.
- `[ ]` Messages/i18n: not migrated.

### `PieChart`

- `[ ]` Main content: `src/components/widgets/chart/pie-chart-widget.vue`.
- `[ ]` Settings: `pie-chart-settings`; migrate chart metrics, show mode, time range, filters and display options.
- `[ ]` Store/API/socket: KPI metric queries.
- `[ ]` Messages/i18n: not migrated.

### `Numbers`

- `[ ]` Main content: `src/components/widgets/chart/numbers-widget.vue`.
- `[ ]` Settings: `numbers-settings`; migrate number metrics, time range, filters and display options.
- `[ ]` Store/API/socket: KPI metric queries.
- `[ ]` Messages/i18n: not migrated.

### `UserStatistics`

- `[ ]` Main content: `src/components/widgets/statistics/user-statistics-widget.vue`.
- `[ ]` Settings: `user-statistics-settings`; migrate user statistics parameters, filters and display options.
- `[ ]` Store/API/socket: user statistics queries.
- `[ ]` Messages/i18n: not migrated.

### `AlarmStatistics`

- `[ ]` Main content: `src/components/widgets/statistics/alarm-statistics-widget.vue`.
- `[ ]` Settings: `alarm-statistics-settings`; migrate alarm statistics parameters, filters and display options.
- `[ ]` Store/API/socket: alarm statistics queries.
- `[ ]` Messages/i18n: not migrated.

### `Availability`

- `[ ]` Main content: `src/components/widgets/availability/availability-widget.vue`.
- `[ ]` Settings: `availability-settings`; migrate availability parameters, time range, filters and display options.
- `[ ]` Store/API/socket: availability queries.
- `[ ]` Messages/i18n: not migrated.

### `ExternalDataTable`

- `[ ]` Main content: `src/components/widgets/external-data-table/external-data-table-widget.vue`.
- `[ ]` Settings: `external-data-table-settings`; migrate table selection, columns, filters and display options.
- `[ ]` Store/API/socket: external data table queries.
- `[ ]` Messages/i18n: not migrated.

## Dynamic Feature Extensions

The source project allows additional routes/widgets through `featuresService.get(...)`. The checked-in community source has `src/features/.gitignore`, so no extra feature files are present in this source tree.

Relevant extension points:

- `constants.WIDGET_TYPES`
- `constants.COMPONENTS_BY_WIDGET_TYPES`
- `constants.WIDGET_ICONS`
- `constants.SIDE_BARS`
- `constants.SIDE_BARS_BY_WIDGET_TYPES`
- `constants.WIDGET_TYPES_RULES`
