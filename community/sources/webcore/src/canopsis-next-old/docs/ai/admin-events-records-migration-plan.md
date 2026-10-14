# Admin Events Records Migration Plan

Temporary working plan for migrating `admin-events-records` from Vue 2 to Vue 3. Delete this file after the page is fully migrated and `docs/migration-inventory.md` is updated.

## Scope

- Source page: `/Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next/src/views/admin/events-records.vue`
- Target page: `src/views/admin/events-records-page.vue`
- Route: `admin-events-records`, `/admin/events-records`
- Source project is read-only. Edit only `canopsis-next-vue-3`.

## Migration Style

- Preserve source component style by default. Options API stays Options API; Composition API stays Composition API.
- Adapt only what is required by Vue 3, Vuetify 4, Vite, Pinia stores, vue-i18n 11, and existing migrated wrappers.
- Components must use hooks from `src/hooks/store/modules`; hooks wrap Pinia stores; only stores call `request` / API routes.
- Do not recreate the old `events-record-current` hook/store as a separate module. The Vue 3 `events-record` store already owns current status, recordings and resendings.
- Prefer existing wrappers (`c-layout`, `c-advanced-data-table`, `c-action-btn`, `c-copy-btn`, `c-event-filter-patterns-field`, etc.) over raw Vuetify controls.
- Do not start the dev server unless explicitly asked.

## Existing Vue 3 Pieces

Already present in the Vue 3 project:

- Route name/path constants for `admin-events-records`.
- Router entry with correct permission metadata, currently pointing to `PlaceholderPage`.
- Events record permission constants and documentation links.
- `MODALS.eventsRecord` and `MODALS.applyEventFilter` constants.
- `src/constants/events-record.js` with `SINGLE_EVENT_RESEND_DELAY`.
- EN/FR `events-record` messages.
- EN/FR modal and page-header messages for events records.
- `src/stores/modules/events-record.js` Pinia store with:
  - list fetching;
  - events fetching;
  - JSON export creation and polling fetch;
  - playback start/stop;
  - record/event deletion;
  - current status fetching/start/stop;
  - derived `recordings`, `resendings`, `limit`, `recordingsById`, `resendingsById`.
- `src/hooks/store/modules/events-record.js` exposing the Pinia store methods and refs.
- Common JSON treeview component.
- File polling hook used by healthcheck/KPI migrations.

## Source Inventory

### Page

- `src/views/admin/events-records.vue`

### Main List

- `components/other/events-record/events-records-header.vue`
- `components/other/events-record/events-records-list.vue`
- `components/other/events-record/partials/events-record-download-btn.vue`

### Record Details Modal

- `components/modals/events-record/events-record.vue`
- `components/modals/events-record/apply-event-filter.vue`
- `components/other/events-record/events-record-events-header.vue`
- `components/other/events-record/events-record-events-filter.vue`
- `components/other/events-record/events-record-events-list.vue`
- `components/other/events-record/partials/events-record-events-list-expand-panel.vue`

### Hooks And Helpers

- `components/other/events-record/hooks/recording.js`
- `components/other/events-record/hooks/resending.js`
- `components/other/events-record/hooks/current-polling.js`
- `components/other/events-record/hooks/export-json.js`
- `helpers/entities/events-record/url.js`

### Old Store Mapping

- Old `useEventsRecordCurrent()` maps to new `useEventsRecord()`:
  - `recordings` -> `eventsRecordRecordings`
  - `resendings` -> `eventsRecordResendings`
  - `limit` -> `eventsRecordLimit`
  - `recordingsById` -> `eventsRecordRecordingsById`
  - `resendingsById` -> `eventsRecordResendingsById`
  - `fetchEventsRecordCurrent` -> `fetchCurrentEventsRecord`
  - `startEventsRecordCurrent` -> `startCurrentEventsRecord`
  - `stopEventsRecordCurrent` -> `stopCurrentEventsRecord`
- Old `useEventsRecord()` action names map to new names:
  - `playbackEventsRecordEvents` -> `playbackEventsRecord`
  - `stopPlaybackEventsRecordEvents` -> `stopEventsRecordPlayback`
  - `bulkRemoveEventsRecordEvent` -> `bulkRemoveEventsRecordEvents`

## Migration Steps

### 1. Helpers And Hooks

- [x] Add `src/helpers/entities/events-record/url.js` with `getEventsRecordFileUrl()`.
- [x] Port `useEventsRecordExportJson()` and adapt it to the new store action payload names.
- [x] Port `useEventRecordCurrentPolling()` and use `fetchCurrentEventsRecord()` from the merged Vue 3 hook.
- [x] Port `useEventsRecordRecording()`:
  - use `MODALS.applyEventFilter`;
  - use `startCurrentEventsRecord()`;
  - use `stopCurrentEventsRecord()`;
  - preserve the optional event-pattern behavior.
- [x] Port `useEventsRecordResending()`:
  - use `playbackEventsRecord()`;
  - use `stopEventsRecordPlayback()`;
  - refresh current status with `fetchCurrentEventsRecord()`;
  - preserve the one-event immediate resend behavior and duration modal for multiple events.

