# BarChart Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/chart/bar-chart-widget.vue`

Target runtime: `src/components/widgets/chart/bar-chart-widget.vue`

The widget is the KPI alarm bar chart (`WIDGET_TYPES.barChart`, component name
`bar-chart-widget`). It is part of the `/view/:id` widget system and reuses the
metrics API, KPI toolbar filters and migrated chart rendering partials.

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
  - `src/components/widgets/chart/partials/bar-chart-metrics.vue`
  - `src/components/widgets/chart/partials/chart-loader.vue`
  - `src/components/widgets/partials/kpi-widget-filters.vue`
- `[x]` Existing target metrics store exposes vector metric fetch and KPI alarm export actions.

## Step 1: Helpers

- `[x]` Register BarChart parameter form conversion in `src/helpers/entities/widget/form.js`.
- `[x]` Reuse existing BarChart query conversion in `src/helpers/entities/widget/query.js`.

## Step 2: Settings Sidebar

- `[x]` Add `src/components/sidebars/chart/bar-chart.vue`.
- `[x]` Reuse `src/components/sidebars/chart/form/bar-chart-widget-form.vue`.
- `[x]` Restore chart preset selection with preset help text.
- `[x]` Preserve permission behavior for chart filters.

## Step 3: Runtime

- `[x]` Add `src/components/widgets/chart/bar-chart-widget.vue`.
- `[x]` Preserve source Options API style.
- `[x]` Fetch vector alarm metrics by widget id and current KPI query.
- `[x]` Apply metric labels/colors from widget metric presets.
- `[x]` Preserve toolbar controls for interval, sampling and user filters.
- `[x]` Preserve CSV/PNG export actions.

## Step 4: View Integration

- `[x]` Register `BarChartWidget` in `widget-wrapper.vue`.
- `[x]` Register `BarChartSettings` in `sidebars/index.js`.
- `[x]` Enable edit/duplicate actions in `widget-wrapper-menu.vue`.

## Step 5: Verification

- `[x]` Targeted ESLint passes for BarChart files.
- `[x]` Production build passes.
- `[x]` Ready For Testing entry includes:
  - create BarChart widget
  - edit and save BarChart settings
  - filter by KPI interval/sampling/user filter
  - render metrics
  - export PNG/CSV
