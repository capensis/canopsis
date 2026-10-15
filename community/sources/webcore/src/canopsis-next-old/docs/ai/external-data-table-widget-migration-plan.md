# External Data Table Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/external-data-table/external-data-table-widget.vue`

Target runtime: `src/components/widgets/external-data-table/external-data-table-widget.vue`

The widget is the External Data Table viewer (`WIDGET_TYPES.externalDataTable`,
component name `external-data-table-widget`). It is part of the `/view/:id`
widget system. Most admin/supporting infrastructure was already migrated before
this widget slice: Pinia store, record/import modals, `external-data-table-records-list`
and the admin page.

The widget is thin orchestration over the shared records list. It uses
`c-advanced-data-table` through that list, not a dedicated table chain.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` Widget form helpers and query converters.
- `[x]` Complete settings sidebar.
- `[x]` Shared hooks: `useWidgetUserPreference`, `useWidgetPeriodicRefresh`, `useFilePollingWithPending`.
- `[x]` Widget runtime and widget-specific hooks.
- `[x]` View integration: `widget-wrapper`, settings sidebar registry, wrapper menu.
- `[x]` Final: Ready For Testing update.

## Step 1: Widget Form Helpers And Query Layer

- `[x]` Port `helpers/entities/widget/forms/external-data-table.js`.
- `[x]` Port `helpers/entities/external-data-table/query.js`.
- `[x]` Wire into `helpers/entities/widget/form.js` and `helpers/entities/widget/query.js`.

## Step 2: Settings Sidebar

- `[x]` Port `sidebars/external-data-table/external-data-table.vue` with `useWidgetSettings()`.
- `[x]` Port `field-columns-without-template`.
- `[x]` Reuse table selector, periodic refresh, default sort, density, resize behavior, export CSV form.
- `[x]` Export from `components/sidebars/index.js`.

## Step 3: Shared Hooks

- `[x]` Port `useWidgetUserPreference` into `hooks/store/modules/user-preference.js`.
- `[x]` Port `hooks/periodic-refresh.js` and `hooks/widget/periodic-refresh.js`.
- `[x]` Port `useFilePollingWithPending` into `hooks/polling.js`.

## Step 4: Widget Runtime

- `[x]` Port `hooks/external-data-table-widget.js`.
- `[x]` Port `external-data-table-widget.vue`.
- `[x]` Reuse `useExternalDataTableRecordsList` and `external-data-table-records-list`.

## Step 5: View Integration

- `[x]` Register in `widget-wrapper.vue`.
- `[x]` Enable edit/duplicate in `widget-wrapper-menu.vue`.

## Step 6: Verification

Verify table selection, column auto-population, settings save/reopen, empty table
import CTA, record CRUD, bulk delete, column drag/resize persistence, CSV export,
periodic refresh, density toggle and view action gating.
