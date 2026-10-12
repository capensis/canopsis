# AlarmStatistics Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/statistics/alarm-statistics-widget.vue`

Target runtime: `src/components/widgets/statistics/alarm-statistics-widget.vue`

The widget is the Pro KPI alarm statistics table (`WIDGET_TYPES.alarmStatistics`,
component name `alarm-statistics-widget`). It belongs to the `/view/:id` widget
system. The runtime is a thin permission wrapper around shared
`statistics-widget`; UserStatistics uses the same stack with
`KPI_RATING_SETTINGS_TYPES.user`.

Prerequisite: shared statistics stack from
`docs/ai/user-statistics-widget-migration-plan.md` (form helpers, query
converters, `statistics-widget`, `statistics-form` and column/main-parameter
fields). Do not duplicate those files here.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` Step 1: confirm shared helper registration covers AlarmStatistics.
- `[x]` Step 2: complete AlarmStatistics settings sidebar.
- `[x]` Step 3: AlarmStatistics runtime wrapper.
- `[x]` Step 4: view integration.
- `[x]` Step 5: verification.

## Step 0: Inventory And Dependency Audit

### Source runtime

- `src/components/widgets/statistics/alarm-statistics-widget.vue` — Options API wrapper; alarm interval/filter mixins; renders `statistics-widget`.
- Shared runtime: `src/components/widgets/statistics/partials/statistics-widget.vue` (owned by the UserStatistics plan).

Pass `KPI_RATING_SETTINGS_TYPES.entity` into `statistics-widget` (source wrappers omit `type`, but fetch uses `this.type`).

### Source settings

- `src/components/sidebars/statistics/alarm-statistics.vue` — `SIDE_BARS.alarmStatisticsSettings`; `KPI_RATING_SETTINGS_TYPES.entity`.
- Shared form: `statistics-form.vue`, `main-parameter.vue`, `statistics-columns.vue`, `statistics-column.vue`.
- Entity-only behavior already in the shared form:
  - custom criteria radio (`KPI_ENTITY_RATING_SETTINGS_CUSTOM_CRITERIA`) with column name + entity patterns
  - `STATISTICS_WIDGETS_ENTITY_METRICS` for column metrics
  - query omits `criteria` when custom; sends `entity_patterns` from `mainParameter.patterns`

### Source mixins

- `src/mixins/permissions/widgets/statistics/alarm/{interval,filters}.js`

### Target already present

- Constants: `WIDGET_TYPES.alarmStatistics`, `SIDE_BARS.alarmStatisticsSettings`, `KPI_RATING_SETTINGS_TYPES.entity`, `STATISTICS_WIDGETS_ENTITY_METRICS`.
- i18n: `settings.alarmStatisticsSettings`, `settings.mainParameter.custom`, permission/modal titles.
- `create-widget.vue` already lists AlarmStatistics next to UserStatistics.
- Shared UI/stores listed in the UserStatistics plan.

### Target missing (this widget only)

- `src/mixins/permissions/widgets/statistics/alarm/{interval,filters}.js`
- AlarmStatistics filter/interval permission hooks for the sidebar
- `src/components/widgets/statistics/alarm-statistics-widget.vue`
- `src/components/sidebars/statistics/alarm-statistics.vue`
- Wrapper/menu/sidebar registrations

## Step 1: Helpers

- `[x]` Confirm `statisticsWidgetParametersToForm` / `formToStatisticsWidgetParameters` are registered for `WIDGET_TYPES.alarmStatistics` in `src/helpers/entities/widget/form.js`.
- `[x]` Confirm `convertStatisticsWidgetParametersToQuery` / `convertStatisticsUserPreferenceToQuery` are registered for AlarmStatistics in `src/helpers/entities/widget/query.js`.
- `[x]` Confirm custom criteria query shape: no `criteria` when `isCustomCriteria`, `entity_patterns` from `mainParameter.patterns`.

## Step 2: Settings Sidebar

- `[x]` Add `src/components/sidebars/statistics/alarm-statistics.vue` with `useWidgetSettings()` and `KPI_RATING_SETTINGS_TYPES.entity`.
- `[x]` Reuse shared `statistics-form` (title, periodic refresh, main parameter with custom criteria, entity metric columns, header, default time range, filters).
- `[x]` Add alarm-statistics permission mixins for Options API runtime.
- `[x]` Add filter/interval permission hooks for the Composition API sidebar (`USER_PERMISSIONS.business.alarmStatistics`).

## Step 3: Runtime

- `[x]` Add `src/components/widgets/statistics/alarm-statistics-widget.vue`. Keep Options API.
- `[x]` Wrap shared `statistics-widget` with alarm interval/filter mixins and `type: KPI_RATING_SETTINGS_TYPES.entity`.
- `[x]` Preserve toolbar interval / user filter (no sampling).
- `[x]` Preserve custom-criteria first column (`columnName`) and pattern-based grouping.

## Step 4: View Integration

- `[x]` Register `AlarmStatisticsWidget` in `widget-wrapper.vue`.
- `[x]` Register `AlarmStatisticsSettings` in `sidebars/index.js`.
- `[x]` Add `WIDGET_TYPES.alarmStatistics` to `MIGRATED_SETTINGS_WIDGET_TYPES` in `widget-wrapper-menu.vue`.
- `[x]` Update `docs/migration-inventory.md`, this plan, the UserStatistics plan and `docs/ai/view-migration-plan.md`.

## Step 5: Verification

- `[x]` Targeted ESLint passes for AlarmStatistics files.
- `[x]` Production build passes.
- `[x]` Ready For Testing entry includes:
  - create AlarmStatistics widget
  - edit and save settings (entity main parameter, custom criteria + patterns, columns, header, interval, filters)
  - required-columns validation
  - toolbar interval / user filter
  - table headers and formatted metric cells
  - custom-criteria grouping
  - periodic refresh
