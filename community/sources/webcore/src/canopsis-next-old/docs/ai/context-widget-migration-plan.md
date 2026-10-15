# Context Widget Migration Plan

Source runtime: `../canopsis-next/src/components/widgets/context/entities-list.vue`

Target runtime: `src/components/widgets/context/entities-list.vue`

The widget is the Context explorer (`WIDGET_TYPES.context`, component name
`entities-list-widget`). It is part of the `/view/:id` widget system. This plan
covers the complete widget flow: settings, rendering, query state, table,
columns, row/mass actions, expand panels, create/edit entity and service
modals, exports, pattern-preview reuse and verification.

The old Context table uses `c-advanced-data-table`. Keep that table. Do not copy
the dedicated AlarmList table chain. Context does not use virtual scroll, column
drag/resize, sticky header, density or Ctrl-drag. Harden only the table
contracts this widget uses: expand, select-all, mass-actions teleport, toolbar
slot, top/bottom pagination and custom `item.<field>` cells.

The settings sidebar must be migrated completely. Every setting from the old
sidebar and its full dependency slice must be ported; TODOs, stubs, disabled
controls and visual placeholders are not acceptable substitutes.

Do not port normalizr, `SCHEMA_EMBEDDED_KEY` or `entities/removeFromStore`.
Entity list state must follow the AlarmList Pinia pattern: plain records-by-id
plus per-widget `allIds` / `meta` / `pending` / `fetchingParams`. After
pbehavior delete, refresh the local pbehavior list only.

Defer AI chat in `create-service`. Follow `docs/ai/architecture.md`.

## Status

- `[x]` Step 0: inventory and dependency audit.
- `[x]` First implementation scope: complete widget settings sidebar.
- `[x]` Next: widget foundation and default configuration.
- `[x]` Next: query state, pagination and entity store integration.
- `[x]` Next: dedicated Context table shell on `c-advanced-data-table`.
- `[x]` Next: columns and cell formatting.
- `[x]` Next: create/edit entity and service modals.
- `[x]` Next: row actions, mass actions and create-service FAB.
- `[x]` Next: expand panel and details.
- `[x]` Toolbar, filters, CSV export and pattern-preview entities-list modal.
- `[x]` Final: view integration and Ready For Testing update.

## Step 0: Inventory And Dependency Audit

### Source runtime

- `widgets/context/entities-list.vue`
- `widgets/context/partials/entities-list-table-with-pagination.vue`
- `widgets/context/partials/entities-list-expand-panel.vue`
- `widgets/context/actions/actions-panel.vue`
- `widgets/context/actions/mass-actions-panel.vue`
- `widgets/context/actions/context-fab.vue`
- `widgets/context/manage-infos.vue` (entity form, not the table)

### Column and cell rendering

- `columns-formatting/entity-column-cell.vue`
- `columns-formatting/entity-column-state.vue`
- `columns-formatting/entity-column-pbehavior-info.vue`
- `columns-formatting/entity-column-event-statistics.vue`

### Expand panel tabs

- `partials/expand-panel-tabs/infos-tab.vue`
- `partials/expand-panel-tabs/impact-depends-tab.vue`
- `partials/expand-panel-tabs/tree-of-dependencies-tab.vue`
- `partials/expand-panel-tabs/impact-chain-dependencies-tab.vue`
- `partials/expand-panel-tabs/entity-alarms-list-table.vue`

### Settings and supporting modals

- `sidebars/context/context.vue`
- `sidebars/context/form/fields/context-entities-types-filter.vue`
- `sidebars/form/fields/tree-of-dependencies-settings.vue`
- `modals/entity/create-entity.vue` and `other/entity/form/*`
- `modals/entity/create-entity-info.vue`
- `modals/entity/entities-list.vue`
- `modals/entity/entities-root-cause-diagram.vue`
- `modals/service/create-service.vue` and `other/service/form/*`
- `other/entity/entity-comments.vue` and comment list/item partials

Already registered but unused by Context: `MODALS.pbehaviorPlanning`,
`MODALS.variablesHelp`, `MODALS.confirmation`, `MODALS.serviceDependencies`.

### Cross-cutting source dependencies

- entity Vuex module widget buckets, CRUD, context-graph, bulk enable/disable,
  export create/fetch;
- service Vuex CRUD (`create` / `update` / `remove` / `fetchItemWithoutStore`);
- `helpers/entities/widget/forms/context.js`;
- `helpers/entities/entity/query.js`;
- `helpers/entities/entity/url.js`;
- `helpers/entities/entity/form.js` (already present);
- `helpers/entities/service/form.js`;
- widget fetch-query, filter-select, advanced-search saved items and export
  mixins (already present);
