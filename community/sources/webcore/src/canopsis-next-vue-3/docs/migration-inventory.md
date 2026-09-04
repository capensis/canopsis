# Canopsis Next Migration Inventory

Source project: `/Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next`

This file tracks migration progress for top-level route pages and widget types from `canopsis-next`.

Status legend:

- `[x]` migrated
- `[~]` partially migrated
- `[ ]` not migrated

## Ready For Testing

Update this section every time a page, widget or shared feature becomes ready
for QA. Keep entries grouped by product area so QA can test
top-down without reading the full inventory.

### Administration

- `/admin/rights` - rights matrix, lightweight checkboxes and sticky actions.
- `/admin/users` - users table, create/edit user modal and shared tokens list/revoke.
- `/admin/roles` - roles table, create/edit/duplicate modal, template permissions, single delete and mass delete.
- `/admin/broadcast-messages` - active messages, list and create/edit flow.
- `/admin/playlists` - list, create/edit/duplicate, manage tabs, expand panel and actions.
- `/admin/planning` - pbehavior types, reasons and exception dates tabs with lists, expand panels, create/edit/import and mass actions.
- `/admin/kpi` - graphs, filters, rating settings and collection settings tabs.
- `/admin/events-records` - current status, events-record modal, apply-event-filter modal, actions and export.
- `/admin/external-auth-tokens` - external auth token list, expand panel, CRUD/delete modals, request fields, token variables and template testing tab.
- `/admin/external-data-tables` - data tables, records, import/export and table/record modals.
- `/admin/entity-infos-properties` - list, type filter, create/edit/duplicate and delete/bulk delete.
- `/admin/icons` - icon list, preview, create/edit upload modal and delete.
- `/admin/maps` - maps list, expand previews, geo/Mermaid/tree/flowchart map modals and nested editor required validation.
- `/admin/tags` - list, expand panel, create/edit modal and tag form pattern tabs.
- `/admin/views-import-export` - group/view selection, JSON export/import and import modal.
- `/admin/notifications` - notification settings form and save flow.
- `/admin/comment-templates` - list, create/edit modal and fields form.
- `/admin/jobs` - ticket status jobs tabs/table, filters, expand panel, play/pause/repeat actions and edit job modal/form.
- `/admin/template-testing` - data/tests tabs, CRUD modals and shared test variables/template studio stack.
- `/admin/widget-templates` - fully migrated: list, template type selector, create/edit modal, columns/quick-actions/sort/text template forms and widget settings consumers.
- `/admin/storage-settings` - storage settings, archive/clean confirmation actions and duration units.
- `/admin/state-settings` - list, state setting modals and create/edit/duplicate/remove actions.
- `/admin/healthcheck` - system status, connectors, connector history modal, graphs, parameters and engines metrics.
- `/admin/anomaly-monitored-connectors` - list, CRUD modal, history modal parity from healthcheck connector tiles and mass enable/disable/remove actions.
- `/admin/llms` - LLM list, important notes, create/edit form and prompts history expand panel.

### Exploitation

- `/exploitation/remediation` - instructions/configurations/jobs/statistics/instruction-stats tabs, execution components, approval modal and template-testing variables output editor.
- `/exploitation/pbehaviors` - list, expand tabs, advanced search, bulk actions, create/edit/duplicate through planning calendar, recurrence editor and recurrence calendar.
- `/exploitation/event-filters` - list, search/pagination, expand tabs, summary/patterns/external data/failures, create/edit/duplicate modal form, delete/export and bulk actions.
- `/exploitation/scenarios` - list, read-only expand content, create/edit/duplicate modal, test-query execution UI and DB export.
- `/exploitation/meta-alarm-rule` - list, expand patterns, bulk/delete/export actions and create/edit/duplicate modal form stack.
- `/exploitation/idle-rules` - list, expand, delete/bulk actions and create/edit/duplicate form with pbehavior action parameters.
- `/exploitation/flapping-rules` - list, expand content, delete/bulk actions and create/edit/duplicate alarm-status-rule modal.
- `/exploitation/resolve-rules` - list, expand content, delete/bulk actions and create/edit/duplicate alarm-status-rule modal.
- `/exploitation/link-rules` - list, mass actions, DB export, create/edit/duplicate modal, link rule form and Monaco advanced editor.
- `/exploitation/snmp-rules` - list, MIB upload, mass actions and create/edit/duplicate/delete SNMP rule modal.
- `/exploitation/declare-ticket-rules` - list, CRUD/mass actions, rule modal, test-query preview and ticket action modals.
- `/exploitation/dynamic-infos` - list, advanced search, Information/Patterns expand tabs, templates list, create/edit/duplicate modal with General/Patterns/Templates Studio tabs, add-from-template and mass actions.

### Profile And System

- `/login` - login page, basic/LDAP/CAS/SAML/OAuth2/OpenID auth variants, app-info polling, footer and i18n.
- `/` - home route default view redirect: user default view first, then role default view fallback with legacy info popups.
- `/kiosk-views/:id/:tabId` - standalone kiosk view route with tab widgets and view periodic refresh.
- `/alarms/:id` - standalone alarm route with single-alarm AlarmList table, optional widget configuration and expand panel.
- `/playlist/:id` - standalone playlist player with top-bar playback controls, tab rotation, autoplay and fullscreen handling.
- `/notifications/:tabId` - notifications tabs, socket refresh, approval/rating/event-filter failures and remediation stats alarms timeline.
- `/themes` - theme list, create/edit/duplicate theme modal and theme form.
- `/filters` - personal/corporate patterns page, create/edit/delete/bulk-delete flows and create-pattern modal.
- `/error` - error page polling and redirect logic.

### View Widgets

The main `/view/:id` page and these widget runtimes/settings are ready to test
inside migrated views:

- `AlarmsList`
- `Context`
- `ServiceWeather`
- `StatsCalendar`
- `Text`
- `Counter`
- `Junit`
- `Map`
- `BarChart`
- `LineChart`
- `PieChart`
- `Numbers`
- `UserStatistics`
- `AlarmStatistics`
- `Availability`
- `ExternalDataTable`

### Shared Features

- AI chat sidebar/form integration, pattern generation/update flow, model
  selector, suggestions, restart/cancel states and `/admin/llms` chat history
  replay modal.

## Immediate Tasks

- `[x]` Keep `normalizr` out of the target migration path. The Vue 3 project does not depend on `normalizr`; use plain normalized-by-id state, Pinia stores and explicit mapper helpers instead of restoring the old Vuex/normalizr shape.
  - Do not port the old `entities/removeFromStore` cleanup mechanism. Embedded metadata was only used internally by the old normalizr plugin for mark-and-sweep and parent cleanup.
  - Single pbehavior delete paths now use API delete plus local list refresh. `AlarmsList`, `Context` and `ServiceWeather` pbehavior tabs refresh locally instead of restoring global embedded cleanup.
- `[x]` Port remaining old entity store modules to Pinia facades.
  - Simple list/detail/bulk modules are wrapped with `createCRUDStore` where the old Vuex module was only CRUD around an API route.
  - Request-only or widget-scoped modules stay as explicit Pinia stores: `metrics` keeps per-widget buckets for aggregated/group/vector/availability data and meta, `template-testing` also covers old `copy/vars`, `snmp` is split into `snmp-rule` and `snmp-mib`, and `info` remains covered by `app-info` plus `infos`.
  - Empty old namespace folders `auth-token-job`, `instruction-job` and `webhook-job` do not require Pinia modules.
  - Consumers are validated through the page/widget entries that use those hooks.

## Cross-Cutting Migration Rules

- `[x]` Target dependency baseline for UI migration is Vuetify `^4.1.2` from `package.json`. Do not describe the target as Vuetify 3; adapt old Vuetify 2 markup against the actual Vuetify 4 API and the compatibility wrappers already present in this project.
- `[x]` Every completed migration stage must be recorded in documentation in the same change: update this inventory for page/widget/shared-component state and update the related `docs/ai/*-migration-plan.md` checklist when one exists.
- `[x]` When a migrated area becomes ready for QA, update `Ready For Testing` at the top of this file in the same change.
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

## AI Chat Migration

- `[x]` AI chat stack is migrated and wired into existing consumers:
  `hooks/ai/ai-chat.js`, `hooks/ai/ai-chat-form.js`,
  `helpers/entities/llm/chat/**`, `components/other/llm/chat/**` and
  `components/modals/llm/ai-chat-history.vue`.
- `[x]` Existing AI chat consumers pass sidebar bindings and events, generated
  patterns update scoped form fields, modal submit flows keep chat-history link
  behavior, and `/admin/llms` prompt history can open the chat history modal.
- `[x]` Targeted eslint and `npm run build` pass for the migrated stack.

Detailed plan: `docs/ai/ai-chat-migration-plan.md`.

## Shared Components

