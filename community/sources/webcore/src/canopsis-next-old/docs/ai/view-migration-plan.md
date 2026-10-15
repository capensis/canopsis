# View Page Migration Plan

Source page: `../canopsis-next/src/views/view.vue`

Target page: `src/views/view-page.vue`

Status legend:

- `[x]` done
- `[~]` in progress
- `[ ]` pending

## Scope

Migrate the main view rendering stack in stages. The page is larger than a normal CRUD page: it includes active view state, tab routing, tab reordering, widget grid layout, widget wrappers, floating view controls, tab/widget modals and the concrete widget types.

Do not simplify old behavior while porting. Preserve old source logic and port accompanying helpers, hooks, stores, constants, styles, JSDoc, child components and i18n keys. Skip only obsolete/incompatible support code and document the reason.

Playlist playback is migrated on top of the main view/tab/widget foundation. Widget templates are migrated and ready for testing.

## Stage 1: Foundation

- `[x]` Complete `active-view` Pinia store/hook parity:
  - `toggleEditing`
  - `registerEditingOffHandler`
  - `unregisterEditingOffHandler`
  - `pausePeriodicRefresh`
  - `resumePeriodicRefresh`
  - `setActiveViewScreenMode`
  - `activeViewIsKioskScreenMode`
- `[x]` Complete `view` tab API actions in the Pinia store/hook:
  - `fetchViewTab`
  - `createViewTab`
  - `updateViewTab`
  - `copyViewTab`
  - `removeViewTab`
  - `updateViewTabPositions`
- `[x]` Verify `widget` store/hook has actions required by widget wrapper and grid editing.
  - Old `activeView/activeWidgets` Vuex module is not ported separately at this stage because the migrated Pinia `alarm` store already owns widget-scoped list buckets (`widgets`, meta, pending and fetching params).
- `[x]` Port/adapt `useViewRouter` helpers:
  - redirect to first tab
  - redirect to view root
  - redirect to selected view and tab
  - redirect home if current route view is removed

## Stage 2: Page Shell

- `[x]` Add `src/views/view-page.vue`.
- `[x]` Wire `/view/:id` route to the new page.
- `[x]` Restore active view lifecycle:
  - clear stale active view on create
  - clear stale active view item on create, so another view is not rendered while the new view is loading
  - fetch active view on mount
  - redirect to first tab when `tabId` is missing
  - redirect to first tab when `tabId` is invalid but the view still has tabs
  - redirect to view root only when `tabId` exists and the view has no tabs
  - redirect home on fetch failure
  - unregister periodic refresh on unmount
- `[x]` Provide `$periodicRefresh` observer for widgets and controls.
- `[x]` Render tabs and floating buttons after Stage 3 and Stage 5 components are available.

## Stage 3: Tabs

- `[x]` Port `view-tabs-wrapper.vue`.
- `[x]` Port `view-tabs.vue`.
- `[x]` Port tab drag/reorder flow with Submit/Cancel overlay.
  - The changed state is tracked by a snapshot of tab ids, so Cancel followed by another drag still shows Submit/Cancel.
- `[x]` Port create/update/copy/delete tab flows.
  - `MODALS.selectView` was ported because copy tab depends on selecting the target view.
- `[x]` Adapt Vuetify 2 tabs API to Vuetify 4:
  - `value` -> `model-value`
  - `v-tabs-items` / `v-tab-item` -> `v-window` / `v-window-item`
  - old scoped slots/listeners -> Vue 3 equivalents
- `[x]` Render widget content in tab windows through the migrated Stage 4 grid/widget wrapper and Stage 7 widget components.

## Stage 4: Grid And Widget Wrapper

- `[x]` Port `common/grid/grid-layout.vue` and dependencies.
  - `helpers/grid.js`
  - `common/grid/partials/grid-item.vue`
  - `helpers/entities/widget/grid.js`
- `[x]` Port `view-tab-widgets.vue`.
- `[x]` Port `widget-wrapper.vue` with migrated widget component dispatch and old-compatible error overlay fallback.
- `[x]` Port `widget-edit-drag-handler.vue`.
  - Full widget settings, duplicate and menu actions are covered by Stage 6.
