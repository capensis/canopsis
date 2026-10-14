# JUnit Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/testing-weather/testing-weather.vue`

Target runtime: `src/components/widgets/testing-weather/testing-weather.vue`

The widget is the JUnit testing weather widget (`WIDGET_TYPES.testingWeather`, component
name `testing-weather-widget`). It is part of the `/view/:id` widget system and displays
test suite cards with history/details drilldowns, alarm-list handoff and storage-backed
settings for reports, screenshots and videos.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` Step 1: widget form helpers and default preparation.
- `[x]` Step 2: settings sidebar.
- `[x]` Step 3: runtime card grid.
- `[x]` Step 4: test-suite modal and details stack.
- `[x]` Step 5: view integration.
- `[x]` Step 6: Ready For Testing update.

## Step 0: Inventory And Dependency Audit

### Source runtime

- `src/components/widgets/testing-weather/testing-weather.vue`
- `src/components/widgets/testing-weather/testing-weather-item.vue`
- `src/components/widgets/testing-weather/test-suite-statistics.vue`
- `src/components/widgets/testing-weather/test-suite-statistic-row.vue`

### Source settings

- `src/components/sidebars/testing-weather/testing-weather.vue`
- `src/components/sidebars/testing-weather/form/fields/storage.vue`
- `src/components/sidebars/testing-weather/form/fields/storages.vue`
- `src/components/sidebars/testing-weather/form/fields/file-name-masks.vue`

### Source modal and test-suite details

- `src/components/modals/test-suite/test-suite.vue`
- `src/components/other/test-suite/test-suite-history.vue`
- `src/components/other/test-suite/partials/test-suite-tabs.vue`
- `src/components/other/test-suite/partials/test-suite-summary-tab.vue`
- `src/components/other/test-suite/partials/test-suite-messages-tab.vue`
- `src/components/other/test-suite/partials/test-suite-gantt-tab.vue`
- `src/components/other/test-suite/partials/test-suite-details-tab.vue`
- `src/components/other/test-suite/partials/test-suite-details.vue`
- `src/components/other/test-suite/partials/test-suite-details-expand-panel.vue`
- `src/components/other/test-suite/partials/test-suite-screenshots-tab.vue`
- `src/components/other/test-suite/partials/test-suite-screenshots.vue`
- `src/components/other/test-suite/partials/test-suite-videos-tab.vue`
- `src/components/other/test-suite/partials/test-suite-videos.vue`
- `src/components/other/test-suite/partials/test-suite-history-panel.vue`
- `src/components/other/test-suite/partials/test-suite-message-panel.vue`
- `src/components/other/test-suite/partials/system-message.vue`
- `src/components/other/test-suite/partials/test-suite-historical-data-months-field.vue`
- `src/components/other/test-suite/partials/test-suite-status-pie-chart.vue`
- `src/components/other/test-suite/partials/test-suite-summary-status-row.vue`
- `src/components/common/chips/c-test-suite-chip.vue`

### Source helpers and stores

- `src/helpers/entities/widget/forms/junit.js`
- `src/helpers/entities/junit/url.js`
- `src/store/modules/entities/test-suite/index.js`
- `src/store/modules/entities/test-suite/history/index.js`
- `src/store/modules/entities/test-suite/entity-gantt/index.js`
- `src/mixins/entities/test-suite.js`
- `src/mixins/entities/test-suite-history.js`

### Target dependencies already present

- `src/constants/junit.js`
- `src/stores/modules/test-suite.js`
- `src/hooks/store/modules/test-suite.js`
- `src/components/common/chart/junit-gantt-chart.vue`
- `src/components/widgets/alarm/entity-gantt/entity-gantt.vue`
- `src/i18n/messages/{en,fr}/test-suite.js`
- Storage setting JUnit form and state-setting JUnit modal stacks are already migrated.

## Step 1: Widget Form Helpers And Defaults

- `[x]` Port `testingWeatherWidgetParametersToForm`.
- `[x]` Port `formToTestingWeatherWidgetParameters`.
- `[x]` Skip `prepareTestingWeatherWidget`; runtime reads parameters as stored, like Map/Text.
- `[x]` Wire the helpers in `src/helpers/entities/widget/form.js`.
- `[x]` Keep storage arrays as keyed form rows and submit them as plain directory strings.
- `[x]` Preserve legacy defaults:
  - `periodic_refresh` from `DEFAULT_PERIODIC_REFRESH`.
  - `directory` as an empty string.
  - `is_api` as `false`.
  - `screenshot_filemask` and `video_filemask` from `TEST_CASE_FILE_MASK`.
  - `report_fileregexp` as an empty string.

Implemented files:

- `src/helpers/entities/widget/forms/junit.js`
- `src/helpers/entities/junit/url.js`
- `src/helpers/entities/widget/form.js`
- `src/mixins/entities/test-suite.js`
- `src/mixins/entities/test-suite-history.js`

## Step 2: Settings Sidebar

- `[x]` Add `src/components/sidebars/testing-weather/testing-weather.vue` with `useWidgetSettings()`.
- `[x]` Port title and periodic refresh fields using the existing sidebar field components.
- `[x]` Port result directory field with create/edit/remove actions.
- `[x]` Port `receiveByApi` switch and disable filesystem directory fields when enabled.
- `[x]` Port screenshot and video directory lists.
- `[x]` Port screenshot, video and report filename mask fields.
- `[x]` Use `MODALS.textFieldEditor` to edit directories and validate with `useTestSuite().validateTestSuiteDirectory`.
- `[x]` Register `testing-weather-settings` in `src/components/sidebars/index.js`.
- `[x]` Add `WIDGET_TYPES.testingWeather` to `MIGRATED_SETTINGS_WIDGET_TYPES` after the sidebar is complete.

Implemented files:

- `src/components/sidebars/testing-weather/testing-weather.vue`
- `src/components/sidebars/testing-weather/form/fields/storage.vue`
- `src/components/sidebars/testing-weather/form/fields/storages.vue`
- `src/components/sidebars/testing-weather/form/fields/file-name-masks.vue`
- `src/components/forms/fields/c-storage-field.vue`
- `src/components/forms/fields/c-storages-field.vue`
- `src/components/forms/fields/c-file-name-mask-field.vue`
- `src/components/sidebars/index.js`

## Step 3: Runtime Card Grid

- `[x]` Add `src/components/widgets/testing-weather/testing-weather.vue`.
- `[x]` Add `testing-weather-item.vue`, `test-suite-statistics.vue` and `test-suite-statistic-row.vue`.
- `[x]` Fetch test suites through `fetchTestSuitesList` with `widgetId` and `limit: MAX_LIMIT`.
- `[x]` Read widget-scoped list and pending state from `getListByWidgetId` and `getPendingByWidgetId`.
- `[x]` Preserve periodic refresh behavior with `widgetPeriodicRefreshMixin` and `widgetFetchQueryMixin`.
- `[x]` Render the loading progress, responsive card grid and `common.noData` state.
- `[x]` Preserve status colors from `TEST_SUITE_COLORS[testSuite.state]`.
- `[x]` Port the mini bar chart timestamp block and fixed test statistics.
- `[x]` Open `MODALS.testSuite` on card click.
- `[x]` Open `MODALS.alarmsList` from the card action with `generatePreparedDefaultAlarmListWidget()` and component alarms fetch by `testSuite.entity_id`.
- `[x]` Guard the alarm-list action with `USER_PERMISSIONS.business.testingWeather.actions.alarmsList`.

Implemented files:

- `src/components/widgets/testing-weather/testing-weather.vue`
- `src/components/widgets/testing-weather/testing-weather-item.vue`
- `src/components/widgets/testing-weather/test-suite-statistics.vue`
- `src/components/widgets/testing-weather/test-suite-statistic-row.vue`
- `src/components/common/chart/c-mini-bar-chart.vue`

## Step 4: Test-Suite Modal And Details Stack

- `[x]` Add `src/components/modals/test-suite/test-suite.vue` and register `MODALS.testSuite`.
- `[x]` Fetch history summaries with `fetchHistoryList` and `limit: MAX_LIMIT`.
- `[x]` Port `test-suite-history.vue` and `test-suite-history-panel.vue`.
- `[x]` Port `test-suite-tabs.vue`; fetch each selected summary with `fetchSummaryWithoutStore`.
- `[x]` Port the summary tab and status pie chart.
- `[x]` Port global message panels with copy-to-clipboard popup behavior.
- `[x]` Port Gantt tab using the existing `junit-gantt-chart.vue`, `fetchItemGanttIntervalsWithoutStore`, months selector and total-time interval.
- `[x]` Port details tab/table, row expand panel and pagination.
- `[x]` Port screenshot and video tabs, including file URL generation from `helpers/entities/junit/url.js`.
- `[x]` Port `c-test-suite-chip.vue` and `c-download-btn.vue`.

Implemented files:

- `src/components/modals/test-suite/test-suite.vue`
- `src/components/other/test-suite/test-suite-history.vue`
- `src/components/other/test-suite/partials/*`
- `src/components/common/chips/c-test-suite-chip.vue`
- `src/components/common/buttons/c-download-btn.vue`
- `src/helpers/file/files.js` (`saveTextFile`)
- `src/components/modals/index.js`

## Step 5: View Integration

- `[x]` Register the runtime component in `src/components/widgets/widget-wrapper.vue`.
- `[x]` No prepare helper; widget is rendered as stored.
- `[x]` Widget create/edit modal already lists JUnit via `TOP_LEVEL_WIDGET_TYPES` and `modals.createWidget.types[WIDGET_TYPES.testingWeather]`.
- `[x]` `SIDE_BARS_BY_WIDGET_TYPES[WIDGET_TYPES.testingWeather]` already points to `SIDE_BARS.testingWeatherSettings`.
- `[x]` Keep old widget type/component name compatibility: `Junit` / `testing-weather-widget`.

Implemented files:

- `src/components/widgets/widget-wrapper.vue`
- `src/components/widgets/partials/widget-wrapper-menu.vue`

## Step 6: Verification

Verify create, edit, duplicate, delete, save/reopen settings, directory validation, API mode
disabled states, periodic refresh, card rendering, mini chart, modal opening, history loading,
summary/details/messages/gantt/screenshots/videos tabs, alarm-list handoff and permission
gating.

Build/static verification:

- `[x]` `npx eslint` on implemented JUnit widget, sidebar, modal and helper files.
- `[x]` `npm run build`.
- `[x]` Ready For Testing entry is updated.

The migration is complete only when the widget has no placeholder branch, Vue 2 listener/slot
API, Vuex mixin dependency, `$validator` dependency or unregistered `c-*` components.

## Documentation Rule

After every completed step, update this file and `docs/migration-inventory.md` in the same
change. Record migrated files, verification results and remaining risks.