- `[x]` `c-fab-btn`: migrated to `src/components/common/buttons/c-fab-btn.vue`.
- `[x]` `c-refresh-btn`: migrated to `src/components/common/buttons/c-refresh-btn.vue` because `c-fab-btn` depends on it.
- `[x]` Modal base dependencies for maintenance mode: `c-modal-wrapper`, `c-form-block`, `c-form-block-row`, `text-editor-field`, `c-color-picker-field`.
- `[x]` Common modals: `confirmation`, `create-maintenance`.
- `[x]` User interface settings modal: settings menu action, modal registration, general/notifications/advanced/login tabs, app-info save flow, notification-settings save flow, rich text variables, logo upload and app title refresh are migrated and checked.
- `[x]` `c-duration-field`: full old duration value/unit field is migrated in the compat field barrel, with `long`, `medium`, `sameWidth`, custom `units`, clearable unit handling and pluralized unit labels adapted to Vue I18n `t`.
- `[x]` `c-payload-textarea-field` parity audit: old usage inventory is recorded in `docs/ai/coding-standards.md`; migrated uses are aligned, including meta alarm rule output template, and `MODALS.payloadTextareaEditor` is restored for simple template editors.
- `[x]` Text editor stack: `text-editor`, `text-editor-field`, `text-editor-with-template`, `c-widget-template-field`, variables menu/list, upload hook, template-vars hook, image/images viewer modals, sidebar text-editor fields and editor modal vee-validate form scopes are migrated.
- `[x]` `c-draggable-list-field`: migrated with Vue 3 `vuedraggable` API; emits `update:value` / `update:model-value` matching its `value` / `modelValue` props.
- `[x]` `c-treeview`: migrated as a Vuetify 4 `VTreeview` compatibility wrapper with legacy `item-key` / `item-text` / `open` / `active` / `value` aliases and old `label` slot `leaf` support. `variables-help` now uses this wrapper instead of a local recursive tree.
- `[x]` `c-advanced-data-table`: migrated on top of Vuetify 4 `v-data-table-server` with legacy-compatible `headers.text/value`, `headers.title/key`, `server-items-length` through `totalItems`, `options`, `sortBy/sortDesc`, `show-select`, `show-expand`, `item`, `item.<field>`, `header`, `header.<field>`, `body.prepend`, `progress`, `toolbar`, `expand` and `mass-actions` behavior. Selection filters disabled rows, keeps selected objects/keys, clears or preserves selected rows through mass-action controls, and resets stale selected/expanded rows when items/options change. Expand rows keep the old custom expand slot, expandable predicate, colspan, transition cleanup and reset behavior. Search and pagination are wired through `c-search`, `c-advanced-search`, `c-table-pagination`, `c-pagination`, `c-items-per-page-field` and shared query helpers. The component is in active use by migrated consumers including users, roles, broadcast messages, template-testing, dynamic-info, event-filter, idle-rules and scenarios. Deep advanced-search editor parity, if needed later, is tracked by the search stack rather than by `c-advanced-data-table`.
- `[x]` `c-patterns-field` migration strategy:
  - `[x]` Scope rule: old component tree is ported without rewriting the main logic; Options API components stay Options API and Composition API hooks stay hooks.
  - `[x]` Phase 1, shared model/helpers: `helpers/entities/pattern/form.js`, `helpers/entities/filter/form.js`, `helpers/entities/pattern/fields/form.js` and required constants are available.
  - `[x]` Phase 2, Pinia stores/hooks: migrated `pattern`, `pbehavior-patterns`, `patterns-fields`, `pattern-entities-optimize`, entity-category and pbehavior dictionary facades. Hooks wrap stores; components do not call request directly.
  - `[x]` Phase 3, base rule editor leaves: migrated `pattern-attribute-field`, `pattern-operator-field`, `pattern-rule-field-date-value`, `pattern-rule-field`, `pattern-group-field`, `pattern-groups-field`, `pattern-editor-field`, `pattern-advanced-editor-field`, `c-pattern-field`, `c-filter-field`.
  - `[x]` Phase 4, domain pattern fields: migrated `c-alarm-patterns-field`, `c-entity-patterns-field`, `c-pbehavior-patterns-field`, `c-event-filter-patterns-field`, `c-service-weather-patterns-field`; added lightweight compatibility fields needed by their dynamic `valueField` configs.
  - `[x]` Phase 5, counters and preview modals: migrated `pattern-counters`, `pattern-count-message`, `pattern-count-alarms-modal`, `pattern-count-entities-modal`. Alarm-list modal body is available through the migrated AlarmsList stack; entities-list and entities-comparison modal bodies are migrated through the Context entities table.
  - `[x]` Phase 6, optimization/suggestions: migrated `pattern-optimization`, `pattern-progress`, `pattern-try-optimization`, `pattern-field-suggestions-wrapper`, `pattern-suggestions`, `pattern-suggestions-list`, `pattern-suggestions-header`, `pattern-suggestion-tab`, `pattern-suggestion-content`, `pattern-fix-with-ai-btn`; API calls go through pattern-entities-optimize store/hook.
  - `[x]` Phase 7, `c-patterns-field` wrapper: migrated `c-patterns-field.vue`; old `v-field` usage is replaced with explicit `:value` / `@update:value` bindings.
  - `[x]` Phase 8, consumers: integrated in migrated pattern consumers, including scenario modal/action patterns, event-filter, dynamic-info, meta-alarm-rule, declare-ticket-rule, pbehavior, remediation, storage/state settings, tags, filters, KPI filters, link rules and alarm-status rules.
  - `[x]` Verification checklist: build and lint pass. Functional QA is tracked through `Ready For Testing`.
- `[x]` Common button components from old `components/common/buttons`: `c-action-btn`, `c-action-fab-btn`, `c-fab-expand-btn`, `c-fab-btn`, `c-refresh-btn`, `c-expand-btn`, `c-btn-with-error`, `c-db-export-btn`, `c-copy-btn`, `c-download-btn`, `c-speed-dial`, `c-grid-btns` and `c-density-btn-toggle` are migrated, exported from the common buttons barrel and used by migrated pages/widgets where needed. Density field integration is migrated through `src/components/sidebars/form/fields/density.vue`.

## Shared Hooks And Mixins

- `[x]` Query hooks: `hooks/query/local-query.js`, `hooks/query/options.js`, `hooks/query/shared.js`.
- `[x]` Query pending hook: `hooks/query/pending.js`; extended with `initialPending`.
- `[x]` Query mixins: `mixins/query/query.js`, `mixins/query/interval.js`.
- `[x]` Query helper dependencies: `helpers/query.js`, `helpers/entities/shared/query.js`, `helpers/collection.js`, `helpers/date/date-intervals.js`.
- `[x]` Groups side bar stack: `groups-side-bar`, nested group/view/playlists/draggable components, navigation editing mode, view/group/playlist/active-view Pinia stores, related mixins and create/edit group/view modal entry points are migrated.
- `[x]` Query-related entity helpers from old `helpers/entities/*/query.js`: migrated with parity for `alarm`, `availability`, `counter`, `entity`, `external-data-table`, `map`, `metric`, `service-weather`, `shared`, `stats` and `widget` query helpers. `Map` widget query conversion is registered in the shared widget query dispatcher.

## Pages

Pages are taken from `src/router.js` in the source project.