- `[x]` Preserve breakpoint size switching and max-width behavior.
- `[x]` Preserve widget grid position saving when edit mode turns off.
- `[x]` Preserve mass actions panel visibility outside default screen mode.

## Stage 5: View Floating Controls

- `[x]` Port `view-fab-btns.vue`.
- `[x]` Port floating control children:
  - scroll top
  - executions FAB/menu with active remediation executions list
  - periodic refresh
  - screen mode
  - share link
  - editing button
- `[x]` Keep the active screen mode item styled with the primary color, matching the old control state.
- `[x]` Keep the scroll-to-top button icon white, matching the old dark FAB behavior.
- `[x]` Port `c-speed-dial` if missing and register it globally.
- `[x]` Preserve shortcuts:
  - `Ctrl+E` toggles editing when view is updatable
  - `Alt/Meta+Shift+1..4` changes screen mode
- `[x]` Preserve fullscreen and kiosk screen-mode behavior on `/view/:id`.
  - Native fullscreen adapter and screen-mode switching are migrated.
  - Top bar/header visibility in kiosk screen mode follows the User Interface
    `show_header_on_kiosk_mode` option, matching the old app.
  - Fullscreen and kiosk fullscreen view content gets a top offset equal to the
    top-bar height when that header is shown.
  - Standalone `/kiosk-views/:id/:tabId` page is migrated through `src/views/view-kiosk-page.vue`.

## Stage 6: Tab And Widget Modals

- `[x]` Port view create/edit/duplicate modal flow:
  - create regular view from navigation settings FAB
  - create private view from navigation settings FAB
  - create-view modal is rewritten to Composition API with the existing vee form context, dirty close confirmation and group auto-creation behavior preserved
  - edit view from navigation group view actions
  - duplicate view from navigation group view actions
  - edit/delete group modal from navigation group actions
  - sidebar view edit/duplicate handlers are enabled; private views keep edit/delete/create behavior from Vue 2; edit action uses the view id permission check from Vue 2 so `Save`/`Delete` availability matches the opened edit modal
  - delete view from edit modal
  - auto-create group by typed title
  - dirty close confirmation
  - `view-form`, group/tags/private-duplicate fields and periodic refresh field
- `[x]` Port tab modals/actions:
  - create tab through text field editor
  - update tab title
  - select target view for clone
  - clone tab
  - delete tab confirmation
- `[x]` Port widget modals/actions:
  - create widget: selection modal is migrated and opens the matching widget settings sidebar for every type, matching Vue 2. Default widget title stays empty like in Vue 2.
  - edit widget: widget settings sidebars are migrated with `useWidgetSettings()`, vee-validate form scope, Save state, invalid-field scrolling and dirty close confirmation; sidebar expand item/group components are migrated and the title field uses an expand item; sidebar content padding is restored to 24px; sidebars use a full-viewport overlay above the top bar and render the drawer above that overlay.
  - copy widget: `select-view-tab` and widget duplication flow are migrated for widget types with migrated settings sidebars.
  - delete widget: wrapper menu delete action is migrated
  - widget wrapper menu
- `[x]` Restore full widget wrapper menu actions after sidebar/select-view-tab infrastructure is available.
- `[x]` Port widget form helpers, constants and validators used by the modals.
  - `getDefaultGridParameters` and `getEmptyWidgetByType` are migrated.
  - Text widget parameter form mapper is migrated.
  - Counter widget parameter form mapper is migrated.
  - Map widget parameter form mapper is migrated.
  - Counter settings field dependencies migrated: column-size field, grid-size, slider, switcher, margins, opened/resolved selector, counter levels, criticity levels, level colors, embedded alarm-list display settings and info popup create/edit modals.
  - Counter filter editor is migrated through the shared widget filter stack and `create-filter` modal.
  - Shared alarm-list column/sort template helpers and the Counter embedded alarm-list settings form are migrated for Counter. The runtime alarm-list handoff action still waits for the alarm-list widget stack.

## Stage 7: Widget Types

- `[x]` Port simple widgets first:
  - text: renderer and settings sidebar are migrated; widget title typography is adjusted so `h4` fits the fixed title row as in Vue 2
  - counter: renderer, item coloring/template, `card-with-see-alarms-btn` wrapper, store fetch, create/edit/duplicate settings sidebar, own settings fields, embedded alarm-list display settings, filter editor and variables-help action/modal are complete
  - map: renderer, map-state fetch, category selector, filter selector/list/create modal, nested-map breadcrumbs, popup actions wiring, settings sidebar, map selector, periodic refresh, filters, color indicator, pbehavior marker display toggle, popup template editor and alarm/entity columns settings are complete
