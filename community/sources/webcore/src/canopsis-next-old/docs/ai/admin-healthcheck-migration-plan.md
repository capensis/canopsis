# Admin Healthcheck Migration Plan

Temporary working plan for migrating `admin-healthcheck` from Vue 2 to Vue 3. Delete this file after the page is fully migrated and the inventory is updated.

## Scope

- Source page: `/Users/verokhov/Workspace/canopsis-pro/community/sources/webcore/src/canopsis-next/src/views/admin/healthcheck.vue`
- Target page: `src/views/admin/healthcheck-page.vue`
- Route: `admin-healthcheck`, `/admin/healthcheck`
- Source project is read-only. Edit only `canopsis-next-vue-3`.

## Migration Style

- Preserve the source component style by default. Do not rewrite Options API components to Composition API.
- Adapt only the parts required by Vue 3, Vuetify 4, Vite, Pinia stores, and updated libraries.
- Rewrite forms only where the source project has already rewritten them or where current Vue 3 form infrastructure requires it.

## Dependencies

Install latest package versions and adapt the migrated code to them:

```bash
npm install chart.js@4.5.1 chartjs-adapter-moment@1.0.1 chartjs-plugin-annotation@3.1.0 chartjs-plugin-datalabels@2.2.0 chartjs-plugin-zoom@2.2.0 cytoscape@3.34.0 cytoscape-cise@2.0.1 cytoscape-dagre@4.0.0 cytoscape-node-html-label@1.2.2 file-saver@2.0.5 jspdf@4.2.1
```

Compatibility notes:

- The old project used `chart.js@3.5.1`; the target should use Chart.js 4.
- Chart.js 4 is ESM-only. The Vue 3 project already uses Vite and `"type": "module"`.
- `chartjs-plugin-annotation@3` requires Chart.js 4.
- Register Chart.js controllers, elements, scales, and plugins explicitly with `Chart.register(...)`.
- Check `chartjs-plugin-zoom` peer warnings after install.
- Cytoscape is already 3.x in the old project, so the main Cytoscape risk is Vue lifecycle/ref cleanup, not core API migration.

## Existing Vue 3 Pieces

Already present in the Vue 3 project:

- Healthcheck API routes in `src/config.js`.
- Healthcheck constants in `src/constants/healthcheck.js`.
- Healthcheck i18n messages.
- `src/stores/modules/healthcheck.js`, currently status-only.
- `src/hooks/store/modules/healthcheck.js`, currently status-only.
- `src/stores/modules/healthcheck-parameters.js`.
- `src/hooks/store/modules/healthcheck-parameters.js`.
- `src/components/other/healthcheck/partials/healthcheck-chips-list.vue`.
- `src/hooks/healthcheck-nodes.js`.
- `src/helpers/entities/healthcheck/color.js`.

## Migration Steps

### 1. Infrastructure

- [x] Install dependencies listed above.
- [x] Port `src/services/cytoscape.js` and register Cytoscape extensions:
  - `cytoscape-node-html-label`
  - `cytoscape-dagre`
  - `cytoscape-cise`
- [x] Port the Cytoscape `network-graph` wrapper used by the system status graph.
- [x] Port the minimal chart wrapper/integration layer for Chart.js 4.
- Port chart plugin registration:
  - annotation
  - datalabels
  - zoom
  - moment adapter
- Port `src/helpers/file/pdf.js` if PDF export is still needed.
- Check existing `src/helpers/file/files.js` with `file-saver`.

### 2. Stores And Hooks

- [x] Extend `src/stores/modules/healthcheck.js` with:
  - `fetchHealthcheckEnginesWithoutStore`
  - `fetchHealthcheckStatusWithoutStore`
  - `fetchHealthcheckEnginesOrderWithoutStore`
- [x] Extend `src/hooks/store/modules/healthcheck.js` accordingly.
- Verify `healthcheck-parameters` store/hook supports:
  - fetching parameters
  - updating parameters