| Status | Route name                                     | Path                                  | Source component                                             | Target component                                                  | Notes                                                                                                                                                                                                                                                                                                                                                                        |
| ------ | ---------------------------------------------- | ------------------------------------- | ------------------------------------------------------------ | ----------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `[x]`  | `login`                                        | `/login`                              | `src/views/login.vue`                                        | `src/views/login-page.vue`                                        | Login page, basic/LDAP/CAS/SAML/OAuth2/OpenID auth variants, app-info polling, footer, anonymous route access and i18n are migrated.                                                                                                                                                                                                                                         |
| `[x]`  | `home`                                         | `/`                                   | `src/views/home.vue`                                         | `src/views/home-page.vue`                                         | Default-view redirect is migrated in `app.vue`: user default view first, then role default view fallback with legacy info popups. The standalone `home-page.vue` is an intentional empty fallback while redirect state is resolving.                                                                                                                                          |
| `[x]`  | `view`                                         | `/view/:id`                           | `src/views/view.vue`                                         | `src/views/view-page.vue`                                         | Foundation, page shell, tabs, grid shell, widget wrapper shell, floating controls including executions menu, create-tab flow, basic create-widget modal, migrated widgets/settings and view screen-mode kiosk behavior are migrated. Playlist playback now uses the same migrated view/tab/widget stack.             |
| `[x]`  | `view-kiosk`                                   | `/kiosk-views/:id/:tabId`             | `src/views/view-kiosk.vue`                                   | `src/views/view-kiosk-page.vue`                                   | Standalone kiosk page is migrated with active view loading, tab fallback, kiosk screen mode, widget rendering and view periodic refresh.                                                                                                                                                                                                                                      |
| `[x]`  | `alarms`                                       | `/alarms/:id`                         | `src/views/alarm.vue`                                        | `src/views/alarm-page.vue`                                        | Standalone alarm page is migrated with one-alarm fetch, optional `widgetId` query configuration, AlarmList table rendering, expand panel and legacy technical-view route permission.                                                                                                                                                                                          |
| `[x]`  | `admin-rights`                                 | `/admin/rights`                       | `src/views/admin/permissions.vue`                            | `src/views/admin/rights-page.vue`                                 | Main rights table, lightweight checkbox, sticky actions, route/store wiring, i18n and visual parity are migrated.                                                                                                                                                                                                                                                            |
| `[x]`  | `admin-users`                                  | `/admin/users`                        | `src/views/admin/users.vue`                                  | `src/views/admin/users-page.vue`                                  | Users tab/table, create/edit user form modal, shared tokens tab/list/revoke flow, route/store wiring and i18n are migrated.                                                                                                                                                                                                                                                  |
| `[x]`  | `admin-roles`                                  | `/admin/roles`                        | `src/views/admin/roles.vue`                                  | `src/views/admin/roles-page.vue`                                  | Main list, search/pagination, single delete, mass delete, create/edit/duplicate role modal, template permissions, route/store wiring and i18n are migrated. Full permissions matrix editing belongs to `/admin/rights`.                                                                                                                                                    |
| `[x]`  | `admin-broadcast-messages`                     | `/admin/broadcast-messages`           | `src/views/admin/broadcast-messages.vue`                     | `src/views/admin/broadcast-messages-page.vue`                     | Main list, active messages, create/edit form, route/store wiring and i18n are migrated.                                                                                                                                                                                                                                                                                      |
| `[x]`  | `admin-playlists`                              | `/admin/playlists`                    | `src/views/admin/playlists.vue`                              | `src/views/admin/playlists-page.vue`                              | Main list, refresh, copy/play actions, delete/bulk actions, create/edit/duplicate playlist modal, manage tabs modal and expand panel are migrated.                                                                                                                                                                                                                           |
| `[x]`  | `admin-planning-administration`                | `/admin/planning`                     | `src/views/admin/planning.vue`                               | `src/views/admin/planning-page.vue`                               | Tabs, lists, expand panels, create buttons and pbehavior type/reason/exception create/edit/import modal forms are migrated.                                                                                                                                                                                                                                                  |
| `[x]`  | `admin-healthcheck`                            | `/admin/healthcheck`                  | `src/views/admin/healthcheck.vue`                            | `src/views/admin/healthcheck-page.vue`                            | Main page tabs are wired for system status, connectors, connector history modal, graphs, parameters and engines metrics.                                                                                                                                                                                                                                                     |
| `[x]`  | `admin-kpi`                                    | `/admin/kpi`                          | `src/views/admin/kpi.vue`                                    | `src/views/admin/kpi-page.vue`                                    | Graphs, filters, rating settings, collection settings tabs, route/store wiring and i18n are migrated.                                                                                                                                                                                                                                                                        |
| `[x]`  | `admin-events-records`                         | `/admin/events-records`               | `src/views/admin/events-records.vue`                         | `src/views/admin/events-records-page.vue`                         | Main page, current status, events-record modal, apply-event-filter modal, record actions, export flows, route/store wiring and i18n are migrated.                                                                                                                                                                                                                            |
| `[x]`  | `admin-anomaly-monitored-connectors`           | `/admin/anomaly-monitored-connectors` | `src/views/admin/anomaly-monitored-connectors.vue`           | `src/views/admin/anomaly-monitored-connectors-page.vue`           | Main page, list, search/pagination, single create/edit/delete modal, history modal parity from healthcheck connector tiles and mass enable/disable/remove actions are migrated.                                                                                                                                                                                               |
| `[x]`  | `admin-template-testing`                       | `/admin/template-testing`             | `src/views/admin/template-testing.vue`                       | `src/views/admin/template-testing-page.vue`                       | Main data/tests tabs, lists, create/edit data modal, edit/remove test modals and shared test variables/template studio stack are migrated.                                                                                                                                                                                                                                    |
| `[x]`  | `admin-jobs-management`                        | `/admin/jobs`                         | `src/views/admin/jobs.vue`                                   | `src/views/admin/jobs-page.vue`                                   | Main jobs tabs/table, filters, expand panel, play/pause/repeat actions, edit modal/form and routing are migrated.                                                                                                                                                                                                                                                            |
| `[x]`  | `admin-custom-objects-external-auth-tokens`    | `/admin/external-auth-tokens`         | `src/views/admin/custom-objects/external-auth-tokens.vue`    | `src/views/admin/custom-objects/external-auth-tokens-page.vue`    | Main list, expand panel, create/edit/delete modals, request form fields, token variables, route, webhook-token-rule CRUD and template testing tab parity are migrated.                                                                                                                                                                                                        |
| `[x]`  | `admin-custom-objects-external-data-tables`    | `/admin/external-data-tables`         | `src/views/admin/custom-objects/external-data-tables.vue`    | `src/views/admin/custom-objects/external-data-tables-page.vue`    | Page, list, table/record/import/export modals, store actions, route wiring and i18n are migrated.                                                                                                                                                                                                                                                                            |
| `[x]`  | `admin-custom-objects-entity-infos-properties` | `/admin/entity-infos-properties`      | `src/views/admin/custom-objects/entity-infos-properties.vue` | `src/views/admin/custom-objects/entity-infos-properties-page.vue` | Main page, list, type filter, create/edit/duplicate form modal, delete, bulk delete actions, route/store wiring and i18n are migrated.                                                                                                                                                                                                                                       |
| `[x]`  | `admin-custom-objects-icons`                   | `/admin/icons`                        | `src/views/admin/custom-objects/icons.vue`                   | `src/views/admin/custom-objects/icons-page.vue`                   | Main page, icons list, custom icon preview/registration, create/edit upload modal, delete action, route/store wiring and i18n are migrated.                                                                                                                                                                                                                                  |
| `[x]`  | `admin-custom-objects-maps`                    | `/admin/maps`                         | `src/views/admin/custom-objects/maps.vue`                    | `src/views/admin/custom-objects/maps-page.vue`                    | Main page, CRUD table, search/pagination, single delete, mass delete wiring, expand panel shell, geo/Mermaid/tree/flowchart expand previews, create/edit/duplicate map modals, map form/editor stack and nested editor required validation are migrated.                                                                                                                       |
| `[x]`  | `admin-custom-objects-tags`                    | `/admin/tags`                         | `src/views/admin/custom-objects/tags.vue`                    | `src/views/admin/custom-objects/tags-page.vue`                    | Main page, list, expand panel, create/edit modal, tag form with pattern tabs, route/store wiring and i18n are migrated.                                                                                                                                                                                                                                                      |
| `[x]`  | `admin-custom-objects-llms`                    | `/admin/llms`                         | `src/views/admin/custom-objects/llms.vue`                    | `src/views/admin/custom-objects/llms-page.vue`                    | Main page, list, important notes banner, create/edit form and prompts history expand panel are migrated.                                                                                                                                                                                                        |
| `[x]`  | `admin-settings-views-import-export`           | `/admin/views-import-export`          | `src/views/admin/settings/views-import-export.vue`           | `src/views/admin/settings/views-import-export-page.vue`           | Main page, group/view selection, JSON export, JSON import parsing, import modal with draggable groups/views, route/store wiring and i18n are migrated.                                                                                                                                                                                                                       |
| `[x]`  | `admin-settings-notifications`                 | `/admin/notifications`                | `src/views/admin/settings/notifications-settings.vue`        | `src/views/admin/settings/notifications-settings-page.vue`        | Main page, notification settings form, fetch/update flow, routing and i18n are migrated.                                                                                                                                                                                                                                                                                     |
| `[x]`  | `admin-settings-comment-templates`             | `/admin/comment-templates`            | `src/views/admin/settings/comment-templates.vue`             | `src/views/admin/settings/comment-templates-page.vue`             | Main page, list table, create/edit modal, fields form, CRUD flow, routing and i18n are migrated.                                                                                                                                                                                                                                                                             |
| `[x]`  | `admin-settings-widget-templates`              | `/admin/widget-templates`             | `src/views/admin/settings/widget-templates.vue`              | `src/views/admin/settings/widget-templates-page.vue`              | Fully migrated: main page, list table, type selector modal, create/edit modal, columns/quick-actions/sort/text template forms, route/store wiring, i18n and widget settings consumers are migrated.                                                                                                                                                                           |
| `[x]`  | `admin-settings-storage-settings`              | `/admin/storage-settings`             | `src/views/admin/settings/storage-settings.vue`              | `src/views/admin/settings/storage-settings-page.vue`              | Main page, data-storage form mapper, storage settings form blocks, history messages, archive/clean confirmation actions, route/store wiring and i18n are migrated.                                                                                                                                                                                                           |
| `[x]`  | `admin-settings-state-settings`                | `/admin/state-settings`               | `src/views/admin/settings/state-settings.vue`                | `src/views/admin/settings/state-settings-page.vue`                | List table, create/edit/duplicate/remove actions, state setting modals, route/store wiring and i18n are migrated.                                                                                                                                                                                                                                                            |
| `[x]`  | `exploitation-pbehaviors`                      | `/exploitation/pbehaviors`            | `src/views/exploitation/pbehaviors.vue`                      | `src/views/exploitation/pbehaviors-page.vue`                      | Main page, list table, expand tabs, advanced search, delete/export, bulk actions, alarm filtering, create/edit/duplicate list actions, pbehavior planning calendar, `create-pbehavior`, recurrence rule editor and recurrence calendar modals are migrated.                                                                                                                     |
| `[x]`  | `exploitation-event-filters`                   | `/exploitation/event-filters`         | `src/views/exploitation/event-filters.vue`                   | `src/views/exploitation/event-filters-page.vue`                   | Main list, expand content, failures tab, remove/export, bulk actions, create/edit/duplicate modal form, `create-event-filter` registration, event-filter form stack and store hooks are migrated.                                                                                                                                                                            |
| `[x]`  | `exploitation-snmp-rules`                      | `/exploitation/snmp-rules`            | `src/views/exploitation/snmp-rules.vue`                      | `src/views/exploitation/snmp-rules-page.vue`                      | Main list, search/pagination, MIB upload, mass actions, create/edit/duplicate/delete modal, SNMP rule form, route/store wiring and i18n are migrated.                                                                                                                                                                                                                        |
| `[x]`  | `exploitation-dynamic-infos`                   | `/exploitation/dynamic-infos`         | `src/views/exploitation/dynamic-infos.vue`                   | `src/views/exploitation/dynamic-infos-page.vue`                   | Page shell, dynamic infos list with advanced search, Information/Patterns expand tabs, dynamic info templates list, create/edit modal with General/Patterns/Templates Studio tabs, add-from-template modal, route/store wiring, i18n and dynamic-info/template mass actions are migrated.                                                                                    |
| `[x]`  | `exploitation-meta-alarm-rules`                | `/exploitation/meta-alarm-rule`       | `src/views/exploitation/meta-alarm-rules.vue`                | `src/views/exploitation/meta-alarm-rules-page.vue`                | Main list, read-only expand patterns, route, helper, bulk/delete/export actions and create/edit/duplicate modal form stack are migrated.                                                                                                                                                                                                                                      |
| `[x]`  | `exploitation-scenarios`                       | `/exploitation/scenarios`             | `src/views/exploitation/scenarios.vue`                       | `src/views/exploitation/scenarios-page.vue`                       | Main list/read-only expand content, scenario store, create/edit/duplicate modal, test-query execution UI and row DB export button are migrated.                                                                                                                                                                                                                              |
| `[x]`  | `exploitation-idle-rules`                      | `/exploitation/idle-rules`            | `src/views/exploitation/idle-rules.vue`                      | `src/views/exploitation/idle-rules-page.vue`                      | Main list, expand, delete/bulk actions and create/edit/duplicate form are migrated, including pbehavior action parameters.                                                                                                                                                                                                                                                   |
| `[x]`  | `exploitation-flapping-rules`                  | `/exploitation/flapping-rules`        | `src/views/exploitation/flapping-rules.vue`                  | `src/views/exploitation/flapping-rules-page.vue`                  | Main list, read-only expand, delete/bulk actions and create/edit/duplicate form are migrated.                                                                                                                                                                                                                                                                                |
| `[x]`  | `exploitation-resolve-rules`                   | `/exploitation/resolve-rules`         | `src/views/exploitation/resolve-rules.vue`                   | `src/views/exploitation/resolve-rules-page.vue`                   | Main list, read-only expand, delete/bulk actions and create/edit/duplicate form are migrated.                                                                                                                                                                                                                                                                                |
| `[x]`  | `exploitation-declare-ticket-rules`            | `/exploitation/declare-ticket-rules`  | `src/views/exploitation/declare-ticket-rules.vue`            | `src/views/exploitation/declare-ticket-rules-page.vue`            | Page, list, CRUD actions, mass actions, create/edit/duplicate rule modal, test-query preview, declare/associate/execute ticket modals, alarm-widget flows and required helpers are migrated.                                                                                                                                                                                  |
| `[x]`  | `exploitation-link-rules`                      | `/exploitation/link-rules`            | `src/views/exploitation/link-rules.vue`                      | `src/views/exploitation/link-rules-page.vue`                      | Page, list, mass actions, DB export, create/edit/duplicate modal, link rule form and Monaco-based advanced editor are migrated.                                                                                                                                                                                                                                              |
| `[x]`  | `exploitation-remediation`                     | `/exploitation/remediation`           | `src/views/exploitation/remediation.vue`                     | `src/views/exploitation/remediation-page.vue`                     | Instructions, configurations, jobs, statistics, instruction stats, execution components, approval modal and template-testing variables output editor are migrated.                                                                                                                                                                                                            |
| `[x]`  | `playlist`                                     | `/playlist/:id`                       | `src/views/playlist.vue`                                     | `src/views/playlist-page.vue`                                     | Standalone playlist player is migrated with playlist fetch, tab fetch, top-bar playback controls, timer rotation, autoplay, fullscreen and migrated `view-tab-widgets` rendering.                                                                                                                                                                                           |
| `[x]`  | `profile-patterns`                             | `/filters`                            | `src/views/profile/patterns.vue`                             | `src/views/profile/patterns-page.vue`                             | Page shell, personal/corporate pattern lists, corporate permission tabs, create/edit/delete/bulk-delete flows, route and `create-pattern` modal are migrated.                                                                                                                                                                                                                |
| `[x]`  | `profile-themes`                               | `/themes`                             | `src/views/profile/themes.vue`                               | `src/views/profile/themes-page.vue`                               | Page, list, route, create/edit/duplicate theme modal, theme form, compact color picker menu, theme form helpers, route/store wiring and i18n are migrated.                                                                                                                                                                                                                   |
| `[x]`  | `notifications`                                | `/notifications/:tabId`               | `src/views/notification/notifications.vue`                   | `src/views/notification/notifications-page.vue`                   | Page shell, tab routing, socket refresh, approve/rate/event-filter-failure tabs, active notification handling, rating modal and remediation stats alarms timeline expand tab are migrated.                                                                                                                                                                                   |
| `[x]`  | `error`                                        | `/error`                              | `src/views/error.vue`                                        | `src/views/error-page.vue`                                        | Migrated with polling and redirect logic.                                                                                                                                                                                                                                                                                                                                    |

## Page Blocks

Each page is tracked by functional blocks. Add concrete component/file names inside the block as we start migrating that page.

### `login`

- `[x]` Main content: `src/views/login.vue` -> `src/views/login-page.vue`; login card, basic/LDAP/CAS/SAML/OAuth2/OpenID auth variants and footer are migrated.
- `[x]` Modal windows: none in the old page.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: app-info polling, auth feature flags, footer/description app-info and login flow are migrated.
- `[x]` Routing/access: `/login`, `hideHeader`, anonymous access and protected-route redirect handling are migrated.
- `[x]` Messages/i18n: login EN/FR messages are present and wired.

### `home`

- `[x]` Main content: `src/views/home.vue` -> `src/views/home-page.vue`; home redirects to the user default view first, then role default view fallback. `home-page.vue` is an intentional empty fallback while redirect state is resolving.
- `[x]` Modal windows: none in the old page.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: current-user default view and role default view are reused from auth state.
- `[x]` Routing/access: `/` route exists and redirects to user default view, or role default view when the user default view is missing/forbidden.
- `[x]` Messages/i18n: legacy default-view redirect info messages are wired.

