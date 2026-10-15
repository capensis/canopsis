# Admin KPI Migration Plan

Temporary working plan for migrating `admin-kpi` from Vue 2 to Vue 3. Delete this file after the page is fully migrated and `docs/migration-inventory.md` is updated.

## Scope

- Source page: `/Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next/src/views/admin/kpi.vue`
- Target page: `src/views/admin/kpi-page.vue`
- Route: `admin-kpi`, `/admin/kpi`
- Source project is read-only. Edit only `canopsis-next-vue-3`.

## Migration Style

- Preserve source component style by default. Options API stays Options API; Composition API stays Composition API.
- Adapt only what is required by Vue 3, Vuetify 4, Vite, Pinia stores, vue-i18n 11, and updated chart libraries.
- Rewrite forms to the current `src/hooks/form` architecture where old vee-validate v2 / `v-field` contracts no longer fit.
- Components must use hooks from `src/hooks/store/modules`; hooks wrap Pinia stores; only stores call `request` / API routes.
- Prefer existing wrappers (`c-layout`, `c-select`, `c-autocomplete`, `c-text-field`, `c-advanced-data-table`, etc.) over raw Vuetify controls.
- Do not start the dev server unless explicitly asked.

## Existing Vue 3 Pieces

Already present in the Vue 3 project:

- Route name/path constants for `admin-kpi`.
- Router entry and top-bar administration menu item.
- KPI permissions in `src/constants/permission.js`.
- KPI constants in `src/constants/kpi.js`.
- KPI EN/FR messages in `src/i18n/messages/*/kpi.js`.
- KPI API routes in `src/config.js`:
  - `API_ROUTES.filters`
  - `API_ROUTES.ratingSettings`
  - `API_ROUTES.bulkRatingSettings`
  - `API_ROUTES.metrics.settings`
  - `API_ROUTES.metrics.sli`
  - `API_ROUTES.metrics.rating`
  - `API_ROUTES.metrics.alarms`
  - metric export routes
- `patterns-fields` store/hook supports KPI filter pattern fields.
- `metrics` store/hook already exposes SLI, rating and alarms metric requests.
- `rating-settings` store/hook exists.
- `metrics-settings` store/hook exists.
- Chart.js 4 infrastructure from the healthcheck migration is available and should be reused where possible.

## Source Inventory

### Page

- `src/views/admin/kpi.vue`

### Charts

- `components/other/kpi/charts/kpi-charts.vue`
- `components/other/kpi/charts/kpi-alarms.vue`
- `components/other/kpi/charts/kpi-rating.vue`
- `components/other/kpi/charts/kpi-sli.vue`
- `components/other/kpi/charts/partials/kpi-error-overlay.vue`
- `components/other/kpi/charts/partials/kpi-alarms-filters.vue`
- `components/other/kpi/charts/partials/kpi-alarms-chart.vue`
- `components/other/kpi/charts/partials/kpi-rating-filters.vue`
- `components/other/kpi/charts/partials/kpi-rating-chart.vue`
- `components/other/kpi/charts/partials/kpi-sli-filters.vue`
- `components/other/kpi/charts/partials/kpi-sli-chart.vue`
- `components/other/kpi/charts/form/fields/kpi-rating-criteria-field.vue`
- `components/other/kpi/charts/form/fields/kpi-rating-metric-field.vue`
- `components/other/kpi/charts/form/fields/kpi-sli-show-mode-field.vue`

### KPI Filters

- `components/other/kpi/filters/kpi-filters.vue`
- `components/other/kpi/filters/kpi-filters-list.vue`
- `components/other/kpi/filters/partials/kpi-filters-expand-item.vue`
- `components/other/kpi/filters/form/kpi-filter-form.vue`
- `components/other/kpi/filters/form/kpi-filter-patterns-form.vue`
- `components/modals/kpi/create-kpi-filter.vue`

### Rating Settings

- `components/other/kpi/rating-settings/kpi-rating-settings.vue`
- `components/other/kpi/rating-settings/kpi-rating-settings-list.vue`

### Collection Settings

- `components/other/kpi/collection-settings/kpi-collection-settings.vue`
- `components/other/kpi/collection-settings/partials/kpi-collection-settings-basic-metrics-list.vue`
- `components/other/kpi/collection-settings/form/kpi-collection-settings-optional-metrics-form.vue`

### Stores, Mixins, Helpers

- `store/modules/entities/filter`
- `store/modules/entities/rating-settings`
- `store/modules/entities/metrics-settings`
- `store/modules/entities/metrics`
- `store/modules/entities/patterns-fields`
- `mixins/entities/filter`
- `mixins/entities/rating-settings`
- `mixins/entities/metrics-settings`
- `mixins/permissions/technical/kpi-filters`
- `mixins/permissions/technical/kpi-rating-settings`
- `mixins/permissions/technical/kpi-collection-settings`
- `helpers/entities/metrics-settings/form`
- filter / pattern helpers used by KPI filter forms

## Migration Steps

### 1. Stores And Hooks

- [x] Add or verify KPI filter Pinia store:
  - use `createCRUDStore()` if it matches old API behavior;
  - route: `API_ROUTES.filters`;
  - expose list, previous params, list without store, create, update, remove.
- [x] Add or verify `src/hooks/store/modules/filter.js` with domain names:
  - `filters`, `filtersMeta`, `filtersPending`, `filtersFetchingParams`;
  - `fetchFiltersList`, `fetchFiltersListWithPreviousParams`, `fetchFiltersListWithoutStore`;
  - `createFilter`, `updateFilter`, `removeFilter`.
