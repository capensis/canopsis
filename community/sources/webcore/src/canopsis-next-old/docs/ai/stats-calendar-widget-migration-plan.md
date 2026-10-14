# StatsCalendar Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/stats/calendar/stats-calendar.vue`

Target runtime: `src/components/widgets/stats/calendar/stats-calendar.vue`

The widget is the alarm statistics calendar (`WIDGET_TYPES.statsCalendar`, component name
`stats-calendar-widget`). It is part of the `/view/:id` widget system and reuses the
existing alarm list endpoint plus the migrated alarm-list modal for drilldown.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` Step 1: widget form helpers and query converter.
- `[x]` Step 2: complete settings sidebar.
- `[x]` Step 3: runtime calendar and alarm-list drilldown.
- `[x]` Step 4: view integration.
- `[x]` Step 5: Ready For Testing update.

## Step 0: Inventory And Dependency Audit

### Source runtime

- `src/components/widgets/stats/calendar/stats-calendar.vue`
- `src/components/widgets/stats/calendar/stats-calendar-menu.vue`

### Source settings

- `src/components/sidebars/stats/stats-calendar.vue`
- `src/components/sidebars/stats/form/stats-calendar-advanced.vue`
- `src/components/sidebars/stats/form/fields/criticity-levels.vue`
- `src/components/sidebars/stats/form/fields/levels-colors-selector.vue`

### Source helpers

- `src/helpers/entities/widget/forms/stats-calendar.js`
- `src/helpers/entities/stats/query.js`

### Target dependencies already present

- `c-calendar` and `helpers/calendar/calendar.js`.
- Pinia alarm store with `fetchListWithoutStore`.
- `generatePreparedDefaultAlarmListWidget()`.
- Alarm-list modal and `alarms-list-modal-form`.
- Shared widget settings hook, filter field, opened/resolved field, switcher field.
- Criticity level and color selector fields.

## Step 1: Widget Form Helpers And Query

- `[x]` Port `statsCalendarWidgetParametersToForm`.
- `[x]` Port `formToStatsCalendarWidgetParameters`.
- `[x]` Port `prepareStatsCalendarWidget`.
- `[x]` Register StatsCalendar in `widgetParametersToForm` and `formToWidgetParameters`.
- `[x]` Port and register `convertStatsCalendarWidgetToQuery`.

Implemented files:

- `src/helpers/entities/widget/forms/stats-calendar.js`
- `src/helpers/entities/stats/query.js`
- `src/helpers/entities/widget/form.js`
- `src/helpers/entities/widget/query.js`

## Step 2: Settings Sidebar

- `[x]` Add `stats-calendar-settings` sidebar with `useWidgetSettings()`.
- `[x]` Keep title, alarm-list modal settings, filters, opened/resolved,
  consider pbehaviors, criticity thresholds and colors.
- `[x]` Fetch infos and widget templates like Counter because the embedded alarm-list
  settings can use template-backed columns.
- `[x]` Register the sidebar in `src/components/sidebars/index.js`.

Implemented files:

- `src/components/sidebars/stats/stats-calendar.vue`
- `src/components/sidebars/index.js`

## Step 3: Runtime Calendar

- `[x]` Add `stats-calendar.vue` runtime component.
- `[x]` Add `stats-calendar-menu.vue` for grouped filter events in month view.
- `[x]` Preserve source component style: runtime and menu stay Options API.
- `[x]` Fetch alarms with `fetchListWithoutStore` for the visible calendar range.
- `[x]` Preserve filter splitting: no filters fetches one collection; configured filters
  fetch one collection per filter and group by filter in the calendar.
- `[x]` Preserve `considerPbehaviors` by adding the active pbehavior pattern exclusion.
- `[x]` Preserve criticity color thresholds and grouped month events.
- `[x]` Open the migrated alarm-list modal for day/hour drilldown.
- `[x]` Do not add periodic refresh; the legacy widget refreshes on calendar/query changes only.

Implemented files:

- `src/components/widgets/stats/calendar/stats-calendar.vue`
- `src/components/widgets/stats/calendar/stats-calendar-menu.vue`

## Step 4: View Integration

- `[x]` Register `stats-calendar-widget` in `widget-wrapper.vue`.
- `[x]` Prepare widget alarm-list columns before rendering.
- `[x]` Enable edit and duplicate in `widget-wrapper-menu.vue`.

Implemented files:

- `src/components/widgets/widget-wrapper.vue`
- `src/components/widgets/partials/widget-wrapper-menu.vue`

## Step 5: Verification

Verify create, edit, duplicate, delete, save/reopen settings, filters, opened/resolved,
consider pbehaviors, criticity colors, month/week/day navigation, grouped month menu,
alarm-list drilldown and API error overlay.

Build/static verification:

- `[x]` `npx eslint` on the implemented StatsCalendar files and touched widget/sidebar helper files.
- `[x]` `npm run build`.
- `[x]` Ready For Testing entry is updated.

The migration is complete only when the widget has no placeholder branch, Vue 2
listener/slot API, `$validator` dependency or unregistered `c-*` components.

## Documentation Rule

After every completed step, update this file and `docs/migration-inventory.md` in the
same change. Record migrated files, verification results and remaining risks.
