# UserStatistics Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/statistics/user-statistics-widget.vue`

Target runtime: `src/components/widgets/statistics/user-statistics-widget.vue`

The widget is the Pro KPI user statistics table (`WIDGET_TYPES.userStatistics`,
component name `user-statistics-widget`). It belongs to the `/view/:id` widget
system. The runtime is a thin permission wrapper around shared
`statistics-widget`; AlarmStatistics uses the same stack with
`KPI_RATING_SETTINGS_TYPES.entity`.

Migrate this widget before AlarmStatistics. Shared files listed here must not be
reimplemented in `docs/ai/alarm-statistics-widget-migration-plan.md`.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` Step 1: widget form/query helper registration.
- `[x]` Step 2: complete settings sidebar and shared form fields.
- `[x]` Step 3: shared runtime table and UserStatistics wrapper.
- `[x]` Step 4: view integration.
- `[x]` Step 5: verification.

## Step 0: Inventory And Dependency Audit

### Source runtime

- `src/components/widgets/statistics/user-statistics-widget.vue` — Options API wrapper; mixins for interval/filter permissions; renders `statistics-widget`.
- `src/components/widgets/statistics/partials/statistics-widget.vue` — Options API; KPI toolbar, group metrics table, rating-settings first-column labels.

Source `statistics-widget.vue` calls `fetchMainRatingSettingsList` with `this.type`, but the wrappers do not pass `type`. Pass `KPI_RATING_SETTINGS_TYPES.user` from the UserStatistics wrapper so the first-column header can resolve rating-setting labels.

### Source settings

- `src/components/sidebars/statistics/user-statistics.vue` — `SIDE_BARS.userStatisticsSettings`; `KPI_RATING_SETTINGS_TYPES.user`.
- `src/components/sidebars/statistics/form/statistics-form.vue` — title, periodic refresh, main parameter, columns, table header, default time range, entity filters.
- `src/components/sidebars/statistics/form/fields/main-parameter.vue` — rating-settings radios; custom criteria (column name + entity patterns) only for entity type.
- `src/components/sidebars/statistics/form/fields/statistics-columns.vue` — required movable column list; rating settings fetched with `type: entity`, `main: false`.
- `src/components/sidebars/statistics/form/fields/statistics-column.vue` — metric, custom label, split-by-infos for `STATISTICS_WIDGETS_USER_METRICS_WITH_ENTITY_TYPE`.

### Source helpers / mixins / store

- `src/helpers/entities/widget/forms/statistics.js`
- `src/helpers/entities/metric/query.js` — `convertStatisticsWidgetParametersToQuery`, `convertStatisticsUserPreferenceToQuery`
- `src/helpers/entities/metric/form.js` — `isCustomCriteria`
- `src/mixins/entities/group-metrics.js` — Vuex `groupMetrics` list/pending/meta + fetch
- `src/mixins/permissions/widgets/statistics/user/{interval,filters}.js`

### Target already present

- Constants: `WIDGET_TYPES.userStatistics`, `SIDE_BARS.userStatisticsSettings`, `KPI_RATING_SETTINGS_TYPES`, `KPI_RATING_CRITERIA`, `KPI_ENTITY_RATING_SETTINGS_CUSTOM_CRITERIA`, `STATISTICS_WIDGETS_*`.
- i18n: `settings.userStatisticsSettings`, `settings.mainParameter`, `settings.statisticsWidgetColumn`, `kpi.statisticsWidgets.metrics`, permission/modal widget titles.
- `create-widget.vue` already lists UserStatistics / AlarmStatistics as Pro children.
- Pinia `metrics` store: `fetchGroupMetricsList`, group list/pending/meta by widget id.
- Pinia `rating-settings` store: `fetchListWithoutStore`.
- Shared UI: `kpi-widget-filters`, `c-advanced-data-table`, `c-progress-overlay`, `kpi-rating-metric-field`, `c-card-iterator-field` / `c-draggable-list-field`, widget settings fields (`title`, `periodic-refresh`, `quick-date-interval-type`, `filters`, `filters-list`).
- Widget query mixins: `widgetFetchQueryMixin`, `widgetFilterSelectMixin`, `queryIntervalFilterMixin`, `widgetPeriodicRefreshMixin`.

### Target missing

- `src/helpers/entities/widget/forms/statistics.js`
- `isCustomCriteria` in `src/helpers/entities/metric/form.js`
- Statistics query converters in `src/helpers/entities/metric/query.js` and registration in `widget/query.js` / `widget/form.js`
- `src/mixins/entities/group-metrics.js` (or Pinia `mapState` / `mapActions` on `useMetricsStore`, same pattern as Numbers)
- Permission mixins `src/mixins/permissions/widgets/statistics/user/{interval,filters}.js`
- Permission hooks for the settings sidebar (`userStatistics` filter/interval), same pattern as `hooks/permissions/widgets/chart/filters.js`
- All `components/widgets/statistics/**` and `components/sidebars/statistics/**`
- Wrapper/menu/sidebar registrations

## Step 1: Helpers

- `[x]` Port `src/helpers/entities/widget/forms/statistics.js` (`statisticsWidgetParametersToForm`, `formToStatisticsWidgetParameters`, column/main-parameter converters).
- `[x]` Port `isCustomCriteria` into `src/helpers/entities/metric/form.js`.
- `[x]` Port `convertStatisticsWidgetParametersToQuery` and `convertStatisticsUserPreferenceToQuery` into `src/helpers/entities/metric/query.js`.
- `[x]` Register UserStatistics (and AlarmStatistics) in `src/helpers/entities/widget/form.js`.
- `[x]` Register UserStatistics (and AlarmStatistics) in `src/helpers/entities/widget/query.js` for both widget and user-preference converters.

## Step 2: Settings Sidebar

- `[x]` Add `src/components/sidebars/statistics/user-statistics.vue` with `useWidgetSettings()` and `KPI_RATING_SETTINGS_TYPES.user`.
- `[x]` Add `src/components/sidebars/statistics/form/statistics-form.vue`. Keep Options API. Replace `v-field` / `formMixin` with `v-model` + `localForm` like `numbers-widget-form.vue`.
- `[x]` Add `main-parameter.vue`, `statistics-columns.vue`, `statistics-column.vue`.
  - Use `c-layout`, `c-select`, `c-enabled-field`, `c-card-iterator-field`.
  - Fetch rating settings through Pinia `useRatingSettingsStore` (`mapActions`).
  - Custom criteria UI stays gated on entity type (AlarmStatistics); still port it in this shared field.
  - Split toggle stays limited to `STATISTICS_WIDGETS_USER_METRICS_WITH_ENTITY_TYPE`.
  - Column metrics: `STATISTICS_WIDGETS_USER_METRICS` vs `STATISTICS_WIDGETS_ENTITY_METRICS` from `type`.
- `[x]` Add user-statistics permission mixins for Options API runtime.
- `[x]` Add filter/interval permission hooks for the Composition API sidebar.

## Step 3: Runtime

- `[x]` Add `src/components/widgets/statistics/partials/statistics-widget.vue`. Keep Options API.
- `[x]` Reuse `kpi-widget-filters` (interval + user filter; no sampling).
- `[x]` Fetch group metrics by widget id through Pinia metrics store.
- `[x]` Fetch main rating settings without store; pass `type` from the wrapper.
- `[x]` Build headers from rating-setting labels / `KPI_RATING_CRITERIA` map / column labels, with `kpi.statisticsWidgets.metrics` → `alarm.metrics` → `user.metrics` fallback.
- `[x]` Flatten group-metric cells with `convertMetricValueToString` and `DATETIME_FORMATS.refreshFieldFormat`.
- `[x]` Clamp interval `from` to `groupMetricsMeta.min_date` in the system timezone.
- `[x]` Add `src/components/widgets/statistics/user-statistics-widget.vue` wrapping the shared runtime with user interval/filter mixins and `type: KPI_RATING_SETTINGS_TYPES.user`.

## Step 4: View Integration

- `[x]` Register `UserStatisticsWidget` in `widget-wrapper.vue`.
- `[x]` Register `UserStatisticsSettings` in `sidebars/index.js`.
- `[x]` Add `WIDGET_TYPES.userStatistics` to `MIGRATED_SETTINGS_WIDGET_TYPES` in `widget-wrapper-menu.vue`.
- `[x]` Update `docs/migration-inventory.md`, this plan, the AlarmStatistics plan and `docs/ai/view-migration-plan.md`.

## Step 5: Verification

- `[x]` Targeted ESLint passes for UserStatistics / shared statistics files.
- `[x]` Production build passes.
- `[x]` Ready For Testing entry includes:
  - create UserStatistics widget
  - edit and save settings (main parameter, columns, header, interval, filters)
  - required-columns validation
  - toolbar interval / user filter
  - table headers and formatted metric cells
  - periodic refresh