- Verify `tech-metrics` store/hook supports the engines metrics tab.
- Verify `anomaly-monitored-connector` store/hook supports connectors tab needs.
- Replace Vuex mixins and `createNamespacedHelpers` usage with Pinia hooks.

### 3. Page Shell

- [x] Create `src/views/admin/healthcheck-page.vue`.
- [x] Wire `admin-healthcheck` route to the real page.
- [x] Port top-level state:
  - `activeTab`
  - `pending`
  - `services`
  - `enginesGraph`
  - `enginesParameters`
  - `hasInvalidEnginesOrder`
  - `maxQueueLength`
  - `maxMessagesLength`
  - `hasServerError`
- [x] Port access checks for:
  - `USER_PERMISSIONS.technical.healthcheck`
  - `USER_PERMISSIONS.technical.anomalyMonitoredConnector`
  - `USER_PERMISSIONS.technical.techmetrics`
- [x] Join and leave `SOCKET_ROOMS.healthcheck`.
- Keep tabs lazy by default. Avoid eager mounting unless runtime behavior requires it.
- Avoid unpleasant height animation if the page has the same issue as planning.

### 4. System Status And Network Graph

- [x] Port:
  - `healthcheck-network-graph.vue`
  - `healthcheck-engine-information.vue`
  - `partials/healthcheck-engine-instance-diagram.vue`
- [x] Replace Vue 2 lifecycle and refs with Vue 3 equivalents.
- [x] Ensure Cytoscape instance cleanup on unmount.
- [x] Ensure graph resize/re-render works when the tab becomes visible.
- [x] Port graph helpers from:
  - `helpers/charts/healthcheck.js`
  - old healthcheck node mixin behavior not already covered by `useHealthcheckNodes`.
- [x] Port and register modals:
  - `MODALS.healthcheckEngine`
  - `MODALS.healthcheckEnginesChainReference`
- [x] Runtime-check graph rendering and modal behavior with live healthcheck data.

### 5. Connectors Tab

- [x] Port:
  - `healthcheck-connectors.vue`
  - `partials/healthcheck-connectors-list.vue`
  - `partials/healthcheck-connector-tile.vue`
  - `hooks/healthcheck-connectors-socket.js`
- [x] Replace Vue 2 `set` usage with ordinary Vue 3 reactive updates or immutable replacement.
- [x] Wire the real Connectors tab into `healthcheck-page.vue` instead of the placeholder.
- [x] Port and register the connector history modal before restoring tile click behavior:
  - `components/modals/anomaly-monitored-connector/anomaly-monitored-connector-history.vue`
  - `components/modals/anomaly-monitored-connector/partials/anomaly-monitored-connector-history-body.vue`
  - `components/modals/anomaly-monitored-connector/partials/anomaly-monitored-connector-history-chart.vue`
  - `components/modals/anomaly-monitored-connector/partials/anomaly-monitored-connector-history-graph-placeholder.vue`
  - `components/modals/anomaly-monitored-connector/partials/anomaly-monitored-connector-history-summary-stats.vue`
  - direct dependencies such as `c-hours-interval-field` and the old line chart wrapper.
- [x] Verify socket updates merge without losing pagination/meta state.
- [x] Verify disabled/not-enough-data connector states.

### 6. Parameters Tab

- [x] Port:
  - `healthcheck-parameters.vue`
  - `form/healthcheck-form.vue`
  - `form/fields/healthcheck-enabled-limit-field.vue`
  - `form/fields/healthcheck-engine-instance-field.vue`
  - `helpers/entities/healthcheck/form.js`
- [x] Replace legacy `v-field` and vee-validate v2 field wiring with local Vue 3-compatible field updates.
- [x] Preserve the old page flow: fetch parameters on mount, edit limits and engine instances, submit, show success popup.
- [x] Avoid immediate validation caused by initial form hydration.
- [x] Runtime-check submit behavior with live API data.

### 7. Graphs Tab