### `view`

- `[x]` Migration plan: staged plan is tracked in `docs/ai/view-migration-plan.md`.
- `[x]` Main content: `src/views/view.vue` -> `src/views/view-page.vue`; active view lifecycle, periodic refresh provider, tabs, grid shell, widget wrapper, floating controls including executions menu, tab/widget creation entry points, migrated widget types/settings and view screen-mode kiosk behavior are migrated.
- `[x]` Tabs: view tab navigation, active tab query, tab reorder submit/cancel overlay, create/update/copy/delete tab flows and `select-view` modal are migrated.
- `[x]` Grid system: layout shell, breakpoint handling, drag/resize rules, edit mode and widget grid position saving are migrated.
- `[x]` Grid system content: `view-tab-widgets`, `widget-wrapper`, `widget-edit-drag-handler`, `c-alert-overlay`, `window-size-field` and migrated widget rendering are complete. AlarmList uses the alarm Pinia store for widget-scoped query/pagination/pending/meta state; toolbar, pagination, columns/cells, row actions, custom checkbox selection column, mass-actions panel, column drag/resize/reset controls, expand panel tabs, sticky header and sticky horizontal scroll are complete.
- `[x]` Modal windows: `select-view`, create/edit/duplicate view, create/update/copy/delete tab, create-widget, widget duplication, widget delete and widget settings sidebar modal-like flows are migrated.
- `[x]` Sidebars/drawers: migrated widget settings sidebars use `useWidgetSettings`; shared settings stacks, dirty-close confirmation and widget settings forms are migrated. The old widget settings mixin is removed.
- `[x]` State/API/socket: active-view foundation, view tab store actions, widget store action verification and view router helpers are migrated. View refresh now force-updates the active tab query nonce so widget lists (alarms list, details, etc.) refetch with the view.
- `[x]` Routing/access: `/view/:id` now points to `view-page.vue`; old route only required login, so route meta matches the source project.
- `[x]` Messages/i18n: existing view, modal, widget, common and settings messages are wired for the migrated view stack.

### `view-kiosk`

- `[x]` Standalone page: `src/views/view-kiosk.vue` -> `src/views/view-kiosk-page.vue`.
- `[x]` Main content: renders `view-tab-widgets` for the selected tab in kiosk screen mode.
- `[x]` Modal windows: none in the old page.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: active view fetch, tab fallback and view-level periodic refresh are migrated.
- `[x]` Routing/access: `/kiosk-views/:id/:tabId`, simple navigation, view permission by route id.
- `[x]` Messages/i18n: no dedicated old page messages.

### `alarms`

- `[x]` Main content: `src/views/alarm.vue` -> `src/views/alarm-page.vue`; renders the migrated AlarmList table for the requested alarm.
- `[x]` Modal windows: old page has no dedicated modals; row/expand action modals come from the migrated AlarmList stack.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: one-alarm fetch through the alarm Pinia store, optional `widgetId` fetch through the widget store and default AlarmList widget fallback are migrated.
- `[x]` Routing/access: `/alarms/:id`, technical view permission.
- `[x]` Messages/i18n: no dedicated old page messages; AlarmList messages are reused.

### `admin-rights`

- `[x]` Main content: `src/views/admin/permissions.vue` -> `src/views/admin/rights-page.vue`; tabs, search, expand/collapse, roles matrix, lightweight checkbox and sticky submit/cancel are migrated and visually checked.
- `[x]` Modal windows: none in the old page.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: permission and role Pinia stores/hooks are migrated for the rights page.
- `[x]` Routing/access: `/admin/rights`, permission management access.
- `[x]` Messages/i18n: existing permission messages are wired.

### `admin-users`

- `[x]` Main content: `src/views/admin/users.vue`; users tab, users table and shared tokens tab/table are migrated and checked.
- `[x]` Modal windows: confirmation modal and create/edit user form modal are migrated; shared token revoke uses confirmation like the old page.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: user CRUD and share-token list/remove stores/hooks are migrated and checked for the users page.
- `[x]` Routing/access: `/admin/users`, user management access.
- `[x]` Messages/i18n: existing user, role, permission, modal and common messages are wired.

### `admin-roles`

- `[x]` Main content: `src/views/admin/roles.vue` -> `src/views/admin/roles-page.vue`; roles table with search, pagination, edit, duplicate, remove and mass remove actions is migrated.
- `[x]` Modal windows: `create-role` modal and role form are migrated with template permissions, type, description, theme, expiration settings and default view. Full permissions matrix editing is handled on `/admin/rights`, not in the role modal.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: role Pinia store/hook is reused and extended with previous-params fetching.
- `[x]` Routing/access: `/admin/roles`, role management access.
- `[x]` Messages/i18n: existing role, permission, modal and common messages are wired.

### `admin-broadcast-messages`

- `[x]` Main content: `src/views/admin/broadcast-messages.vue` -> `src/views/admin/broadcast-messages-page.vue`; table, search, pagination, status and actions are migrated.
- `[x]` Modal windows: `create-broadcast-message` modal and form are migrated with new `useForm`.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: broadcast message Pinia store/hook, active broadcast messages, mark-as-read and maintenance edit integration are migrated.
- `[x]` Routing/access: `/admin/broadcast-messages`, broadcast message access.
- `[x]` Messages/i18n: existing broadcast message and modal messages are wired.

### `admin-playlists`

- `[x]` Main content: `src/views/admin/playlists.vue` -> `src/views/admin/playlists-page.vue`; route is wired to a real page and the playlists table is migrated with search, pagination, play link, copy link, create/edit/duplicate, delete, bulk actions and read-only expand panel.
- `[x]` Modal windows: `create-playlist` and `manage-playlist-tabs` are migrated with `useForm`; selected tabs are stored as tab objects in form state and converted to `_id` list for API submit.
- `[x]` Sidebars/drawers: no playlist sidebars/drawers are required for this page.
- `[x]` State/API/socket: playlist Pinia CRUD store/hook is reused; current user refresh runs after single delete and create/update.
- `[x]` Routing/access: `/admin/playlists`, playlist management access.
- `[x]` Messages/i18n: existing playlist/common modal messages are wired.

### `admin-planning-administration`

- `[x]` Main content: `src/views/admin/planning.vue` -> `src/views/admin/planning-page.vue`; route is wired to a real page with pbehavior type/reason/exception tabs, migrated list tables, read-only expand panels and create/import FAB actions. Type, reason and exception date edit actions are wired; duplicate is intentionally absent because the old page did not have it. Type expand shows description/type tabs with icon/color, reason expand shows description, exception expand shows description/periods tabs with formatted dates and type metadata.
- `[x]` Modal windows: `create-pbehavior-type`, `create-pbehavior-reason`, `create-pbehavior-exception` and `import-pbehavior-exception` are migrated and registered with `useForm`. The forms keep the old `c-form-block` layout, modal date picker variant, primary date header, adjacent months, and no premature validation on open.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: pbehavior type/reason/exception Pinia CRUD stores/hooks are reused for list, delete and bulk hide/unhide/remove.
- `[x]` Routing/access: `/admin/planning`, grouped planning permission and per-tab read permissions.
- `[x]` Messages/i18n: existing pbehavior/common/modal messages are wired.

### `admin-healthcheck`

- `[x]` Main content: `src/views/admin/healthcheck.vue` -> `src/views/admin/healthcheck-page.vue`; system status, connectors, graphs, parameters and engines metrics tabs are wired to the real page.
- `[x]` Modal windows: healthcheck engine modals and anomaly monitored connector history modal used by the Connectors tab are migrated.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: top-bar healthcheck chips, admin page data fetching, connector states pagination/socket flow, graph metric flows and connector history modal data flow are migrated.
- `[x]` Routing/access: `/admin/healthcheck`, healthcheck `can` access.
- `[x]` Messages/i18n: existing healthcheck messages are wired.

### `admin-kpi`

- `[x]` Main content: `src/views/admin/kpi-page.vue` is wired; Graphs tab has SLI, alarm metrics and rating metrics fetch/render/export paths; KPI filters, rating settings and collection settings tabs are migrated.
- `[x]` Modal windows: `create-kpi-filter` is migrated. The old KPI page has no separate rating/collection modals; rating settings and collection settings are inline tab forms.
- `[x]` Sidebars/drawers: not used by this page.
- `[x]` State/API/socket: filter CRUD, rating settings list/bulk update, metrics settings fetch/update, healthcheck socket availability, KPI graph metric fetch and CSV/PNG export paths are wired.
- `[x]` Routing/access: `/admin/kpi`, KPI `can` access.
- `[x]` Messages/i18n: KPI messages already exist.

### `admin-events-records`

- Working plan: `docs/ai/admin-events-records-migration-plan.md`.
- `[x]` Main content: `src/views/admin/events-records-page.vue`.
- `[x]` Modal windows: `events-record`, `apply-event-filter`.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: events-record Pinia store/hook, current status polling, recording/resending/export flows.
- `[x]` Routing/access: `/admin/events-records`, events record `can` access.
- `[x]` Messages/i18n: base events-record/modal/page-header messages already present.

### `admin-anomaly-monitored-connectors`

- `[x]` Main content: `src/views/admin/anomaly-monitored-connectors.vue` -> `src/views/admin/anomaly-monitored-connectors-page.vue`; list, search, pagination, refresh, create/edit/delete entry points and mass enable/disable/remove actions are migrated.
- `[x]` Modal windows: `create-anomaly-monitored-connector` modal/form and `anomaly-monitored-connector-history` modal from healthcheck connector tiles are migrated.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: Pinia store/hook, CRUD, enabled patch, states/history fetches and bulk actions are migrated.
- `[x]` Routing/access: `/admin/anomaly-monitored-connectors`, anomaly monitored connector access.
- `[x]` Messages/i18n: base page-header, modal and anomaly monitored connector messages already present.

### `admin-template-testing`

- `[x]` Main content: `src/views/admin/template-testing.vue` -> `src/views/admin/template-testing-page.vue`; data and tests tabs with search, pagination, type filters and row actions are migrated.
- `[x]` Modal windows: `create-template-testing-data` and `create-template-testing-test` are migrated with `useForm`.
- `[x]` Sidebars/drawers: none in the old route page; shared test variables/template studio stack is migrated for consumer forms.
- `[x]` State/API/socket: template testing Pinia store/hook reused for data/tests CRUD, list fetching and test variables validation.
- `[x]` Routing/access: `/admin/template-testing`, template testing `can` access.
- `[x]` Messages/i18n: existing template testing and modal messages are wired.

### `admin-jobs-management`

- `[x]` Main content: `src/views/admin/jobs.vue`; migrated to `src/views/admin/jobs-page.vue` with tabs, table, filters, row expansion and play/pause/repeat actions.
- `[x]` Modal windows: edit ticket status job modal/form is migrated with declare-ticket check-status settings.
- `[x]` Sidebars/drawers: none in the migrated jobs page.
- `[x]` State/API/socket: existing ticket status job Pinia store/hook and bulk play/pause endpoints are reused.
- `[x]` Routing/access: `/admin/jobs`, job management `can` access.
- `[x]` Messages/i18n: existing jobs/modal messages are reused.

### `admin-custom-objects-external-auth-tokens`