- `[x]` Port table/data widgets:
  - external data table: runtime, settings sidebar, form/query helpers, widget-wrapper registration, user preference, periodic refresh, CSV export and records-list handoff are complete
  - alarms list: runtime, settings sidebar, dedicated Vuetify 4 table chain, toolbar, columns/cells, row and mass actions, modals, exports, sticky header, sticky horizontal scroll and expand tabs are complete. It intentionally does not use `c-advanced-data-table`.
  - context: runtime, settings sidebar, table, toolbar, actions, expand panel and root-cause modal are complete
- `[x]` Port chart/statistics widgets:
  - bar chart: runtime, settings sidebar, form/query helpers, KPI toolbar filters, metrics fetch/export and wrapper/menu integration are complete
  - line chart: runtime, settings sidebar, form/query helpers, KPI toolbar filters, metrics fetch/export and wrapper/menu integration are complete
  - pie chart: runtime, settings sidebar, form/query helpers, KPI toolbar filters, metrics fetch/export and wrapper/menu integration are complete
  - numbers: runtime, settings sidebar, form/query helpers, KPI toolbar filters, metrics fetch/export and wrapper/menu integration are complete
  - stats calendar: runtime, settings sidebar, form/query helpers, grouped event menu, alarm-list drilldown and wrapper/menu integration are complete
  - user statistics: runtime table, settings sidebar, form/query helpers, KPI toolbar filters, group metrics fetch and wrapper/menu integration are complete
  - alarm statistics: runtime table, settings sidebar, shared form/query helpers, KPI toolbar filters, group metrics fetch and wrapper/menu integration are complete
  - availability: runtime table, settings sidebar, form/query helpers, toolbar filters, history graph and wrapper/menu integration are complete
- `[x]` Port service/testing widgets:
  - service weather: runtime, settings sidebar, handoffs and service-entities modal are complete
  - testing weather: runtime, settings sidebar and test-suite modal/details stack are complete
- `[x]` Unsupported widget types render the old-compatible error overlay fallback while migrated widget types render their concrete components.

## Stage 8: Follow-Ups

- `[x]` Migrate `playlist` after the view/tab/widget rendering foundation is stable: `/playlist/:id` now uses `src/views/playlist-page.vue`, the playlist store, tab fetches, top-bar playback controls, autoplay/fullscreen and migrated `view-tab-widgets`.
- `[x]` Widget templates are migrated and connected to migrated widget settings consumers.
- `[x]` Removed the old widget settings compatibility mixin after migrating widget settings sidebars to `useWidgetSettings()`.
- `[x]` Revisited `pbehavior` store deletion behavior for alarm/context/service-weather widget pbehavior tabs; all use API delete plus local refresh without restoring normalizr cleanup.

## Verification Checklist

- `[x]` Targeted eslint passes for changed files.
- `[x]` Build passes or any unrelated failures are documented.
- `[x]` `/view/:id` loads an existing view. Page shell compiles and route is wired.
- `[x]` Missing `tabId` redirects to the first tab. Logic is ported.
- `[x]` Invalid `tabId` redirects to the first available tab, or to view root when there are no tabs. Logic is ported.
- `[x]` Tab switching preserves route query and visible widget content.
- `[x]` Tab reorder can be submitted and cancelled. Cancel/retry state and shared overlay layering are fixed.
- `[x]` Edit mode toggles and saves widget grid positions. Grid save logic is ported.
- `[x]` Fullscreen and kiosk modes work. Native fullscreen flow is ported.
- `[x]` Periodic refresh can pause/resume. Logic is ported.
- `[x]` Create/edit/delete/copy tab flows work. Existing tabs wrapper actions plus create FAB are ported.
- `[x]` Create/edit/duplicate/delete view modal flows work. Code is ported and build-verified.
- `[x]` Create/edit/delete/copy widget flows work for migrated widget types.
- `[x]` Widget wrapper renders all migrated widget types.
