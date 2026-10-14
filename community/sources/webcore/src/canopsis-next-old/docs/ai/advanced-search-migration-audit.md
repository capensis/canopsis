# Advanced Search Migration Audit

Source runtime: `../canopsis-next/src/components/common/search`

Target runtime: `src/components/common/search`

This file tracks the advanced-search migration as an isolated scope because the
AlarmList toolbar depends on the full old behavior, not only on the visible list
of selectable attributes.

## Rules

- Port the whole advanced-search slice, including nested chips, input behavior,
  history, validators, hooks, list rendering helpers and serialization helpers.
- Do not replace object-valued chips with plain text chips. Duration and range
  chips keep their original model shapes.
- Keep Vue 3 model events explicit: advanced-search rule partials pass
  `modelValue` and listen to `update:modelValue`.
- Adapt only framework-specific APIs: Vuetify activator props, Vue 3 emits,
  and `t` instead of removed `tc`.

## Component Checklist

- `[x]` `c-advanced-search.vue`: wrapper ported with Vue 3 model wiring.
- `[x]` `c-advanced-search-field.vue`: field shell, history actions, submit,
  reset and validator bridge ported.
- `[x]` `advanced-search-rules.vue`: rule collection and input focusing ported.
- `[x]` `advanced-search-rule.vue`: dynamic chip `bind/on` wiring restored and
  switched to `modelValue`.
- `[x]` `advanced-search-chip.vue`: attribute/operator/value list chip restored
  with `variables-list`, disabled message, no-data slot and lazy items.
- `[x]` `advanced-search-duration-chip.vue`: restored as a compound chip with
  `{ value, unit }` model and unit selector.
- `[x]` `advanced-search-range-value-date-chip.vue`: restored as a compound
  from/to date chip with date-time picker menus.
- `[x]` `advanced-search-range-value-period-chip.vue`: restored as a compound
  from/to quick range chip with disabled invalid ranges.
- `[x]` `advanced-search-history-btn.vue`: old-style history dropdown restored.
- `[x]` `advanced-search-history-item-btns.vue`: pinned/remove actions restored.

## Dependencies

- `[x]` `variables-list.vue`: restored nested `c-list` rendering used by
  advanced-search menus.
- `[x]` `c-list` and `v-list-item-mask`: ported and globally registered for the
  variable/search menu tree.
- `[x]` `useQuickDateIntervalRange`: ported for period range chips.
- `[x]` `useDateRangeAllowedDates`: available for date range chips.
- `[x]` `helpers/search/advanced-search.js`: checked for form-to-pattern
  conversion; service-only form fields such as `filled`, `rangeValueDate` and
  `rangeValuePeriod` are stripped before pattern serialization.
- `[x]` `hooks/advanced-search/attributes-map*.js`: restored the full old
  alarm/entity/pbehavior/dynamic-info attribute maps, including store-backed
  value fetching, enum values, aliases, allow/disallow flags and lazy value
  display settings for Connector, Component, Resource, Tags, State, Status,
  Meta, users, entity category/type and pbehavior fields.

## Current Notes

- The old vee-validate `$validator` per-component integration is intentionally
  replaced by the target `useAdvancedSearchValidator` bridge.
- Differences caused only by Vuetify 4 and Vue 3 APIs are expected; simplified
  behavior is not.