- `[x]` Main content: `src/views/admin/custom-objects/external-auth-tokens.vue`; migrated to `src/views/admin/custom-objects/external-auth-tokens-page.vue` with list, search, pagination, expand panel and linked-rule delete protection.
- `[x]` Modal windows: create/edit external auth token and confirmation phrase delete modals are migrated.
- `[x]` Form fields: request URL, timeout, repeat request, skip verify, auth, headers, payload, token, token variables switch and token expiration are migrated. Supporting shared fields now include `c-duration-field`, `c-retry-field`, `c-payload-text-field`, `c-payload-textarea-field`, request pair fields and payload variable menus.
- `[x]` Remaining form parity: shared `c-form-general-patterns-tabs` / test-variables stack is migrated and template testing tab behavior is restored in `external-auth-token-form.vue`.
- `[x]` Sidebars/drawers: none in the migrated page.
- `[x]` State/API/socket: existing webhook-token-rule Pinia store/hook and template-vars endpoint are reused.
- `[x]` Routing/access: `/admin/external-auth-tokens`, external auth token access.
- `[x]` Messages/i18n: existing external auth token and modal messages are reused.

### `admin-custom-objects-external-data-tables`

- `[x]` Main content: `src/views/admin/custom-objects/external-data-tables.vue`; migrated to `src/views/admin/custom-objects/external-data-tables-page.vue` with list, pagination, linked-rule delete protection, removed-from-config warning and records expand panel.
- `[x]` Modal windows: create/edit table, create/edit/duplicate record and import records modals are migrated.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: external data table Pinia store now uses CRUD helper for table CRUD and exposes import, export and record actions through migrated hooks.
- `[x]` Routing/access: `/admin/external-data-tables`, external data table access.
- `[x]` Messages/i18n: existing external data and modal messages are reused.

### `admin-custom-objects-entity-infos-properties`

- `[x]` Main content: `src/views/admin/custom-objects/entity-infos-properties.vue` -> `src/views/admin/custom-objects/entity-infos-properties-page.vue`; list, type filter, search, pagination and row actions are migrated.
- `[x]` Modal windows: `create-entity-info-property` modal and form are migrated; `c-entity-info-property-key-field` and type field are migrated.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: entity info property Pinia store/hook already exists and is wired to page CRUD and bulk remove.
- `[x]` Routing/access: `/admin/entity-infos-properties`, entity info property access.
- `[x]` Messages/i18n: existing old keys are reused.

### `admin-custom-objects-icons`

- `[x]` Main content: `src/views/admin/custom-objects/icons.vue` -> `src/views/admin/custom-objects/icons-page.vue`; icons list, pagination, refresh, custom SVG preview and actions are migrated.
- `[x]` Modal windows: `create-icon` modal, icon form and SVG file selector are migrated.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: material icon font/list generator and custom icon registry are migrated; page uses icon Pinia store hook for multipart CRUD and socket refresh parity is covered.
- `[x]` Routing/access: `/admin/icons`, icon access.
- `[x]` Messages/i18n: existing old keys are reused.

### `admin-custom-objects-maps`

- `[x]` Main content: `src/views/admin/custom-objects/maps.vue` -> `src/views/admin/custom-objects/maps-page.vue`; table/search/pagination/actions, expand panel shell, geo/Mermaid/tree-of-dependencies/flowchart expand previews, common flowchart stack and cluster graph editor stack are migrated.
- `[x]` Modal windows: `create-map`, `create-geo-map`, `create-flowchart-map`, `create-mermaid-map`, `create-tree-of-dependencies-map` are migrated and registered with the Vue 3 modal plugin. Their forms/editors now use explicit `v-model:*`/`update:*` events instead of the old `v-field` directive and restore nested editor required-rule validation through the new validator APIs.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: existing Pinia CRUD map store/hook is used, including bulk remove and `fetchMapStateWithoutStore`; Mermaid, geo, tree and flowchart preview helpers, popup support and form/editor-specific map helpers are migrated.
- `[x]` Routing/access: `/admin/maps`, map access.
- `[x]` Messages/i18n: existing old map/geomap/modal keys are reused.

### `admin-custom-objects-tags`

- `[x]` Main content: `src/views/admin/custom-objects/tags.vue` -> `src/views/admin/custom-objects/tags-page.vue`; tags list, search/pagination, expand panel and delete/bulk actions are migrated.
- `[x]` Modal windows: `create-tag` modal and general/patterns tag form are migrated.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: tag Pinia CRUD store/hook are used for list/create/update/delete and bulk actions.
- `[x]` Routing/access: `/admin/tags`, tag access.
- `[x]` Messages/i18n: existing tag/pattern messages are wired.

### `admin-custom-objects-llms`

- `[x]` Main content: `src/views/admin/custom-objects/llms.vue` -> `src/views/admin/custom-objects/llms-page.vue`; LLM list, important notes banner, create/edit form and prompts history expand panel are migrated.
- `[x]` Modal windows: `create-llm` modal and form are migrated.
- `[x]` Sidebars/drawers: none required for the migrated page.
- `[x]` State/API/socket: LLM Pinia store/hook and form helper are migrated for page CRUD.
- `[x]` Routing/access: `/admin/llms`, LLM access.
- `[x]` Messages/i18n: existing LLM messages are wired.

### `admin-settings-views-import-export`

- `[x]` Main content: `src/views/admin/settings/views-import-export.vue` -> `src/views/admin/settings/views-import-export-page.vue`; page shell, select all, group/view selection panel, JSON export and JSON import file selection are migrated.
- `[x]` Modal windows: `import-groups-and-views` modal is migrated with draggable imported groups, imported views and result groups.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: existing view Pinia store/hook is used for fetching groups, export and import API calls; payload parity is checked.
- `[x]` Routing/access: `/admin/views-import-export`, view import/export `can` access.
- `[x]` Messages/i18n: existing import/export modal and page header messages are wired.

### `admin-settings-notifications`

- `[x]` Main content: `src/views/admin/settings/notifications-settings.vue` -> `src/views/admin/settings/notifications-settings-page.vue`; page shell, loading state, form and submit flow are migrated.
- `[x]` Modal windows: not used by this page.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: existing notification settings Pinia store/hook is used for fetch/update.
- `[x]` Routing/access: `/admin/notifications`, notification settings `can` access.
- `[x]` Messages/i18n: existing notification settings messages are wired.

### `admin-settings-comment-templates`

- `[x]` Main content: `src/views/admin/settings/comment-templates.vue` -> `src/views/admin/settings/comment-templates-page.vue`; page shell, list table, query pagination, create/edit/delete actions are migrated.
- `[x]` Modal windows: `create-comment-template` modal is migrated with name and fields form.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: existing comment template Pinia store/hook is used for list/create/update/remove.
- `[x]` Routing/access: `/admin/comment-templates`, comment template access.
- `[x]` Messages/i18n: existing comment template messages are wired.

### `admin-settings-widget-templates`

- `[x]` Main content: `src/views/admin/settings/widget-templates.vue` -> `src/views/admin/settings/widget-templates-page.vue`; page shell, list table, query pagination, create/edit/delete actions are migrated.
- `[x]` Modal windows: `select-widget-template-type` and `create-widget-template` modals are migrated.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: existing widget template Pinia store/hook is used for list/create/update/remove.
- `[x]` Routing/access: `/admin/widget-templates`, widget template access.
- `[x]` Messages/i18n: existing widget template messages are wired.
- `[x]` Widget settings integration: widget templates are consumed in all migrated settings that matched old `widgetTemplatesMixin` usage.

### `admin-settings-storage-settings`

- `[x]` Main content: `src/views/admin/settings/storage-settings.vue` -> `src/views/admin/settings/storage-settings-page.vue`; page wrapper, storage settings component, right-aligned history action buttons, right-aligned submit/cancel actions and all storage form blocks are migrated.
- `[x]` Modal windows: existing `confirmation` and `confirmation-phrase` modals are reused. Additional archive disabled/unlinked forms are migrated for confirmation phrase actions, with explicit vue-i18n 11 message mapping for title/text/phrase fields and old day/week/month/year archive duration units.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: existing data-storage Pinia store/hook and entity archive/clean hooks are used.
- `[x]` Routing/access: `/admin/storage-settings`, storage settings access.
- `[x]` Messages/i18n: existing storage-setting and confirmation phrase messages are wired.

### `admin-settings-state-settings`

- `[x]` Main content: `src/views/admin/settings/state-settings.vue` -> `src/views/admin/settings/state-settings-page.vue`; list table, pagination, create/edit/duplicate/remove actions are migrated.
- `[x]` Modal windows: `create-state-setting` and `create-junit-state-setting` are migrated and registered. `state-setting-inherited-entity-pattern` is registered and opened from the Context root-cause / state-settings summary **See filter pattern** action.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: state-setting CRUD store/hook are reused.
- `[x]` Routing/access: `/admin/state-settings`, state setting access.
- `[x]` Messages/i18n: reused existing state-setting and modal messages.

### `exploitation-pbehaviors`

- `[x]` Main content: `src/views/exploitation/pbehaviors.vue` -> `src/views/exploitation/pbehaviors-page.vue`; list table, advanced search, pagination, expand tabs (patterns, entities, comments, recurrence), delete/export action, bulk enable/disable/remove, alarm filtering button and create/edit/duplicate row actions are migrated.
- `[x]` Modal windows: `pbehaviorPlanning`, `create-pbehavior`, `create-recurrence-rule`, `pbehavior-recurrence-rule` and recurrent changes confirmation are migrated and registered. Recurrence rule form fields and recurrence calendar timespan display are migrated.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: pbehavior store/hooks are reused for list fetch, create/update/delete, bulk create/update/delete, DB export, recurrence timespans, bulk actions and alarm filtering socket room. Do not port old `entities/removeFromStore` for single pbehavior delete; use API delete plus list refresh.
- `[x]` Routing/access: `/exploitation/pbehaviors`, pbehavior access.
- `[x]` Messages/i18n: reused existing pbehavior messages.

### `exploitation-event-filters`

- `[x]` Main content: `src/views/exploitation/event-filters.vue` -> `src/views/exploitation/event-filters-page.vue`; list, search, pagination, create/edit/duplicate/remove/export actions, DB export button, expand tabs, summary, patterns, external data and failures tab are migrated.
- `[x]` Modal windows: `create-event-filter` is registered and adapted to the new `useForm` layer, wrappers, validation, confirmable close and AI chat compatibility provider. `event-filter-general-form`, `event-filter-enrichment-form` and `event-filter-change-entity-form` use `useFormModels`.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: event-filter Pinia store/hook, failures API, bulk actions, DB export facade and supporting external-data/pbehavior lookup facades are migrated for the page.
- `[x]` Routing/access: `/exploitation/event-filters`, event filter access.
- `[x]` Messages/i18n: existing event-filter messages are wired.

### `exploitation-snmp-rules`

- `[x]` Main content: `src/views/exploitation/snmp-rules.vue` -> `src/views/exploitation/snmp-rules-page.vue`; list, search, pagination, MIB upload button, row actions and mass enable/disable/remove are migrated.
- `[x]` Modal windows: `create-snmp-rule` modal and SNMP rule form tree are migrated with `useForm`; severity field, MIB object append action and shared add button translations are aligned.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: existing SNMP rule and SNMP MIB Pinia stores/hooks are used for list CRUD, bulk actions and MIB upload.
- `[x]` Routing/access: `/exploitation/snmp-rules`, SNMP rule access.
- `[x]` Messages/i18n: existing SNMP rule/modal messages are wired.

