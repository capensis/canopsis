# LineChart Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/chart/line-chart-widget.vue`

Target runtime: `src/components/widgets/chart/line-chart-widget.vue`

The widget is the KPI alarm line chart (`WIDGET_TYPES.lineChart`, component name
`line-chart-widget`). It belongs to the `/view/:id` widget system and reuses the
metrics API, KPI toolbar filters, shared metric preset form and migrated chart rendering partials.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` Step 1: widget form/query helper registration.
- `[x]` Step 2: complete settings sidebar.
- `[x]` Step 3: runtime chart and metric export.
- `[x]` Step 4: view integration.
- `[x]` Step 5: Ready For Testing update.

## Step 0: Inventory And Dependency Audit

- `[x]` Source runtime is Options API and must stay Options API.
- `[x]` Source settings sidebar uses the widget settings mixin; target settings sidebars use `useWidgetSettings()`.
- `[x]` Existing target chart partials are reusable:
  - `src/components/widgets/chart/partials/line-chart-metrics.vue`
  - `src/components/widgets/chart/partials/chart-loader.vue`
  - `src/components/widgets/partials/kpi-widget-filters.vue`
- `[x]` Existing target metrics store exposes vector metric fetch and KPI alarm export actions.
- `[x]` Shared chart metric preset fields are available from the BarChart migration.

## Step 1: Helpers

- `[x]` Register LineChart parameter form conversion in `src/helpers/entities/widget/form.js`.
- `[x]` Reuse existing LineChart query conversion in `src/helpers/entities/widget/query.js`.

## Step 2: Settings Sidebar

- `[x]` Add `src/components/sidebars/chart/line-chart.vue`.
- `[x]` Reuse `src/components/sidebars/chart/form/line-chart-widget-form.vue`.
- `[x]` Preserve permission behavior for chart filters.

## Step 3: Runtime

- `[x]` Add `src/components/widgets/chart/line-chart-widget.vue`.
- `[x]` Preserve source Options API style.
- `[x]` Fetch vector alarm metrics by widget id and current KPI query.
- `[x]` Apply metric labels/colors from widget metric presets.
- `[x]` Preserve toolbar controls for interval, sampling and user filters.
- `[x]` Preserve CSV/PNG export actions.

## Step 4: View Integration

- `[x]` Register `LineChartWidget` in `widget-wrapper.vue`.
- `[x]` Register `LineChartSettings` in `sidebars/index.js`.
- `[x]` Enable edit/duplicate actions in `widget-wrapper-menu.vue`.

## Step 5: Verification

- `[x]` Targeted ESLint passes for LineChart files.
- `[x]` Production build passes.
- `[x]` Ready For Testing entry includes:
  - create LineChart widget
  - edit and save LineChart settings
  - filter by KPI interval/sampling/user filter
  - render metrics
  - export PNG/CSV
