# AlarmsList Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/alarm/alarms-list.vue`

Target runtime: `src/components/widgets/alarm/alarms-list.vue`

The widget is part of the `/view/:id` widget system. This plan covers the
complete widget flow: rendering, query state, sockets, table behavior, row
actions, mass actions, expand panels, settings, modals, exports and
verification.

`c-advanced-data-table` is intentionally out of scope for this migration. The
AlarmList table will be implemented as a dedicated Vuetify 4 table component,
preserving the old AlarmList-specific behavior without waiting for the shared
table migration.

The settings sidebar must be migrated completely. Every setting from the old
sidebar and its full dependency slice must be ported; TODOs, stubs, disabled
controls and visual placeholders are not acceptable substitutes.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` First implementation scope: complete widget settings sidebar.
- `[x]` Next: widget foundation and default configuration.
- `[x]` Next: query state, pagination and alarm store integration.
- `[x]` Periodic refresh and socket rooms.
- `[x]` Dedicated AlarmList table shell: dedicated wrapper, top/bottom
  pagination, empty-column UI, selection, Ctrl-drag and row behavior.
- `[x]` Columns, headers and cell formatting.
- `[x]` Next: alarm rows and selection.
- `[x]` Next: single-alarm actions.
- `[x]` Next: mass actions.
- `[x]` Next: expand panel and details.
- `[x]` Toolbar and filters runtime scope.
- `[x]` Next: modals and exports.
- `[x]` Final: view integration and Ready For Testing update.

## Step 0: Inventory And Dependency Audit

### Source runtime

- `widgets/alarm/alarms-list.vue`
- `widgets/alarm/partials/alarms-list-table.vue`
- `widgets/alarm/partials/alarms-list-table-with-pagination.vue`
- `widgets/alarm/partials/alarms-list-row.vue`
- `widgets/alarm/partials/alarms-list-row-bookmark-icon.vue`
- `widgets/alarm/partials/alarms-list-row-instructions-icon.vue`
- `widgets/alarm/partials/alarms-list-remediation-instructions-filter.vue`
- `widgets/alarm/partials/alarms-list-remediation-instructions-filter-fields.vue`
- `widgets/alarm/partials/hooks/alarms-list-remediation-instructions-filter.js`

### Column and cell rendering

- `headers-formatting/alarm-header-cell.vue`
- `headers-formatting/alarm-header-priority.vue`
- `columns-formatting/alarm-column-cell.vue`
- `columns-formatting/alarm-column-cell-popup-body.vue`
- `columns-formatting/alarm-column-value.vue`
- `columns-formatting/alarm-column-value-status.vue`
- `columns-formatting/alarm-column-value-extra-details.vue`
- all components under `columns-formatting/extra-details/`
- `hooks/extra-details-tooltips.js`
- `partials/alarm-status-chip-with-relations.vue`

### Actions and expanded content

- `actions/actions-panel.vue`
- `actions/mass-actions-panel.vue`
- `expand-panel/alarms-expand-panel.vue`
- `expand-panel/alarms-expand-panel-btn.vue`
- `expand-panel/alarms-expand-panel-tab-item-wrapper.vue`
- `expand-panel/alarms-expand-panel-more-infos.vue`
- `expand-panel/alarms-expand-panel-children.vue`
- `expand-panel/alarms-expand-panel-remediation.vue`
- `expand-panel/alarms-expand-panel-entity-enrichments.vue`
- `entity-gantt/entity-gantt.vue`
- `timeline/` and all timeline components
- alarm action forms under `forms/`

### Settings and supporting modals

- `sidebars/alarm/alarm.vue`
- `sidebars/alarm/form/alarms-list-modal.vue`
- all fields under `sidebars/alarm/form/fields/`
- `sidebars/alarm/form/fast-pbehavior-form.vue`
- `modals/alarm/alarms-list.vue`
- alarm event, ticket, pbehavior, meta-alarm, chart and live-reporting modals
- `modals/alarm/info-popup-setting/`

### Cross-cutting source dependencies

- alarm entity store and alarm details store;
- alarm query, URL, list, form, formatting and icon helpers;
- widget alarm form mapper and column/sort/template mappers;
- widget fetch-query, periodic-refresh, filter-select and export mixins;
- advanced-search saved items;
- alarm socket mixin;
- alarm Handlebars precompiled templates;
- alarm columns, rows, sticky, resizing and dragging mixins;
- alarm action-panel and PDF-export mixins;
- alarm and remediation permissions;
- entity info properties and alarm column filters;
- pbehavior, service dependency, declared ticket and remediation components;
- alarm constants, widget constants, API routes, feature-service extension points;
- alarm/widget settings translations and export templates.

### Existing target coverage

Already present in the target project:

- Pinia alarm store and `useAlarm()` facade;
- widget alarm form mapper;
- partial alarm general list;
- a subset of alarm extra-details renderers;
- info-popup modal and form;
- partial AlarmList settings fields;
- widget constants containing `AlarmsList` and `alarmSettings`.

Still missing for the complete widget:

- runtime AlarmList component;
- dedicated AlarmList table and pagination wrapper;
- alarm row, headers and most cell formatters;
- single and mass action panels;
- expand panel and all its tabs;
- alarm socket/details query integration;
- AlarmList settings sidebar registration and full form;
- action and export modal registration/parity;
- complete alarm widget i18n and permission wiring;
- view widget-wrapper integration;
- Ready For Testing update.

## First Implementation Scope: Complete Settings Sidebar

The first implementation scope is the complete settings sidebar. It is
intentionally executed before the runtime table so that the final shape of
`widget.parameters`, widget templates, column mappings and settings defaults
is known before the renderer consumes them.

Use the complete checklist in `Complete Settings Sidebar Scope` below. The
sidebar is not split across later runtime stages and no settings group is
marked complete independently.

## Next Scope: Widget Foundation

- `[x]` Register the runtime component in `widget-wrapper.vue`.
- `[x]` Register and connect `alarmSettings` sidebar.
- `[x]` Enable Create, Edit, Duplicate and Delete entry points for the
  `AlarmsList` settings sidebar. Browser save/reopen verification remains in
  the final settings verification pass.
- `[x]` Port default widget parameters, default columns and default group
  columns through `generateDefaultAlarmListWidgetForm`,
  `generateDefaultAlarmListWidget` and `generatePreparedDefaultAlarmListWidget`.
- `[x]` Port `widgetToForm` and `formToWidget` mappings, including templates.
- `[x]` Preserve all old typedefs and JSDoc in migrated helpers touched by this
  scope.
- `[x]` Keep Options API components as Options API and Composition API components as
  Composition API. Do not introduce mixed API components.

Implemented files:

- `src/components/widgets/alarm/alarms-list.vue`
- `src/components/widgets/widget-wrapper.vue`
- `src/components/widgets/partials/widget-wrapper-menu.vue`
- `src/helpers/entities/widget/form.js`
- `src/helpers/entities/widget/forms/alarm.js`

Verification:

- `npx eslint src/helpers/entities/widget/form.js src/helpers/entities/widget/forms/alarm.js src/components/widgets/widget-wrapper.vue src/components/widgets/partials/widget-wrapper-menu.vue src/components/widgets/alarm/alarms-list.vue`

## Step 2: Query, Pagination And Store

- `[x]` Reuse widget-scoped alarm list state in the Pinia alarm store:
  `widgets[widgetId].allIds`, `meta`, `pending` and `fetchingParams`.
- `[x]` Connect the runtime widget to `fetchAlarmsList`, `getListByWidgetId`,
  `getMetaByWidgetId`, `getPendingByWidgetId` and
  `getFetchingParamsByWidgetId`.
- `[x]` Preserve page, items per page, sorting, search and total count at the
  query/store boundary.
- `[x]` Port AlarmList widget/user-preference query mapping into
  `helpers/entities/alarm/query.js`, including active columns, opened/resolved
  defaults, locked filter, correlation and pbehavior origin.
- `[x]` Port the shared widget query layer from the old Vuex `query` module to
  Pinia and widget mixins: `queryId`, `query` getter/setter, table `options`,
  `tabQueryNonce` and `prepareWidgetQuery`.
- `[x]` Remove the temporary local `data.query` state from `AlarmsList`; runtime
  query now follows the old shared widget-query flow.
- `[x]` Keep obsolete request cancellation and pending state per widget through
  the existing Pinia alarm store request-cancelling key.
- `[x]` Preserve toolbar updates for category, correlation, bookmarks, tags and
  selected filters in the runtime UI.
- `[x]` Restore advanced-search toolbar behavior as a full slice, including
  nested rules, object-valued duration/range chips, input behavior, history and
  helper serialization. See `docs/ai/advanced-search-migration-audit.md`.
- `[x]` Preserve remediation instruction filters in the runtime toolbar UI,
  including lazy fetching for selected instructions, the first page, search and
  pagination instead of loading the full instruction list.
- `[x]` Preserve toolbar alarm/entity/pbehavior patterns through advanced-search
  query serialization and export query forwarding.
- `[x]` Preserve live-reporting range and time field in the runtime UI.
- `[x]` Restore toolbar CSV export button, export polling mixin, alarm export
  download URL helper and export query helpers with old JSDoc.
- `[x]` Preserve group-alarm and alarm-details queries in Step 9 together with
  the expand panel migration. Periodic list refresh now also refetches opened
  expand-panel details.
- `[x]` Refresh expanded rows after data changes and drop rows that disappeared
  together with expanded row state. List fetch still refetches opened details;
  `synchronizeExpandedAlarms` drops expanded ids that left the page, and the
  expand panel removes its details query on unmount.
- `[x]` Do not port normalizr embedded cleanup or `entities/removeFromStore`.

## Step 3: Refresh And Sockets

- `[x]` Connect periodic refresh to widget visibility.
- `[x]` Pause refresh where the old widget paused it, including
  `pausePeriodicRefreshOnExpandPanel`.
- `[x]` Join and leave alarm socket rooms by visible alarm ids.
- `[x]` Join and leave alarm-details rooms by expanded details queries.
- `[x]` Recreate subscriptions on `visible`, `liveWatching`, edit mode and data
  changes.
- `[x]` Update alarm and details records in Pinia from socket events.
- `[x]` Remove all listeners on unmount.

Implemented files:

- `src/mixins/widget/periodic-refresh.js`
- `src/mixins/widget/alarms-socket.js`
- `src/stores/modules/alarm.js`
- `src/helpers/entities/alarm/list.js`
- `src/components/widgets/alarm/alarms-list.vue`
- `src/components/widgets/alarm/partials/alarms-list-table.vue`

Verification:

- `npx eslint src/mixins/widget/periodic-refresh.js src/mixins/widget/alarms-socket.js src/stores/modules/alarm.js src/helpers/entities/alarm/list.js src/components/widgets/alarm/alarms-list.vue src/components/widgets/alarm/partials/alarms-list-table.vue`

## Step 4: Dedicated AlarmList Table

- `[x]` Build a dedicated Vuetify 4 table wrapper, without
  `c-advanced-data-table`. Runtime uses `alarms-list-table-with-pagination` ->
  `alarms-list-table` -> `alarms-list-row`.
- `[x]` Preserve server-side pagination, sorting, loading and empty-column
  states. Empty columns render `c-empty-data-table-columns`.
- `[x]` Preserve old top and bottom pagination layout:
  `c-pagination type="top"` above the table, hidden Vuetify default footer, and
  `c-table-pagination` with the old per-page values below the table.
- `[x]` Preserve custom header rendering.
- `[>]` Preserve custom row and expanded-row rendering in Step 6 and Step 9.
- `[x]` Preserve selectable rows, select-all behavior and disabled rows.
- Preserve expanded rows, transition behavior and stable row keys.
- `[x]` Preserve dense, medium-dense and ultra-dense runtime density selection.
- `[x]` Preserve sticky table header, sticky top toolbar with density/pagination
  controls and sticky horizontal scrollbar through the dedicated AlarmList table
  chain and the migrated `c-horizontal-scrollbar` common component.
- `[x]` Preserve virtual rendering/lazy booting where enabled. Delayed
  viewport-nearest row booting and `isVirtualScrollEnabled` off-screen row
  placeholders are migrated through `helpers/render`, `TableIntersectionObserver`
  and the alarm table/row chain. Nested children and alarm-list modal tables
  stay eager.
- `[x]` Preserve column drag, resize and reset controls.
- Preserve parent-alarm and correlation-child modes.
- Expose explicit Vue 3 props, events, slots and `emits`.

## Step 5: Columns And Cells

- `[x]` Port alarm runtime columns used by `AlarmsList`.
- `[>]` Port group and service-dependency columns with the expand-panel and
  service-dependency scopes.
- `[x]` Port labels, aliases, sorting rules and used-property tracking through
  the alarm widget form mapper.
- `[x]` Port column templates, HTML sanitizing and color indicators.
- `[x]` Port header priority and custom header cells.
- `[x]` Port state/status/relation rendering.
- `[x]` Port all extra-details chips, badges, tooltips and popup bodies.
- `[x]` Render extra-details tooltip HTML through explicit sanitized
  `c-simple-tooltip html` mode, matching the old HTML tooltip behavior without
  enabling HTML for plain text tooltips globally.
- `[x]` Port tags and tag click/remove events.
- `[x]` Port alarm links helper slice, link chips and link width calculation
  helper.
- `[x]` Port responsive column widths, backend column filters and entity info
  properties with column resize/filter controls. Runtime column drag, resize
  and wrap/truncate behavior are migrated. Backend `alarm-column-filters` and
  entity info property aliases are wired through the alarm columns mixins.

Implemented files:

- `src/helpers/entities/alarm/list.js`
- `src/helpers/entities/link/list.js`
- `src/components/widgets/alarm/headers-formatting/alarm-header-cell.vue`
- `src/components/widgets/alarm/headers-formatting/alarm-header-priority.vue`
- `src/components/widgets/alarm/columns-formatting/alarm-column-value.vue`
- `src/components/widgets/alarm/columns-formatting/alarm-column-cell.vue`
- `src/components/widgets/alarm/columns-formatting/alarm-column-cell-popup-body.vue`
- `src/components/widgets/alarm/columns-formatting/alarm-column-value-status.vue`
- `src/components/widgets/alarm/columns-formatting/alarm-column-value-extra-details.vue`
- `src/components/widgets/alarm/columns-formatting/extra-details/*`
- `src/components/widgets/alarm/hooks/extra-details-tooltips.js`
- `src/components/widgets/alarm/partials/alarm-status-chip-with-relations.vue`
- `src/components/common/tooltip/c-simple-tooltip.vue`
- `src/components/common/chips/c-alarm-chip.vue`
- `src/components/common/chips/c-alarm-actions-chips.vue`
- `src/components/common/chips/c-alarm-state-chip.vue`
- `src/components/common/chips/c-alarm-tags-chips.vue`
- `src/components/common/chips/c-alarm-links-chips.vue`
- `src/components/common/chips/c-links-chips.vue`
- `src/components/common/chips/c-alarm-pbehavior-chip.vue`
- `src/components/common/table/color-indicator-wrapper.vue`
- `src/components/common/icon/c-no-events-icon.vue`

Verification:

- `npx eslint ...` for the full columns/header/extra-details/chips helper
  scope.
- `npm run build`

## Step 6: Rows And Selection

- `[x]` Port row selection and disabled selection rules.
- `[x]` Port select-all and Ctrl-drag selection.
- `[x]` Port bookmarks and remediation instruction icons.
- Port expand button and expanded state.
- `[x]` Port inline actions shell and connect the row action column to the
  AlarmList table.
- `[x]` Port action availability for row actions. The base action list,
  permission filtering, quick-action ordering and action pending state are
  connected; modal integration is completed in Step 7 and Step 12.
- Port row highlighting, search context and parent-alarm context.
- `[x]` Port resize handles, density classes and lazy/virtual booting. Resize
  handles, column drag/drop, lock/reset toolbar buttons, user-preference
  column width/position persistence, delayed row booting and virtual scroll
  are migrated.
- Preserve all row events with explicit Vue 3 event names.

Implemented files:

- `src/helpers/actions-panel.js`
- `src/helpers/entities/alarm/form.js`
- `src/helpers/entities/remediation/instruction-execution/form.js`
- `src/components/common/actions-panel/actions-panel.vue`
- `src/components/common/actions-panel/actions-panel-btn.vue`
- `src/components/common/actions-panel/actions-panel-menu.vue`
- `src/components/common/actions-panel/index.js`
- `src/components/widgets/alarm/actions/actions-panel.vue`
- `src/components/widgets/alarm/partials/alarms-list-row-bookmark-icon.vue`
- `src/components/widgets/alarm/partials/alarms-list-row-instructions-icon.vue`
- `src/components/widgets/alarm/alarms-list.vue`
- `src/hooks/table/columns-settings.js`
- `src/helpers/render.js`
- `src/helpers/filter.js`
- `src/services/table-intersection-observer.js`
- `src/mixins/widget/columns/alarm.js`
- `src/mixins/widget/columns/entity-info-property.js`
- `src/mixins/widget/rows/alarm-selecting.js`
- `src/mixins/entities/associative-table/index.js`
- `src/mixins/entities/associative-table/alarm-columns-filters.js`
- `src/mixins/entities/entity-info-property.js`
- `src/components/widgets/alarm/partials/alarms-list-row.vue`
- `src/components/widgets/alarm/partials/alarms-list-table.vue`
- `src/components/widgets/alarm/partials/alarms-list-table-with-pagination.vue`

Verification:

- `npx eslint src/helpers/actions-panel.js src/helpers/entities/alarm/form.js src/helpers/entities/remediation/instruction-execution/form.js src/components/common/actions-panel/actions-panel.vue src/components/common/actions-panel/actions-panel-btn.vue src/components/common/actions-panel/actions-panel-menu.vue src/components/common/actions-panel/index.js src/components/widgets/alarm/actions/actions-panel.vue src/components/widgets/alarm/alarms-list.vue`
- `npm run build`

## Step 7: Single-Alarm Actions

- `[x]` Port base row action panel rendering, quick action sorting, inline/menu
  split, permissions, loading state and common action panel UI.
- `[x]` Port action handlers for existing infrastructure: fast ack/cancel,
  bookmark add/remove, variable helper, dynamic links, declare ticket execution
  handoff, pbehavior planning/fast pbehavior and remediation instruction
  execution are wired to the current target stores/modals.
- `[x]` Event, ticket, comment, snooze, change-state, history, meta-alarm and
  remediation execute modals are registered and opened from the row action
  panel. PDF export is wired from the row action panel and variables-help.

Remaining action parity:

- alarm details (Step 9 expand panel).

## Step 8: Mass Actions

- `[x]` Port the mass-actions panel into the dedicated table toolbar through
  `c-mass-actions-panel` and the widget grid teleport target.
- `[x]` Preserve selected alarms, available actions and disabled rows.
- `[x]` Preserve pending state per action.
- `[x]` Preserve keep-selected-after-action behavior.
- `[x]` Port all bulk alarm, ticket, pbehavior, bookmark, meta-alarm and
  remediation actions. Bulk action selection/filtering, store dispatches and
  shared event/meta/remediation modal registrations are migrated with the
  single-action modal stack.
- `[x]` Port follow-up refresh and selection clearing.
- Verify panel z-index, top-bar overlap and no clipping.

Implemented files:

- `src/components/widgets/alarm/actions/mass-actions-panel.vue`
- `src/components/widgets/alarm/alarms-list.vue`
- `src/helpers/entities/remediation/instruction-execution/form.js`

Verification:

- `npx eslint src/helpers/entities/remediation/instruction-execution/form.js src/components/widgets/alarm/alarms-list.vue src/components/widgets/alarm/actions/mass-actions-panel.vue src/components/widgets/alarm/partials/alarms-list-row-bookmark-icon.vue src/components/widgets/alarm/partials/alarms-list-row-instructions-icon.vue`
- `npm run build`

## Step 9: Expand Panel

- `[x]` Port conditional expand panel tabs, permissions and tab ordering:
  More infos, Timeline, Charts, Declared tickets, Pbehaviors, Child alarms,
  Track source, Impact chain, Entity gantt, Remediation, Availability and
  Entity enrichments.
- `[x]` Port remediation executions tab with executions table, status icons and
  nested execution steps.
- `[x]` Port entity enrichments tab with entity-info logs table, search,
  interval toolbar and pagination.
- `[x]` Port service dependency helper slice with old JSDoc/typedefs and adapt
  Track source / Impact chain rendering to the target Vuetify 4 tree-table
  surface. The old `c-treeview-data-table` wrapper was not ported because it is
  a Vuetify 2 table abstraction; the target keeps the same lazy dependency
  loading, load-more and nested service-dependencies modal behavior.
- `[x]` Port availability tab with old quick interval and show-type controls.
- `[x]` Port entity gantt data loading, pagination and Chart.js gantt rendering
  through `junit-gantt-chart` / `horizontal-bar`.
- `[x]` Port chart tab data loading and visual rendering for configured entity
  metrics: `bar-chart-metrics`, `line-chart-metrics`, `numbers-metrics`, interval
  / sampling filters and CSV/PNG export.
- `[x]` Keep explicit local refresh for pbehavior lists. Do not restore the old
  global normalizr deletion mechanism.
- `[x]` Group-alarm and alarm-details query socket refresh is wired through
  `widgetAlarmsSocketMixin`.
- `[x]` Expanded row state refresh/drop after data changes: opened details are
  refetched with the list, disappeared rows are dropped from `expandedAlarmIds`,
  and `removeDetailsQuery` runs when the expand panel unmounts.
- `[x]` Visual chart/gantt rendering for the expand-panel Charts and Entity
  gantt tabs.

Implemented files:

- `src/components/widgets/alarm/expand-panel/alarms-expand-panel.vue`
- `src/components/widgets/alarm/expand-panel/alarms-expand-panel-remediation.vue`
- `src/components/widgets/alarm/expand-panel/alarms-expand-panel-entity-enrichments.vue`
- `src/components/other/remediation/instruction-execute/remediation-instruction-executions-list.vue`
- `src/components/other/remediation/instruction-execute/partials/remediation-instruction-execution-status-icon.vue`
- `src/components/other/remediation/instruction-execute/partials/remediation-instruction-executions-expand-panel.vue`
- `src/components/other/remediation/instruction-execute/partials/remediation-instruction-executions-expand-panel-row.vue`
- `src/components/other/entity/entity-infos-logs-list.vue`
- `src/components/other/entity/entity-availability.vue`
- `src/components/other/availability/partials/availability-bar.vue`
- `src/components/other/availability/partials/availability-bar-chart.vue`
- `src/components/other/availability/partials/availability-bar-chart-information-row.vue`
- `src/components/other/availability/partials/availability-filters.vue`
- `src/components/other/service/partials/service-dependencies.vue`
- `src/components/other/service/partials/service-dependencies-row.vue`
- `src/components/other/service/partials/service-dependencies-expand.vue`
- `src/components/other/service/partials/service-dependencies-entity-cell.vue`
- `src/components/other/service/form/fields/service-dependencies-show-type-field.vue`
- `src/components/other/state-setting/state-settings-summary.vue`
- `src/components/widgets/chart/entity-charts.vue`
- `src/components/widgets/chart/entity-chart-widget.vue`
- `src/components/widgets/chart/partials/bar-chart-metrics.vue`
- `src/components/widgets/chart/partials/line-chart-metrics.vue`
- `src/components/widgets/chart/partials/numbers-metrics.vue`
- `src/components/widgets/chart/partials/numbers-metrics-item.vue`
- `src/components/widgets/chart/partials/chart-loader.vue`
- `src/components/widgets/partials/kpi-widget-filters.vue`
- `src/components/widgets/alarm/entity-gantt/entity-gantt.vue`
- `src/components/common/chart/line-chart.vue`
- `src/components/common/chart/horizontal-bar.vue`
- `src/components/common/chart/junit-gantt-chart.vue`
- `src/mixins/chart/metrics-options.js`
- `src/mixins/widget/chart/metrics-map.js`
- `src/helpers/entities/metric/list.js`
- `src/helpers/entities/metric/form.js`
- `src/helpers/entities/metric/query.js`
- `src/components/modals/service/service-dependencies.vue`
- `src/helpers/entities/service-dependencies/list.js`

Verification:

- `npx eslint` over the full Step 9 file set.
- `npx oxlint` over the full Step 9 file set.

## Step 10: Toolbar And Filters

- `[x]` Port advanced search and saved searches, including old-style history dropdown, pinned/remove actions, client-side infinite scroll, virtualized history rows, and the custom `variables-list -> c-list` selection layer used by attribute/operator/value menus.
- `[x]` Port category and correlation controls.
- `[x]` Port user filter selector and filter management button.
- `[x]` Port tags, bookmarks and remediation instruction filters.
- `[x]` Port live-reporting chip, edit modal and reset behavior.
- `[x]` Port CSV export button and export query flow.
- `[x]` Port toolbar visibility modes, including kiosk `hideToolbar` propagation
  from `widget-wrapper`.
- `[x]` Port runtime density selector above the table and persist it through
  user preferences.
- `[x]` Preserve empty-value behavior of all comboboxes and selectors used in
  the toolbar.
- `[x]` Preserve locked widget filter display in the filter selector selection
  with a lock chip and tooltip.

## Complete Settings Sidebar Scope

The sidebar is migrated as a single scope. Do not mark individual settings
groups as complete while the full sidebar flow is still incomplete. The scope
includes the sidebar component, every nested field, every helper/mixin/hook it
uses, widget-template loading, entity-info loading, permissions, validation,
dirty-state handling, submit mapping, translations.

The sidebar scope is complete only when all settings below can be opened,
edited, validated, saved, reopened and mapped back without losing values.

Current progress:

- `[x]` Register `alarmSettings` in the target sidebar registry.
- `[x]` Connect the sidebar to `useWidgetSettings()`.
- `[x]` Port the full AlarmList parameter form/payload mapper foundation.
- `[x]` Port quick-action, kiosk, availability and grid-range helper foundations.
- `[x]` Port Time filter and meta-alarm/track-source column settings.
- `[x]` Keep AlarmList column settings direct in the sidebar; do not reuse
  `AlarmsListModalForm`, which is reserved for Counter, Stats calendar and
  Service weather settings.
- `[x]` Add direct root-cause and expand-panel More infos settings.
- `[x]` Port the remaining nested settings fields and their modals. The Time
  filter (`liveReporting`) field, direct column settings, the Actions section
  in the old order (Ack, Cancel, comment templates, fast pbehavior, action
  switches, Export PDF, quick actions and quick mass actions), availability
  graph settings, Export CSV, Export PDF template and Charts settings
  (including BarChart, LineChart and Numbers chart forms) are implemented.

- `[x]` Complete the sidebar implementation and move it to Ready For Testing.

Port every group from the old sidebar:

- title and periodic refresh/live watching;
- live reporting;
- default sort and sort templates;
- alarm, group and service-dependency columns;
- column templates, filters, color indicators and HTML;
- drag/resize and wrap/truncate behavior;
- root cause and info popups;
- main filters, opened/resolved, correlation and clear-filter behavior;
- items per page, density, sticky header and horizontal scroll;
- virtual scroll and kiosk options;
- acknowledgement, cancel, comment, pbehavior, snooze and ticket options;
- meta-alarm, OK-state and keep-selected options;
- PDF template and quick/quick-mass actions;
- expand-panel templates, charts, availability and timeline options;
- CSV export options, instruction fields and optional info attributes.

Use `useWidgetSettings()` for the new Composition API settings root. The legacy
widget settings mixin has been removed after widget settings sidebars were
migrated.

## Step 12: Modals And Exports

- `[x]` Register and adapt all alarm event, ticket, pbehavior, meta-alarm and
  chart modals used by single and mass actions.
- `[x]` Register live-reporting and info-popup modals.
- `[x]` Port CSV export request/poll/download flow.
- `[x]` Port PDF export and alarm-detail comment loading.
- `[x]` Port remediation instruction execution popups, including paused
  execution startup fetch/socket room and the runtime-template "See details"
  link component.
- `[x]` Preserve dirty-close confirmation, vee-validate behavior, progress
  overlays, z-index and Vuetify 4 slots for the registered action modals.

Implemented action modals:

- `src/components/modals/alarm/create-event.vue`
- `src/components/modals/alarm/create-ack-event.vue`
- `src/components/modals/alarm/create-snooze-event.vue`
- `src/components/modals/alarm/create-change-state-event.vue`
- `src/components/modals/common/create-comment-event.vue`
- `src/components/modals/alarm/create-remove-associated-ticket-event.vue`
- `src/components/modals/alarm/link-to-meta-alarm.vue`
- `src/components/modals/alarm/remove-alarms-from-meta-alarm.vue`
- `src/components/modals/alarm/alarms-list.vue`
- `src/components/modals/remediation/execute-remediation-instruction.vue`
- `src/components/modals/remediation/execute-remediation-simple-instruction.vue`
- `src/components/common/popups/c-remediation-instruction-execution-see-details.vue`

Intentionally skipped:

- `confirmAckWithTicket` — defined in the old source, never shown from AlarmList
  action panels.
- History modal expand panel and mass-actions teleport — owned by Step 9 / the
  dedicated table wrapper.

Verification:

- `npx eslint` on the registered action-modal files, forms, fields, helpers and
  remediation execute slice
- `npm run build`

## Step 13: Integration And Verification

Verify creation, editing, duplication, deletion, rendering, pagination,
sorting, search, filters, tags, bookmarks, live reporting, exports, sockets,
periodic refresh, selection, mass actions, every expand tab, settings,
permissions, kiosk mode, responsive layouts, dark/light themes and API errors.

- `[x]` Check gantt in the alarms list table expand panel (Entity gantt tab:
  intervals, pagination, tooltips and layout inside the expanded row).

The migration is complete only when the widget has no placeholder branches,
Vue 2 listener/slot API, `tc`, `$validator`/legacy injection warnings,
normalizr cleanup dependency or unregistered `c-*` components.

## Documentation Rule

After every completed step, update this file and `docs/migration-inventory.md`
in the same change. Record migrated files, intentionally skipped legacy code,
verification results and remaining risks.