### `exploitation-dynamic-infos`

- `[x]` Main content: `src/views/exploitation/dynamic-infos.vue` -> `src/views/exploitation/dynamic-infos-page.vue`; tabs, dynamic infos list with `c-advanced-data-table`/advanced-search, Information/Patterns expand tabs with infos table and readonly patterns form, templates list and mass actions are migrated.
- `[x]` Modal windows: `create-dynamic-info` and `create-dynamic-info-template` are migrated for create/edit with `vee-validate` form; General/Patterns/Templates Studio tabs, infos item types and AI chat sidebar match the old modal; `add-dynamic-info-infos-from-template` is migrated for selecting names from templates.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: dynamic info Pinia CRUD store/hook, dynamic info templates hook and bulk/mass action UI are migrated.
- `[x]` Routing/access: `/exploitation/dynamic-infos`, dynamic info access.
- `[x]` Messages/i18n: existing `dynamic-info` and dynamic info template modal/form messages are wired.

### `exploitation-meta-alarm-rules`

- `[x]` Main content: `src/views/exploitation/meta-alarm-rules.vue` -> `src/views/exploitation/meta-alarm-rules-page.vue`; list, search, pagination, read-only pattern expand panel, DB export, delete/bulk actions and create/edit/duplicate actions are migrated.
- `[x]` Modal windows: `src/components/modals/meta-alarm-rule/create-meta-alarm-rule.vue` and related form stack are migrated.
- `[x]` Sidebars/drawers: not used by this page.
- `[x]` State/API/socket: Pinia CRUD store/hook and meta alarm rule form helper are migrated.
- `[x]` Routing/access: `/exploitation/meta-alarm-rule`, meta alarm rule access.
- `[x]` Messages/i18n: existing meta alarm rule messages are wired.

### `exploitation-scenarios`

- `[x]` Main content: `src/views/exploitation/scenarios.vue` -> `src/views/exploitation/scenarios-page.vue`; table, search, pagination, row expand, action general tab, page create/edit/duplicate entry points and row DB export button are migrated.
- `[x]` Modal windows: `create-scenario` modal is migrated for create/edit/duplicate, including general fields, scenario actions, webhook/pbehavior/action parameters, declare-ticket mapping, action patterns and template variables tab. Confirmation modal is reused for remove.
- `[x]` Sidebars/drawers: none required for the migrated page.
- `[x]` State/API/socket: scenario Pinia store and hook are migrated for list/create/update/remove/bulk enable/bulk disable/bulk remove/test execution API wrappers; test execution UI is migrated through `scenario-test-query`.
- `[x]` Routing/access: `/exploitation/scenarios`, scenario access.
- `[x]` Messages/i18n: existing `scenario` and `modals.createScenario` messages are wired.

### `exploitation-idle-rules`

- `[x]` Main content: `src/views/exploitation/idle-rules.vue` -> `src/views/exploitation/idle-rules-page.vue`; list, search, pagination, summary/patterns expand, delete/bulk actions and create/edit/duplicate entry points are migrated.
- `[x]` Modal windows: `create-idle-rule` form is migrated, including pbehavior action parameters.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: idle rules Pinia store/hook are migrated for list/create/update/delete/bulk actions.
- `[x]` Routing/access: `/exploitation/idle-rules`, idle rules access.
- `[x]` Messages/i18n: existing `idle-rules` messages are wired.

### `exploitation-flapping-rules`

- `[x]` Main content: `src/views/exploitation/flapping-rules.vue` -> `src/views/exploitation/flapping-rules-page.vue`; shared alarm-status rules list, search, pagination, expand, delete/bulk actions and create/edit/duplicate entry points are migrated.
- `[x]` Modal windows: `create-alarm-status-rule` modal and form are migrated with general/patterns tabs.
- `[x]` Sidebars/drawers: none required for the migrated page.
- `[x]` State/API/socket: flapping rules Pinia store/hook are migrated for list/delete/bulk actions.
- `[x]` Routing/access: `/exploitation/flapping-rules`, flapping rules access.
- `[x]` Messages/i18n: existing `flapping-rule` and page header messages are wired.

### `exploitation-resolve-rules`

- `[x]` Main content: `src/views/exploitation/resolve-rules.vue` -> `src/views/exploitation/resolve-rules-page.vue`; shared alarm-status rules list, search, pagination, expand, delete/bulk actions and create/edit/duplicate entry points are migrated.
- `[x]` Modal windows: `create-alarm-status-rule` modal and form are migrated with general/patterns tabs.
- `[x]` Sidebars/drawers: none required for the migrated page.
- `[x]` State/API/socket: resolve rules Pinia store/hook are migrated for list/delete/bulk actions.
- `[x]` Routing/access: `/exploitation/resolve-rules`, resolve rules access.
- `[x]` Messages/i18n: existing `resolve-rule` and page header messages are wired.

### `exploitation-declare-ticket-rules`

- `[x]` Main content: `src/views/exploitation/declare-ticket-rules-page.vue`, `src/components/other/declare-ticket/declare-ticket-rules-list.vue`.
- `[x]` Modal windows: `create-declare-ticket-rule`, `create-declare-ticket-event`, `create-associate-ticket-event`, `execute-declare-tickets`.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: rule store/hook, bulk actions, test execution, execution sockets and live alarm-widget execution parity are wired and checked.
- `[x]` Routing/access: `/exploitation/declare-ticket-rules`, declare ticket rule access.
- `[x]` Messages/i18n: existing declare-ticket/modal/alarm messages are wired.

### `exploitation-link-rules`

- `[x]` Main content: `src/views/exploitation/link-rules-page.vue`, `src/components/other/link-rule/link-rules-list.vue`.
- `[x]` Modal windows: `src/components/modals/link-rule/create-link-rule.vue`.
- `[x]` Sidebars/drawers: none required for the migrated page.
- `[x]` State/API/socket: Pinia store/hook already available in `src/stores/modules/link-rule.js` and `src/hooks/store/modules/link-rule.js`.
- `[x]` Routing/access: `/exploitation/link-rules`, link rule access.
- `[x]` Messages/i18n: link rule and modal messages already present.

### `exploitation-remediation`

- `[x]` Main content: `src/views/exploitation/remediation.vue` -> `src/views/exploitation/remediation-page.vue`; remediation configurations, instructions, jobs, statistics and instruction stats tabs are migrated.
- `[x]` Modal windows: remediation configuration, instruction/job, execution and instruction approval modals are migrated.
- `[x]` Sidebars/drawers: none required for the migrated page.
- `[x]` State/API/socket: remediation Pinia stores/hooks are migrated for page list/actions, statistics, instruction stats, execution flows and approval flow; backend bulk route for configurations remains a backend/API TODO outside the UI migration.
- `[x]` Routing/access: `/exploitation/remediation`, grouped remediation permission.
- `[x]` Messages/i18n: existing remediation messages are reused.

### `playlist`

- `[x]` Main content: `src/views/playlist.vue` -> `src/views/playlist-page.vue`; renders the selected playlist tab through the migrated `view-tab-widgets` stack and shows the central play overlay before playback starts.
- `[x]` Top-bar controls: timer, previous, play/pause, next and fullscreen controls are teleported into the top bar.
- `[x]` Modal windows: none in the old page.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: playlist fetch, per-tab fetch, timer rotation, autoplay and view-level periodic refresh provider are migrated.
- `[x]` Routing/access: `/playlist/:id`, playlist read permission by route id.
- `[x]` Messages/i18n: existing common and playlist messages are reused.

### `profile-patterns`

- `[x]` Main content: `src/views/profile/patterns-page.vue`, `src/components/other/pattern/patterns.vue`, `corporate-patterns.vue`, `patterns-list.vue`.
- `[x]` Modal windows: `create-pattern` modal is migrated.
- `[x]` Sidebars/drawers: none required for the migrated page.
- `[x]` State/API/socket: pattern and corporate-pattern Pinia stores/hooks are wired.
- `[x]` Routing/access: `/filters`, authenticated access plus corporate pattern CRUD permission for the corporate tab/create action.
- `[x]` Messages/i18n: existing pattern/modal messages are reused.

### `profile-themes`

- `[x]` Main content: `src/views/profile/themes.vue` -> `src/views/profile/themes-page.vue`; themes list, search, pagination, row actions and bulk delete are migrated.
- `[x]` Modal windows: `create-theme` modal, theme form, compact color picker menu and preview components are migrated.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: existing theme Pinia CRUD store/hook and base Vuetify theme runtime are used; active user theme refresh is wired.
- `[x]` Routing/access: `/themes`, profile theme permission.
- `[x]` Messages/i18n: existing theme/modal/common messages are reused.

### `notifications`

- `[x]` Main content: `src/views/notification/notifications.vue` -> `src/views/notification/notifications-page.vue`; route page, tab badges, refresh FAB, approve/rate/event-filter-failure tabs and remediation stats alarms timeline expand tab are migrated.
- `[x]` Modal windows: `rate` modal/form is migrated for instruction rating. Existing remediation instruction approval/edit and event-filter modals are reused from their migrated slices.
- `[x]` Sidebars/drawers: none in the old page.
- `[x]` State/API/socket: existing notification socket room and remediation instruction/stats/event-filter Pinia hooks are wired.
- `[x]` Routing/access: `/notifications/:tabId`, authenticated access.
- `[x]` Messages/i18n: existing notifications/remediation/modal/common messages are reused.
- `[x]` Remediation instruction stats alarms timeline expand tab is migrated, including executions list, horizontal alarm timeline, failed-executions switch and per-execution steps expand panel.

### `error`

- `[x]` Main content: `src/views/error.vue` -> `src/views/error-page.vue`.
- `[x]` Modal windows: none expected.
- `[x]` Sidebars/drawers: none expected.
- `[x]` State/API/socket: app-info polling and redirect logic migrated.
- `[x]` Routing/access: `/error`, `hideHeader`, query props.
- `[x]` Messages/i18n: error page uses existing app/common messages.

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
- `[x]` Permission waiting: current implementation fetches current user when needed. App-info loading sets popup default close times from `popup_timeout` both from app bootstrap and from route access guard fallback.
- `[x]` Permission algorithm: target uses `checkUserAnyAccess`, grouped permissions and CRUD action behavior.
- `[x]` Guard failure redirect: source redirects failed `beforeResolve` checks to `home`.

### Post-navigation Rules

- `[x]` Entity sweep: old `entities/sweep` belonged to the Vuex/normalizr embedded-entity mark-and-sweep plugin and is intentionally not ported. Pinia stores keep explicit local list/detail state and refresh migrated consumers directly after mutations.
- `[x]` Initial error popup: on router ready, if initial route has `errorMessage`, show persistent error popup.
- `[x]` Scroll/restoration behavior: source router did not define `scrollBehavior`; target intentionally keeps Vue Router default behavior and page-level scroll helpers handle local scroll cases.

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