- `mixins/entities/context-entity.js`;
- `mixins/entities/service.js`;
- `mixins/widget/actions-panel/context.js`;
- `mixins/permissions/widgets/context/{filters,category,pbehavior}.js`;
- `c-entity-advanced-search`, `c-entity-category-field`, filter selector/list;
- `c-speed-dial` for the create-service FAB;
- entity comments store (already present);
- Context i18n (`i18n/messages/{en,fr}/context.js`, already present).

### Existing target coverage

Already present in the target project:

- widget constants: `WIDGET_TYPES.context`, `SIDE_BARS.contextSettings`,
  `COMPONENTS_BY_WIDGET_TYPES`;
- Pinia `entity` store/hook with widget-scoped list buckets, entityBasics CRUD,
  context-graph, export, bulk enable/disable and archive/clean/check-state
  methods;
- Pinia `service` store/hook with dependency/impact fetches, bulk remove,
  and create/update/remove/fetch-item;
- Pinia `entity-comment` store/hook;
- `entityToForm` / `formToEntity`;
- `c-entity-advanced-search`, `c-entity-category-field`, filter selector/list;
- widget query mixins, export mixin, column-filter and entity-info-property
  mixins;
- settings fields reused by the Context sidebar: title, default sort, columns,
  filters, grid-range-size, switcher, export CSV, root-cause, availability
  graph, charts form;
- expand-panel reuse: `entity-charts`, `entity-availability`,
  `pbehaviors-simple-list`, `service-dependencies`,
  `alarms-list-table-with-pagination`;
- `generatePreparedDefaultContextWidget()`;
- Context i18n files.

Still missing for the complete widget:

- Ready For Testing update (Step 9).

## First Implementation Scope: Complete Settings Sidebar

The first implementation scope is the complete settings sidebar. Execute it
before the runtime table so `widget.parameters`, templates, column mappings and
defaults are known before the renderer consumes them.

Use `useWidgetSettings()` like AlarmList / Map / Counter. Keep Options API
leaves as Options API. Replace `v-field` with explicit Vue 3 bindings. Use
`c-*` field wrappers.

The sidebar is not complete until every setting can be opened, edited,
validated, saved, reopened and mapped back without losing values.

Port every group from the old sidebar:

- `[x]` title;
- `[x]` default sort column;
- `[x]` widget columns and entity column templates;
- `[x]` service-dependencies columns with color indicator and templates;
- `[x]` tree-of-dependencies show type;
- `[x]` availability graph settings;
- `[x]` root-cause settings;
- `[x]` active alarms columns and templates;
- `[x]` resolved alarms columns and templates;
- `[x]` main filters (`with-alarm`, `with-entity`, `with-pbehavior`,
  `entity-counters-type`);
- `[x]` selected entity types filter;
- `[x]` expand grid range size;
- `[x]` CSV export form (`ENTITIES_TYPES.entity`, without infos attributes);
- `[x]` default inherited pbehavior;
- `[x]` keep selected after action;
- `[x]` charts form.

Implemented files:

- `src/helpers/entities/widget/forms/context.js`
- `src/helpers/entities/widget/form.js`
- `src/components/sidebars/form/fields/default-sort-column.vue`
- `src/components/sidebars/form/fields/tree-of-dependencies-settings.vue`
- `src/components/sidebars/context/form/fields/context-entities-types-filter.vue`
- `src/components/sidebars/context/context.vue`
- `src/components/sidebars/index.js`
- `src/components/widgets/partials/widget-wrapper-menu.vue`

Verification:

- `npx eslint` on the implemented files above.
- `[x]` Complete the sidebar implementation and move it to Ready For Testing.

## Next Scope: Widget Foundation

- `[x]` Register the runtime component in `widget-wrapper.vue` and call
  `prepareContextWidget`.
- `[x]` Register `contextSettings` in `components/sidebars/index.js`.
- `[x]` Add `WIDGET_TYPES.context` to `MIGRATED_SETTINGS_WIDGET_TYPES`.
- `[x]` Port `contextWidgetParametersToForm`, `formToContextWidgetParameters` and
  `prepareContextWidget` with the original JSDoc/typedefs.
- `[x]` Wire those mappers into `widgetParametersToForm` / `formToWidgetParameters`.
- `[x]` Replace the `generatePreparedDefaultContextWidget()` stub with
  `prepareContextWidget(generateDefaultContextWidget())`.
- Keep Options API runtime components as Options API.

Implemented files:

- `src/components/widgets/context/entities-list.vue`
- `src/components/widgets/widget-wrapper.vue`

