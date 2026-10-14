# Service Weather Widget Migration Plan

**Project:** `canopsis-next-vue-3` (Vue 3 port)

**Legacy source:** `canopsis-next/src/components/widgets/service-weather/`

**Full dual-track doc:** `canopsis-next/.cursor/rules/service-weather-migration.mdc`

**Active track:** Vue 3 port — Phase 5 next

---

## Progress Summary

| Phase | Description | Status |
|-------|-------------|--------|
| 0 | Form helpers, Pinia store, permission mixins | ✅ Done |
| 1 | Settings sidebar + service-weather form fields | ✅ Done |
| 2 | Widget runtime (tile grid + toolbar) | ✅ Done |
| 3 | Alarm list / root cause / variables help handoffs | ✅ Done |
| 4 | Service-entities modal stack | ✅ Done |
| 5 | Ready For Testing update | ✅ Done |

**Last verified:** 2026-08-18 (codebase audit)

---

## Phase 0 — Done

- [x] `src/helpers/entities/widget/forms/service-weather.js` — form mappers + `prepareServiceWeatherWidget`
- [x] Wired in `src/helpers/entities/widget/form.js`
- [x] `src/stores/modules/weather-service-widget.js`
- [x] `src/hooks/store/modules/weather-service-widget.js`
- [x] `src/mixins/permissions/widgets/service-weather/filters.js`
- [x] `src/mixins/permissions/widgets/service-weather/category.js`

---

## Phase 1 — Done

- [x] `src/components/sidebars/service-weather/service-weather.vue`
- [x] Custom settings fields under `src/components/sidebars/service-weather/form/fields/` (counters-selector, sort-column, modal-type, action/no-action-required settings)
- [x] `c-service-weather-icon-field.vue`
- [x] `c-service-weather-state-counter-field.vue`
- [x] `c-enabled-icon-field.vue`
- [x] Registered in `src/components/sidebars/index.js`

---

## Phase 2 — Done

- [x] `src/helpers/entities/service-weather/query.js` — registered in `src/helpers/entities/widget/query.js`
- [x] `src/components/widgets/service-weather/service-weather.vue` — toolbar, fetch, grid, periodic refresh
- [x] `src/components/widgets/service-weather/service-weather-item.vue`
- [x] Partials: `alarm-counter`, `alarm-state-counters`, `alarm-pbehavior-counters`, `impact-state-indicator`
- [x] Registered in `src/components/widgets/widget-wrapper.vue` with `prepareServiceWeatherWidget`

#### Post-Phase 2 fixes

- [x] `service-weather-item__content` / `__main` flex layout — main tile area fills width beside counter columns

---

## Phase 3 — Done

Handoffs wired in widget runtime:

- [x] Alarm list modal via `generatePreparedDefaultAlarmListWidget()` + `fetchServiceAlarmsWithoutStore`
- [x] Root cause via `MODALS.entitiesRootCauseDiagram`
- [x] Variables help via `MODALS.variablesHelp` on tile toolbar (`service-weather-item.vue`)
- [x] Handlebars `tags` helper registered in `service-weather-item.vue`

---

## Phase 4 — Done

### Modal

- [x] `src/components/modals/service/service-entities.vue`
- [x] `src/components/modals/service/create-service-pause-event.vue`
- [x] Registered in `src/components/modals/index.js`

### Partials ported

- [x] `service-template.vue`
- [x] `service-entities-list.vue`
- [x] `service-entity.vue`
- [x] `service-entity-header.vue`
- [x] `service-entity-info-tab.vue`
- [x] `service-entity-template.vue`
- [x] `service-entity-actions.vue`
- [x] `service-entity-alarm-instruction-menu.vue`
- [x] `service-entity-tree-of-dependencies-tab.vue`

### Supporting files

- [x] `src/mixins/widget/actions-panel/common.js`
- [x] `src/mixins/widget/actions-panel/service-entity.js`
- [x] `src/mixins/handlebars/links-helper-creator.js`
- [x] `src/components/widgets/service-weather/service-pause-event-form.vue`
- [x] `src/helpers/entities/entity/actions.js`
- [x] `src/helpers/entities/entity/icons.js`

**Build:** `yarn build` passes (2026-08-18)

#### Post-Phase 4 fixes (2026-08-18)

- [x] Entity expansion panels default collapsed
- [x] Entities list pagination fix
- [x] Hide gray switcher compact density (alarms list correlation style)
- [x] `c-service-weather-patterns-field.vue` → Composition API

---

## Phase 5 — Done

- [x] Add widget to view, configure all settings, save/reopen
- [x] Category / filter / hide-grey user preference
- [x] Periodic refresh
- [x] Tile rendering, counters, modal type routing
- [x] Alarm list / root cause / variables help modals
- [x] Service-entities modal: entity list, actions, pbehaviors tab, comments tab
- [x] Pause/play entity actions
- [x] Mass actions with entitiesActionsInQueue

---

## Track A reference (canopsis-next)

Separate migration: Options API → Composition API within Vue 2.7.

**Status:** nothing migrated yet. See `canopsis-next/.cursor/rules/service-weather-migration.mdc`.