| Status | Widget type         | Source component name        | Source file                                                                 | Icon                     | Edition / scope                                                                                                                                            | Notes                                                                                                                                                                                                                 |
| ------ | ------------------- | ---------------------------- | --------------------------------------------------------------------------- | ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `[x]`  | `AlarmsList`        | `alarms-list-widget`         | `src/components/widgets/alarm/alarms-list.vue`                              | `view_list`              | Core, top-level                                                                                                                                            | Step 0 audit, complete settings sidebar scope, default widget configuration, shared widget query layer, toolbar, columns/cells, pagination, base row/mass actions, action modals, alarm-list modal handoff, PDF export and column drag/resize/reset are implemented. Table structure now follows the legacy split: `alarms-list-table-with-pagination` -> `alarms-list-table` -> `alarms-list-row`. Row actions restore the legacy shared panel `closeMenu()` contract, dropdown close behavior and action-panel async booting. Expand row state, details fetch/store path, expand button, More infos, Timeline, Charts, Tickets declared, Pbehaviors, Children, Track source, Impact chain, Entity gantt, Remediation, Availability and Entity enrichments tabs are migrated. Children tab now reuses the table chain, keeps parent alarm context for child rows and only appears when the current opened/resolved filter has matching consequences. Periodic refresh, live watching, alarm/details socket rooms, delayed row booting, virtual scroll, empty-column UI, Ctrl-drag selection, backend column filters, entity info property aliases, sticky header, sticky horizontal scroll and toolbar are migrated. `c-advanced-data-table` is intentionally excluded. |
| `[x]`  | `Context`           | `entities-list-widget`       | `src/components/widgets/context/entities-list.vue`                          | `view_list`              | Core, top-level                                                                                                                                            | Plan: `docs/ai/context-widget-migration-plan.md`. Settings sidebar, form mappers, widget-wrapper registration, widget-scoped entity store, query converters, list fetch, `c-advanced-data-table` shell, entity column cells, create/edit entity/service/info modals, row/mass actions, create-service FAB, expand panel, toolbar/search/category/filters/idle-entities (`c-enabled-field`), CSV export, pattern entities-list modal and root-cause diagram are complete. |
| `[x]`  | `ServiceWeather`    | `service-weather-widget`     | `src/components/widgets/service-weather/service-weather.vue`                | `view_module`            | Core, top-level                                                                                                                                            | Phases 0–5 are complete: settings sidebar, form helpers, weather-service Pinia store, permission mixins, query converters, widget runtime (tile grid + toolbar + counters), widget-wrapper registration, alarm-list/root-cause/variables-help handoffs and service-entities modal stack. See `docs/ai/service-weather-widget-migration-plan.md`. |
| `[x]`  | `StatsCalendar`     | `stats-calendar-widget`      | `src/components/widgets/stats/calendar/stats-calendar.vue`                  | `calendar_today`         | Pro, top-level                                                                                                                                             | Plan: `docs/ai/stats-calendar-widget-migration-plan.md`. Runtime calendar, grouped event menu, form/query helpers, settings sidebar, alarm-list drilldown and widget-wrapper/menu integration are complete. |
| `[x]`  | `Text`              | `text-widget`                | `src/components/widgets/text/text.vue`                                      | `view_headline`          | Core, top-level                                                                                                                                            | Renderer and settings sidebar are migrated.                                                                                                                                                                           |
| `[x]`  | `Counter`           | `counter-widget`             | `src/components/widgets/counter/counter.vue`                                | `numbers`                | Renderer, settings sidebar, embedded alarm-list display settings, filter editor, variables-help action and alarm-list modal handoff are migrated. |                                                                                                                                                                                                                       |
| `[x]`  | `Junit`             | `testing-weather-widget`     | `src/components/widgets/testing-weather/testing-weather.vue`                | `view_module`            | Pro, top-level                                                                                                                                             | Plan: `docs/ai/junit-widget-migration-plan.md`. Runtime cards, settings sidebar, test-suite modal/details stack and widget-wrapper/menu integration are complete.                           |
| `[x]`  | `Map`               | `map-widget`                 | `src/components/widgets/map/map.vue`                                        | `location_on`            | Pro, top-level                                                                                                                                             | Renderer, settings sidebar, map selector, category filter, filter selector/list/create modal, breadcrumbs, map-state fetch, columns and alarm-list modal handoff are migrated.                  |
| `[x]`  | `BarChart`          | `bar-chart-widget`           | `src/components/widgets/chart/bar-chart-widget.vue`                         | `stacked_bar_chart`      | KPI/chart child widget                                                                                                                                     | Plan: `docs/ai/bar-chart-widget-migration-plan.md`. Runtime chart, KPI toolbar filters, metrics fetch/export, form/query helpers, settings sidebar and widget-wrapper/menu integration are complete. |
| `[x]`  | `LineChart`         | `line-chart-widget`          | `src/components/widgets/chart/line-chart-widget.vue`                        | `show_chart`             | KPI/chart child widget                                                                                                                                     | Plan: `docs/ai/line-chart-widget-migration-plan.md`. Runtime chart, KPI toolbar filters, metrics fetch/export, form/query helpers, settings sidebar and widget-wrapper/menu integration are complete. |
| `[x]`  | `PieChart`          | `pie-chart-widget`           | `src/components/widgets/chart/pie-chart-widget.vue`                         | `donut_small`            | KPI/chart child widget                                                                                                                                     | Plan: `docs/ai/pie-chart-widget-migration-plan.md`. Runtime chart, KPI toolbar filters, metrics fetch/export, form/query helpers, settings sidebar and widget-wrapper/menu integration are complete. |
| `[x]`  | `Numbers`           | `numbers-widget`             | `src/components/widgets/chart/numbers-widget.vue`                           | `functions`              | KPI/chart child widget                                                                                                                                     | Plan: `docs/ai/numbers-widget-migration-plan.md`. Runtime metrics, KPI toolbar filters, metrics fetch/export, form/query helpers, settings sidebar and widget-wrapper/menu integration are complete. |
| `[x]`  | `UserStatistics`    | `user-statistics-widget`     | `src/components/widgets/statistics/user-statistics-widget.vue`              | `perm_identity`          | Pro, statistics child widget                                                                                                                               | Plan: `docs/ai/user-statistics-widget-migration-plan.md`. Runtime table, settings sidebar, form/query helpers, KPI toolbar filters, group metrics fetch and widget-wrapper/menu integration are complete. AlarmStatistics reuses the shared statistics stack. |
| `[x]`  | `AlarmStatistics`   | `alarm-statistics-widget`    | `src/components/widgets/statistics/alarm-statistics-widget.vue`             | `notification_important` | Pro, statistics child widget                                                                                                                               | Plan: `docs/ai/alarm-statistics-widget-migration-plan.md`. Runtime table, settings sidebar, shared form/query helpers, KPI toolbar filters, group metrics fetch and widget-wrapper/menu integration are complete. |
| `[x]`  | `Availability`      | `availability-widget`        | `src/components/widgets/availability/availability-widget.vue`               | `published_with_changes` | Pro, top-level                                                                                                                                             | Plan: `docs/ai/availability-widget-migration-plan.md`. Runtime table, settings sidebar, form/query helpers, toolbar filters, history graph and widget-wrapper/menu integration are complete. |
| `[x]`  | `ExternalDataTable` | `external-data-table-widget` | `src/components/widgets/external-data-table/external-data-table-widget.vue` | `database`               | Core, top-level                                                                                                                                            | Widget runtime, settings sidebar, form/query helpers, widget-wrapper registration, user preference, periodic refresh, CSV export and records-list handoff are complete. Reuses admin-layer store, modals and `external-data-table-records-list`. |

## Widget Blocks

Each widget is tracked by main rendering logic and settings/sidebar logic. Source setting sidebars are mapped from `SIDE_BARS_BY_WIDGET_TYPES` in `src/constants/widget.js`.

### `AlarmsList`

- `[x]` Main content: `src/components/widgets/alarm/alarms-list.vue`; runtime foundation is registered in `widget-wrapper.vue`, renders prepared columns, uses shared widget query state through the Pinia `query` store, includes the migrated toolbar controls, and delegates the table to the legacy-style component chain `partials/alarms-list-table-with-pagination.vue`, `partials/alarms-list-table.vue` and `partials/alarms-list-row.vue`. Base row rendering, mass action panels, column drag/resize/reset, sticky header, sticky top toolbar, sticky horizontal scrollbar and pagination live in that table chain. Alarm event, ticket, comment, snooze, change-state, history, meta-alarm, remediation execute, remediation execution popups and alarm-list handoff modals are registered. PDF export is wired. Row actions now keep the shared action panel close contract used by legacy `afterSubmit()` and use the migrated `c-booted-placeholder-loader`/`AsyncBooting` path. Expand panel now has row state, details fetching through Pinia, expand button, all legacy tabs, local details query lifecycle, Pbehaviors calendar modal, Children tab backed by `alarms-list-table-with-pagination`, remediation executions, service dependency tree, availability, Entity gantt and enrichments data loading. Children/consequences rows keep the parent alarm context for templates/actions/instruction icons, and tab visibility follows `opened_children`/`closed_children` for the current widget mode. Periodic refresh, live watching, alarm/details socket rooms, pause-on-expand, delayed row booting, virtual scroll, empty-column UI, Ctrl-drag selection, backend column filters, entity info property aliases, sticky header, sticky horizontal scroll and toolbar are wired.
- `[x]` Expand panel tab switching is complete for migration scope. The current non-animated switch is intentional because the Vuetify `v-window` horizontal transition leaves the previous tab visually hanging inside the table expand row.
- `[x]` Entity gantt in the alarms list table expand panel is migrated and checked: intervals, pagination, tooltips and layout inside the expanded row.
- `[x]` Settings: `alarm-settings`; widget parameters form, columns, filters, actions, quick actions, fast pbehavior, Export PDF, Export CSV, Charts and expand panel settings are migrated and checked.
- `[x]` Store/API/socket: alarm list query/pagination/user-preference foundation is migrated through the Pinia alarm/query stores. Periodic refresh, live watching, alarm/alarm-details socket rooms and opened-details refresh on list fetch are migrated. When migrating the alarm expand pbehavior tab, do not rely on old embedded entity cleanup; `pbehaviors-simple-list` should refresh its own list after pbehavior delete.
- `[x]` Messages/i18n: existing `alarm` and settings keys are wired.

### `Context`

- `[x]` Main content: `src/components/widgets/context/entities-list.vue`. Runtime foundation is registered in `widget-wrapper.vue` with `prepareContextWidget`. Query state, user preference, category/filter/no-events updates and list fetch (`with_flags: true` when widget columns exist) are wired. Table shell is `partials/entities-list-table-with-pagination.vue` on `c-advanced-data-table`: empty-column UI, top/bottom pagination, select-all, expand, mass-actions grid teleport, `keepSelectedAfterAction` and actions column unless `hideActions`. Entity column cells are migrated (`entity-column-cell` with Vue 3 `<component :is>`, state/pbehavior/event-statistics, date filters, links categories, root-cause click). Create/edit entity, entity-info and service modals and their form stacks are registered. Row actions (edit entity/service, duplicate service, delete, pbehavior, variables help), mass actions (delete with partial-deletable warning, pbehavior, mass enable/disable) and the create-service FAB are migrated; the list refreshes through `fetchContextEntitiesListWithPreviousParams` after each successful action. Expand panel tabs are migrated (availability, pbehaviors with local refresh after delete, impact/depends, infos, charts, tree of dependencies, impact chain, active/resolved alarms, comments). Toolbar advanced search, category, user filters, idle-entities `c-enabled-field` (same compact inline switch as AlarmsList correlation), CSV export, pattern entities-list modal and root-cause diagram modal are migrated. Root-cause follow-ups are done: Vue 3 node fill/icons, layout `animate: false` + `cy.stop()` so the second layout is not blocked by the first fit animation, fit/pan checked, inherited-state summary **See filter pattern** via `MODALS.stateSettingInheritedEntityPattern`.
- `[x]` Settings: `context-settings`; title, default sort, widget/service-dependency/active/resolved/export columns, tree-of-dependencies show type, availability, root-cause, filters, entity types, expand grid range, inherited pbehavior, keep-selected and charts are migrated and checked.
- `[x]` Store/API/socket: Pinia entity store now has widget-scoped list state, request-cancelling `fetchList` / `fetchListWithPreviousParams`, entityBasics CRUD, context-graph, export create/fetch, bulk enable/disable and `fetchStateSettingWithoutStore`. Pinia service store now has create/update/remove/`fetchItemWithoutStore` plus existing dependency/impact/bulk-remove. Query converters are registered for `WIDGET_TYPES.context`. Expand pbehavior tab refreshes its local list after delete instead of old `entities/removeFromStore`.
- `[x]` Messages/i18n: `src/i18n/messages/{en,fr}/context.js` and `state-setting.js` (`seeFilterPattern`, `stateIsInheritFrom`) are wired.