## Step 2: Query, Pagination And Store

- `[x]` Extend the Pinia entity store with widget-scoped state like the alarm store:
  `entitiesById`, `widgets[widgetId].allIds`, `meta`, `pending`,
  `fetchingParams`.
- `[x]` Add request-cancelling `fetchList`, `fetchListWithPreviousParams`, CRUD
  (`create` / `update` / `remove` on `API_ROUTES.entityBasics`),
  `fetchBasicEntityWithoutStore`, `fetchContextGraphWithoutStore`,
  `createContextExport`, `fetchContextExport`, `bulkEnable`, `bulkDisable`.
- `[x]` Keep existing archive/clean/check-state-setting methods.
- `[x]` Expose the new methods through `useEntity()` and
  `mixins/entities/context-entity.js`.
- `[x]` Port `convertContextWidgetToQuery` and `convertContextUserPreferenceToQuery`
  into `helpers/entities/entity/query.js` and register them in
  `helpers/entities/widget/query.js`.
- `[x]` Fetch only when `widget.parameters.widgetColumns.length` is non-zero, with
  `with_flags: true`.
- `[x]` Reuse `widgetFetchQueryMixin` for query, user preference, category, filter
  and no-events updates.
- `[x]` Do not port normalizr cleanup.

Implemented files:

- `src/stores/modules/entity.js`
- `src/hooks/store/modules/entity.js`
- `src/helpers/entities/entity/query.js`
- `src/helpers/entities/widget/query.js`
- `src/mixins/entities/context-entity.js`
- `src/mixins/permissions/widgets/context/category.js`
- `src/mixins/permissions/widgets/context/filters.js`
- `src/components/widgets/context/entities-list.vue`

Intentionally skipped:

- `cleanEntitiesData` from the old mixin (not present on the old entity
  module; do not invent it).
- `widgetFilterSelectMixin` (not in the Vue 3 project). Filter/category
  updates are inlined like AlarmList / Map.

Verification:

- `npx eslint` on the implemented files above.

## Step 3: Table Shell

- `[x]` Port `entities-list-table-with-pagination.vue` onto `c-advanced-data-table`.
- `[x]` Preserve empty-column UI through `c-empty-data-table-columns`.
- `[x]` Preserve top `c-pagination` and bottom `c-table-pagination`.
- `[x]` Preserve select-all, expand, mass-actions grid teleport and
  `keepSelectedAfterAction`.
- `[x]` Preserve actions column unless `hideActions`.
- `[x]` Load entity-info properties and alarm column filters on mount.
- `[x]` `openRootCauseDiagram` is wired from the state cell click.
- `[x]` If expand, selection or mass-actions cannot match old parity on
  `c-advanced-data-table`, stop and document the gap. Do not silently switch
  to a dedicated table.

Did not stop: expand, select-all and mass-actions teleport already work on
`c-advanced-data-table` (same contracts as users-list / records-list). Grid
position uses `$vuetify.display.name` instead of old `$mq`.

Cell formatters are migrated. Row/mass action bodies and the expand-panel
tree stay in later steps. Expand and mass-actions slots are registered so
those contracts are live; their inner components are still empty placeholders.

Implemented files:

- `src/components/widgets/context/partials/entities-list-table-with-pagination.vue`
- `src/components/widgets/context/entities-list.vue`

Verification:

- `npx eslint` on the implemented files above.

## Step 4: Columns And Cells

- `[x]` Port `entity-column-cell` with Vue 3 `<component :is>`. Do not keep the old
  `v-bind.is` / `v-on` dynamic host.
- `[x]` Port state, pbehavior-info and event-statistics cells.
- `[x]` Reuse `c-enabled`, `c-no-events-icon`, `c-links-chips`, `c-ellipsis`,
  `c-alarm-chip`.
- `[x]` Preserve date filters, column-filter map, links category columns and
  root-cause click on state.

Column filters use `FILTERS` from `helpers/filter.js` instead of Vue 2
`$options.filters`.

Implemented files:

- `src/mixins/widget/columns-filters.js`
- `src/components/widgets/context/columns-formatting/entity-column-cell.vue`
- `src/components/widgets/context/columns-formatting/entity-column-state.vue`
- `src/components/widgets/context/columns-formatting/entity-column-pbehavior-info.vue`
- `src/components/widgets/context/columns-formatting/entity-column-event-statistics.vue`
- `src/components/widgets/context/partials/entities-list-table-with-pagination.vue`

Verification:

- `npx eslint` on the implemented files above.

## Step 5: Create Entity And Service Modals

