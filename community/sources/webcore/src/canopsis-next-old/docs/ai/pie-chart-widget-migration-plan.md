# PieChart Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/chart/pie-chart-widget.vue`

Target runtime: `src/components/widgets/chart/pie-chart-widget.vue`

The widget is the KPI alarm pie chart (`WIDGET_TYPES.pieChart`, component name
`pie-chart-widget`). It belongs to the `/view/:id` widget system and reuses the
aggregated metrics API, KPI toolbar filters, shared metric preset form and migrated
Chart.js rendering wrappers.

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
- `[x]` Target lacked PieChart-specific form fields, common chart wrapper and metrics partial.
- `[x]` Target metrics store exposes aggregated metric fetch and KPI aggregate export actions.
- `[x]` Shared chart metric preset fields are available from the BarChart migration.

## Step 1: Helpers

- `[x]` Add PieChart parameter form conversion in `src/helpers/entities/widget/forms/pie-chart.js`.
- `[x]` Register PieChart parameter form conversion in `src/helpers/entities/widget/form.js`.
- `[x]` Add PieChart query conversion in `src/helpers/entities/metric/query.js`.
- `[x]` Register PieChart runtime and user preference query conversion in `src/helpers/entities/widget/query.js`.

## Step 2: Settings Sidebar

- `[x]` Add `src/components/sidebars/chart/pie-chart.vue`.
- `[x]` Add `src/components/sidebars/chart/form/pie-chart-widget-form.vue`.
- `[x]` Add PieChart-specific show mode and aggregate function fields.
- `[x]` Preserve permission behavior for chart filters.

## Step 3: Runtime

- `[x]` Add `src/components/widgets/chart/pie-chart-widget.vue`.
- `[x]` Preserve source Options API style.
- `[x]` Fetch aggregated alarm metrics by widget id and current KPI query.
- `[x]` Apply metric labels/colors from widget metric presets.
- `[x]` Preserve toolbar controls for interval, sampling and user filters.
- `[x]` Preserve CSV/PNG export actions.
- `[x]` Add Chart.js pie wrapper and empty-pie rendering plugin.

## Step 4: View Integration

- `[x]` Register `PieChartWidget` in `widget-wrapper.vue`.
- `[x]` Register `PieChartSettings` in `sidebars/index.js`.
- `[x]` Enable edit/duplicate actions in `widget-wrapper-menu.vue`.

## Step 5: Verification

- `[x]` Targeted ESLint passes for PieChart files.
- `[x]` Production build passes.
- `[x]` Ready For Testing entry includes:
  - create PieChart widget
  - edit and save PieChart settings
  - filter by KPI interval/sampling/user filter
  - render metrics
  - export PNG/CSV