### `ServiceWeather`

- `[x]` Main content: `src/components/widgets/service-weather/service-weather.vue` and `service-weather-item.vue` are migrated and registered in `widget-wrapper.vue` with `prepareServiceWeatherWidget`. Runtime includes category/filter/hide-grey toolbar, `widgetFetchQueryMixin` + `widgetPeriodicRefreshMixin`, weather-service Pinia list fetch, responsive column grid, state/pbehavior counter columns, impact indicator, modal-type routing (more infos / alarm list / both), and tile layout fixes. Partials: `alarm-counter`, `alarm-state-counters`, `alarm-pbehavior-counters`, `impact-state-indicator`, `service-pause-event-form`.
- `[x]` Settings: `service-weather-settings`; service display, counters, pause/event settings, linked alarm list settings, templates, grid, margins, modal type and action-required settings are migrated.
- `[x]` Store/API: weather-service widget list Pinia store is wired to `/api/v4/weather-services`. Service-entity store facade supports the service-entities modal. Query converters registered in `helpers/entities/widget/query.js`. No dedicated socket layer (parity with legacy).
- `[x]` Messages/i18n: existing `service-weather.js` keys are wired; no new keys added in this step.
- `[x]` Modals: `MODALS.serviceEntities` (`service-entities.vue` + partials under `other/service/partials/`), `create-service-pause-event.vue`, alarm-list and root-cause handoffs from widget runtime.

### `StatsCalendar`

- `[x]` Main content: `src/components/widgets/stats/calendar/stats-calendar.vue`; runtime uses `c-calendar`, fetches alarms for the visible range through the alarm Pinia `fetchListWithoutStore`, groups alarms by day/hour, renders grouped month events and opens the migrated alarm-list modal for drilldown.
- `[x]` Settings: `stats-calendar-settings`; title, linked alarm-list modal settings, filters, opened/resolved, consider pbehaviors, criticity levels and colors are migrated with `useWidgetSettings`.
- `[x]` Store/API/socket: query converter is registered for StatsCalendar and drilldown uses the alarm-list modal handoff. No dedicated store/socket layer is required by the legacy widget.
- `[x]` Messages/i18n: reused existing settings, modal, error and common keys.

### `Text`

- `[x]` Main content: `src/components/widgets/text/text.vue`; runtime template renderer is migrated.
- `[x]` Settings: `text-settings`; template editor sidebar is migrated with `useWidgetSettings`.
- `[x]` Store/API/socket: text widget does not require a dedicated runtime store; broader widget refresh parity is covered by the view stack.
- `[x]` Messages/i18n: existing settings keys are wired.

### `Counter`

- `[x]` Main content: `src/components/widgets/counter/counter.vue`; runtime renderer, state color/icon logic, block template context, counter store fetch and variables-help action are migrated.
- `[x]` Settings: `counter-settings`; levels, margins, columns, height, block template, opened/resolved selector, filter editor, embedded alarm-list display settings and info popup create/edit modals are migrated.
- `[x]` Store/API/socket: counter queries and alarm-list handoff modal are migrated.
- `[x]` Messages/i18n: existing settings keys are wired.

### `Junit`

- Working plan: `docs/ai/junit-widget-migration-plan.md`.
- `[x]` Main content: `src/components/widgets/testing-weather/testing-weather.vue` and `testing-weather-item.vue` are migrated and registered in `widget-wrapper.vue`. Runtime includes periodic refresh, widget-scoped test-suite list fetch, status-colored cards, mini chart, statistics and alarm-list handoff.
- `[x]` Settings: `testing-weather-settings`; title, periodic refresh, result directory, receive-by-API, screenshot/video directories and filename masks are migrated with `useWidgetSettings`.
- `[x]` Store/API/socket: Pinia `test-suite` store is reused for widget list, directory validation, summary, gantt, details and history. No dedicated socket layer (parity with legacy).
- `[x]` Modals: `MODALS.testSuite` (`test-suite.vue` + history/details stack under `other/test-suite/`) and alarm-list handoff from widget runtime.
- `[x]` Messages/i18n: existing `test-suite` and settings keys are wired.

### `Map`

- `[x]` Main content: `src/components/widgets/map/map.vue`; renderer, map-state fetch, category selector, filter selector/list/create modal, nested-map navigation breadcrumbs, preview component dispatch for geo/flowchart/mermaid/tree maps and alarm-list modal handoff are migrated.
- `[x]` Settings: `map-settings`; map selector, periodic refresh, filter editor, color indicator, pbehavior marker display toggle, entity info popup editor and alarm/entity columns are migrated.
- `[x]` Store/API/socket: map queries use the migrated Pinia map store and service-alarm route helper is restored; periodic refresh is local to the widget.
- `[x]` Messages/i18n: existing `map` and settings keys are wired.

### `BarChart`

- `[x]` Main content: `src/components/widgets/chart/bar-chart-widget.vue` is migrated and registered in `widget-wrapper.vue`. Runtime preserves Options API, uses `widgetFetchQueryMixin`, KPI toolbar interval/sampling/filter controls, vector metrics from the Pinia metrics store, migrated `bar-chart-metrics`, metric color/label presets and PNG/CSV export actions.
- `[x]` Settings: `bar-chart-settings`; chart metrics, sampling, time range, filters and display options are migrated through `bar-chart-widget-form.vue`.
- `[x]` Store/API/socket: KPI vector metric queries and alarm metric export use the existing metrics Pinia store.
- `[x]` Messages/i18n: existing BarChart settings, permissions and modal messages are available.

### `LineChart`

- `[x]` Main content: `src/components/widgets/chart/line-chart-widget.vue` is migrated and registered in `widget-wrapper.vue`. Runtime preserves Options API, uses `widgetFetchQueryMixin`, KPI toolbar interval/sampling/filter controls, vector metrics from the Pinia metrics store, migrated `line-chart-metrics`, metric color/label presets and PNG/CSV export actions.
- `[x]` Settings: `line-chart-settings`; chart metrics, sampling, time range, filters and display options are migrated through `line-chart-widget-form.vue`.
- `[x]` Store/API/socket: KPI vector metric queries and alarm metric export use the existing metrics Pinia store.
- `[x]` Messages/i18n: existing LineChart settings, permissions and modal messages are available.

### `PieChart`

- `[x]` Main content: `src/components/widgets/chart/pie-chart-widget.vue`.
- `[x]` Settings: `pie-chart-settings`; migrate chart metrics, show mode, time range, filters and display options.
- `[x]` Store/API/socket: KPI metric queries.
- `[x]` Messages/i18n: existing PieChart settings, permissions and modal messages are available.

### `Numbers`

- `[x]` Main content: `src/components/widgets/chart/numbers-widget.vue`.
- `[x]` Settings: `numbers-settings`; migrate number metrics, time range, filters and display options.
- `[x]` Store/API/socket: KPI metric queries.
- `[x]` Messages/i18n: existing Numbers settings, permissions and modal messages are available.

### `UserStatistics`

- `[x]` Migration plan: staged plan is tracked in `docs/ai/user-statistics-widget-migration-plan.md`. This widget owns the shared statistics stack used by AlarmStatistics.
- `[x]` Main content: `src/components/widgets/statistics/user-statistics-widget.vue` and shared `partials/statistics-widget.vue`.
- `[x]` Settings: `user-statistics-settings`; title, periodic refresh, user main parameter, columns, table header, default time range and filters.
- `[x]` Store/API/socket: group metrics by widget id and rating-settings fetch without store.
- `[x]` Messages/i18n: existing `settings`, `kpi.statisticsWidgets`, `user.metrics` and `alarm.metrics` keys are wired with the runtime.

### `AlarmStatistics`

- `[x]` Migration plan: staged plan is tracked in `docs/ai/alarm-statistics-widget-migration-plan.md`. Reuses the shared statistics stack from the UserStatistics plan.
- `[x]` Main content: `src/components/widgets/statistics/alarm-statistics-widget.vue`.
- `[x]` Settings: `alarm-statistics-settings`; entity main parameter including custom criteria, entity metric columns, table header, default time range and filters.
- `[x]` Store/API/socket: same group metrics / rating-settings path as UserStatistics, with `KPI_RATING_SETTINGS_TYPES.entity`.
- `[x]` Messages/i18n: existing settings/kpi/metric keys are wired with the runtime.

### `Availability`

- `[x]` Migration plan: staged plan is tracked in `docs/ai/availability-widget-migration-plan.md`.
- `[x]` Main content: `src/components/widgets/availability/availability-widget.vue` with toolbar, list, expand graph and alarm tabs.
- `[x]` Settings: `availability-settings`; columns, interval, display parameter, show type, filters and export CSV.
- `[x]` Store/API/socket: availability list/history/export through the existing metrics Pinia store.
- `[x]` Messages/i18n: existing `availability`, settings and permission keys are wired with the runtime.

### `ExternalDataTable`

- `[x]` Main content: `src/components/widgets/external-data-table/external-data-table-widget.vue`. Runtime orchestrates table fetch, user preference (dense, items per page, column drag/resize settings), periodic refresh, CSV export polling and delegates rendering/actions to migrated `external-data-table-records-list`. Column tag edits at runtime update the table definition via `updateExternalDataTable`.
- `[x]` Settings: `external-data-table-settings`; title, periodic refresh, table selector, column names, resize/drag behavior, default sort, items per page, density and CSV export settings are migrated through `useWidgetSettings()`.
- `[x]` Store/API/socket: reuses consolidated Pinia `external-data-table` store for table fetch, record CRUD, export create/status and column-tag updates. Widget query converters registered for pagination/sort defaults.
- `[x]` Messages/i18n: reuses `external-data.js`, `settings.js` and modal strings already present in target.

## Dynamic Feature Extensions

The source project allows additional routes/widgets through `featuresService.get(...)`. The checked-in community source has `src/features/.gitignore`, so no extra feature files are present in this source tree.

Relevant extension points:

- `constants.WIDGET_TYPES`
- `constants.COMPONENTS_BY_WIDGET_TYPES`
- `constants.WIDGET_ICONS`
- `constants.SIDE_BARS`
- `constants.SIDE_BARS_BY_WIDGET_TYPES`
- `constants.WIDGET_TYPES_RULES`