Row actions and the FAB cannot ship without these modals.

- `[x]` Add service Pinia CRUD: `create`, `update`, `remove`,
  `fetchItemWithoutStore`. Expose them through `useService()` and
  `mixins/entities/service.js`. Widget-scoped weather list getters were not
  invented.
- `[x]` Port `helpers/entities/service/form.js` and
  `helpers/entities/entity-info/form.js`.
- `[x]` Port `entity-form`, `entity-general-form`, `entity-manage-infos-form`,
  `entity-info-form` and `widgets/context/manage-infos.vue`.
- `[x]` Port service form stack (`service-form`, `service-general-form`,
  `service-manage-infos-form`, `service-entity-patterns-form`) and
  `entity-state-setting`.
- `[x]` Port `create-entity`, `create-entity-info` and `create-service`.
- `[x]` Skip `ai-chat-sidebar` / `useAiChatForm` in `create-service`.
- `[x]` Register the modals in `components/modals/index.js`. `main.js` already
  has max-width entries.

Implemented files:

- `src/stores/modules/service.js`
- `src/hooks/store/modules/service.js`
- `src/mixins/entities/service.js`
- `src/helpers/entities/service/form.js`
- `src/helpers/entities/entity-info/form.js`
- `src/components/forms/fields/entity/c-entity-type-field.vue`
- `src/components/forms/fields/entity/c-impact-level-field.vue`
- `src/components/forms/fields/compat/c-mixed-field.vue`
- `src/components/forms/fields/compat/c-mixed-input-field.vue`
- `src/components/other/state-setting/entity-state-setting.vue`
- `src/components/other/entity/form/entity-form.vue`
- `src/components/other/entity/form/entity-general-form.vue`
- `src/components/other/entity/form/entity-manage-infos-form.vue`
- `src/components/other/entity/form/entity-info-form.vue`
- `src/components/widgets/context/manage-infos.vue`
- `src/components/other/service/form/service-form.vue`
- `src/components/other/service/form/service-general-form.vue`
- `src/components/other/service/form/service-manage-infos-form.vue`
- `src/components/other/service/form/service-entity-patterns-form.vue`
- `src/components/modals/entity/create-entity.vue`
- `src/components/modals/entity/create-entity-info.vue`
- `src/components/modals/service/create-service.vue`
- `src/components/modals/index.js`

Verification:

- `npx eslint` on the implemented files above.

## Step 6: Row Actions, Mass Actions And FAB

- `[x]` Port `mixins/widget/actions-panel/context.js` and both action panels.
- `[x]` Row actions: edit entity/service, duplicate service, delete, pbehavior,
  variables help.
- `[x]` Mass actions: delete, pbehavior, mass enable, mass disable, with the
  partial-deletable warning.
- `[x]` Port `context-fab`. The old FAB only creates services. `c-fab-expand-btn`
  is a page-level FAB, so `c-speed-dial` was ported for the inline toolbar
  control.
- `[x]` Refresh the widget list through `fetchContextEntitiesListWithPreviousParams`
  after every successful action.

Implemented files:

- `src/mixins/widget/actions-panel/context.js`
- `src/components/common/buttons/c-speed-dial.vue`
- `src/components/widgets/context/actions/actions-panel.vue`
- `src/components/widgets/context/actions/mass-actions-panel.vue`
- `src/components/widgets/context/actions/context-fab.vue`
- `src/components/widgets/context/partials/entities-list-table-with-pagination.vue`
- `src/components/widgets/context/entities-list.vue`

Verification:

- `npx eslint` on the implemented files above.

## Step 7: Expand Panel

- `[x]` Port `entities-list-expand-panel.vue` with the old tab order and
  conditions, using Vue 3 `v-tabs` / `v-window` named values. Window
  transitions are disabled like AlarmList expand (table expand rows break
  the horizontal animation).
- `[x]` Availability, if `availability.enabled` (reuse `entity-availability`).
- `[x]` Pbehaviors, if permission (reuse `pbehaviors-simple-list`, local
  refresh after delete).
- `[x]` Impact/depends, if the entity is not a service.
- `[x]` Infos.
- `[x]` Charts, if widget charts and `filtered_perf_data` exist (reuse
  `entity-charts`).
- `[x]` Tree of dependencies, if service (reuse `service-dependencies`).
- `[x]` Impact chain (reuse `service-dependencies` with `impact`).
- `[x]` Active alarm (reuse `alarms-list-table-with-pagination`).
- `[x]` Resolved alarms (same table, `resolved`).
- `[x]` Comments, if permission and type is service or resource.
- `[x]` `$periodicRefresh` registration on impact/depends, entity-alarms and
  entity-comments. `view-page.vue` already provides the observer.