### 2. Page Shell And Routing

- [x] Create `src/views/admin/events-records-page.vue`.
- [x] Wire `src/router/index.js` from `PlaceholderPage` to the real page.
- [x] Port local state:
  - `eventsRecords`;
  - `meta`;
  - query/options;
  - pending state.
- [x] Fetch list through `fetchEventsRecordsListWithoutStore()`.
- [x] Fetch current status after list fetch and during polling.
- [x] Preserve `prepareEventsRecords()` behavior with Vue 3-safe array replacement/assignment.
- [x] Preserve record removal confirmation.
- [x] Open `MODALS.eventsRecord` with `{ eventsRecord, fetchList, removeAction }`.

### 3. Main List Components

- [x] Port `events-records-header.vue`.
- [x] Port `events-records-list.vue`.
- [x] Adapt `c-advanced-data-table` usage to the current Vue 3 table API.
- [x] Replace Vue 2 filters with date/duration/json helpers or existing migrated global equivalents.
- [x] Adapt Vuetify tooltip activator slots from `{ on }` to `{ props }`.
- [x] Preserve recording/resending visual states:
  - recording progress tooltip;
  - resending blinking icon;
  - disabled view/delete/download actions while recording where the old UI did so.
- [x] Preserve advanced pagination and total item behavior.

### 4. Details Modal

- [x] Add `MODALS.eventsRecord` async component registration.
- [x] Port `components/modals/events-record/events-record.vue`.
- [x] Replace `v-tabs-items` / `v-tab-item` with `v-window` / `v-window-item`.
- [x] Preserve the automatic switch from Events tab to Pattern tab while the record is still recording.
- [x] Preserve title formatting with record start date.
- [x] Preserve record deletion and stop-recording confirmation flows.
- [x] Preserve nested events query:
  - page;
  - limit;
  - `event_pattern` serialized as JSON.

### 5. Event Filter Modal

- [x] Add `MODALS.applyEventFilter` async component registration.
- [x] Port `components/modals/events-record/apply-event-filter.vue`.
- [x] Reuse the existing Vue 3 `c-event-filter-patterns-field` API.
- [x] Keep `required` configurable so event recording can start without a pattern.
- [x] Keep AI chat sidebar integration if the current Vue 3 project still exposes compatible hooks/components.
- [x] Preserve the delayed `afterSubmit()` call after applying the filter because the old backend flow depended on it.

### 6. Nested Events Table

- [x] Port `events-record-events-header.vue`.
- [x] Port `events-record-events-filter.vue`.
- [x] Port `events-record-events-list.vue`.
- [x] Port `events-record-events-list-expand-panel.vue`.
- [x] Preserve mass actions:
  - resend selected events;
  - export selected events;
  - remove selected events.
- [x] Preserve row actions:
  - resend one event;
  - copy event JSON;
  - export one event;
  - remove one event.
- [x] Preserve expanded row JSON preview with base64 decoding fallback.
- [x] Adapt table selection slots/events to the current `c-advanced-data-table` implementation.

### 7. Export And Download

- [x] Port `events-record-download-btn.vue`.
- [x] Use `useFilePolling()` for export generation.
- [x] Open the generated file with `getEventsRecordFileUrl(fileData._id ?? fileData.id)` if backend response shape differs.
- [x] Support full record export, one-event export, selected-events export, and filtered export.

### 8. Styling And Vuetify 4 Adaptation

- [x] Replace old Vuetify 2 utility classes where necessary:
  - `grey--text text--darken-1`;
  - `blue darken-3`;
  - `outlined` button prop;
  - tooltip positioning props.
- [x] Keep expand-panel background behavior consistent with already migrated planning/table components:
  - secondary wrapper background outside;
  - card background inside.
- [x] Avoid hard-coded table fixes unless the current `c-advanced-data-table` requires them.

### 9. Verification

- [x] Run targeted eslint for migrated files.
- [x] Run `npm run build`.
- [x] Runtime check `/admin/events-records`:
  - list loads and paginates correctly;
  - current recording/resending status updates without losing table expansion state;
  - launch recording opens the filter modal and allows optional pattern;
  - stop recording works;
  - details modal opens for finished records;
  - Events tab filter modal applies/resets filters;
  - Pattern tab renders readonly event pattern;
  - selected and single event resend works;
  - record/event export opens generated JSON download URL;
  - delete record/event/selected events works;
  - console has no Vue/Vuetify/i18n errors.

## Open Checks Before Implementation

- [x] Confirm whether AI chat sidebar dependencies are already fully migrated. If not, keep the modal functional without breaking the sidebar path and document the gap.
- [x] Confirm the current `c-advanced-data-table` selection slot names and selected key payload.
- [x] Confirm date/duration/json formatting helpers used by recently migrated components.
- [x] Confirm whether `bulkRemoveEventsRecordEvents()` should receive an array directly or `{ data: array }` in all call sites.
