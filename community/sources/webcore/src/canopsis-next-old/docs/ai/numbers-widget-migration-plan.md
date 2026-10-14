# Numbers Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/chart/numbers-widget.vue`

Target runtime: `src/components/widgets/chart/numbers-widget.vue`

The widget is the KPI alarm numbers chart (`WIDGET_TYPES.numbers`, component name
`numbers-widget`). It belongs to the `/view/:id` widget system and reuses the
aggregated metrics API, KPI toolbar filters, shared metric preset form and migrated
number metric rendering partials.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` Step 1: widget form/query helper registration.
- `[x]` Step 2: complete settings sidebar.
- `[x]` Step 3: runtime metrics and export.
- `[x]` Step 4: view integration.
- `[x]` Step 5: verification.

## Step 0: Inventory And Dependency Audit

- `[x]` Source runtime is Options API and must stay Options API.
- `[x]` Source settings sidebar uses the widget settings mixin; target settings sidebars use `useWidgetSettings()`.
- `[x]` Target already has:
  - `src/helpers/entities/widget/forms/numbers-chart.js`
  - `src/components/sidebars/chart/form/numbers-widget-form.vue`
  - `src/components/widgets/chart/partials/numbers-metrics.vue`
  - `src/components/widgets/chart/partials/numbers-metrics-item.vue`
- `[x]` Target is missing:
  - `src/components/widgets/chart/numbers-widget.vue`
  - `src/components/sidebars/chart/numbers.vue`
  - Numbers chart permission mixins
  - central form/helper/sidebar/wrapper/menu registrations
- `[x]` Target metrics store exposes aggregated metric fetch and KPI aggregate export actions.
- `[x]` Shared chart metric preset fields are available from the chart migrations.

## Step 1: Helpers

- `[x]` Register Numbers parameter form conversion in `src/helpers/entities/widget/form.js`.
- `[x]` Reuse existing Numbers query conversion in `src/helpers/entities/widget/query.js`.
- `[x]` Confirm user preference query conversion is registered for Numbers.

## Step 2: Settings Sidebar

- `[x]` Add `src/components/sidebars/chart/numbers.vue`.
- `[x]` Ensure `src/components/sidebars/chart/form/numbers-widget-form.vue` supports preset selection.
- `[x]` Preserve permission behavior for chart filters, interval and sampling.

## Step 3: Runtime

- `[x]` Add `src/components/widgets/chart/numbers-widget.vue`.
- `[x]` Preserve source Options API style.
- `[x]` Fetch aggregated alarm metrics by widget id and current KPI query.
- `[x]` Apply metric labels from widget metric presets.
- `[x]` Preserve toolbar controls for interval, sampling and user filters.
- `[x]` Preserve automatic value font sizing with `ResizeObserver`.
- `[x]` Preserve CSV export action.

## Step 4: View Integration

- `[x]` Register `NumbersWidget` in `widget-wrapper.vue`.
- `[x]` Register `NumbersSettings` in `sidebars/index.js`.
- `[x]` Enable edit/duplicate actions in `widget-wrapper-menu.vue`.
- `[x]` Update migration inventory and view migration docs.

## Step 5: Verification

- `[x]` Targeted ESLint passes for Numbers files.
- `[x]` Production build passes.
- `[x]` Ready For Testing entry includes:
  - create Numbers widget
  - edit and save Numbers settings
  - filter by KPI interval/sampling/user filter
  - render metrics and trends
  - export CSV