- `[x]` `$tc` replaced with `$t(..., count)`.

Implemented files:

- `src/mixins/permissions/widgets/context/pbehavior.js`
- `src/mixins/permissions/widgets/entity-comment.js`
- `src/mixins/permissions/technical/exploitation/pbehavior.js`
- `src/components/widgets/context/partials/expand-panel-tabs/infos-tab.vue`
- `src/components/widgets/context/partials/expand-panel-tabs/impact-depends-tab.vue`
- `src/components/widgets/context/partials/expand-panel-tabs/tree-of-dependencies-tab.vue`
- `src/components/widgets/context/partials/expand-panel-tabs/impact-chain-dependencies-tab.vue`
- `src/components/widgets/context/partials/expand-panel-tabs/entity-alarms-list-table.vue`
- `src/components/other/entity/entity-comments.vue`
- `src/components/other/entity/partials/entity-comments-list.vue`
- `src/components/other/entity/partials/entity-comment.vue`
- `src/components/widgets/context/partials/entities-list-expand-panel.vue`
- `src/components/widgets/context/partials/entities-list-table-with-pagination.vue`

Verification:

- `npx eslint` on the implemented files above.

## Step 8: Toolbar, Export And Pattern Preview

- `[x]` Port the toolbar: `c-entity-advanced-search`, category field, filter
  selector/list (`with-alarm`, `with-entity`, `with-pbehavior`,
  `entity-counters-type`), no-events checkbox, FAB, CSV export.
- `[x]` Port `getContextExportDownloadFileUrl`.
- `[x]` Preserve kiosk `hideActions` / `hideMassSelection` from
  `widget-wrapper`.
- `[x]` Port `modals/entity/entities-list.vue` and register it so pattern
  counters can pass `generatePreparedDefaultContextWidget()`. Query init uses
  `convertContextWidgetToQuery` instead of the old alarm query converter.
- `[x]` Port `modals/entity/entities-root-cause-diagram.vue` with the
  state-click flow, including `entity-dependencies-by-state-settings`,
  `entity-network-graph` and `useEntityNetworkGraph`.
- `[x]` Idle-entities toolbar control uses `c-enabled-field` (compact inline
  switch), same as AlarmsList correlation.
- `[x]` Root-cause inherited-state summary: `{name} state is inherit from`
  plus **See filter pattern**. `MODALS.stateSettingInheritedEntityPattern`
  is registered and opens the read-only pattern step.
- `[x]` Root-cause node HTML labels: Vue 3 `state-*` classes only set text
  color, so node fill uses `getEntityColor()` / CSS vars; icons are
  flex-centered (Vue 2 `v-btn__content` flex is gone).
- `[x]` Root-cause layout: Cytoscape constructor plus `resetLayout` both
  run dagre. Vue 2 `animate: 'end'` left the first fit animation running
  so the second fit was ignored. Layout is `animate: false`; `cy.stop()`
  runs before layout and `fit()`.

Implemented files:

- `src/helpers/entities/entity/url.js`
- `src/helpers/entities/entity/cytoscape.js`
- `src/helpers/entities/map/list.js`
- `src/hooks/charts/entity-network-graph.js`
- `src/components/common/chart/entity-network-graph.vue`
- `src/components/other/entity/entity-dependencies-by-state-settings.vue`
- `src/components/other/state-setting/state-settings-summary.vue`
- `src/components/modals/state-setting/state-setting-inherited-entity-pattern.vue`
- `src/components/modals/entity/entities-list.vue`
- `src/components/modals/entity/entities-root-cause-diagram.vue`
- `src/components/modals/index.js`
- `src/components/widgets/context/entities-list.vue`

Verification:

- `npx eslint` on the implemented files above.
- `[x]` Root-cause graph fit/pan after modal open is visually checked.

## Step 9: Integration And Verification

Verify creation, editing, duplication, deletion, rendering, pagination,
sorting, search, category, filters, no-events, CSV export, selection, mass
actions, every expand tab, settings save/reopen, permissions, kiosk
`hideActions`, pattern entities modal, root-cause modal, dark/light themes
and API errors.

- `[x]` Ready For Testing entry is updated.

The migration is complete only when the widget has no placeholder branches,
Vue 2 listener/slot API, `tc`, `$validator`/legacy injection warnings,
normalizr cleanup dependency or unregistered `c-*` components.

## Documentation Rule

After every completed step, update this file and `docs/migration-inventory.md`
in the same change. Record migrated files, intentionally skipped legacy code,
verification results and remaining risks.