- [x] Roll back the first Graphs attempt after the new interval filter broke the tab at runtime.
- [x] Restore the Graphs tab to a placeholder in `healthcheck-page.vue` so the healthcheck page stays usable.
- [x] Remove the broken Graphs-specific implementation files from the first attempt:
  - `healthcheck-graphs.vue`
  - `partials/healthcheck-last-hour.vue`
  - `partials/healthcheck-last-hour-filters.vue`
  - `partials/healthcheck-last-hour-graph.vue`
  - `partials/healthcheck-history-filters.vue`
  - `partials/healthcheck-history-graph.vue`
  - `components/common/chart/limited-time-line-chart.vue`
  - `components/common/chart/chart-export-actions.vue`
  - `helpers/charts/canvas.js`
  - `forms/fields/compat/c-quick-date-interval-field.vue`
- [x] Restart the Graphs migration from the old source, in smaller verified steps.
- [x] Port container:
  - `healthcheck-graphs.vue`
- [x] Port filters used by both graph blocks:
  - `form/fields/healthcheck-event-types-field.vue`
  - `form/fields/healthcheck-connector-names-field.vue`
- [x] Port Last Hour block:
  - `partials/healthcheck-last-hour.vue`
  - `partials/healthcheck-last-hour-filters.vue`
  - `partials/healthcheck-last-hour-graph.vue`
- [x] Port History block:
  - `partials/healthcheck-history-filters.vue`
  - `partials/healthcheck-history-graph.vue`
- [x] Add Chart.js 4 infrastructure:
  - `components/common/chart/limited-time-line-chart.vue`
  - `components/common/chart/chart-export-actions.vue`
  - `helpers/charts/canvas.js`
  - Chart.js registration for line charts, time scale, filler, tooltip, legend, annotation, and zoom.
- [x] Preserve the old `options.plugins.limit` shape through the Chart.js annotation plugin.
- [x] Adapt old Chart.js options to Chart.js 4:
  - replace removed/renamed scale options;
  - replace `grid.drawBorder` with `border.display`;
  - replace `grid.borderWidth` with `border.width`;
  - replace `grid.borderColor` with `border.color`;
  - replace `time.stepSize` with `ticks.stepSize` where needed;
  - check time scale callbacks, which now receive timestamps;
  - register all controllers, elements, scales, and plugins explicitly.
- [x] Reuse existing `message-rate-stats` Pinia store/hook for REST stats.
- [x] Preserve Last Hour behavior:
  - query minute interval;
  - join `SOCKET_ROOMS.messageRates` with filters;
  - refetch/rejoin when filters change;
  - export PNG as `HEALTHCHECK_LAST_HOUR_FILENAME_PREFIX + day(from-to)`.
- [x] Preserve History behavior:
  - default interval `QUICK_RANGES.last30Days`;
  - request a swipe-offset interval;
  - use minute sampling for <= 1 day and hour sampling otherwise;
  - apply `deleted_before` minimum and adjust interval if needed;
  - pan/zoom updates the interval;
  - export PNG as `HEALTHCHECK_HISTORY_FILENAME_PREFIX + from-to`.
- [x] Port/verify file export support:
  - `saveFile` helper or existing equivalent;
  - canvas-to-blob PNG export from Chart.js instance.
- [x] Runtime-check:
  - Last Hour fetch, socket updates, filters, and PNG export;
  - History fetch, quick ranges, pan, zoom buttons, wheel/drag zoom, min boundary, and PNG export;
  - dark/light chart colors and tooltip labels.
- [x] Static checks after rollback:
  - `npx eslint` for `healthcheck-page.vue`, restored compat quick interval type field, and date interval helper;
  - `npm run build`.
- [x] Port the legacy quick interval components before reconnecting Graphs:
  - `components/common/block/c-information-block.vue`;
  - `components/forms/fields/compat/c-quick-date-interval-field.vue`;
  - `components/forms/fields/compat/c-quick-date-interval-type-field.vue`;
  - `helpers/date/date-intervals.js#findQuickRangeValue`.
