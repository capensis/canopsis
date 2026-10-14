# Availability Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/availability/availability-widget.vue`

Target runtime: `src/components/widgets/availability/availability-widget.vue`

The widget is the Pro availability table (`WIDGET_TYPES.availability`,
component name `availability-widget`). It belongs to the `/view/:id` widget
system. Runtime is a toolbar + paginated entity availability list with an
expand panel (history graph, active alarms, resolved alarms).

Shared entity-availability bar UI already migrated for AlarmsList / Context
expand tabs must not be reimplemented here.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` Step 1: widget form/query helper registration.
- `[x]` Step 2: complete settings sidebar and shared form fields.
- `[x]` Step 3: runtime list, graph stack and widget wrapper.
- `[x]` Step 4: view integration.
- `[x]` Step 5: verification.

## Step 0: Inventory And Dependency Audit

### Source runtime

- `src/components/widgets/availability/availability-widget.vue` — Options API;
  filters toolbar + `availability-list`; interval/filter/export permissions;
  periodic refresh; CSV export; query omits `showType`/`displayParameter` when
  comparing if `valueFilter` is empty.
- `src/components/widgets/availability/partials/availability-widget-filters.vue`
  — advanced search, interval (`with-hours`), filter selector, display
  parameter, show type, trend toggle, value filter, CSV export.

Interval conversion uses `convertQueryIntervalToTimestamp` (hour unit), not
metric day-start conversion.

### Source settings

- `src/components/sidebars/availability/availability.vue` —
  `SIDE_BARS.availabilitySettings`; widget templates; permissions.
- `src/components/sidebars/availability/form/availability-form.vue` — title,
  periodic refresh, entity columns, active/resolved alarm columns, default
  time range, display parameter, show type, filters, export CSV.
- Sidebar wrappers around existing radio fields:
  `availability-display-parameter.vue`, `availability-display-show-type.vue`.

### Source shared UI

- `src/components/other/availability/availability-list.vue`
- `partials/availability-list-column-value.vue`
- `partials/availability-list-expand-panel.vue`
- `partials/availability-graph.vue` / `availability-history.vue` /
  `availability-history-filters.vue` / `availability-line-chart.vue`
- Toolbar fields: display-parameter, show-type, value-filter, value-filter-method

### Source helpers / mixins

- `src/helpers/entities/availability/entity.js`
- `src/helpers/entities/availability/query.js`
- `src/mixins/entities/availability.js`
- `src/mixins/permissions/widgets/availability/{filters,interval,export}.js`
  (interval mixin is misnamed `permissionsWidgetsAlarmStatisticsInterval` in
  Vue 2; Vue 3 uses `permissionsWidgetsAvailabilityInterval`)

### Target already present

- Constants, i18n, `WIDGET_TYPES.availability`, `SIDE_BARS.availabilitySettings`
- Form helpers in `src/helpers/entities/widget/forms/availability.js`
  (`availabilityWidgetParametersToForm`, `formToAvailabilityWidgetParameters`,
  `prepareAvailabilityWidget`)
- Pinia metrics store: availability list/pending/meta, fetch, history, export
- Shared bar UI: `availability-bar*`, `availability-filters`, radio fields,
  `entity-availability.vue`, `c-availability-advanced-search`
- `create-widget.vue` already lists Availability as a Pro top-level type
- Widget mixins: fetch-query, filter-select, periodic-refresh, options, export
- `entity-alarms-list-table`, `entity-column-cell`, `line-chart`,
  `chart-export-actions`

### Target missing

- Entity/query helpers and form/query registration
- Permission mixins and hooks
- List / graph / toolbar field stack
- Runtime widget and settings sidebar
- Wrapper / menu / sidebar registration
- `prepareAvailabilityWidget` export-column labels
- `convertQueryIntervalToTimestamp` / `findQuickRangeByInterval`

## Step 1: Helpers

- `[x]` Port `src/helpers/entities/availability/entity.js`.
- `[x]` Port `src/helpers/entities/availability/query.js`.
- `[x]` Add `convertQueryIntervalToTimestamp` and `findQuickRangeByInterval`.
- `[x]` Register Availability in `widget/form.js` and `widget/query.js`.
- `[x]` Add export-column labels to `prepareAvailabilityWidget`.

## Step 2: Settings Sidebar

- `[x]` Add permission mixins and hooks (filter, interval, export).
- `[x]` Add `availability.vue` with `useWidgetSettings()` and widget templates.
- `[x]` Add `availability-form.vue` with `v-model` + `localForm`.
- `[x]` Add display-parameter / show-type settings items.

## Step 3: Runtime

- `[x]` Port list, expand panel, graph/history/line-chart and toolbar fields.
- `[x]` Port `availability-widget.vue` (Options API) and widget filters.
- `[x]` Reuse Pinia metrics availability bucket and export mixin.

## Step 4: View Integration

- `[x]` Register `AvailabilityWidget` in `widget-wrapper.vue` with
  `prepareAvailabilityWidget`.
- `[x]` Register `AvailabilitySettings` in `sidebars/index.js`.
- `[x]` Add `WIDGET_TYPES.availability` to `MIGRATED_SETTINGS_WIDGET_TYPES`.
- `[x]` Update inventory and view-migration-plan docs.

## Step 5: Verification

- `[x]` Targeted ESLint for Availability files.
- `[x]` Production build.
- `[x]` Ready For Testing entry includes:
  - create Availability widget
  - edit and save settings (columns, interval, display, filters, export CSV)
  - toolbar interval / filter / show type / trend / value filter
  - table values, trend arrows, expand graph and alarm tabs
  - CSV export and periodic refresh