- [x] Verify `rating-settings` hook exposes:
  - list refs;
  - `fetchRatingSettingsList`;
  - `fetchRatingSettingsListWithPreviousParams`;
  - `bulkUpdateRatingSettings`.
- [x] Verify `metrics-settings` hook exposes:
  - `fetchMetricsSettingsWithoutStore`;
  - `updateMetricsSettings`.
- [x] Verify `metrics` hook exposes:
  - `fetchSliMetricsWithoutStore`;
  - `fetchRatingMetricsWithoutStore`;
  - `fetchAlarmsMetricsWithoutStore`;
  - export helpers if used by KPI charts.

### 2. Permission Access

- [x] Port or replace:
  - `permissionsTechnicalKpiFiltersMixin`;
  - `permissionsTechnicalKpiRatingSettingsMixin`;
  - `permissionsTechnicalKpiCollectionSettingsMixin`.
- [x] Prefer existing permission hooks/mixins used by nearby migrated pages.
- [x] Keep page tab visibility identical to the old page.

### 3. Page Shell

- [x] Create `src/views/admin/kpi-page.vue`.
- [x] Wire router to the real page if it still points to a placeholder.
- [x] Port top-level state:
  - `activeTab`;
  - `pending`;
  - `timescaleAvailable`.
- [x] Port page lifecycle:
  - fetch healthcheck engines/status;
  - join `SOCKET_ROOMS.healthcheck`;
  - leave room on unmount.
- [x] Replace `v-tabs-items` / `v-tab-item` with `v-window` / `v-window-item`.
- [x] Preserve tab list:
  - Graphs;
  - Filters;
  - Rating settings;
  - Collection settings.
- [x] Preserve FAB behavior:
  - visible only on Filters and Rating settings;
  - create action only for Filters;
  - refresh current tab list.
  - create action opens the KPI filter modal on Filters tab.

### 4. Charts Tab

- [x] Port chart container first:
  - `kpi-charts.vue`;
  - `kpi-error-overlay.vue`;
  - temporary child placeholders are allowed only while wiring the page, then replace with real charts before marking this step complete.
- [x] Port Alarm metrics:
  - filters;
  - chart;
  - metric parameters field dependencies.
- [x] Port Rating metrics:
  - filters;
  - criteria field;
  - metric field;
  - chart.
- [x] Port SLI:
  - [x] filters;
- [x] show mode field;
- [x] chart.
- [x] Adapt Chart.js code to Chart.js 4:
  - explicit `Chart.register(...)`;
  - no reactive proxy around chart instances (`markRaw` if needed);
  - update old scales/options syntax;
  - verify tooltip callbacks and time/category axes.
- [x] Reuse existing chart export helpers where possible.
- [x] Runtime-check graphs with live API data.

Current status:

- [x] Added shared `bar-chart.vue` for Chart.js 4.
- [x] SLI fetch/render path is wired through `useMetrics()`.
- [x] SLI CSV/PNG export behavior is ported.
- [x] Alarm metrics fetch/render/export path is wired through `useMetrics()`.
- [x] Rating metrics fetch/render/export path is wired through `useMetrics()`.
- [x] KPI filter selector fetches filters through `useFilter()`, and KPI filters CRUD plus create/edit modal are migrated.

### 5. KPI Filters CRUD

- [x] Port `kpi-filters.vue`.
- [x] Port `kpi-filters-list.vue` using current `c-advanced-data-table`.
- [x] Port `kpi-filters-expand-item.vue`.
- [x] Port `create-kpi-filter` modal and register it.
- [x] Rewrite KPI filter forms to current form layer:
  - root modal uses `useForm()`;
  - subforms use `useFormContext()` / `useFormModels()`;
  - leaf fields use props and `update:*`;
  - no `v-field`, no component `@input`.
- [x] Preserve old filter payload shape.
- [x] Verify create/edit/duplicate/delete and table pagination/expansion.

### 6. Rating Settings

- [x] Port `kpi-rating-settings.vue`.
- [x] Port `kpi-rating-settings-list.vue`.
- [x] Use existing `rating-settings` hook, not direct store imports.
- [x] Preserve bulk enable/disable behavior and table pagination.
- [x] Runtime-check bulk update payload.

### 7. Collection Settings

- [x] Port `kpi-collection-settings.vue`.
- [x] Port `kpi-collection-settings-basic-metrics-list.vue`.
- [x] Rewrite `kpi-collection-settings-optional-metrics-form.vue` to current form layer.
- [x] Port or verify `metricsSettingsToForm` / form-to-entity helpers.
- [x] Use `metrics-settings` hook.
- [x] Preserve submit disabled state when unchanged.
- [x] Runtime-check fetch, update, validation and success popup.

### 8. I18n, Modals, Barrels

- [x] Register `MODALS.createKpiFilter` target.
- [x] Add component barrels only where the project uses them.
- [x] Verify KPI EN/FR messages cover migrated UI.
- [x] Verify page header and permissions translations.
- [x] Avoid renaming old message keys.

### 9. Verification

- [x] Run targeted eslint for changed KPI files.
- [x] Run `npm run build` after app-code changes compile.
- [x] Runtime-check `/admin/kpi`:
  - permissions and visible tabs;
  - timescale unavailable warning;
  - charts fetch/render/export;
  - filters CRUD;
  - rating-settings bulk update;
  - collection-settings submit;
  - no console errors on tab switch.

## Suggested Execution Order

1. Stores/hooks for KPI filters and permission access helpers.
2. Page shell with tabs and existing route.
3. Charts container and chart blocks.
4. KPI filters CRUD and modal.
5. Rating settings.
6. Collection settings.
7. Inventory update and build.