- [x] Keep the port close to the old source and only adapt Vue 3/Vuetify 4 API:
  - `v-menu` activator uses `{ props }`;
  - fields use `model-value` / `update:model-value` or `update:value`;
  - interval emits both `input` and `update:interval`.
- [x] Runtime-check quick interval in Graphs History before porting chart logic.
- [x] Reconnect Graphs tab with filters only:
  - `healthcheck-graphs.vue`;
  - `partials/healthcheck-last-hour.vue`;
  - `partials/healthcheck-last-hour-filters.vue`;
  - `partials/healthcheck-history-graph.vue`;
  - `partials/healthcheck-history-filters.vue`.
- [x] Keep chart/socket/export logic disabled while validating the interval runtime behavior.
- [x] Restore Chart.js 4 infrastructure after validating the interval layer:
  - `components/common/chart/limited-time-line-chart.vue`;
  - `components/common/chart/chart-export-actions.vue`;
  - `helpers/charts/canvas.js`.
- [x] Restore Last Hour graph logic:
  - REST fetch through `message-rate-stats`;
  - socket join/rejoin for `SOCKET_ROOMS.messageRates`;
  - PNG export;
  - Chart.js line graph and limit annotation.
- [x] Restore History graph logic:
  - REST fetch through `message-rate-stats`;
  - quick interval and filters;
  - swipe-offset query interval;
  - minute/hour sampling;
  - deleted-before boundary;
  - pan/zoom and PNG export.
- [x] Static checks after restoring graphs:
  - `npx eslint` for Graphs/chart files;
  - `npm run build`.

### 8. Engines Metrics Tab

- [x] Port `healthcheck-engines-metrics.vue`.
- [x] Use existing `tech-metrics` store/hook.
- [x] Add missing tech metrics download URL helper.
- [x] Add `useFilePolling` compatibility helper for export generation polling.
- [x] Adapt the table to current Vuetify data-table headers and slots.
- [x] Verify permission gate for `USER_PERMISSIONS.technical.techmetrics` on the page shell.
- [x] Runtime-check enable toggle, dump generation, polling, and download with live API data.

### 9. Modals

- Port and register:
  - `healthcheck-engine`
  - `healthcheck-engines-chain-reference`
- Preserve old modal sizing behavior. The old project configured both as `autoHeight`.
- Verify modal content with real graph node data and engines chain data.

### 10. I18n And Access

- Verify EN/FR keys:
  - `healthcheck`
  - `modals`
  - `page-headers`
  - `permission`
- Verify permissions:
  - technical healthcheck page access
  - healthcheck status access
  - anomaly monitored connector tab
  - tech metrics tab
- Verify broadcast message view mapping still includes `admin-healthcheck`.

### 11. Verification Checklist

- `npx eslint` for all migrated files.
- `npm run build`.
- Runtime checks:
  - page opens from `/admin/healthcheck`;
  - server error state is shown correctly;
  - healthcheck socket updates the system status graph;
  - graph renders nonblank and resizes after tab changes;
  - engine and chain-reference modals open;
  - connectors tab fetches and receives socket updates;
  - parameters tab fetches, edits, validates, and saves;
  - graphs tab renders last-hour and history charts;
  - zoom and export actions work;
  - engines metrics tab renders;
  - permissions hide restricted tabs;
  - tab switching does not introduce distracting height lag.

## Suggested Execution Order

1. Install dependencies.
2. Port infrastructure: Cytoscape service, Chart.js integration, file export helpers.
3. Extend healthcheck stores/hooks.
4. Build the page shell and first System Status tab.
5. Port Network Graph and its modals.
6. Port Connectors tab.
7. Port Parameters tab.
8. Port Graphs tab with Chart.js 4 adaptation.
9. Port Engines Metrics tab.
10. Run full verification and update `docs/migration-inventory.md`.
11. Delete this temporary plan file.
