# Pbehavior Planning Calendar Migration Plan

Temporary working plan for migrating the pbehavior planning calendar from `canopsis-next` to Vue 3 / Vuetify 4. Remove this file after the feature is fully migrated and verified.

## Scope

- Source:
  - `../canopsis-next/src/components/common/calendar/**`
  - `../canopsis-next/src/components/other/pbehavior/calendar/**`
  - `../canopsis-next/src/components/modals/pbehavior/pbehavior-planning.vue`
  - `../canopsis-next/src/components/modals/pbehavior/create-pbehavior.vue`
  - `../canopsis-next/src/components/modals/pbehavior/create-recurrence-rule.vue`
  - `../canopsis-next/src/components/modals/pbehavior/pbehavior-recurrence-rule.vue`
  - `../canopsis-next/src/components/modals/pbehavior/pbehavior-recurrent-changes-confirmation.vue`
- Target:
  - `src/components/common/calendar/**`
  - `src/components/other/pbehavior/calendar/**`
  - `src/components/modals/pbehavior/**`
  - related stores/hooks/helpers/constants/i18n/barrels needed by those components

## Migration Principles

- Keep the external `c-calendar` API compatible with the old consumers.
- Keep Options API components as Options API.
- Preserve old behavior first; adapt only incompatible Vue 2 / Vuetify 2 APIs.
- Port all companion helpers, child components, styles and i18n keys used by the migrated calendar.
- Keep the pbehavior page create/edit/duplicate behavior aligned with the old page: list actions open the planning calendar with the fetched full pbehavior object.

## Phase 1: Inventory And Compatibility Contract

- [x] Compare old `c-calendar` public props, emits, slots and ref methods with Vuetify 4 `VCalendar`.
- [x] Define the target `c-calendar` compatibility contract:
  - props: `events`, `loading`, `readonly`, `hideDetailsMenu`, `type`, `timezone`, `noTimezone`
  - slots: `event`, `form-event`, `loader`, `menu-right`
  - emits: `update:timezone`, `change:event`, `move:event`, `resize:event`, `change:pagination`
  - ref/data API: `focus`, `filled`, `setFocusDate`, `clearPlaceholder`, `prev`, `next`
- [x] Document any Vuetify 4 calendar gaps that require local compatibility code.

## Phase 2: Common Calendar Adapter

- [x] Add `src/components/common/calendar/c-calendar.vue`.
- [x] Port calendar toolbar partials:
  - `calendar-today-btn.vue`
  - `calendar-view-mode.vue`
  - `calendar-pagination.vue`
  - `calendar-pagination-arrow-btn.vue`
  - `calendar-period-picker.vue`
- [x] Adapt old layout to `c-layout` / Vuetify 4 markup.
- [x] Render Vuetify 4 `v-calendar` with old event data shape.
- [x] Preserve event slot output and readable icon/text colors.
- [x] Recreate filled range tracking and emit `change:pagination`.
- [x] Recreate current-time line behavior.

## Phase 3: Calendar Interaction Layer

- [x] Recreate popover positioning behavior using Vuetify 4 menu `target`.
- [x] Recreate create placeholder event behavior.
- [x] Recreate event move handling and emit `move:event`.
- [x] Recreate event resize handling and emit `resize:event`.
- [x] Preserve `hideDetailsMenu`, `readonly`, `dragging`, `resizing`, and `creating` states.

## Phase 4: Pbehavior Planning Calendar Dependencies

- [x] Port `helpers/calendar/calendar.js` with JSDoc.
- [x] Ensure pbehavior form calendar helpers are present:
  - `calendarEventToPbehaviorForm`
  - `formToCalendarEvent`
  - `isFullDayEvent`
- [x] Ensure pbehavior timespan helper exists:
  - `pbehaviorToTimespanRequest`
- [x] Ensure pbehavior timespan and type stores/hooks expose methods used by the calendar.
- [x] Port `pbehavior-planning-calendar-legend.vue`.
- [x] Port `pbehavior-create-event.vue` and its companion dependencies.

## Phase 5: Pbehavior Planning Calendar Component

- [x] Add `src/components/other/pbehavior/calendar/pbehavior-planning-calendar.vue`.
- [x] Replace Vuex mixins/actions with Pinia hooks.
- [x] Preserve props and `update:*` emits for:
  - `pbehaviorsById`
  - `addedPbehaviorsById`
  - `removedPbehaviorsById`
  - `changedPbehaviorsById`
  - `timezone`
- [x] Preserve recurrence change behavior:
  - selected occurrence creates exdate plus standalone pbehavior
  - all occurrences shifts original pbehavior
  - modal confirmation decides selected vs all
- [x] Preserve timezone and `shownUserTimezone` behavior.

## Phase 6: Modals

- [x] Port `pbehavior-recurrent-changes-confirmation.vue`.
- [x] Port `pbehavior-planning.vue`.
- [x] Port and register `create-pbehavior`, `create-recurrence-rule` and `pbehavior-recurrence-rule`.
- [x] Port recurrence rule form and fields used by the editor modal.
- [x] Register pbehavior modals in `src/components/modals/index.js`.
- [x] Wire dialog sizing in the existing modal props map if needed.
- [x] Enable create/edit/duplicate actions on `pbehaviors-page.vue`.

## Phase 7: Verification

- [x] Run targeted ESLint for all migrated calendar files.
- [x] Run `npm run build -- --logLevel error`.
- [x] Manually test `/exploitation/pbehaviors` create action opens planning modal.
- [x] Manually test planning calendar mode switching, pagination, timezone switch, events loading, popover open/close, create placeholder, move and resize.
- [x] Manually test recurrent change confirmation for selected and all occurrences.
- [x] Update `docs/migration-inventory.md`.

## Known Risks

- Vuetify 4 `VCalendar` does not preserve the old Vuetify 2 internal extension surface.
- The old custom timed event height patch cannot be copied with `extends: VCalendar`; it must be reproduced in the wrapper or with CSS/slot rendering.
- Drag/resize depends on DOM geometry and may need browser verification.
- `pbehavior-create-event.vue` may pull in form dependencies that belong to the deferred forms pass; only port the minimum needed for the modal to open and preserve behavior.
