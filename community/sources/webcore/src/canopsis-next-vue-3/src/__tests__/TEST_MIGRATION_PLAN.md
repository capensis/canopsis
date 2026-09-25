# Test Migration Plan

Generated: 2026-09-15T05:08:26.082Z
Updated: 2026-09-25T08:29:58+05:00

## Current Scope

- Old unit specs: 448
- Migrated non-snapshot specs: 352
- Removed snapshot-only specs: 94
- Latest Vitest tests: 1796
- Latest passing tests: 1673
- Latest failing tests: 113
- Latest skipped tests: 10
- Suites failed to load: 19

## Recently migrated

Modal specs: vee-validate submit, `modals.hide({ id })`, widget filter list via Pinia/axios mocks.

Alarm acknowledgement modal: folded the obsolete `confirm-ack-with-ticket` spec into the
current `create-ack-event` contract and migrated its action flags, explicit modal-id close,
cancel, and error-popup assertions.

Associate-ticket modal: migrated form mapping, explicit modal-id close, cancel, and
error-popup assertions to the Vue 3 contract.

Change-state modal: migrated initial-state selection, form submission, updated values,
explicit modal-id close, cancel, and error-popup assertions to the Vue 3 contract.

Declare-ticket modal: migrated explicit modal-id close assertions while preserving the
confirmation-modal and error-popup behavior.

Snooze modal: migrated default payload, note-required configuration, explicit modal-id
close, cancel, and error-popup assertions to the Vue 3 contract.

Meta-alarm link modal: migrated default form, request mapping with alarm ids, empty eligible
alarm handling, and explicit modal-id close assertions.

Comment modal: restored the missing `AlarmCommentTemplateForm` import and migrated plain
comments, structured template comments, error-popup, and explicit modal-id close behavior.

Text-editor modal: migrated modal-id close and Vue 3 form contract assertions, including
required configuration and action error handling.

Common controls and availability batch: migrated action/expand buttons, chips, table
pagination, variables menu, flowchart color/stroke fields, availability value filter, and
availability bar to Vue 3 `update:*` contracts and public component state assertions.

Availability/actions batch: migrated availability filters/history/line-chart event contracts,
the mass-actions teleport, selection, style, and `update:modelValue` assertions, and removed
the obsolete mass-actions panel spec targeting a component that no longer exists.

Entity/filter/rule batch: migrated entity availability query contracts, filter tile/list events,
entity form updates, declare-ticket webhook state, and removed obsolete KPI chart specs targeting
missing components.

KPI batch: migrated KPI filter/rating-settings pagination contracts, rating criteria selection,
alarm/SLI filter `update:query` contracts, and SLI show-mode updates; removed three obsolete
chart specs targeting missing components.

Map editor batch: migrated 10 map form/list/preview specs from Vue 2 wrapper internals and
legacy input assertions to explicit Vue 3 props/emits, scoped-slot stubs, and public required-field
state; updated Mermaid syntax validation for its asynchronous API.

Map/meta-alarm/pattern batch: migrated 10 remaining map preview, meta-alarm form, pattern list,
and pbehavior comment specs to Vue 3 events, shared form-context assertions, Pinia hook mocks,
and focused scoped-slot stubs.

Variables-help and entity modals: registered the migrated object treeview, rewrote variable
copy/export and entity form tests for current public APIs, and removed the stale title-buttons
spec for a component no longer present.

Service pause modal: migrated default form, submit/cancel, explicit modal-id close, and
error-popup assertions; removed stale modal wrapper specs for components no longer present.

Alarm status rule form: rewrote the parent-form suite around its current section-rendering
contract; field-level behavior remains covered by the dedicated general-form suite.

- components/modals/alarm/info-popup-setting/info-popup-setting.spec.js
- components/modals/common/filters-list.spec.js
- components/modals/alarm/create-ack-event.spec.js
- components/modals/alarm/create-associate-ticket-event.spec.js
- components/modals/alarm/create-change-state-event.spec.js
- components/modals/alarm/create-declare-ticket-event.spec.js
- components/modals/alarm/create-snooze-event.spec.js
- components/modals/alarm/link-to-meta-alarm.spec.js
- components/modals/common/create-comment-event.spec.js
- components/modals/common/text-editor.spec.js
- components/modals/common/variables-help.spec.js
- components/modals/entity/create-entity.spec.js
- components/modals/service/create-service-pause-event.spec.js

## Migration Rules

- Keep all runnable tests under `src/__tests__/unit`; do not colocate tests with components.
- Do not restore snapshot assertions or `__snapshots__` directories.
- Rewrite behavior tests against the new Vue 3/Vitest contracts instead of preserving Vue 2 implementation details.
- Prefer adapting shared test utilities/stubs when many old tests fail for the same Vue 2 to Vue 3 reason.
- For components whose public contract changed, rewrite the spec around the new props, emits, hooks, and store facade.

## Recommended Rewrite Order

1. Stabilize shared test harness: Vue Test Utils helpers, form model helpers, modal/sidebar/store hook bridges, common stubs.
2. Rewrite form-field specs by domain: pattern fields, date/time fields, JSON/mixed fields, entity fields.
3. Rewrite sidebar settings specs: alarm, context, service-weather, map, chart widgets.
4. Rewrite widget behavior specs: alarms list, alarm table, actions panels, availability, service weather.
5. Rewrite view/modal specs where API changed from options API stores to Pinia/composition hooks.
6. Finish store/helper specs and remove remaining compatibility shims that are no longer needed.

## High-Priority Failure Areas

- `toEmitInput` failures: old `input/change` expectations need mapping to `update:*`, form context, or rewritten expectations.
- Store action failures: old store mocks must be bridged or specs rewritten to mock Pinia/composition hooks directly.
- Missing child props/methods: old VTU v1 wrapper APIs and shallow stubs need Vue 3 equivalents.
- Modal close/open failures: specs should use `useModals`/`useInnerModal` behavior and explicit modal ids.
- Component contract changes: old tests that inspect internals should be rewritten around visible behavior and public API.

## Snapshot-Only Specs Removed

- components/common/block/c-collapse-panel.spec.js
- components/common/chart/horizontal-bar.spec.js
- components/common/chips/c-alarm-extra-details-chip.spec.js
- components/common/chips/c-alarm-pbehavior-chip.spec.js
- components/common/chips/c-alarm-tags-chips.spec.js
- components/common/chips/c-pattern-operator-chip.spec.js
- components/common/groups/c-density-btn-toggle.spec.js
- components/common/icons/c-no-events-icon.spec.js
- components/common/links/c-links-list.spec.js
- components/common/overlay/c-progress-overlay.spec.js
- components/common/runtime-template/c-compiled-template.spec.js
- components/common/runtime-template/c-runtime-template.spec.js
- components/common/table/c-treeview-data-table.spec.js
- components/common/table/color-indicator-wrapper.spec.js
- components/common/text-editor/variables-list.spec.js
- components/form/fields/pattern/pattern-operator-information.spec.js
- components/icons/c-enabled.spec.js
- components/layout/navigation/partials/groups-side-bar/groups-side-bar-playlists.spec.js
- components/modals/alarm/alarms-list.spec.js
- components/other/alarm-status-rule/form/partials/alarm-status-rule-patterns-form.spec.js
- components/other/availability/availability-list.spec.js
- components/other/availability/partials/availability-bar-chart-information-row.spec.js
- components/other/availability/partials/availability-bar-chart.spec.js
- components/other/availability/partials/availability-list-column-value.spec.js
- components/other/availability/partials/availability-list-expand-panel.spec.js
- components/other/declare-ticket/declare-ticket-rules-list.spec.js
- components/other/declare-ticket/declared-tickets-list.spec.js
- components/other/declare-ticket/form/declare-ticket-rule-patterns-form.spec.js
- components/other/dynamic-info/form/fields/dynamic-info-patterns-form.spec.js
- components/other/event-filter/event-filters-list.spec.js
- components/other/healthcheck/healthcheck-network-graph.spec.js
- components/other/idle-rule/form/idle-rule-patterns-form.spec.js
- components/other/kpi/charts/kpi-charts.spec.js
- components/other/kpi/filters/kpi-filters-list.spec.js
- components/other/kpi/filters/partials/kpi-filters-expand-item.spec.js
- components/other/login/cas-login.spec.js
- components/other/login/login-footer.spec.js
- components/other/login/oauth-login.spec.js
- components/other/login/partials/third-party-login.spec.js
- components/other/login/saml-login.spec.js
- components/other/map/partials/maps-list-expand-flowchart-item.spec.js
- components/other/map/partials/maps-list-expand-geomap-item.spec.js
- components/other/map/partials/maps-list-expand-mermaid-item.spec.js
- components/other/map/partials/mermaid-code-preview.spec.js
- components/other/map/partials/tree-of-dependencies-preview.spec.js
- components/other/meta-alarm-rule/form/meta-alarm-rule-patterns-form.spec.js
- components/other/pbehavior/pbehaviors/calendar/pbehavior-planning-calendar.spec.js
- components/other/pbehavior/pbehaviors/form/pbehavior-patterns-form.spec.js
- components/other/pbehavior/pbehaviors/partials/pbehavior-comments.spec.js
- components/other/pbehavior/pbehaviors/partials/pbehavior-exceptions-list.spec.js
- components/other/pbehavior/pbehaviors/partials/pbehavior-recurrence-rule.spec.js
- components/other/pbehavior/pbehaviors/partials/pbehaviors-list-expand-item.spec.js
- components/other/pbehavior/pbehaviors/pbehaviors-list.spec.js
- components/other/scenario/form/scenario-action-patterns-form.spec.js
- components/other/service/partials/service-dependencies-entity-cell.spec.js
- components/other/service/partials/service-entity-template.spec.js
- components/other/service/partials/service-entity-tree-of-dependencies-tab.spec.js
- components/other/tag/form/tag-patterns-form.spec.js
- components/other/tag/partials/tags-list-expand-panel.spec.js
- components/table/c-empty-data-table-columns.spec.js
- components/widgets/alarm/alarm-general-list.spec.js
- components/widgets/alarm/column-formatting/alarm-column-value-extra-details.spec.js
- components/widgets/alarm/column-formatting/alarm-column-value-status.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-ack.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-canceled.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-children.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-close-delay-value.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-close-delay.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-last-comment.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-parents.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-pbehavior.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-snooze.spec.js
- components/widgets/alarm/column-formatting/extra-details/extra-details-ticket.spec.js
- components/widgets/alarm/expand-panel/alarms-expand-panel-children.spec.js
- components/widgets/alarm/expand-panel/alarms-expand-panel-more-infos.spec.js
- components/widgets/alarm/headers-formatting/alarm-header-cell.spec.js
- components/widgets/alarm/headers-formatting/alarm-header-priority.spec.js
- components/widgets/alarm/timeline/alarm-timeline-days.spec.js
- components/widgets/alarm/timeline/alarm-timeline-step-expand.spec.js
- components/widgets/alarm/timeline/alarm-timeline-step-icon.spec.js
- components/widgets/alarm/timeline/alarm-timeline-step-title.spec.js
- components/widgets/chart/partials/bar-chart-metrics.spec.js
- components/widgets/chart/partials/line-chart-metrics.spec.js
- components/widgets/chart/partials/numbers-metrics-item.spec.js
- components/widgets/chart/partials/numbers-metrics.spec.js
- components/widgets/chart/partials/pie-chart-metrics.spec.js
- components/widgets/context/columns-formatting/entity-column-event-statistics.spec.js
- components/widgets/context/columns-formatting/entity-column-pbehavior-info.spec.js
- components/widgets/service-weather/alarm-pbehavior-counters.spec.js
- components/widgets/service-weather/alarm-state-counters.spec.js
- components/widgets/service-weather/impact-state-indicator.spec.js
- components/widgets/text/text.spec.js
- components/widgets/widget-wrapper.spec.js
- views/login.spec.js

## Runnable Non-Snapshot Tests

### components/common/actions-panel

#### components/common/actions-panel/actions-panel.spec.js
Status: passed 0, failed 4, skipped 0

- [test] actions-panel > Method into list called after trigger click on action item button. Size 'xl'
- [test] actions-panel > Fast pbehavior add action method called after trigger click
- [test] actions-panel > Fast pbehavior remove action method called after trigger click
- [test] actions-panel > Method into dropdown called after trigger click on action item button. Size 'xl'

#### components/common/actions-panel/c-mass-actions-panel.spec.js
Status: passed 8, failed 0, skipped 0

- [test] c-mass-actions-panel > Portal target is massActionsPanel when modal is not injected
- [test] c-mass-actions-panel > Portal target includes modal id when modal is injected
- [test] c-mass-actions-panel > Message displays selected count
- [test] c-mass-actions-panel > clearSelected emits clear:selected when close button clicked
- [test] c-mass-actions-panel > clearSelected emits clear:selected when called directly
- [test] c-mass-actions-panel > Style computed from grid parameters
- [test] c-mass-actions-panel > Keep-selected field emits input with true
- [test] c-mass-actions-panel > Keep-selected field emits input with false

### components/common/block

#### components/common/block/c-form-general-patterns-tabs.spec.js
Status: passed 0, failed 10, skipped 0

- [test] c-form-general-patterns-tabs > General and patterns tabs are rendered by default
- [test] c-form-general-patterns-tabs > Patterns tab is rendered before general tab when reverse is true
- [test] c-form-general-patterns-tabs > General tab is hidden when hideGeneral is true
- [test] c-form-general-patterns-tabs > Active tab switches to patterns when hideGeneral becomes true
- [test] c-form-general-patterns-tabs > Custom patterns and additional labels are rendered
- [test] c-form-general-patterns-tabs > Test query tab is disabled when disabledTestQueryTooltip is set
- [test] c-form-general-patterns-tabs > Template vars lists are fetched on mount when template testing tab is available
- [test] c-form-general-patterns-tabs > Template vars lists are not fetched on mount when template testing tab is unavailable
- [test] c-form-general-patterns-tabs > Testing tab is rendered when user has access and type is defined
- [test] c-form-general-patterns-tabs > Visible tabs contain expected tab ids

### components/common/buttons

#### components/common/buttons/c-action-btn.spec.js
Status: passed 7, failed 0, skipped 0

- [test] c-action-btn > Action button with default type and custom icon.
- [test] c-action-btn > Action button with default type and custom color.
- [test] c-action-btn > Action button without default type and all custom props.
- [test] c-action-btn > Check loading property.
- [test] c-action-btn > Check disabled property.
- [test] c-action-btn > Check button slot.
- [test] c-action-btn > Click event working correctly.

#### components/common/buttons/c-expand-btn.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-expand-btn > Expand button expand is worked.
- [test] c-expand-btn > Expand button collapse is worked.

### components/common/chart

#### components/common/chart/chart-export-actions.spec.js
Status: passed 2, failed 0, skipped 0

- [test] chart-export-actions > Export csv event emitted
- [test] chart-export-actions > Export png event emitted with correct blob

#### components/common/chart/junit-gantt-chart.spec.js
Status: passed 2, failed 0, skipped 0

- [test] junit-gantt-chart > Pagination next page event
- [test] junit-gantt-chart > Pagination rows per page event

### components/common/chips

#### components/common/chips/c-alarm-action-chip.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-chip > Renders `c-alarm-tag-chip` with default props

#### components/common/chips/c-alarm-actions-chips.spec.js
Status: passed 4, failed 0, skipped 0

- [test] c-alarm-actions-chips > Should emit `select` event on `click`
- [test] c-alarm-actions-chips > Should emit `select` event on `click` with custom itemValue
- [test] c-alarm-actions-chips > Should emit `select` event on `click` with returnObject
- [test] c-alarm-actions-chips > Should not emit `select` when clicking an already active item

#### components/common/chips/c-alarm-chip.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-alarm-chip > Click event emitted

#### components/common/chips/c-alarm-state-chip.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-alarm-state-chip > Click emitted after trigger click on chip

### components/common/flowchart

#### components/common/flowchart/fields/flowchart-color-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] flowchart-color-field > Value changed to first palette color after trigger select field with true
- [test] flowchart-color-field > Value changed to transparent after trigger select field with false
- [test] flowchart-color-field > Color picker modal showed after trigger button

#### components/common/flowchart/fields/flowchart-line-type-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] flowchart-line-type-field > Value changed after trigger select field

#### components/common/flowchart/fields/flowchart-number-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] flowchart-number-field > Value changed after trigger select field

#### components/common/flowchart/fields/flowchart-stroke-type-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] flowchart-stroke-type-field > Value changed after trigger select field

#### components/common/flowchart/flowchart-editor.spec.js
Status: passed 0, failed 29, skipped 0

- [test] flowchart-editor > Shape selected after mouse events triggered
- [test] flowchart-editor > Second shape selected after mouse events triggered with ctrl
- [test] flowchart-editor > Shape unselected after mouse events triggered on already selected shape with ctrl
- [test] flowchart-editor > Shapes unselected after mouse events triggered on already selected shape without ctrl
- [test] flowchart-editor > Selected shapes cleared and new selected after mouse events triggered without ctrl
- [test] flowchart-editor > Shapes selected after mouse events triggered on container
- [test] flowchart-editor > Shapes selected after mouse events triggered on container with shift
- [test] flowchart-editor > Shapes removed after keyboard event triggered
- [test] flowchart-editor > Shapes moved up after keyboard event triggered
- [test] flowchart-editor > Shapes moved right after keyboard event triggered
- [test] flowchart-editor > Shapes moved up after keyboard event triggered
- [test] flowchart-editor > Shapes moved left after keyboard event triggered
- [test] flowchart-editor > Shapes copied and pasted after keyboard event triggered
- [test] flowchart-editor > Shape moved after mouse event triggered
- [test] flowchart-editor > Shapes moved after mouse event triggered
- [test] flowchart-editor > Connected line moved after move connected shape
- [test] flowchart-editor > Line connected after mouseover with point
- [test] flowchart-editor > Line unconnected after mouseleave
- [test] flowchart-editor > Line unconnected after move connected shape
- [test] flowchart-editor > Shape: ${...} resized after trigger selection
- [test] flowchart-editor > Shape: ${...} resized after trigger selection
- [test] flowchart-editor > Shape: %s resized after trigger selection
- [test] flowchart-editor > <dynamic test name>
- [test] flowchart-editor > Shape: %s text changed after trigger
- [test] flowchart-editor > <dynamic test name>
- [test] flowchart-editor > View box changed after trigger mousemove event
- [test] flowchart-editor > View box scaled after trigger mousewheel event
- [test] flowchart-editor > View box moved horizontal after trigger mousewheel event with shift
- [test] flowchart-editor > View box moved vertical after trigger mousewheel event with alt

#### components/common/flowchart/flowchart-properties.spec.js
Status: passed 8, failed 0, skipped 0

- [test] flowchart-properties > Fill changed after trigger color field
- [test] flowchart-properties > Stroke changed after trigger color field
- [test] flowchart-properties > Stroke width changed after trigger number field
- [test] flowchart-properties > Stroke type changed after trigger stroke type field
- [test] flowchart-properties > Stroke type changed after trigger stroke type field
- [test] flowchart-properties > Text color changed after trigger color field
- [test] flowchart-properties > Text background color changed after trigger color field
- [test] flowchart-properties > text size changed after trigger number field

#### components/common/flowchart/flowchart-sidebar.spec.js
Status: passed 21, failed 0, skipped 0

- [test] flowchart-sidebar > Rect shape added after trigger button
- [test] flowchart-sidebar > Rounded rect shape added after trigger button
- [test] flowchart-sidebar > Square shape added after trigger button
- [test] flowchart-sidebar > Rhombus shape added after trigger button
- [test] flowchart-sidebar > Circle shape added after trigger button
- [test] flowchart-sidebar > Ellipse shape added after trigger button
- [test] flowchart-sidebar > Parallelogram shape added after trigger button
- [test] flowchart-sidebar > Process shape added after trigger button
- [test] flowchart-sidebar > Document shape added after trigger button
- [test] flowchart-sidebar > Storage shape added after trigger button
- [test] flowchart-sidebar > Curve shape added after trigger button
- [test] flowchart-sidebar > Curve arrow line shape added after trigger button
- [test] flowchart-sidebar > Bidirectional arrow curve line shape added after trigger button
- [test] flowchart-sidebar > Line shape added after trigger button
- [test] flowchart-sidebar > Arrow line shape added after trigger button
- [test] flowchart-sidebar > Bidirectional arrow line shape added after trigger button
- [test] flowchart-sidebar > Text shape added after trigger button
- [test] flowchart-sidebar > Textbox shape added after trigger button
- [test] flowchart-sidebar > Image shape added after trigger button
- [test] flowchart-sidebar > Image shape added and cropped after trigger button
- [test] flowchart-sidebar > Image shape added and cropped after trigger button

#### components/common/flowchart/flowchart.spec.js
Status: passed 0, failed 4, skipped 0

- [test] flowchart > Shapes added after trigger flowchart sidebar
- [test] flowchart > Shapes added after trigger flowchart editor
- [test] flowchart > Shapes added after trigger flowchart properties
- [test] flowchart > Points updated after trigger flowchart properties

### components/common/list

#### components/common/list/c-list.spec.js
Status: passed 6, failed 0, skipped 0

- [test] c-list > Item selected after click on list item
- [test] c-list > Item selected after click on list item (returnObject is enabled)
- [test] c-list > Submenu opened when mouseover list item (childrenKey: %s)
- [test] c-list > <dynamic test name>
- [test] c-list > Submenu closed after mouseover on other list item
- [test] c-list > Sub item selected after trigger input on nested list

### components/common/pagination

#### components/common/pagination/c-pagination.spec.js
Status: passed 6, failed 0, skipped 0

- [test] c-pagination > Pagination hidden without total
- [test] c-pagination > Pagination on the top. Check prev page button.
- [test] c-pagination > Pagination on the top. Check prev page button is disabled, when first page.
- [test] c-pagination > Pagination on the top. Check next page button.
- [test] c-pagination > Pagination on the top. Check next page button is disabled, when last page.
- [test] c-pagination > Pagination on the bottom. Check next page button is disabled, when last page.

#### components/common/pagination/c-table-pagination.spec.js
Status: passed 4, failed 0, skipped 0

- [test] c-table-pagination > Pagination hidden without items
- [test] c-table-pagination > Update pagination page
- [test] c-table-pagination > Update pagination items per page (without page changing)
- [test] c-table-pagination > Update pagination items per page (with page changing)

### components/common/runtime-template

#### components/common/runtime-template/c-request-helper.spec.js
Status: passed 7, failed 0, skipped 0

- [test] c-request-helper > Does not send request when helperData is not available
- [test] c-request-helper > Sends request successfully and updates template
- [test] c-request-helper > Sends request with custom method and options
- [test] c-request-helper > Handles unauthorized error
- [test] c-request-helper > Handles timeout error
- [test] c-request-helper > Handles other errors
- [test] c-request-helper > Cancels request on component unmount

### components/common/search

#### components/common/search/c-advanced-search-field.spec.js
Status: passed 4, failed 0, skipped 0

- [test] c-advanced-search-field > Emit submit when submit method called and validation passes
- [test] c-advanced-search-field > Emit reset when reset method called
- [test] c-advanced-search-field > Emit remove:search when history item removed
- [test] c-advanced-search-field > Emit toggle-pin:search when history item pin toggled

#### components/common/search/c-search-field.spec.js
Status: passed 0, failed 0, skipped 10

- [test] c-search-field > Not empty value was pass into props and it was applied to input field
- [test] c-search-field > Input value was update, after prop change and it was applied to input field
- [test] c-search-field > Set value into input element
- [test] c-search-field > Keyup without enter key on input element
- [test] c-search-field > Keyup with enter key on input element
- [test] c-search-field > Submit search button is the first button
- [test] c-search-field > Clear search button is the second button
- [test] c-search-field > Click on submit button
- [test] c-search-field > Click on clear button
- [test] c-search-field > Set value into combobox element

### components/common/text-editor

#### components/common/text-editor/text-editor.spec.js
Status: passed 6, failed 4, skipped 0

- [test] text-editor > Value changed after change props
- [test] text-editor > Value changed after trigger editor
- [test] text-editor > Value changed after trigger variables
- [test] text-editor > Menu showed after trigger variables button
- [test] text-editor > Value changed after trigger variables with caret in variable
- [test] text-editor > Value changed after trigger variables with selected variable
- [test] text-editor > Image uploaded after trigger image control
- [test] text-editor > Image as url uploaded after trigger image control
- [test] text-editor > Image not uploaded after trigger image control with large file
- [test] text-editor > Image not uploaded after trigger paste event

#### components/common/text-editor/variables-menu.spec.js
Status: passed 2, failed 0, skipped 0

- [test] variables-menu > Input event emitted after trigger variables
- [test] variables-menu > Close event emitted after trigger menu

### components/form/date-interval-selector.spec.js

#### components/form/date-interval-selector.spec.js
Status: passed 10, failed 0, skipped 0

- [test] date-interval-selector > Start time updated after trigger tstart field
- [test] date-interval-selector > Stop time updated after trigger tstart field
- [test] date-interval-selector > Start object value updated after trigger update object value
- [test] date-interval-selector > Stop object value updated after trigger update object value
- [test] date-interval-selector > Range updated after trigger quick range field
- [test] date-interval-selector > Range updated after trigger quick range field with custom type
- [test] date-interval-selector > Range didn't updated after trigger quick range field with previous value
- [test] date-interval-selector > Dates prepared after trigger prepare callback
- [test] date-interval-selector > Dates prepared after trigger prepare callback with rounded hours
- [test] date-interval-selector > Dates prepared after trigger prepare callback with string value

### components/form/fields

#### components/form/fields/ pbehavior/c-pbehavior-reason-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-pbehavior-reason-field > Reasons fetched after mount
- [test] c-pbehavior-reason-field > Reason changed after trigger select field

#### components/form/fields/ pbehavior/c-pbehavior-type-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-pbehavior-type-field > Type changed after trigger select field

#### components/form/fields/alarm/c-alarm-metric-aggregate-function-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-alarm-metric-aggregate-function-field > Value changed after trigger the input

#### components/form/fields/alarm/c-alarm-metric-parameters-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-alarm-metric-parameters-field > Value changed after trigger the input

#### components/form/fields/alarm/c-alarm-state-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-alarm-state-field > State type changed after trigger select field

#### components/form/fields/alarm/c-alarm-status-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-alarm-status-field > Value changed after trigger the input

#### components/form/fields/alarm/c-alarm-tag-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] c-alarm-tag-field > fetchAlarmTags was called
- [test] c-alarm-tag-field > fetchAlarmTagsLabels was called
- [test] c-alarm-tag-field > Value changed after trigger the input

#### components/form/fields/c-array-mixed-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] c-array-text-field > Empty string added after click on add button
- [test] c-array-text-field > Value changed after trigger mixed field
- [test] c-array-text-field > Field removed after click on remove button

#### components/form/fields/c-change-state-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-change-state-field > State changed after trigger the state field
- [test] c-change-state-field > Output changed after trigger the textarea field

#### components/form/fields/c-column-size-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-column-size-field > Renders `c-column-size-field` with required props

#### components/form/fields/c-date-interval-field.spec.js
Status: passed 4, failed 0, skipped 0

- [test] c-date-interval-field > From value changed after trigger from date picker
- [test] c-date-interval-field > To value changed after trigger from date picker
- [test] c-date-interval-field > Allowed dates on the "from" picker works
- [test] c-date-interval-field > Allowed dates on the "to" picker works

#### components/form/fields/c-description-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-description-field > Value changed after trigger the textarea

#### components/form/fields/c-draggable-list-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] c-draggable-list-field > Item position changed after trigger draggable with moved event
- [test] c-draggable-list-field > Filter position changed after trigger draggable with added event
- [test] c-draggable-list-field > Filter position changed after trigger draggable with removed event

#### components/form/fields/c-enabled-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-enabled-field > Value set in the input
- [test] c-enabled-field > Value changed in the input

#### components/form/fields/c-event-filter-type-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-event-filter-type-field > Value changed after trigger the text field

#### components/form/fields/c-id-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-id-field > Value changed after trigger the text field

#### components/form/fields/c-infos-attribute-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-infos-attribute-field > Dictionary changed after trigger the dictionary select
- [test] c-infos-attribute-field > Field changed after trigger the field select

#### components/form/fields/c-input-type-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-input-type-field > Input type changed after trigger select field

#### components/form/fields/c-json-field.spec.js
Status: passed 16, failed 0, skipped 0

- [test] c-json-field > Object value as prop
- [test] c-json-field > Input event after set value
- [test] c-json-field > Payload json value as prop with variables
- [test] c-json-field > Payload json value as prop without variables
- [test] c-json-field > v-validate works correctly on valid json
- [test] c-json-field > v-validate works correctly on invalid json and blur event
- [test] c-json-field > v-validate works correctly on invalid json and blur event (validate on `button`)
- [test] c-json-field > v-validate works correctly on invalid json and parse click
- [test] c-json-field > v-validate reset by previous value
- [test] c-json-field > v-validate reset by another value
- [test] c-json-field > v-validate works correctly on valid payload json value
- [test] c-json-field > v-validate works correctly on invalid payload json value
- [test] c-json-field > Set new value by prop
- [test] c-json-field > Set new value by prop (same value)
- [test] c-json-field > Set the same values
- [test] c-json-field > Check reset click

#### components/form/fields/c-mixed-field.spec.js
Status: passed 7, failed 0, skipped 0

- [test] c-mixed-field > Input type changed after trigger input type field with string value
- [test] c-mixed-field > Input type changed after trigger input type field with number value
- [test] c-mixed-field > Input type changed after trigger input type field with null value
- [test] c-mixed-field > Input type changed after trigger input type field with boolean value
- [test] c-mixed-field > Input type changed after trigger input type field with array value
- [test] c-mixed-field > Value changed on the first field after remove selected type
- [test] c-mixed-field > Value change on undefined after remove all types

#### components/form/fields/c-mixed-input-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] c-mixed-input-field > Value changed to empty string after trigger the input with null value
- [test] c-mixed-input-field > Value changed after trigger the input with number value
- [test] c-mixed-input-field > Value changed after trigger the select with items

#### components/form/fields/c-number-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-number-field > Value changed after trigger the input

#### components/form/fields/c-quick-date-interval-field.spec.js
Status: passed 11, failed 0, skipped 0

- [test] c-quick-date-interval-field > Value changed after trigger date interval field
- [test] c-quick-date-interval-field > Value changed after trigger quick date interval type field
- [test] c-quick-date-interval-field > Value changed after trigger quick date interval type field with short format
- [test] c-quick-date-interval-field > Value not changed after trigger date interval field with custom quick range
- [test] c-quick-date-interval-field > Ranges filtered with accumulatedBefore and min
- [test] c-quick-date-interval-field > Allowed from date function work correctly with accumulated before
- [test] c-quick-date-interval-field > Allowed from date function work correctly with min
- [test] c-quick-date-interval-field > Allowed from date function work correctly with to date
- [test] c-quick-date-interval-field > Allowed to date function work correctly with accumulated before
- [test] c-quick-date-interval-field > Allowed to date function work correctly with now
- [test] c-quick-date-interval-field > Allowed to date function work correctly with from date

#### components/form/fields/c-quick-date-interval-type-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-quick-date-interval-type-field > Value changed after trigger select field

#### components/form/fields/c-records-per-page-field.spec.js
Status: passed 8, failed 0, skipped 0

- [test] c-items-per-page-field > Default items is equal to PAGINATION_PER_PAGE_VALUES
- [test] c-items-per-page-field > Default items property applied to select element
- [test] c-items-per-page-field > Custom items property applied to select element
- [test] c-items-per-page-field > Default value property is equal to PAGINATION_LIMIT
- [test] c-items-per-page-field > Default value property applied to select element
- [test] c-items-per-page-field > Custom value property applied to select element
- [test] c-items-per-page-field > Custom items and value properties applied to select element
- [test] c-items-per-page-field > Set value into select element

#### components/form/fields/c-sampling-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-sampling-field > Value changed after trigger the input

#### components/form/fields/card-iterator/c-movable-card-iterator-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-movable-card-iterator-field > Cards reordered after trigger draggable list
- [test] c-movable-card-iterator-field > Cards reordered to first position after trigger draggable list

#### components/form/fields/color/c-color-indicator-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-color-indicator-field > The value did set in the input
- [test] c-color-indicator-field > Value changed after trigger the input

#### components/form/fields/column/c-columns-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] c-columns-field > Column added after trigger add event
- [test] c-columns-field > Column changed after trigger column field
- [test] c-columns-field > Column removed after trigger remove event

#### components/form/fields/date-picker/c-date-picker-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] c-date-picker-field > Value changed after trigger date picker
- [test] c-date-picker-field > Value cleared after trigger text field
- [test] c-date-picker-field > Change event emitted after trigger date picker

#### components/form/fields/date-time-picker/c-date-time-interval-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-date-time-interval-field > From changed after trigger from date time picker field
- [test] c-date-time-interval-field > To changed after trigger to date time picker field

#### components/form/fields/date-time-picker/date-time-picker-menu.spec.js
Status: passed 3, failed 0, skipped 0

- [test] date-time-picker-menu > Value changed after trigger input on the date time picker
- [test] date-time-picker-menu > Menu opened after trigger button
- [test] date-time-picker-menu > Menu closed after trigger close event on the date time picker

#### components/form/fields/date-time-picker/date-time-picker-text-field.spec.js
Status: passed 6, failed 0, skipped 0

- [test] date-time-picker-text-field > String value changed after trigger text field
- [test] date-time-picker-text-field > Date object not updated if text field is focused
- [test] date-time-picker-text-field > Date object updated if text field is blurred
- [test] date-time-picker-text-field > Errors added after validate date time text field
- [test] date-time-picker-text-field > Error added after validate date time text field with failed preparer
- [test] date-time-picker-text-field > Object field updated after trigger date time picker button

#### components/form/fields/date-time-picker/date-time-picker.spec.js
Status: passed 5, failed 0, skipped 0

- [test] date-time-picker > Time updated after trigger time picker field
- [test] date-time-picker > Date updated after trigger date picker
- [test] date-time-picker > Value changed after trigger apply button without changes
- [test] date-time-picker > Value changed after trigger apply button with changes
- [test] date-time-picker > Close emitted after trigger cancel button

#### components/form/fields/date-time-picker/date-time-splitted-picker-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] date-time-splitted-picker-field > Date changed after trigger date field
- [test] date-time-splitted-picker-field > Time changed after trigger time field

#### components/form/fields/date-time-picker/time-picker-field.spec.js
Status: passed 7, failed 0, skipped 0

- [test] time-picker-field > Combobox filter return true with exist text
- [test] time-picker-field > Combobox filter return false with excluded text
- [test] time-picker-field > Value changed after trigger change with empty value on the combobox
- [test] time-picker-field > List scrolled correctly without value
- [test] time-picker-field > List scrolled correctly with value
- [test] time-picker-field > Value changed after trigger change with valid value on the combobox
- [test] time-picker-field > Value changed after trigger change with round hours on the combobox

#### components/form/fields/date-time-splitted-range-picker-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] date-time-splitted-range-picker-field > Start changed after trigger date field
- [test] date-time-splitted-range-picker-field > End changed after trigger time field

#### components/form/fields/duration/c-duration-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] c-duration-field > Value changed after trigger text field
- [test] c-duration-field > Unit changed after trigger select field
- [test] c-duration-field > Value cleared after trigger select field without value

#### components/form/fields/duration/c-splitted-duration-field.spec.js
Status: passed 6, failed 0, skipped 0

- [test] c-splitted-duration-field > Minutes changed after trigger minutes field with value less then 1 hour
- [test] c-splitted-duration-field > Minutes changed after trigger minutes field with value greater then 1 hour
- [test] c-splitted-duration-field > Hours changed after trigger minutes field with value less then 1 hour
- [test] c-splitted-duration-field > Hours changed after trigger hours field with value greater then 1 day
- [test] c-splitted-duration-field > Days changed after trigger days field
- [test] c-splitted-duration-field > Value changed to maxValue after trigger with big value

#### components/form/fields/entity/c-entity-category-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-entity-category-field > Check fetch list call
- [test] c-entity-category-field > Category changed after trigger select field

#### components/form/fields/entity/c-entity-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-entity-field > Entities fetched after mount
- [test] c-entity-field > Value changed after trigger the input

#### components/form/fields/map/c-map-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] c-map-field > Filters fetched after mount
- [test] c-map-field > Map changed after trigger select field

#### components/form/fields/pattern/c-filter-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] c-filter-field > Filters fetched after mount if pending is false
- [test] c-filter-field > Filters not fetched after mount if pending is true
- [test] c-filter-field > Filter changed after trigger select field

#### components/form/fields/pattern/c-pattern-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] c-pattern-field > Patterns fetched after mount
- [test] c-pattern-field > Pbehavior patterns fetched after mount
- [test] c-pattern-field > Value changed after trigger the select

#### components/form/fields/pattern/c-patterns-field.spec.js
Status: passed 11, failed 0, skipped 0

- [test] c-patterns-field > Alarm pattern changed after trigger alarm patterns field
- [test] c-patterns-field > Entity pattern changed after trigger entity patterns field
- [test] c-patterns-field > Pbehavior pattern changed after trigger pbehavior patterns field
- [test] c-patterns-field > Event filter pattern changed after trigger event filter patterns field
- [test] c-patterns-field > Shows try optimization button when pattern has regexp infos and form has changes
- [test] c-patterns-field > Triggers optimization when try optimization button is clicked
- [test] c-patterns-field > Shows optimization suggestions when optimization succeeds with suggestions
- [test] c-patterns-field > Applies suggestion when apply suggestion event is triggered
- [test] c-patterns-field > Cancels optimization when cancel optimization event is triggered
- [test] c-patterns-field > Expands entity panel when optimization suggestions are received
- [test] c-patterns-field > Does not expand entity panel when optimization suggestions are empty

#### components/form/fields/pattern/pattern-advanced-editor-field.spec.js
Status: passed 5, failed 0, skipped 0

- [test] pattern-advanced-editor-field > Patterns invalid with wrong structure
- [test] pattern-advanced-editor-field > Patterns invalid with disabled attribute
- [test] pattern-advanced-editor-field > Patterns invalid with invalid infos attribute
- [test] pattern-advanced-editor-field > Patterns invalid with not exist attribute
- [test] pattern-advanced-editor-field > Patterns valid with valid rules

#### components/form/fields/pattern/pattern-attribute-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] pattern-attribute-field > Value changed after trigger the input

#### components/form/fields/pattern/pattern-editor-field.spec.js
Status: passed 4, failed 0, skipped 0

- [test] pattern-editor-field > Pattern id changed to custom after trigger input event on the pattern field
- [test] pattern-editor-field > Pattern changed after trigger input event on the pattern field
- [test] pattern-editor-field > Pattern changed to custom after click on the edit button
- [test] pattern-editor-field > Pattern groups changed after trigger advanced editor

#### components/form/fields/pattern/pattern-group-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] pattern-group-field > Rules updated after trigger input event on pattern rules field
- [test] pattern-group-field > Group removed after trigger input event on pattern rules field without rules

#### components/form/fields/pattern/pattern-groups-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] pattern-groups-field > Group removed after trigger remove event on the pattern group field
- [test] pattern-groups-field > Group updated after trigger update event on the pattern group field
- [test] pattern-groups-field > Group added after click on the add button

#### components/form/fields/pattern/pattern-operator-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] pattern-operator-field > Value changed after trigger the select

#### components/form/fields/pattern/pattern-progress.spec.js
Status: passed 4, failed 0, skipped 0

- [test] pattern-progress > Should emit close when close button is clicked in failed state
- [test] pattern-progress > Should emit try-again when try again button is clicked in failed state
- [test] pattern-progress > Should emit cancel when cancel button is clicked in progress state
- [test] pattern-progress > Should render custom labels and message when corresponding props are set

#### components/form/fields/pattern/pattern-rule-field-date-value.spec.js
Status: passed 39, failed 0, skipped 0

- [test] pattern-rule-field-date-value > Mount c-date-time-interval-field when operator is inRangeDates
- [test] pattern-rule-field-date-value > Mount c-quick-date-interval-type-field when operator is within
- [test] pattern-rule-field-date-value > Mount c-quick-date-interval-type-field when operator is olderThan
- [test] pattern-rule-field-date-value > Mount c-quick-date-interval-type-range-field when operator is inRangePeriod
- [test] pattern-rule-field-date-value > Mount nothing when operator is not date-related
- [test] pattern-rule-field-date-value > Passes correct props to c-date-time-interval-field
- [test] pattern-rule-field-date-value > Passes correct props to c-quick-date-interval-type-field
- [test] pattern-rule-field-date-value > Passes correct props to c-quick-date-interval-type-range-field
- [test] pattern-rule-field-date-value > Returns true for isDateRange when operator is inRangeDates
- [test] pattern-rule-field-date-value > Returns false for isDateRange when operator is not inRangeDates
- [test] pattern-rule-field-date-value > Returns true for isIntervalRange when operator is inRangePeriod
- [test] pattern-rule-field-date-value > Returns false for isIntervalRange when operator is not inRangePeriod
- [test] pattern-rule-field-date-value > Returns true for isInterval when operator is within
- [test] pattern-rule-field-date-value > Returns true for isInterval when operator is olderThan
- [test] pattern-rule-field-date-value > Returns false for isInterval when operator is not within or olderThan
- [test] pattern-rule-field-date-value > Returns true for isAllowedFromDate when no toTimestamp is set
- [test] pattern-rule-field-date-value > Returns true for isAllowedFromDate when from date is before to date
- [test] pattern-rule-field-date-value > Returns false for isAllowedFromDate when from date is after to date
- [test] pattern-rule-field-date-value > Returns true for isAllowedToDate when no fromTimestamp is set
- [test] pattern-rule-field-date-value > Returns true for isAllowedToDate when to date is after from date
- [test] pattern-rule-field-date-value > Returns false for isAllowedToDate when to date is before from date
- [test] pattern-rule-field-date-value > Updates model with default date form range when operator changes
- [test] pattern-rule-field-date-value > Does not update model when operator remains the same
- [test] pattern-rule-field-date-value > Uses default intervalRanges when not provided
- [test] pattern-rule-field-date-value > Uses default disabled value when not provided
- [test] pattern-rule-field-date-value > Does not require name prop
- [test] pattern-rule-field-date-value > Emits input event when c-date-time-interval-field value changes
- [test] pattern-rule-field-date-value > Emits input event when c-quick-date-interval-type-field value changes
- [test] pattern-rule-field-date-value > Emits input event when c-quick-date-interval-type-range-field value changes
- [test] pattern-rule-field-date-value > Returns intervalRanges when value.type is not custom
- [test] pattern-rule-field-date-value > Returns intervalRanges with custom range item when value.type is custom
- [test] pattern-rule-field-date-value > Returns intervalRanges when value.from is not custom
- [test] pattern-rule-field-date-value > Returns intervalRanges with custom range item when value.from is custom
- [test] pattern-rule-field-date-value > Returns intervalRanges when value.to is not custom
- [test] pattern-rule-field-date-value > Returns intervalRanges with custom range item when value.to is custom
- [test] pattern-rule-field-date-value > Creates correct custom range item for minutes
- [test] pattern-rule-field-date-value > Creates correct custom range item for hours with plural form
- [test] pattern-rule-field-date-value > Creates correct custom range item for single day
- [test] pattern-rule-field-date-value > Handles missing duration gracefully

#### components/form/fields/pattern/pattern-rule-field.spec.js
Status: passed 24, failed 0, skipped 0

- [test] pattern-rule-field > Attribute changed after trigger input event on the attribute field
- [test] pattern-rule-field > Pattern attribute field has return-object prop set to true
- [test] pattern-rule-field > Pattern attribute field emits input when updateAttribute is triggered
- [test] pattern-rule-field > Attribute with alias changed after trigger input event on the attribute field
- [test] pattern-rule-field > Attribute with definedType sets fieldType on updateAttribute
- [test] pattern-rule-field > Operator changed after trigger input event on the operator field
- [test] pattern-rule-field > Value changed after trigger input event on the value field
- [test] pattern-rule-field > Field and dictionary changed after trigger input event on the infos attribute field
- [test] pattern-rule-field > Field changed after trigger input event on the extra infos attribute field
- [test] pattern-rule-field > Range operator changed after trigger input event on the operator field
- [test] pattern-rule-field > Range type changed after trigger input event on the date value field
- [test] pattern-rule-field > Duration changed after trigger input event on the date duration field
- [test] pattern-rule-field > IsAnyInfosRuleOrAlias computed returns true when rule alias is true
- [test] pattern-rule-field > IsAnyInfosRuleOrAlias computed returns true when type is infos
- [test] pattern-rule-field > IsAnyInfosRuleOrAlias computed returns true when type is extra infos
- [test] pattern-rule-field > IsAnyInfosRuleOrAlias computed returns false when rule alias is false and type is not infos
- [test] pattern-rule-field > isInfosValueField computed returns true when rule field is value
- [test] pattern-rule-field > isInfosValueField computed returns true when rule alias is true
- [test] pattern-rule-field > isInfosValueField computed returns false when rule alias is false and field is not value
- [test] pattern-rule-field > isDateRule computed returns true when rule fieldType is timestamp
- [test] pattern-rule-field > isDateRule computed returns true when type is date
- [test] pattern-rule-field > isDateRule computed returns false when type is not date and fieldType is not timestamp
- [test] pattern-rule-field > notDefinedType computed returns true when item definedType differs from rule fieldType
- [test] pattern-rule-field > notDefinedType computed returns false when item definedType equals rule fieldType

#### components/form/fields/pattern/pattern-rules-field.spec.js
Status: passed 21, failed 0, skipped 0

- [test] pattern-rules-field > Rule removed after trigger remove event on the pattern Rule field
- [test] pattern-rules-field > Rule updated after trigger update event on the pattern Rule field
- [test] pattern-rules-field > Operator and value cleared after update rule with new attribute
- [test] pattern-rules-field > Value changed to array after update rule with array operator
- [test] pattern-rules-field > Rule added after click on the add button
- [test] pattern-rules-field > Add button is hidden when readonly prop is true
- [test] pattern-rules-field > Remove buttons are disabled when readonly prop is true
- [test] pattern-rules-field > Remove buttons are disabled when disabled prop is true
- [test] pattern-rules-field > Validation rule is attached on component creation
- [test] pattern-rules-field > Validation rule is detached on component destruction
- [test] pattern-rules-field > Validation getter filters out disabled rules
- [test] pattern-rules-field > hasRulesErrors returns true when validator has errors for the field
- [test] pattern-rules-field > Error message is displayed on add button when hasRulesErrors is true
- [test] pattern-rules-field > Operator is automatically set when only one operator is available
- [test] pattern-rules-field > getUpdatedRule method clears operator when not in available operators
- [test] pattern-rules-field > Rule props are calculated correctly with attribute options
- [test] pattern-rules-field > String rule props include one-of operators and keep opened value
- [test] pattern-rules-field > rulesMap computed property maps attributes correctly
- [test] pattern-rules-field > Default value is applied when changing attribute
- [test] pattern-rules-field > Custom name prop is used in validation
- [test] pattern-rules-field > Pattern rule field receives correct props from getRuleProps
- [test] pattern-rules-field > Disabled state is propagated correctly through readonly and disabled props
- [test] pattern-rules-field > Add button uses first attribute when adding new rule
- [test] pattern-rules-field > Value is converted when operator changes
- [test] pattern-rules-field > beforeDestroy hook calls detachMinValueRule

#### components/form/fields/periodic-refresh-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] periodic-refresh-field > Enabled changed after trigger the enabled field
- [test] periodic-refresh-field > Duration changed after trigger the duration field

#### components/form/fields/service-weather/c-service-weather-state-counter-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] c-service-weather-state-counter-field > Value changed after trigger the input

#### components/form/fields/state-criticity-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] state-criticity-field > Value changed after trigger click on the button

### components/form/pattern-form.spec.js

#### components/form/pattern-form.spec.js
Status: passed 4, failed 0, skipped 0

- [test] pattern-form > Title changed after trigger text field
- [test] pattern-form > Alarm pattern changed after trigger alarm patterns field
- [test] pattern-form > Entity pattern changed after trigger entity patterns field
- [test] pattern-form > Pbehavior pattern changed after trigger pbehavior patterns field

### components/form/patterns-form.spec.js

#### components/form/patterns-form.spec.js
Status: passed 2, failed 0, skipped 0

- [test] patterns-form > Title changed after trigger text field
- [test] patterns-form > Patterns changed after trigger patterns field

### components/layout/navigation

#### components/layout/navigation/partials/groups-side-bar/group-panel.spec.js
Status: passed 1, failed 0, skipped 0

- [test] group-panel > Change event emitted after trigger button

#### components/layout/navigation/partials/groups-side-bar/group-view-panel.spec.js
Status: passed 2, failed 0, skipped 0

- [test] group-view-panel > Change event emitted after trigger edit button
- [test] group-view-panel > Duplicate event emitted after trigger duplicate button

#### components/layout/navigation/partials/groups-side-bar/groups-side-bar-group-view.spec.js
Status: passed 2, failed 0, skipped 0

- [test] groups-side-bar-group-view > Duplicate view modal showed after trigger duplicate event
- [test] groups-side-bar-group-view > Change view modal showed after trigger change event

#### components/layout/navigation/partials/groups-side-bar/groups-side-bar-group.spec.js
Status: passed 2, failed 0, skipped 0

- [test] groups-side-bar-group > Edit group modal showed after trigger group panel
- [test] groups-side-bar-group > Change views order after trigger draggable list

#### components/layout/navigation/partials/groups-side-bar/groups-side-bar.spec.js
Status: passed 6, failed 0, skipped 0

- [test] groups-side-bar > Side bar opened after trigger navigation
- [test] groups-side-bar > Groups fetched after mount
- [test] groups-side-bar > Change groups order after trigger submit button
- [test] groups-side-bar > Error popup showed after trigger submit button with error
- [test] groups-side-bar > Cancel groups order after trigger cancel button
- [test] groups-side-bar > Editing mode changed after trigger settings button

### components/modals/alarm

#### components/modals/alarm/create-ack-event.spec.js
Status: passed 5, failed 0, skipped 0

- [test] create-ack-event > submits default form values and closes the modal by id
- [test] create-ack-event > passes the associate-ticket action flag
- [test] create-ack-event > passes the declare-ticket action flag when declare rules exist
- [test] create-ack-event > shows an error popup and keeps the modal open for unmapped action errors
- [test] create-ack-event > closes without an action when cancel is clicked

#### components/modals/alarm/create-associate-ticket-event.spec.js
Status: passed 4, failed 0, skipped 0

- [test] create-associate-ticket-event > submits the default form and closes the modal by id
- [test] create-associate-ticket-event > maps updated form values before submitting
- [test] create-associate-ticket-event > shows an error popup and keeps the modal open for unmapped action errors
- [test] create-associate-ticket-event > closes without an action when cancel is clicked

#### components/modals/alarm/create-change-state-event.spec.js
Status: passed 6, failed 0, skipped 0

- [test] create-change-state-event > initializes the state from a single alarm
- [test] create-change-state-event > uses the major state for multiple alarms
- [test] create-change-state-event > submits the form and closes the modal by id
- [test] create-change-state-event > submits updated state-field values
- [test] create-change-state-event > shows an error popup and keeps the modal open for unmapped action errors
- [test] create-change-state-event > closes without an action when cancel is clicked

#### components/modals/alarm/create-declare-ticket-event.spec.js
Status: passed 4, failed 0, skipped 0

- [test] create-declare-ticket-event > Form submitted after trigger submit button
- [test] create-declare-ticket-event > Confirmation modal showed after submit with exist tickets
- [test] create-declare-ticket-event > Error popup showed after trigger submit button with action errors
- [test] create-declare-ticket-event > Modal hidden after trigger cancel button

#### components/modals/alarm/create-snooze-event.spec.js
Status: passed 4, failed 0, skipped 0

- [test] create-snooze-event > submits default snooze values and closes the modal by id
- [test] create-snooze-event > forwards the note-required setting to the snooze form
- [test] create-snooze-event > shows an error popup and keeps the modal open for unmapped action errors
- [test] create-snooze-event > closes without an action when cancel is clicked

#### components/modals/alarm/edit-live-reporting.spec.js
Status: passed 9, failed 0, skipped 0

- [test] edit-live-reporting > Default parameters applied to form
- [test] edit-live-reporting > Config parameters applied to form
- [test] edit-live-reporting > Form submitted after trigger submit button
- [test] edit-live-reporting > Form didn't submitted after trigger submit button with error
- [test] edit-live-reporting > Form submitted after trigger submit button without action
- [test] edit-live-reporting > Errors added after trigger submit button with action errors
- [test] edit-live-reporting > Error popup showed after trigger submit button with action errors
- [test] edit-live-reporting > Modal submitted with correct data after trigger form
- [test] edit-live-reporting > Modal hidden after trigger cancel button

#### components/modals/alarm/info-popup-setting/add-info-popup.spec.js
Status: passed 9, failed 0, skipped 0

- [test] add-info-popup > Default parameters applied to form
- [test] add-info-popup > Config parameters applied to form
- [test] add-info-popup > Form submitted after trigger submit button
- [test] add-info-popup > Form didn't submitted after trigger submit button with error
- [test] add-info-popup > Form submitted after trigger submit button without action
- [test] add-info-popup > Errors added after trigger submit button with action errors
- [test] add-info-popup > Error popup showed after trigger submit button with action errors
- [test] add-info-popup > Modal submitted with correct data after trigger form
- [test] add-info-popup > Modal hidden after trigger cancel button

#### components/modals/alarm/info-popup-setting/info-popup-setting.spec.js
Status: passed 6, failed 0, skipped 0

- [test] info-popup-setting > Form submitted after trigger submit button
- [test] info-popup-setting > Form submitted after trigger submit button without action
- [test] info-popup-setting > Popup added after trigger add button
- [test] info-popup-setting > Popup edited after trigger edit button
- [test] info-popup-setting > Popup removed after trigger delete button
- [test] info-popup-setting > Modal hidden after trigger cancel button

#### components/modals/alarm/link-to-meta-alarm.spec.js
Status: passed 4, failed 0, skipped 0

- [test] link-to-meta-alarm > initializes the default link form
- [test] link-to-meta-alarm > maps the form and alarm ids before submitting
- [test] link-to-meta-alarm > submits an empty alarm list when no alarms are eligible for linking
- [test] link-to-meta-alarm > closes without an action when cancel is clicked

### components/modals/common

#### components/modals/common/confirmation.spec.js
Status: passed 4, failed 0, skipped 0

- [test] confirmation > Submit action called after trigger submit button
- [test] confirmation > Hidden action called after trigger submit button without action
- [test] confirmation > Cancel action called after trigger cancel button
- [test] confirmation > Cancel action called after destroy

#### components/modals/common/create-comment-event.spec.js
Status: passed 6, failed 0, skipped 0

- [test] create-comment-event > submits a plain comment and closes the modal by id
- [test] create-comment-event > converts template fields to struct_comment
- [test] create-comment-event > omits empty template fields from struct_comment
- [test] create-comment-event > passes configured templates to the view model
- [test] create-comment-event > shows an error popup and keeps the modal open for unmapped action errors
- [test] create-comment-event > closes without an action when cancel is clicked

#### components/modals/common/create-filter.spec.js
Status: passed 8, failed 0, skipped 0

- [test] create-filter > Form submitted without fields after trigger submit button
- [test] create-filter > Form submitted with all fields after trigger submit button
- [test] create-filter > Form didn't submitted after trigger submit button with error
- [test] create-filter > Form submitted after trigger submit button without action
- [test] create-filter > Errors added after trigger submit button with action errors
- [test] create-filter > Error popup showed after trigger submit button with action errors
- [test] create-filter > Modal submitted with correct data after trigger form
- [test] create-filter > Modal hidden after trigger cancel button

#### components/modals/common/filters-list.spec.js
Status: passed 6, failed 0, skipped 0

- [test] filters-list > Filters fetched after mount
- [test] filters-list > Filter created after trigger filters list
- [test] filters-list > Filter edited after trigger filters list
- [test] filters-list > Filter removed after trigger filters list
- [test] filters-list > Filters priority changed after trigger filters list
- [test] filters-list > Filters priority didn't change after trigger filters list with error

#### components/modals/common/text-editor.spec.js
Status: passed 9, failed 0, skipped 0

- [test] text-editor > Form submitted with empty string after trigger submit button
- [test] text-editor > Form submitted with correct value after trigger submit button
- [test] text-editor > Form didn't submitted after trigger submit button with error
- [test] text-editor > Form submitted after trigger submit button without action
- [test] text-editor > Validation rules applied to form from config
- [test] text-editor > Errors added after trigger submit button with action errors
- [test] text-editor > Error popup showed after trigger submit button with action errors
- [test] text-editor > Modal submitted with correct data after trigger form
- [test] text-editor > Modal hidden after trigger cancel button

#### components/modals/common/text-field-editor.spec.js
Status: passed 10, failed 0, skipped 0

- [test] text-field-editor > Form submitted with empty string after trigger submit button
- [test] text-field-editor > Form submitted with correct value after trigger submit button
- [test] text-field-editor > Form didn't submitted after trigger submit button with error
- [test] text-field-editor > Form submitted after trigger submit button without action
- [test] text-field-editor > Validation rules applied to form from config
- [test] text-field-editor > Errors added after trigger submit button with field error
- [test] text-field-editor > Errors added after trigger submit button with common error
- [test] text-field-editor > Errors added after trigger submit button with error message
- [test] text-field-editor > Modal submitted with correct data after trigger form
- [test] text-field-editor > Modal hidden after trigger cancel button

#### components/modals/common/variables-help.spec.js
Status: passed 3, failed 0, skipped 0

- [test] variables-help > copies a variable path successfully
- [test] variables-help > shows an error when copying a variable path fails
- [test] variables-help > exports a tree item as JSON

### components/modals/entity

#### components/modals/entity/create-entity.spec.js
Status: passed 7, failed 0, skipped 0

- [test] create-entity > Form submitted after trigger submit button
- [test] create-entity > submits the default form through the Vue 3 form contract
- [test] create-entity > Form submitted after trigger submit button without action
- [test] create-entity > Errors added after trigger submit button with action errors
- [test] create-entity > Error popup showed after trigger submit button with action errors
- [test] create-entity > Modal submitted with correct data after trigger form
- [test] create-entity > Modal hidden after trigger cancel button

### components/modals/pattern

#### components/modals/pattern/create-pattern.spec.js
Status: passed 10, failed 0, skipped 0

- [test] create-pattern > Form submitted after trigger submit button
- [test] create-pattern > Form didn't submitted after trigger submit button with error
- [test] create-pattern > Form submitted after trigger submit button without action
- [test] create-pattern > Errors added after trigger submit button with action errors
- [test] create-pattern > Error popup showed after trigger submit button with action errors
- [test] create-pattern > Modal submitted with type: "$type" data after trigger form
- [test] create-pattern > <dynamic test name>
- [test] create-pattern > Modal hidden after trigger cancel button

### components/modals/pbehavior

#### components/modals/pbehavior/create-pbehavior.spec.js
Status: passed 7, failed 0, skipped 0

- [test] create-pbehavior > Form submitted after trigger submit button
- [test] create-pbehavior > Form didn't submitted after trigger submit button with error
- [test] create-pbehavior > Form submitted after trigger submit button without action
- [test] create-pbehavior > Errors added after trigger submit button with action errors
- [test] create-pbehavior > Error popup showed after trigger submit button with action errors
- [test] create-pbehavior > Modal submitted with correct data after trigger form
- [test] create-pbehavior > Modal hidden after trigger cancel button

### components/modals/service

#### components/modals/service/create-service-pause-event.spec.js
Status: passed 5, failed 0, skipped 0

- [test] create-service-pause-event > initializes the default pause form
- [test] create-service-pause-event > submits updated form values and closes by id
- [test] create-service-pause-event > closes without an action
- [test] create-service-pause-event > shows an error popup and keeps the modal open for unmapped action errors
- [test] create-service-pause-event > closes when cancel is clicked

#### components/modals/service/create-service.spec.js
Status: passed 7, failed 0, skipped 0

- [test] create-service > Form submitted after trigger submit button
- [test] create-service > Form didn't submitted after trigger submit button with error
- [test] create-service > Form submitted after trigger submit button without action
- [test] create-service > Errors added after trigger submit button with action errors
- [test] create-service > Error popup showed after trigger submit button with action errors
- [test] create-service > Modal submitted with correct data after trigger form
- [test] create-service > Modal hidden after trigger cancel button

### components/modals/tag

#### components/modals/tag/create-tag.spec.js
Status: passed 7, failed 0, skipped 0

- [test] create-tag > Form submitted after trigger submit button
- [test] create-tag > Form didn't submitted after trigger submit button with error
- [test] create-tag > Form submitted after trigger submit button without action
- [test] create-tag > Errors added after trigger submit button with action errors
- [test] create-tag > Error popup showed after trigger submit button with action errors
- [test] create-tag > Modal submitted with correct data after trigger form
- [test] create-tag > Modal hidden after trigger cancel button

### components/modals/view

#### components/modals/view/create-view.spec.js
Status: passed 10, failed 0, skipped 0

- [test] create-view > View created after trigger submit button
- [test] create-view > Create error popup showed after trigger submit button with error
- [test] create-view > View updated after trigger submit button
- [test] create-view > Update error popup showed after trigger submit button with error
- [test] create-view > View duplicated after trigger submit button
- [test] create-view > Duplicate error popup showed after trigger submit button with error
- [test] create-view > Action called after trigger submit button with action
- [test] create-view > View removed after trigger remove button with action
- [test] create-view > Remove error popup showed after trigger remove button with error
- [test] create-view > Cancel action called after trigger cancel button

### components/other/alarm-status-rule

#### components/other/alarm-status-rule/form/alarm-status-rule-form.spec.js
Status: passed 2, failed 0, skipped 0

- [test] alarm-status-rule-form > General form is rendered in general tab
- [test] alarm-status-rule-form > hides patterns for a default rule

#### components/other/alarm-status-rule/form/alarm-status-rule-general-form.spec.js
Status: passed 5, failed 0, skipped 0

- [test] alarm-status-rule-general-form > Name changed after trigger name field
- [test] alarm-status-rule-general-form > Duration changed after trigger duration field
- [test] alarm-status-rule-general-form > Priority changed after trigger priority field
- [test] alarm-status-rule-general-form > Freq limit changed after trigger freq limit field when flapping is true
- [test] alarm-status-rule-general-form > Description changed after trigger description field

### components/other/availability

#### components/other/availability/form/fields/availability-display-parameter-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] availability-display-parameter-field > Show type changed after trigger radio field

#### components/other/availability/form/fields/availability-display-parameter-radio-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] availability-display-parameter-radio-field > Show type changed after trigger radio field

#### components/other/availability/form/fields/availability-show-type-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] availability-show-type-field > Show type changed after trigger radio field

#### components/other/availability/form/fields/availability-show-type-radio-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] availability-show-type-radio-field > Show type changed after trigger radio field

#### components/other/availability/form/fields/availability-value-filter-field.spec.js
Status: passed 4, failed 0, skipped 0

- [test] availability-value-filter-field > Value filter changed after trigger availability value filter method field without value
- [test] availability-value-filter-field > Value filter changed after trigger availability value filter method field with value
- [test] availability-value-filter-field > Value changed after trigger percent field
- [test] availability-value-filter-field > Value changed after trigger duration field

#### components/other/availability/form/fields/availability-value-filter-method-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] availability-value-filter-method-field > Method changed after trigger radio field

#### components/other/availability/partials/availability-bar.spec.js
Status: passed 2, failed 0, skipped 0

- [test] availability-bar > Update query emitted after trigger filters
- [test] availability-bar > Show type changed after trigger filters

#### components/other/availability/partials/availability-filters.spec.js
Status: passed 2, failed 0, skipped 0

- [test] availability-filters > Interval changed after trigger quick interval field
- [test] availability-filters > Show type changed after trigger show type field

#### components/other/availability/partials/availability-graph.spec.js
Status: passed 2, failed 4, skipped 0

- [test] availability-graph > History fetched after mount
- [test] availability-graph > History fetched after trigger update:sampling on availability history
- [test] availability-graph > History exported as png after trigger export:png on availability history
- [test] availability-graph > Error popup showed after trigger export:png with error on availability history
- [test] availability-graph > Widget exported after trigger export button
- [test] availability-graph > Error popup showed after trigger export button with error

#### components/other/availability/partials/availability-history-filters.spec.js
Status: passed 2, failed 0, skipped 0

- [test] availability-history-filters > Sampling changed after trigger sampling field
- [test] availability-history-filters > Show type changed after trigger show type field

#### components/other/availability/partials/availability-history.spec.js
Status: passed 2, failed 0, skipped 0

- [test] availability-history > update:sampling emitted after trigger filters
- [test] availability-history > Show type changed after trigger filters

#### components/other/availability/partials/availability-line-chart.spec.js
Status: passed 2, failed 0, skipped 0

- [test] availability-line-chart > Export csv event emitted
- [test] availability-line-chart > Export png event emitted

### components/other/comment-template

#### components/other/comment-template/form/alarm-comment-template-form.spec.js
Status: passed 1, failed 6, skipped 0

- [test] alarm-comment-template-form > Renders default comment field when no template is selected
- [test] alarm-comment-template-form > Renders template fields when template is selected
- [test] alarm-comment-template-form > Does not render template select when templates array is empty
- [test] alarm-comment-template-form > Clears errors when template is changed
- [test] alarm-comment-template-form > Template select field emits input with correct template object
- [test] alarm-comment-template-form > Fields are bound correctly with v-field directive
- [test] alarm-comment-template-form > Field labels use field.label or fallback to field.name

### components/other/declare-ticket

#### components/other/declare-ticket/form/associate-ticket-event-form.spec.js
Status: passed 5, failed 0, skipped 0

- [test] associate-ticket-event-form > System name changed after trigger name field
- [test] associate-ticket-event-form > Ticket changed after trigger ticket id field
- [test] associate-ticket-event-form > Ticket url changed after trigger ticket url field
- [test] associate-ticket-event-form > Ticket url title changed after trigger ticket url field
- [test] associate-ticket-event-form > Mapping changed after trigger custom fields field

#### components/other/declare-ticket/form/declare-ticket-rule-form.spec.js
Status: passed 2, failed 3, skipped 0

- [test] declare-ticket-rule-form > General form is rendered in general tab
- [test] declare-ticket-rule-form > Enabled changed after trigger enabled field
- [test] declare-ticket-rule-form > Form fields changed after trigger input event on general form
- [test] declare-ticket-rule-form > Patterns fields changed after trigger input event on patterns form
- [test] declare-ticket-rule-form > Test query form receives form prop

#### components/other/declare-ticket/form/declare-ticket-rule-general-form.spec.js
Status: passed 5, failed 0, skipped 0

- [test] declare-ticket-rule-general-form > Name changed after trigger name field
- [test] declare-ticket-rule-general-form > System name changed after trigger name field
- [test] declare-ticket-rule-general-form > Emit trigger success field changed after trigger emit trigger success field
- [test] declare-ticket-rule-general-form > Emit trigger fail field changed after trigger emit trigger fail field
- [test] declare-ticket-rule-general-form > Webhooks field changed after trigger webhooks field

### components/other/entity

#### components/other/entity/entity-availability.spec.js
Status: passed 4, failed 0, skipped 0

- [test] entity-availability > Availability fetched after mount
- [test] entity-availability > Availability fetched after trigger update interval with different time range
- [test] entity-availability > Availability fetched after trigger update interval
- [test] entity-availability > Fetching availability data for a different entity ID and time range

#### components/other/entity/form/entity-form.spec.js
Status: passed 11, failed 0, skipped 0

- [test] entity-form > General form is rendered in general tab
- [test] entity-form > Manage infos form is rendered in additional tab
- [test] entity-form > Name changed after trigger general form
- [test] entity-form > Description changed after trigger general form
- [test] entity-form > Available state changed after trigger general form
- [test] entity-form > Entity type changed after trigger general form
- [test] entity-form > Impact level changed after trigger general form
- [test] entity-form > Coordinates changed after trigger general form
- [test] entity-form > Enabled changed after trigger general form
- [test] entity-form > Prepare function passed to general form
- [test] entity-form > Infos changed after trigger manage infos form

### components/other/event-filter

#### components/other/event-filter/form/event-filter-form.spec.js
Status: passed 2, failed 0, skipped 0

- [test] event-filter-form > General form is rendered in general tab
- [test] event-filter-form > Enabled changed after trigger enabled field
- [test] event-filter-form > General fields changed after trigger general form
- [test] event-filter-form > Patterns changed after trigger patterns form

#### components/other/event-filter/form/event-filter-general-form.spec.js
Status: passed 3, failed 0, skipped 0

- [test] event-filter-general-form > Description changed after trigger description field
- [test] event-filter-general-form > Enrichment form is rendered for enrichment type
- [test] event-filter-general-form > Change entity form is rendered for change entity type

### components/other/filter

#### components/other/filter/filters-list.spec.js
Status: passed 4, failed 0, skipped 0

- [test] filters-list > Create filter modal opened after trigger add button
- [test] filters-list > Edit filter modal opened after trigger edit button
- [test] filters-list > Confirmation modal opened after trigger delete button
- [test] filters-list > Should send updated filters array on dragging

#### components/other/filter/partials/filter-selector.spec.js
Status: passed 1, failed 0, skipped 0

- [test] filter-selector > Filter changed after trigger the select field

#### components/other/filter/partials/filter-tile.spec.js
Status: passed 2, failed 0, skipped 0

- [test] filter-tile > Edit event emitted after trigger edit button
- [test] filter-tile > Delete event emitted after trigger delete button

### components/other/idle-rule

#### components/other/idle-rule/form/idle-rule-form.spec.js
Status: passed 2, failed 0, skipped 0

- [test] idle-rule-form > General form is rendered in general tab
- [test] idle-rule-form > IDLE Rule enabled after trigger enabled field
- [test] idle-rule-form > IDLE Rule fields changed after trigger general form
- [test] idle-rule-form > IDLE Rule patterns changed after trigger patterns form

#### components/other/idle-rule/form/idle-rule-general-form.spec.js
Status: passed 4, failed 0, skipped 0

- [test] idle-rule-general-form > Name changed after trigger name field
- [test] idle-rule-general-form > Description changed after trigger description field
- [test] idle-rule-general-form > Action type field is rendered for alarm type
- [test] idle-rule-general-form > Action type field is hidden for entity type

### components/other/kpi

#### components/other/kpi/charts/partials/kpi-alarms-filters.spec.js
Status: passed 6, failed 0, skipped 0

- [test] kpi-alarms-filters > Query changed after trigger a quick interval field
- [test] kpi-alarms-filters > Query changed after trigger a sampling field
- [test] kpi-alarms-filters > Query changed after trigger a sampling field to hour with large interval diff
- [test] kpi-alarms-filters > Query changed after trigger a sampling field to hour with normal interval diff
- [test] kpi-alarms-filters > Query changed after trigger a parameters field
- [test] kpi-alarms-filters > Query changed after trigger a filters field

#### components/other/kpi/charts/partials/kpi-alarms-chart.spec.js
Status: passed 2, failed 0, skipped 0

- [test] kpi-alarms-chart > Export csv event emitted
- [test] kpi-alarms-chart > Export png event emitted

#### components/other/kpi/charts/partials/kpi-sli-filters.spec.js
Status: passed 4, failed 0, skipped 0

- [test] kpi-sli-filters > Query changed after trigger a quick interval field
- [test] kpi-sli-filters > Query changed after trigger a sampling field
- [test] kpi-sli-filters > Query changed after trigger a show mode field
- [test] kpi-sli-filters > Query changed after trigger a filters field

#### components/other/kpi/charts/partials/kpi-rating-criteria-field.spec.js
Status: passed 4, failed 0, skipped 0

- [test] kpi-rating-criteria-field > selects the first available rating setting when mandatory
- [test] kpi-rating-criteria-field > clears a mandatory value when no rating settings are available
- [test] kpi-rating-criteria-field > keeps the current value when it is still available
- [test] kpi-rating-criteria-field > emits a changed criteria value

#### components/other/kpi/charts/partials/kpi-rating-chart.spec.js
Status: passed 2, failed 0, skipped 0

- [test] kpi-rating-chart > Export csv event emitted
- [test] kpi-rating-chart > Export png event emitted

#### components/other/kpi/charts/partials/kpi-rating-filters.spec.js
Status: passed 5, failed 0, skipped 0

- [test] kpi-rating-filters > Query changed after trigger a quick interval field
- [test] kpi-rating-filters > Query changed after trigger a filters field
- [test] kpi-rating-filters > Query changed after trigger a criteria field
- [test] kpi-rating-filters > Query changed after change criteria to entity criteria
- [test] kpi-rating-filters > Query changed after trigger a parameters field

#### components/other/kpi/charts/partials/kpi-rating-metric-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] kpi-rating-metric-field > Metric changed after trigger select field

#### components/other/kpi/charts/partials/kpi-sli-show-mode-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] kpi-sli-show-mode-field > Unit changed after trigger select field

#### components/other/kpi/filters/kpi-filters.spec.js
Status: passed 2, failed 0, skipped 0

- [test] kpi-filters > Filters fetched after mount
- [test] kpi-filters > Filters fetched after change query
- [test] kpi-filters > Patterns modal showed after trigger edit
- [test] kpi-filters > Patterns modal showed after trigger duplicate
- [test] kpi-filters > Patterns modal showed after trigger remove

#### components/other/kpi/rating-settings/kpi-rating-settings-list.spec.js
Status: passed 1, failed 0, skipped 0

- [test] kpi-rating-settings-list > Rating settings changed after enable rows

#### components/other/kpi/rating-settings/kpi-rating-settings.spec.js
Status: passed 2, failed 0, skipped 0

- [test] kpi-rating-settings > Rating settings fetched after mount
- [test] kpi-rating-settings > Rating settings fetched after change query

### components/other/login

#### components/other/login/basic-login.spec.js
Status: passed 0, failed 4, skipped 0

- [test] basic-login > Form submitted after trigger submit button
- [test] basic-login > Error shown after trigger submit button with error
- [test] basic-login > Form submitted after trigger submit button with redirect query
- [test] basic-login > Form submitted after trigger submit button with default view

#### components/other/login/form/login-form.spec.js
Status: passed 0, failed 2, skipped 0

- [test] login-form > Username changed after trigger text field
- [test] login-form > Password changed after trigger password field

### components/other/map

#### components/other/map/form/fields/add-location-btn.spec.js
Status: passed 0, failed 1, skipped 0

- [test] add-location-btn > Value updated after click button

#### components/other/map/form/fields/flowchart-editor.spec.js
Status: passed 0, failed 6, skipped 0

- [test] flowchart-editor > Background changed after trigger flowchart
- [test] flowchart-editor > Shapes changed after trigger flowchart
- [test] flowchart-editor > Linked point removed when shape removed
- [test] flowchart-editor > Points changed after trigger points editor
- [test] flowchart-editor > Add on click mode enabled after trigger add location btn
- [test] flowchart-editor > Form re-validated after change form with error

#### components/other/map/form/fields/geomap-editor.spec.js
Status: passed 12, failed 0, skipped 0

- [test] geomap-editor > Add on click mode enabled after trigger add location btn
- [test] geomap-editor > Point didn't add after click with add on click disabled
- [test] geomap-editor > Form re-validated after change form with error
- [test] geomap-editor > Point added after double click
- [test] geomap-editor > Point dialog menu didn't changed after double click when dialog already opened
- [test] geomap-editor > Point edited after open form dialog by contextmenu
- [test] geomap-editor > Point added after open form dialog by contextmenu
- [test] geomap-editor > Point removed form point dialog
- [test] geomap-editor > Point edited after open form dialog by point double click
- [test] geomap-editor > Point removed after trigger remove in contextmenu
- [test] geomap-editor > Coordinates changed after change coordinates by point dialog
- [test] geomap-editor > Marker moved after dragend

#### components/other/map/form/fields/mermaid-code-editor.spec.js
Status: passed 0, failed 3, skipped 0

- [test] mermaid-code-editor > Value changed after trigger code editor
- [test] mermaid-code-editor > Error added after update value
- [test] mermaid-code-editor > Error removed after input valid value

#### components/other/map/form/fields/mermaid-editor.spec.js
Status: passed 0, failed 5, skipped 0

- [test] mermaid-editor > Code changed after trigger code editor
- [test] mermaid-editor > Theme changed after trigger theme field
- [test] mermaid-editor > Points changed after trigger points editor
- [test] mermaid-editor > Add on click mode enabled after trigger add location btn
- [test] mermaid-editor > Form re-validated after change form with error

#### components/other/map/form/fields/mermaid-points-editor.spec.js
Status: passed 12, failed 0, skipped 0

- [test] mermaid-points-editor > Contextmenu opened after right click
- [test] mermaid-points-editor > Contextmenu didn't change when menu already opened
- [test] mermaid-points-editor > Point contextmenu didn't change when menu already opened
- [test] mermaid-points-editor > Point added after double click
- [test] mermaid-points-editor > Point dialog menu didn't changed after double click when dialog already opened
- [test] mermaid-points-editor > Point didn't add after click without prop
- [test] mermaid-points-editor > Point added after click with prop
- [test] mermaid-points-editor > Point edited after open form dialog by contextmenu
- [test] mermaid-points-editor > Point edited after open form dialog by point double click
- [test] mermaid-points-editor > Point removed after trigger remove in contextmenu
- [test] mermaid-points-editor > Point moved after trigger mouse events
- [test] mermaid-points-editor > Point didn't moving when dialog opened

#### components/other/map/form/fields/mermaid-theme-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] mermaid-theme-field > Value changed after trigger select

#### components/other/map/form/fields/point-contextmenu.spec.js
Status: passed 0, failed 3, skipped 0

- [test] mermaid-contextmenu > Add point event emitted after trigger click on add point
- [test] mermaid-contextmenu > Edit point event emitted after trigger click on edit point
- [test] mermaid-contextmenu > Remove point event emitted after trigger click on remove point

#### components/other/map/form/fields/point-form-dialog.spec.js
Status: passed 9, failed 0, skipped 0

- [test] point-form-dialog > Point dialog submitted with changes
- [test] point-form-dialog > Point dialog submitted without changes
- [test] point-form-dialog > Point dialog submitted with valid data after update prop
- [test] point-form-dialog > Point dialog didn't submitted with errors
- [test] point-form-dialog > Point dialog closed after trigger cancel button
- [test] point-form-dialog > Point dialog closed after trigger close button
- [test] point-form-dialog > Point removed after trigger delete button
- [test] point-form-dialog > Cancel emitted after click outside and confirm close
- [test] point-form-dialog > Submit emitted after click outside and confirm save

#### components/other/map/form/flowchart-map-form.spec.js
Status: passed 2, failed 0, skipped 0

- [test] flowchart-map-form > Name changed after trigger name field
- [test] flowchart-map-form > Parameters changed after trigger flowchart editor field

#### components/other/map/form/geomap-map-form.spec.js
Status: passed 2, failed 0, skipped 0

- [test] geomap-map-form > Name changed after trigger name field
- [test] geomap-map-form > Parameters changed after trigger mermaid editor field

#### components/other/map/form/mermaid-map-form.spec.js
Status: passed 0, failed 2, skipped 0

- [test] mermaid-map-form > Name changed after trigger name field
- [test] mermaid-map-form > Parameters changed after trigger mermaid editor field

#### components/other/map/maps-list.spec.js
Status: passed 0, failed 3, skipped 0

- [test] maps-list > Remove event emitted after trigger click on the remove button
- [test] maps-list > Update event emitted after trigger click on the edit button
- [test] maps-list > Duplicate event emitted after trigger click on the duplicate button

#### components/other/map/partials/flowchart-points-preview.spec.js
Status: passed 2, failed 1, skipped 0

- [test] flowchart-points-preview > Show alarms emitted after trigger point popup dialog
- [test] flowchart-points-preview > Show map emitted after trigger point popup dialog
- [test] flowchart-points-preview > Point popup dialog closed after trigger close event

#### components/other/map/partials/flowchart-preview.spec.js
Status: passed 0, failed 2, skipped 0

- [test] flowchart-preview > Shapes color changed by point with color indicator impact state
- [test] flowchart-preview > Shapes color changed by point with color indicator state

#### components/other/map/partials/geomap-preview.spec.js
Status: passed 3, failed 1, skipped 0

- [test] geomap-preview > Show alarms emitted after trigger point popup dialog
- [test] geomap-preview > Map fitted when points updated
- [test] geomap-preview > Show map emitted after trigger point popup dialog
- [test] geomap-preview > Point popup dialog closed after trigger close event

#### components/other/map/partials/mermaid-point-marker.spec.js
Status: passed 0, failed 1, skipped 0

- [test] mermaid-point-marker > Listeners applied to point icon

#### components/other/map/partials/mermaid-points-preview.spec.js
Status: passed 2, failed 1, skipped 0

- [test] mermaid-points-preview > Show alarms emitted after trigger point popup dialog
- [test] mermaid-points-preview > Show map emitted after trigger point popup dialog
- [test] mermaid-points-preview > Point popup dialog closed after trigger close event

#### components/other/map/partials/mermaid-preview.spec.js
Status: passed 3, failed 0, skipped 0

- [test] mermaid-preview > Show map emitted after trigger mermaid points preview
- [test] mermaid-preview > Show alarms emitted after trigger mermaid points preview
- [test] mermaid-preview > Panzoom reset after map prop updated

### components/other/meta-alarm-rule

#### components/other/meta-alarm-rule/form/meta-alarm-rule-form.spec.js
Status: passed 2, failed 2, skipped 0

- [test] meta-alarm-rule-form > General form is rendered in general tab
- [test] meta-alarm-rule-form > Enabled changed after trigger enabled field
- [test] meta-alarm-rule-form > General fields updated after trigger general form
- [test] meta-alarm-rule-form > Parameters changed after trigger parameters form

#### components/other/meta-alarm-rule/form/meta-alarm-rule-general-form.spec.js
Status: passed 5, failed 1, skipped 0

- [test] meta-alarm-rule-general-form > ID changed after trigger id field
- [test] meta-alarm-rule-general-form > Output template changed after trigger description field
- [test] meta-alarm-rule-general-form > Name changed after trigger text field
- [test] meta-alarm-rule-general-form > Enabled changed after trigger enabled field
- [test] meta-alarm-rule-general-form > Tags changed after trigger tags form
- [test] meta-alarm-rule-general-form > Infos changed after trigger infos form

#### components/other/meta-alarm-rule/form/meta-alarm-rule-parameters-form.spec.js
Status: passed 2, failed 3, skipped 0

- [test] meta-alarm-rule-parameters-form > Corel changed after trigger corel form
- [test] meta-alarm-rule-parameters-form > Threshold changed after trigger threshold form
- [test] meta-alarm-rule-parameters-form > Time based changed after trigger time based field
- [test] meta-alarm-rule-parameters-form > Value paths changed after trigger value paths field
- [test] meta-alarm-rule-parameters-form > Patterns changed after trigger patterns field

### components/other/pattern

#### components/other/pattern/corporate-patterns.spec.js
Status: passed 0, failed 5, skipped 0

- [test] corporate-patterns > Filters fetched after mount
- [test] corporate-patterns > Filters fetched after change query
- [test] corporate-patterns > Edit event emitted after trigger edit event on patterns list
- [test] corporate-patterns > Remove selected event emitted after trigger remove selected event on patterns list
- [test] corporate-patterns > Remove event emitted after trigger remove event on patterns list

#### components/other/pattern/patterns-list.spec.js
Status: passed 2, failed 1, skipped 0

- [test] patterns-list > Selected items removed after trigger mass remove button
- [test] patterns-list > Remove event emitted after trigger click on the remove button
- [test] patterns-list > Update event emitted after trigger click on the remove button

#### components/other/pattern/patterns.spec.js
Status: passed 0, failed 5, skipped 0

- [test] patterns > Filters fetched after mount
- [test] patterns > Filters fetched after change query
- [test] patterns > Edit event emitted after trigger edit event on patterns list
- [test] patterns > Remove selected event emitted after trigger remove selected event on patterns list
- [test] patterns > Remove event emitted after trigger remove event on patterns list

### components/other/pbehavior

#### components/other/pbehavior/pbehaviors/fields/pbehavior-comment-field.spec.js
Status: passed 1, failed 1, skipped 0

- [test] pbehavior-comment-field > Message changed after trigger textarea
- [test] pbehavior-comment-field > Remove event emitted after trigger remove button

#### components/other/pbehavior/pbehaviors/fields/pbehavior-comments-field.spec.js
Status: passed 0, failed 3, skipped 0

- [test] pbehavior-comments-field > Comment added after trigger create button
- [test] pbehavior-comments-field > Comment changed after trigger comment field
- [test] pbehavior-comments-field > Comment removed after trigger comment field

#### components/other/pbehavior/pbehaviors/fields/pbehavior-exception-field.spec.js
Status: passed 5, failed 2, skipped 0

- [test] pbehavior-exception-field > Begin date changed after trigger date time picker
- [test] pbehavior-exception-field > Begin date changed after trigger date time picker with defined type
- [test] pbehavior-exception-field > End date changed after trigger date time picker
- [test] pbehavior-exception-field > Type changed after trigger type field
- [test] pbehavior-exception-field > Delete event emitted after trigger button
- [test] pbehavior-exception-field > Full day enabled after trigger checkbox
- [test] pbehavior-exception-field > Full day disabled after trigger checkbox

#### components/other/pbehavior/pbehaviors/fields/pbehavior-recurrence-rule-exceptions-field.spec.js
Status: passed 1, failed 2, skipped 0

- [test] pbehavior-recurrence-rule-exceptions-field > Exceptions list shown when exceptions exist
- [test] pbehavior-recurrence-rule-exceptions-field > Exceptions list hidden when exceptions are empty
- [test] pbehavior-recurrence-rule-exceptions-field > Passes exdates to exceptions field

#### components/other/pbehavior/pbehaviors/fields/pbehavior-recurrence-rule-field.spec.js
Status: passed 2, failed 0, skipped 0

- [test] pbehavior-recurrence-rule-field > Rrule created after trigger create button
- [test] pbehavior-recurrence-rule-field > Rrule removed after trigger remove button

#### components/other/pbehavior/pbehaviors/form/pbehavior-form.spec.js
Status: passed 3, failed 0, skipped 0

- [test] pbehavior-form > General form is rendered in general tab
- [test] pbehavior-form > General options updated after trigger general pbehavior form
- [test] pbehavior-form > Filter updated after trigger pbehavior filter field

#### components/other/pbehavior/pbehaviors/form/pbehavior-general-form.spec.js
Status: passed 21, failed 0, skipped 0

- [test] pbehavior-general-form > Name changed after trigger name field
- [test] pbehavior-general-form > Inherited changed after trigger enabled field
- [test] pbehavior-general-form > Start on trigger enabled after trigger enabled field
- [test] pbehavior-general-form > Start on trigger disabled after trigger enabled field
- [test] pbehavior-general-form > Duration changed after trigger duration field
- [test] pbehavior-general-form > Start time changed after trigger date time picker field
- [test] pbehavior-general-form > Stop time changed after trigger date time picker field
- [test] pbehavior-general-form > Stop time didn't change after trigger date time picker field without value
- [test] pbehavior-general-form > Full day changed after trigger checkbox field
- [test] pbehavior-general-form > Full day didn't change after trigger checkbox field without tstart
- [test] pbehavior-general-form > No ending enabled after trigger checkbox field
- [test] pbehavior-general-form > No ending disabled after trigger checkbox field
- [test] pbehavior-general-form > No ending didn't changed after trigger checkbox field without tstart
- [test] pbehavior-general-form > No ending changed after form type changed field
- [test] pbehavior-general-form > No ending changed after form type changed field with full day
- [test] pbehavior-general-form > No ending didn't change after form type changed field with paused type
- [test] pbehavior-general-form > Reason changed after trigger reason field
- [test] pbehavior-general-form > Type changed after trigger type field
- [test] pbehavior-general-form > Comments updated after trigger pbehavior comments field
- [test] pbehavior-general-form > Color changed after trigger color field
- [test] pbehavior-general-form > Timezone changed after trigger timezone field

#### components/other/pbehavior/pbehaviors/partials/pbehavior-entities.spec.js
Status: passed 0, failed 1, skipped 0

- [test] pbehavior-entities > Pbehavior entities fetched after mount

#### components/other/pbehavior/pbehaviors/partials/pbehavior-recurrence-rule-periods.spec.js
Status: passed 1, failed 6, skipped 0

- [test] pbehavior-recurrence-rule-periods > Pbehavior timespans fetched after mount
- [test] pbehavior-recurrence-rule-periods > Pbehavior timespans fetched after pbehavior updated
- [test] pbehavior-recurrence-rule-periods > Pbehavior timespans fetched after range changed to next week
- [test] pbehavior-recurrence-rule-periods > Pbehavior timespans fetched after range changed to next two week
- [test] pbehavior-recurrence-rule-periods > Pbehavior timespans fetched after range changed to this month
- [test] pbehavior-recurrence-rule-periods > Pbehavior timespans fetched after range changed to next month
- [test] pbehavior-recurrence-rule-periods > Pbehavior recurrence rule modal opened after click action

#### components/other/pbehavior/pbehaviors/partials/pbehaviors-create-action-btn.spec.js
Status: passed 0, failed 1, skipped 0

- [test] pbehaviors-create-action-btn > Pbehavior planning modal opened after trigger button

#### components/other/pbehavior/pbehaviors/partials/pbehaviors-list-alarm-filtering-btn.spec.js
Status: passed 3, failed 2, skipped 0

- [test] pbehaviors-list-alarm-filtering-btn > Should render button with correct structure
- [test] pbehaviors-list-alarm-filtering-btn > Should have required component structure
- [test] pbehaviors-list-alarm-filtering-btn > Should handle component lifecycle correctly
- [test] pbehaviors-list-alarm-filtering-btn > Should trigger modal when button is clicked
- [test] pbehaviors-list-alarm-filtering-btn > Should call runAlarmFiltering when modal action is executed

#### components/other/pbehavior/pbehaviors/pbehaviors-simple-list.spec.js
Status: passed 1, failed 2, skipped 0

- [test] pbehaviors-simple-list > Pbehaviors fetched after mount
- [test] pbehaviors-simple-list > Pbehavior create modal opened after trigger create button
- [test] pbehaviors-simple-list > Pbehavior calendar modal opened after trigger calendar button

### components/other/service

#### components/other/service/form/fields/service-dependencies-show-type-field.spec.js
Status: passed 1, failed 0, skipped 0

- [test] service-dependencies-show-type-field > Value changed after trigger radio group field

#### components/other/service/form/service-form.spec.js
Status: passed 11, failed 2, skipped 0

- [test] service-form > General form is rendered in general tab
- [test] service-form > Entity patterns form is rendered in additional tab
- [test] service-form > Manage infos form is rendered in patterns tab
- [test] service-form > Name changed after trigger general form
- [test] service-form > Category changed after trigger general form
- [test] service-form > Available state changed after trigger general form
- [test] service-form > Impact level changed after trigger general form
- [test] service-form > Coordinates changed after trigger general form
- [test] service-form > Output template changed after trigger general form
- [test] service-form > Enabled changed after trigger enabled field
- [test] service-form > Prepare function passed to general form
- [test] service-form > Patterns changed after trigger entity patterns form
- [test] service-form > Infos changed after trigger manage infos form

#### components/other/service/partials/service-dependencies-expand.spec.js
Status: passed 0, failed 2, skipped 0

- [test] service-dependencies-expand > Show triggered after trigger click
- [test] service-dependencies-expand > Load triggered after trigger click

#### components/other/service/partials/service-dependencies.spec.js
Status: passed 0, failed 5, skipped 0

- [test] service-dependencies > Dependencies fetched after mount
- [test] service-dependencies > Children loaded after trigger load children
- [test] service-dependencies > More dependencies fetched after trigger more button
- [test] service-dependencies > Dependencies modal showed after trigger dependencies button
- [test] service-dependencies > Dependencies modal not showed after trigger dependencies button with impact

#### components/other/service/partials/service-entities-list.spec.js
Status: passed 17, failed 0, skipped 0

- [test] service-entities-list > Entities selected after trigger select event on service entity
- [test] service-entities-list > All entities selected and unselected after trigger checkbox
- [test] service-entities-list > Ack remove action applied after trigger mass ackremove  action
- [test] service-entities-list > Associate ticket action applied after trigger mass associate ticket action
- [test] service-entities-list > Validate action applied after trigger mass validate action
- [test] service-entities-list > Invalidate action applied after trigger mass invalidate action
- [test] service-entities-list > Pause action applied after trigger mass pause action
- [test] service-entities-list > Play action applied after trigger mass play action
- [test] service-entities-list > Cancel action applied after trigger mass cancel action
- [test] service-entities-list > Comment action applied after trigger mass comment action
- [test] service-entities-list > Declare ticket action applied after trigger mass declare ticket action
- [test] service-entities-list > Unavailable action not applied
- [test] service-entities-list > Unavailable action flag removed after trigger
- [test] service-entities-list > Entities list refreshed after trigger refresh on entity element
- [test] service-entities-list > Page updated after trigger pagination
- [test] service-entities-list > Records per page updated after trigger pagination
- [test] service-entities-list > Selected entities cleared after service entities updated

#### components/other/service/partials/service-entity-header.spec.js
Status: passed 1, failed 0, skipped 0

- [test] service-entity-header > Alert removed after trigger alert

#### components/other/service/partials/service-entity-info-tab.spec.js
Status: passed 0, failed 2, skipped 0

- [test] service-entity-info-tab > Action applied after trigger entity actions
- [test] service-entity-info-tab > Instruction executed after trigger entity actions

#### components/other/service/partials/service-entity.spec.js
Status: passed 0, failed 3, skipped 0

- [test] service-entity > Selected update emitted after trigger entity header
- [test] service-entity > Remove unavailable emitted after trigger entity header
- [test] service-entity > Instruction executed after trigger entity info tab

#### components/other/service/partials/service-template.spec.js
Status: passed 1, failed 2, skipped 0

- [test] service-template > Refresh applied after triggers entities list
- [test] service-template > Add action applied after triggers entities list
- [test] service-template > Pagination updated after triggers entities list

### components/other/storage-setting

#### components/other/storage-setting/form/storage-settings-alarm-external-tag-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-alarm-external-tag-form > Alarm external tag delete after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-alarm-form.spec.js
Status: passed 1, failed 1, skipped 0

- [test] storage-settings-alarm-form > Alarm archive after changed after trigger enabled duration field
- [test] storage-settings-alarm-form > Alarm delete after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-archive-entity-disabled-form.spec.js
Status: passed 0, failed 1, skipped 0

- [test] storage-settings-archive-entity-disabled-form > Form changed after trigger with dependencies checkbox

#### components/other/storage-setting/form/storage-settings-archive-entity-unlinked-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-archive-entity-unlinked-form > Duration changed after trigger duration field

#### components/other/storage-setting/form/storage-settings-connector-anomalies-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-connector-anomalies-form > Event anomaly delete after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-events-records-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-events-records-form > Events record delete after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-form.spec.js
Status: passed 10, failed 0, skipped 0

- [test] storage-settings-form > Alarm storage settings changed after trigger alarm settings
- [test] storage-settings-form > Remediation storage settings changed after trigger remediation settings
- [test] storage-settings-form > Pbehavior storage settings changed after trigger pbehavior settings
- [test] storage-settings-form > Junit storage settings changed after trigger junit settings
- [test] storage-settings-form > Health check storage settings changed after trigger health check settings
- [test] storage-settings-form > Webhook storage settings changed after trigger webhook settings
- [test] storage-settings-form > Metrics storage settings changed after trigger metrics settings
- [test] storage-settings-form > Perf data metrics storage settings changed after trigger perf data metrics settings
- [test] storage-settings-form > Entity infos log storage settings changed after trigger entity infos log settings
- [test] storage-settings-form > Connector anomalies storage settings changed after trigger connector anomalies settings

#### components/other/storage-setting/form/storage-settings-health-check-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-health-check-form > Health check delete after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-junit-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-junit-form > Junit delete after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-metrics-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-metrics-form > Metrics delete after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-pbehavior-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-pbehavior-form > Pbehavior delete after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-perf-data-metrics-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-perf-data-metrics-form > Perf data metrics delete after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-remediation-form.spec.js
Status: passed 3, failed 0, skipped 0

- [test] storage-settings-remediation-form > Remediation delete after changed after trigger enabled duration field
- [test] storage-settings-remediation-form > Remediation delete stats after changed after trigger enabled duration field
- [test] storage-settings-remediation-form > Remediation delete mod stats after changed after trigger enabled duration field

#### components/other/storage-setting/form/storage-settings-webhook-form.spec.js
Status: passed 1, failed 0, skipped 0

- [test] storage-settings-webhook-form > Webhook delete after changed after trigger enabled duration field

### components/other/tag

#### components/other/tag/form/tag-form.spec.js
Status: passed 4, failed 0, skipped 0

- [test] tag-form > General form is rendered in general tab
- [test] tag-form > Value changed after trigger general form
- [test] tag-form > Color changed after trigger general form
- [test] tag-form > Tag patterns changed after trigger patterns form

#### components/other/tag/form/tag-general-form.spec.js
Status: passed 2, failed 0, skipped 0

- [test] tag-general-form > Value changed after trigger name field
- [test] tag-general-form > Color changed after trigger color field

#### components/other/tag/tags-list.spec.js
Status: passed 0, failed 4, skipped 0

- [test] tags-list > Selected items removed after trigger mass remove button
- [test] tags-list > Remove event emitted after trigger click on the remove button
- [test] tags-list > Update event emitted after trigger click on the remove button
- [test] tags-list > Duplicate event emitted after trigger click on the duplicate button

### components/other/view

#### components/other/view/form/view-form.spec.js
Status: passed 4, failed 2, skipped 0

- [test] view-form > Title changed after trigger title field
- [test] view-form > Description changed after trigger description field
- [test] view-form > Enabled changed after trigger enabled field
- [test] view-form > Periodic refresh changed after trigger periodic refresh field
- [test] view-form > Tag changed after trigger tag field
- [test] view-form > Group changed after trigger group field

#### components/other/view/partials/view-periodic-refresh-btn.spec.js
Status: passed 3, failed 4, skipped 0

- [test] view-periodic-refresh-btn > computed > Should return correct periodicRefreshFullPaused value when activeViewPeriodicRefreshPaused
- [test] view-periodic-refresh-btn > computed > Should return correct periodicRefreshFullPaused value when isNavigationEditingMode
- [test] view-periodic-refresh-btn > computed > Should return false for periodicRefreshFullPaused when not paused
- [test] view-periodic-refresh-btn > methods > Should handle interval management correctly
- [test] view-periodic-refresh-btn > methods > Should have periodicRefreshPausedWatcher method
- [test] view-periodic-refresh-btn > new periodic refresh pause functionality > Should incorporate activeViewPeriodicRefreshPaused in periodicRefreshFullPaused
- [test] view-periodic-refresh-btn > new periodic refresh pause functionality > Should access activeViewPeriodicRefreshPaused computed property from mixin

#### components/other/view/view-tab-widgets.spec.js
Status: passed 4, failed 2, skipped 0

- [test] view-tab-widgets > Each query removed after destroy
- [test] view-tab-widgets > Register and unregister editing off handler is working (visible true)
- [test] view-tab-widgets > Register and unregister editing off handler is working (visible false)
- [test] view-tab-widgets > Update positions doesn't trigger updateGridPositions without changes
- [test] view-tab-widgets > Reflects active view screen mode from store
- [test] view-tab-widgets > Update positions triggers updateGridPositions with changes

### components/sidebars/alarm

#### components/sidebars/alarm/alarm.spec.js
Status: passed 43, failed 0, skipped 0

- [test] alarm > Create widget with default parameters
- [test] alarm > Duplicate widget without changes
- [test] alarm > Title changed after trigger field title
- [test] alarm > Periodic refresh changed after trigger field periodic refresh
- [test] alarm > Live watching changed after trigger field periodic refresh
- [test] alarm > Default sort column changed after trigger field default sort column
- [test] alarm > Widget columns changed after trigger field columns
- [test] alarm > Widget columns with `alarm.` prefix changed after trigger field columns
- [test] alarm > Widget group columns changed after trigger field columns
- [test] alarm > Service dependencies columns changed after trigger field columns
- [test] alarm > Default elements per page changed after trigger field default elements per page
- [test] alarm > Opened changed after trigger opened and resolved field
- [test] alarm > Filters changed after trigger update:filters on filters field
- [test] alarm > Filter changed after trigger input on filters field
- [test] alarm > Live reporting changed after trigger live reporting field
- [test] alarm > Info popup changed after trigger info popup field
- [test] alarm > More info template changed after trigger text field
- [test] alarm > Grid range changed after trigger grid range field
- [test] alarm > Clear filter disabled changed after trigger switcher field
- [test] alarm > Root cause settings changed after trigger root cause settings field
- [test] alarm > Availability field changed after trigger availability settings field
- [test] alarm > Html enabled on timeline changed after trigger switcher field
- [test] alarm > Open expandable panel on double click alarm changed after trigger switcher field
- [test] alarm > Ack note required changed after trigger switcher field
- [test] alarm > Multi ack enabled changed after trigger switcher field
- [test] alarm > Fast ack output changed after trigger fast ack output field
- [test] alarm > Fast cancel output changed after trigger fast cancel output field
- [test] alarm > Snooze note required changed after trigger switcher field
- [test] alarm > Remove alarms from meta alarms comment required changed after trigger switcher field
- [test] alarm > Export parameters changed after trigger export csv form
- [test] alarm > Sticky header changed after trigger switcher field
- [test] alarm > Sticky horizontal scroll changed after trigger switcher field
- [test] alarm > Invalid periodic refresh converted to valid object
- [test] alarm > Kiosk fields changed after trigger switcher field
- [test] alarm > Actions allowed with state ok changed after trigger switcher field
- [test] alarm > Keep selected after action changed after trigger switcher field
- [test] alarm > Actions allowed with state ok changed after trigger switcher field
- [test] alarm > Charts fields changed after trigger switcher field
- [test] alarm > Quick actions changed after trigger field quick alarm actions
- [test] alarm > Quick mass actions changed after trigger field quick alarm actions
- [test] alarm > Quick actions template changed after trigger field quick alarm actions
- [test] alarm > Quick mass actions template changed after trigger field quick alarm actions
- [test] alarm > Hide mass actions changed after trigger switcher field

#### components/sidebars/alarm/form/alarms-list-modal.spec.js
Status: passed 7, failed 0, skipped 0

- [test] alarms-list-modal > Sort columns changed after trigger sort columns field
- [test] alarms-list-modal > Sort template changed after trigger sort columns template field
- [test] alarms-list-modal > Columns changed after trigger columns field
- [test] alarms-list-modal > Items per page changed after trigger items per page field
- [test] alarms-list-modal > Info popups changed after trigger info popup field
- [test] alarms-list-modal > More info template changed after trigger text editor field
- [test] alarms-list-modal > Show root cause by state click changed after trigger switcher field

#### components/sidebars/alarm/form/fields/fast-action-output.spec.js
Status: passed 2, failed 0, skipped 0

- [test] fast-action-output > Enabled changed after trigger the enabled field
- [test] fast-action-output > Value changed after trigger the text field

#### components/sidebars/alarm/form/fields/info-popup.spec.js
Status: passed 1, failed 0, skipped 0

- [test] info-popup > Info popup setting modal opened after trigger create button

#### components/sidebars/alarm/form/fields/live-reporting.spec.js
Status: passed 3, failed 0, skipped 0

- [test] live-reporting > Live reporting modal opened after trigger create button
- [test] live-reporting > Text editor modal opened after trigger edit button
- [test] live-reporting > Value removed after trigger delete button

#### components/sidebars/alarm/form/fields/opened-resolved-filter.spec.js
Status: passed 3, failed 0, skipped 0

- [test] opened-resolved-filter > Value changed after to %s after trigger a radio field
- [test] opened-resolved-filter > <dynamic test name>

### components/sidebars/availability

#### components/sidebars/availability/availability.spec.js
Status: passed 6, failed 0, skipped 0

- [test] availability-settings > Create widget with default parameters
- [test] availability-settings > Duplicate widget without changes
- [test] availability-settings > All fields changed after input trigger
- [test] availability-settings > Widget columns changed after trigger update:widget-columns-template
- [test] availability-settings > Active columns changed after trigger update:active-alarms-columns-template
- [test] availability-settings > Resolved columns changed after trigger update:resolved-alarms-columns-template

#### components/sidebars/availability/form/availability-form.spec.js
Status: passed 8, failed 0, skipped 0

- [test] availability-form > Title changed after trigger title field
- [test] availability-form > Periodic refresh changed after trigger periodic refresh field
- [test] availability-form > Widget columns changed after trigger widget columns field
- [test] availability-form > Active alarms columns changed after trigger columns field
- [test] availability-form > REsolved alarms columns changed after trigger columns field
- [test] availability-form > Default time range changed after trigger quick interval type field
- [test] availability-form > Default display parameter changed after trigger default display parameter field
- [test] availability-form > Filters changed after trigger filters field

#### components/sidebars/availability/form/fields/availability-display-parameter.spec.js
Status: passed 1, failed 0, skipped 0

- [test] availability-display-parameter > Display parameter type changed after trigger radio field

#### components/sidebars/availability/form/fields/availability-display-show-type.spec.js
Status: passed 1, failed 0, skipped 0

- [test] availability-display-show-type > Show type changed after trigger radio field

#### components/sidebars/availability/form/fields/availability-graph-settings.spec.js
Status: passed 3, failed 0, skipped 0

- [test] availability-graph-settings > Enabled field changed after trigger enabled field
- [test] availability-graph-settings > Default show type changed after trigger availability show type field
- [test] availability-graph-settings > Default time range changed after trigger quick date interval field

### components/sidebars/chart

#### components/sidebars/chart/bar-chart.spec.js
Status: passed 3, failed 0, skipped 0

- [test] bar-chart > Create widget with default parameters
- [test] bar-chart > Duplicate widget without changes
- [test] bar-chart > All fields changed after input trigger

#### components/sidebars/chart/form/bar-chart-widget-form.spec.js
Status: passed 10, failed 0, skipped 0

- [test] bar-chart-widget-form > Title changed after trigger field title
- [test] bar-chart-widget-form > Periodic refresh changed after trigger field periodic refresh
- [test] bar-chart-widget-form > Preset changed after trigger preset field
- [test] bar-chart-widget-form > Metrics changed after trigger field alarm metric presets
- [test] bar-chart-widget-form > Graph type changed after trigger field bar graph type
- [test] bar-chart-widget-form > Chart title changed after trigger field chart title
- [test] bar-chart-widget-form > Filters changed after trigger field filters
- [test] bar-chart-widget-form > Quick date interval type changed after trigger field quick date interval
- [test] bar-chart-widget-form > Sampling changed after trigger field sampling
- [test] bar-chart-widget-form > Comparison changed after trigger field comparison

#### components/sidebars/chart/form/charts-form.spec.js
Status: passed 3, failed 0, skipped 0

- [test] charts-form > Show create %s chart modal
- [test] charts-form > <dynamic test name>

#### components/sidebars/chart/form/fields/alarm-metric-aggregate-function.spec.js
Status: passed 1, failed 0, skipped 0

- [test] alarm-metric-aggregate-function > Value changed after trigger color indicator field

#### components/sidebars/chart/form/fields/alarm-metric-presers.spec.js
Status: passed 1, failed 0, skipped 0

- [test] alarm-metric-presets > Value changed after trigger color indicator field

#### components/sidebars/chart/form/fields/bar-graph-type.spec.js
Status: passed 1, failed 0, skipped 0

- [test] bar-graph-type > Value changed after trigger radio group field

#### components/sidebars/chart/form/fields/pie-show-mode.spec.js
Status: passed 1, failed 0, skipped 0

- [test] pie-show-mode > Value changed after trigger radio group field

#### components/sidebars/chart/form/fields/sampling.spec.js
Status: passed 1, failed 0, skipped 0

- [test] sampling > Value changed after trigger sampling field

#### components/sidebars/chart/form/line-chart-widget-form.spec.js
Status: passed 8, failed 0, skipped 0

- [test] line-chart-widget-form > Title changed after trigger field title
- [test] line-chart-widget-form > Periodic refresh changed after trigger field periodic refresh
- [test] line-chart-widget-form > Metrics changed after trigger field alarm metric presets
- [test] line-chart-widget-form > Chart title changed after trigger field chart title
- [test] line-chart-widget-form > Filters changed after trigger field filters
- [test] line-chart-widget-form > Quick date interval type changed after trigger field quick date interval
- [test] line-chart-widget-form > Sampling changed after trigger field sampling
- [test] line-chart-widget-form > Comparison changed after trigger field comparison

#### components/sidebars/chart/form/numbers-widget-form.spec.js
Status: passed 10, failed 0, skipped 0

- [test] numbers-widget-form > Title changed after trigger field title
- [test] numbers-widget-form > Periodic refresh changed after trigger field periodic refresh
- [test] numbers-widget-form > Preset changed after trigger preset field
- [test] numbers-widget-form > Metrics changed after trigger field alarm metric presets
- [test] numbers-widget-form > Chart title changed after trigger field chart title
- [test] numbers-widget-form > Filters changed after trigger field filters
- [test] numbers-widget-form > Quick date interval type changed after trigger field quick date interval
- [test] numbers-widget-form > Sampling changed after trigger field sampling
- [test] numbers-widget-form > Trend changed after trigger field switcher
- [test] numbers-widget-form > Font size changed after trigger field font size

#### components/sidebars/chart/form/pie-chart-widget-form.spec.js
Status: passed 10, failed 0, skipped 0

- [test] pie-chart-widget-form > Title changed after trigger field title
- [test] pie-chart-widget-form > Periodic refresh changed after trigger field periodic refresh
- [test] pie-chart-widget-form > Preset changed after trigger preset field
- [test] pie-chart-widget-form > Metrics changed after trigger field alarm metric presets
- [test] pie-chart-widget-form > Show mode changed after trigger field pie show mode
- [test] pie-chart-widget-form > Chart title changed after trigger field chart title
- [test] pie-chart-widget-form > Filters changed after trigger field filters
- [test] pie-chart-widget-form > Quick date interval type changed after trigger field quick date interval
- [test] pie-chart-widget-form > Sampling changed after trigger field sampling
- [test] pie-chart-widget-form > Aggregate function changed after trigger field aggregate function

#### components/sidebars/chart/line-chart.spec.js
Status: passed 3, failed 0, skipped 0

- [test] line-chart > Create widget with default parameters
- [test] line-chart > Duplicate widget without changes
- [test] line-chart > All fields changed after input trigger

#### components/sidebars/chart/numbers.spec.js
Status: passed 3, failed 0, skipped 0

- [test] numbers > Create widget with default parameters
- [test] numbers > Duplicate widget without changes
- [test] numbers > All fields changed after input trigger

#### components/sidebars/chart/pie-chart.spec.js
Status: passed 3, failed 0, skipped 0

- [test] pie-chart > Create widget with default parameters
- [test] pie-chart > Duplicate widget without changes
- [test] pie-chart > All fields changed after input trigger

### components/sidebars/context

#### components/sidebars/context/context.spec.js
Status: passed 20, failed 0, skipped 0

- [test] context > Create widget with default parameters
- [test] context > Duplicate widget without changes
- [test] context > Title changed after trigger field title
- [test] context > Default sort column changed after trigger field default sort column
- [test] context > Widget columns changed after trigger field columns
- [test] context > Widget columns with `entity.` prefix changed after trigger field columns
- [test] context > Service dependencies columns changed after trigger field columns
- [test] context > Widget active alarms columns changed after trigger field columns
- [test] context > Widget resolved alarms columns changed after trigger field columns
- [test] context > Widget tree of dependencies settings changed after trigger field tree of deps
- [test] context > Widget tree of dependencies settings changed after trigger field tree of deps
- [test] context > Filters changed after trigger update:filters on filters field
- [test] context > Filter changed after trigger input on filters field
- [test] context > Root cause settings changed after trigger root cause settings field
- [test] context > Availability field changed after trigger availability settings field
- [test] context > Export parameters changed after trigger export csv form
- [test] context > Charts fields changed after trigger switcher field
- [test] context > Grid range changed after trigger grid range field
- [test] context > Default inherited pbehavior changed after trigger default inherited pbehavior switcher
- [test] context > Keep selected after action changed after trigger keep selected after action switcher

### components/sidebars/counter

#### components/sidebars/counter/form/fields/template.spec.js
Status: passed 1, failed 0, skipped 0

- [test] field-template > Value changed after trigger color indicator field

### components/sidebars/form

#### components/sidebars/form/button-field.spec.js
Status: passed 3, failed 0, skipped 0

- [test] button-field > Create event emitted after click on the button
- [test] button-field > Edit event emitted after click on the button
- [test] button-field > Delete event emitted after click on the button

#### components/sidebars/form/color-indicator.spec.js
Status: passed 1, failed 0, skipped 0

- [test] field-color-indicator > Value changed after trigger color indicator field

#### components/sidebars/form/columns.spec.js
Status: passed 1, failed 0, skipped 0

- [test] columns > Columns changed after trigger columns field

#### components/sidebars/form/default-elements-per-page.spec.js
Status: passed 1, failed 0, skipped 0

- [test] default-elements-per-page > Value changed after trigger select field

#### components/sidebars/form/default-sort-column.spec.js
Status: passed 2, failed 0, skipped 0

- [test] default-sort-column > Column changed after trigger select field with columns
- [test] default-sort-column > Order changed after trigger select field with orders

#### components/sidebars/form/export-csv.spec.js
Status: passed 3, failed 0, skipped 0

- [test] export-csv > Separator changed after trigger separator select field
- [test] export-csv > Datetime format changed after trigger datetime format select field
- [test] export-csv > Columns changed after trigger columns field

#### components/sidebars/form/fields/filters-list.spec.js
Status: passed 5, failed 0, skipped 0

- [test] filters-list > Filters updated after trigger input event on filters list
- [test] filters-list > Filters created after trigger add event on filters list
- [test] filters-list > Filters edited after trigger edit event on filters list
- [test] filters-list > Filter deleted after trigger delete event on filters list
- [test] filters-list > Current filter deleted after trigger delete event on filters list

#### components/sidebars/form/fields/grid-range-size.spec.js
Status: passed 1, failed 0, skipped 0

- [test] grid-range-size > Value changed after trigger range slider field

#### components/sidebars/form/fields/quick-date-interval-type.spec.js
Status: passed 1, failed 0, skipped 0

- [test] quick-date-interval-type > Value changed after trigger number field

#### components/sidebars/form/filters.spec.js
Status: passed 2, failed 0, skipped 0

- [test] filters > Selected filters updated after trigger input on the filter selector field
- [test] filters > Filters updated after trigger input event on filters list

#### components/sidebars/form/grid-size.spec.js
Status: passed 1, failed 0, skipped 0

- [test] field-grid-size > Column size changed after trigger field

#### components/sidebars/form/margins.spec.js
Status: passed 4, failed 0, skipped 0

- [test] margins > Top margin changed after trigger slider field
- [test] margins > Right margin changed after trigger slider field
- [test] margins > Bottom margin changed after trigger slider field
- [test] margins > Left margin changed after trigger slider field

#### components/sidebars/form/number.spec.js
Status: passed 1, failed 0, skipped 0

- [test] field-number > Value changed after trigger number field

#### components/sidebars/form/periodic-refresh.spec.js
Status: passed 2, failed 0, skipped 0

- [test] periodic-refresh > Unit as seconds settled created, if unit doesn't exist
- [test] periodic-refresh > Value changed after trigger periodic refresh field
- [test] periodic-refresh > Live watching triggers input on changes

#### components/sidebars/form/slider.spec.js
Status: passed 1, failed 0, skipped 0

- [test] field-slider > Value changed after trigger number field

#### components/sidebars/form/switcher.spec.js
Status: passed 1, failed 0, skipped 0

- [test] switcher > Value changed after trigger switch field

#### components/sidebars/form/text-editor.spec.js
Status: passed 3, failed 0, skipped 0

- [test] text-editor > Text editor modal opened after trigger create button
- [test] text-editor > Text editor modal opened after trigger edit button
- [test] text-editor > Confirmation modal opened after trigger delete button

#### components/sidebars/form/title.spec.js
Status: passed 1, failed 0, skipped 0

- [test] title > Value changed after trigger text field

#### components/sidebars/form/tree-of-dependencies-settings.spec.js
Status: passed 1, failed 0, skipped 0

- [test] tree-of-dependencies-settings > Value changed after trigger radio group field

### components/sidebars/map

#### components/sidebars/map/form/fields/map.spec.js
Status: passed 1, failed 0, skipped 0

- [test] field-map > Info popup setting modal opened after trigger create button

#### components/sidebars/map/map.spec.js
Status: passed 11, failed 0, skipped 0

- [test] map > Create widget with default parameters
- [test] map > Duplicate widget without changes
- [test] map > Title changed after trigger field title
- [test] map > Periodic refresh changed after trigger field periodic refresh
- [test] map > Map changed after trigger field map
- [test] map > Color indicator changed after trigger field color indicator
- [test] map > Under pbehavior enabled changed after trigger field switcher
- [test] map > Filters changed after trigger field filters
- [test] map > Entity info template changed after trigger field text editor
- [test] map > Alarms columns changed after trigger field columns
- [test] map > Entities columns changed after trigger field columns

### components/sidebars/service-weather

#### components/sidebars/service-weather/form/fields/counters-selector.spec.js
Status: passed 4, failed 0, skipped 0

- [test] field-counters-selector > Pbehavior counters enabled after trigger enabled field
- [test] field-counters-selector > Pbehavior counters changed after trigger pbehavior type field
- [test] field-counters-selector > State counters enabled after trigger enabled field
- [test] field-counters-selector > State counters changed after trigger state type field

#### components/sidebars/service-weather/form/fields/modal-type.spec.js
Status: passed 1, failed 0, skipped 0

- [test] field-modal-type > Value changed after trigger radio group

#### components/sidebars/service-weather/form/fields/sort-column.spec.js
Status: passed 2, failed 0, skipped 0

- [test] field-sort-column > Column changed after trigger combobox field
- [test] field-sort-column > Order changed after trigger select field

#### components/sidebars/service-weather/service-weather.spec.js
Status: passed 27, failed 0, skipped 0

- [test] service-weather > Create widget with default parameters
- [test] service-weather > Duplicate widget without changes
- [test] service-weather > Title changed after trigger field title
- [test] service-weather > Periodic refresh changed after trigger field periodic refresh
- [test] service-weather > Filters changed after trigger field filters
- [test] service-weather > Alarms list modal changed after trigger alarms list modal form
- [test] service-weather > Limit changed after trigger number field
- [test] service-weather > Color indicator changed after trigger field color indicator
- [test] service-weather > Sort changed after trigger sort column field
- [test] service-weather > Modal items per page changed after trigger default per page field
- [test] service-weather > Block template changed after trigger field template
- [test] service-weather > Modal template changed after trigger field template
- [test] service-weather > Entity template changed after trigger field template
- [test] service-weather > Column mobile changed after trigger field grid size
- [test] service-weather > Column tablet changed after trigger field grid size
- [test] service-weather > Column desktop changed after trigger field grid size
- [test] service-weather > Margins changed after trigger margins form
- [test] service-weather > Height factor changed after trigger slider field
- [test] service-weather > Counters changed after trigger counters field
- [test] service-weather > Is priority enabled changed after trigger switcher field
- [test] service-weather > Is hide gray enabled changed after trigger switcher field
- [test] service-weather > Is secondary icon enabled changed after trigger switcher field
- [test] service-weather > Is entities actions in queue changed after trigger switcher field
- [test] service-weather > Modal type changed after trigger modal type field
- [test] service-weather > Action required settings changed after trigger action required settings field
- [test] service-weather > Root cause settings changed after trigger root cause settings field
- [test] service-weather > Default inherited pbehavior changed after trigger default inherited pbehavior switcher

### components/sidebars/text

#### components/sidebars/text/text.spec.js
Status: passed 4, failed 0, skipped 0

- [test] text > Create widget with default parameters
- [test] text > Duplicate widget without changes
- [test] text > Title changed after trigger field title
- [test] text > Template changed after trigger field text editor

### components/table/c-ellipsis.spec.js

#### components/table/c-ellipsis.spec.js
Status: passed 11, failed 0, skipped 0

- [test] c-ellipsis > Text letters count less then default max letters
- [test] c-ellipsis > Text letters count more then default max letters
- [test] c-ellipsis > Text letters count less then custom maxLetters
- [test] c-ellipsis > Text letters count more then custom maxLetters
- [test] c-ellipsis > Click on text with text letters count less then default max letters
- [test] c-ellipsis > Click on text with text letters count more then default max letters
- [test] c-ellipsis > hover > Text letters count less then default max letters
- [test] c-ellipsis > hover > Text letters count more then default max letters
- [test] c-ellipsis > hover > Text letters count less then custom maxLetters
- [test] c-ellipsis > hover > Text letters count more then custom maxLetters
- [test] c-ellipsis > hover > Click on activator emits textClicked

### components/widgets/alarm

#### components/widgets/alarm/actions/actions-panel.spec.js
Status: passed 0, failed 30, skipped 0

- [test] actions-panel > quickActions returns [] if widget.parameters.quickActions is undefined
- [test] actions-panel > quickActions returns array if widget.parameters.quickActions is set
- [test] actions-panel > preparedActions sorts actions according to quickActions order
- [test] actions-panel > additionalProps.inlineCount reflects quickActions count
- [test] actions-panel > Ack modal showed after trigger ack action
- [test] actions-panel > Fast ack event sent after trigger fast ack action
- [test] actions-panel > Ack remove modal showed after trigger ack remove action
- [test] actions-panel > Create pbehavior modal showed after trigger pbehavior add action
- [test] actions-panel > Fast pbehavior add creates downtime pbehavior and calls refreshAlarmsList
- [test] actions-panel > Fast pbehavior remove removes downtime pbehavior and calls refreshAlarmsList
- [test] actions-panel > Snooze modal showed after trigger snooze action
- [test] actions-panel > Declare ticket modal showed after trigger declare action
- [test] actions-panel > Associate ticket modal showed after trigger associate ticket action
- [test] actions-panel > Remove associated ticket modal showed after trigger remove associated ticket action
- [test] actions-panel > Change state modal showed after trigger change state action
- [test] actions-panel > Cancel modal showed after trigger cancel action
- [test] actions-panel > Fast cancel event sent after trigger fast cancel action
- [test] actions-panel > Fast cancel and cancel action available when status is flapping and state is ok
- [test] actions-panel > Fast cancel and cancel action available when status is closed and state is ok
- [test] actions-panel > Fast cancel and cancel action doesn't available when status is closed and state is ok and isActionsAllowWithOkState is disabled
- [test] actions-panel > Variables modal showed after trigger variables help action
- [test] actions-panel > History modal showed after trigger history action
- [test] actions-panel > Comment modal showed after trigger comment action
- [test] actions-panel > Remove alarms from manual meta alarm modal showed after trigger remove alarms from manual meta alarm action
- [test] actions-panel > Remove alarms from auto meta alarm modal showed after trigger remove alarms from auto meta alarm action
- [test] actions-panel > Execute instruction alarm modal showed after trigger execute instruction action
- [test] actions-panel > Export PDF action
- [test] actions-panel > Custom action called after trigger button
- [test] actions-panel > Add bookmark request was sent after trigger add bookmark
- [test] actions-panel > Remove bookmark request was sent after trigger remove bookmark

#### components/widgets/alarm/actions/mass-actions-panel.spec.js
Status: passed 17, failed 0, skipped 0

- [test] mass-actions-panel > Create pbehavior modal showed after trigger pbehavior add action
- [test] mass-actions-panel > Fast pbehavior add creates downtime pbehavior and emits clear:items
- [test] mass-actions-panel > Fast pbehavior remove removes downtime pbehavior and emits clear:items
- [test] mass-actions-panel > Ack modal showed after trigger ack action
- [test] mass-actions-panel > Fast ack event sent after trigger fast ack action
- [test] mass-actions-panel > Fast ack event sent after trigger fast ack action without parameters
- [test] mass-actions-panel > Ack remove modal showed after trigger ack remove action
- [test] mass-actions-panel > Cancel modal showed after trigger cancel action
- [test] mass-actions-panel > Fast cancel event sent after trigger fast cancel action
- [test] mass-actions-panel > Fast cancel event sent after trigger fast cancel action without parameters
- [test] mass-actions-panel > Associate ticket modal showed after trigger associate ticket action
- [test] mass-actions-panel > Remove associated ticket modal showed after trigger remove associated ticket action
- [test] mass-actions-panel > Snooze modal showed after trigger snooze action
- [test] mass-actions-panel > Manual meta alarm group modal showed after trigger manual meta alarm group action
- [test] mass-actions-panel > Comment modal showed after trigger comment action
- [test] mass-actions-panel > Change state modal showed after trigger change state action
- [test] mass-actions-panel > inlineCount reflects quickMassActions count

#### components/widgets/alarm/alarms-list.spec.js
Status: passed 28, failed 2, skipped 0

- [test] alarms-list > Query updated after mount
- [test] alarms-list > User preferences not fetched after mount with local widget prop
- [test] alarms-list > Correlation updated after trigger correlation field
- [test] alarms-list > Only bookmark updated after trigger filter by bookmark field
- [test] alarms-list > Tags updated after trigger alarm tag field
- [test] alarms-list > Filter updated after trigger filter field
- [test] alarms-list > Instruction filters updated after trigger filter field
- [test] alarms-list > Interval query removed after click on the button
- [test] alarms-list > Interval modal showed after click on the live reporting button
- [test] alarms-list > Category updated after trigger category field
- [test] alarms-list > Limit updated after trigger table pagination
- [test] alarms-list > Page updated after trigger table pagination
- [test] alarms-list > Widget exported after trigger export button
- [test] alarms-list > Widget exported after trigger export button with invalid structure
- [test] alarms-list > Widget exported after trigger export button without export columns
- [test] alarms-list > Widget exported after trigger export button with long request time
- [test] alarms-list > Error popup showed exported after trigger export button with failed create export
- [test] alarms-list > Error popup showed exported after trigger export button with failed fetch export
- [test] alarms-list > Error popup showed exported after trigger export button with failed status fetch export
- [test] alarms-list > Alarms not fetched after change query without columns
- [test] alarms-list > Alarms fetched after change query
- [test] alarms-list > Alarms fetched after change query nonce
- [test] alarms-list > Periodic started after mount with enabled value
- [test] alarms-list > Interval cleared after update periodic refresh
- [test] alarms-list > Interval cleared after destroy
- [test] alarms-list > Search updated after alarm-advanced-search-field submit trigger
- [test] alarms-list > Patterns updated after alarm-advanced-search-field submit trigger
- [test] alarms-list > Search field cleared after alarm-advanced-search-field reset trigger
- [test] alarms-list > Search pinned after alarm-advanced-search-field toggle-pin:search trigger
- [test] alarms-list > Search removed after alarm-advanced-search-field remove:search trigger

#### components/widgets/alarm/column-formatting/alarm-column-cell-popup-body.spec.js
Status: passed 0, failed 1, skipped 0

- [test] alarm-column-cell-popup-body > Popup template closed after click on the button

#### components/widgets/alarm/column-formatting/alarm-column-cell.spec.js
Status: passed 22, failed 1, skipped 0

- [test] alarm-column-cell > Default filter for date field: "%s" converted value to time
- [test] alarm-column-cell > <dynamic test name>
- [test] alarm-column-cell > Default filter for duration field: "%s" converted value to duration
- [test] alarm-column-cell > <dynamic test name>
- [test] alarm-column-cell > Custom filter for field converted value correctly
- [test] alarm-column-cell > Applies filter class when column has isFilter set to true
- [test] alarm-column-cell > Does not apply filter class when column has isFilter set to false
- [test] alarm-column-cell > Calls selectAdvancedSearchField when filter column is clicked
- [test] alarm-column-cell > Calls selectAdvancedSearchField with sanitized HTML value when column is HTML and filter
- [test] alarm-column-cell > Calls selectAdvancedSearchField with filtered value when column has custom filter
- [test] alarm-column-cell > Does not call selectAdvancedSearchField when column is not a filter
- [test] alarm-column-cell > Calls both custom click handler and selectAdvancedSearchField when column has both
- [test] alarm-column-cell > Handles empty or null values when filter column is clicked
- [test] alarm-column-cell > Handles undefined selectAdvancedSearchField gracefully

#### components/widgets/alarm/column-formatting/alarm-column-value.spec.js
Status: passed 1, failed 0, skipped 0

- [test] alarm-column-value > Click state emitted after trigger click state event on alarm cell

#### components/widgets/alarm/entity-gantt/entity-gantt.spec.js
Status: passed 0, failed 5, skipped 0

- [test] entity-gantt > Gantt intervals fetched after mount
- [test] entity-gantt > Gantt intervals fetched after change the query
- [test] entity-gantt > Error popup showed after fetch gantt intervals with error message
- [test] entity-gantt > Error popup showed after fetch gantt intervals with error description
- [test] entity-gantt > Error popup showed after fetch gantt intervals with error

#### components/widgets/alarm/expand-panel/alarms-expand-panel-entity-enrichments.spec.js
Status: passed 1, failed 1, skipped 0

- [test] alarms-expand-panel-entity-enrichments > Should call fetchEntityInfosLogsListWithoutStore on mount
- [test] alarms-expand-panel-entity-enrichments > Should render entity-infos-logs-list component

#### components/widgets/alarm/expand-panel/alarms-expand-panel-remediation.spec.js
Status: passed 2, failed 0, skipped 0

- [test] alarms-expand-panel-remediation > Check fetchList calling on mounted
- [test] alarms-expand-panel-remediation > Check fetchList calling on options updating

#### components/widgets/alarm/expand-panel/alarms-expand-panel.spec.js
Status: passed 0, failed 1, skipped 0

- [test] alarms-expand-panel > Tab key updated after change moreInfoTemplate

#### components/widgets/alarm/forms/ack-event-form.spec.js
Status: passed 1, failed 1, skipped 0

- [test] ack-event-form > Comment changed after trigger description field
- [test] ack-event-form > Ack resources changed after trigger checkbox field

#### components/widgets/alarm/forms/info-popup-form.spec.js
Status: passed 1, failed 1, skipped 0

- [test] info-popup-form > Column changed after trigger select field
- [test] info-popup-form > Template changed after trigger text editor field

#### components/widgets/alarm/forms/link-meta-alarm-form.spec.js
Status: passed 0, failed 3, skipped 0

- [test] link-meta-alarm-form > Meta alarm changed after trigger text field
- [test] link-meta-alarm-form > Comment changed after trigger description field
- [test] link-meta-alarm-form > Auto resolve changed after trigger enabled field

#### components/widgets/alarm/forms/snooze-event-form.spec.js
Status: passed 2, failed 0, skipped 0

- [test] snooze-event-form > Duration changed after trigger duration field
- [test] snooze-event-form > Output changed after trigger description field

#### components/widgets/alarm/partials/alarm-status-chip-with-relations.spec.js
Status: passed 0, failed 2, skipped 0

- [test] alarm-status-chip-with-relations > Shows modal with correct config when button clicked and entity has upstream
- [test] alarm-status-chip-with-relations > Shows modal with correct config when button clicked and entity has no upstream

#### components/widgets/alarm/partials/alarms-list-remediation-instructions-filter.spec.js
Status: passed 1, failed 2, skipped 0

- [test] alarms-list-remediation-instructions-filter > Emits input with empty object when cleared
- [test] alarms-list-remediation-instructions-filter > Does not emit input if value is not empty
- [test] alarms-list-remediation-instructions-filter > Computes selectedItems and selectedAllFields correctly

#### components/widgets/alarm/partials/alarms-list-row-instructions-icon.spec.js
Status: passed 17, failed 0, skipped 0

- [test] alarms-list-row-instructions-icon > Icon computation > Shows assignment_warning icon when withoutAnyExecution is true
- [test] alarms-list-row-instructions-icon > Icon computation > Shows manual_instruction icon when manual instruction is available
- [test] alarms-list-row-instructions-icon > Icon computation > Shows assignment icon when auto instruction is available
- [test] alarms-list-row-instructions-icon > Icon computation > Defaults to manualAvailable when instruction_execution_icon is not provided
- [test] alarms-list-row-instructions-icon > Icon classes > Applies warning class for withoutAnyExecution
- [test] alarms-list-row-instructions-icon > Icon classes > Applies blinking and dotted classes for running instructions
- [test] alarms-list-row-instructions-icon > Icon classes > Applies failed class for failed instructions
- [test] alarms-list-row-instructions-icon > Icon classes > Applies completed class for successful instructions
- [test] alarms-list-row-instructions-icon > Tooltip content > Shows running manual instruction tooltip
- [test] alarms-list-row-instructions-icon > Tooltip content > Shows running auto instruction tooltip
- [test] alarms-list-row-instructions-icon > Tooltip content > Shows failed manual instruction tooltip
- [test] alarms-list-row-instructions-icon > Tooltip content > Shows failed auto instruction tooltip
- [test] alarms-list-row-instructions-icon > Tooltip content > Shows successful manual instruction tooltip
- [test] alarms-list-row-instructions-icon > Tooltip content > Shows successful auto instruction tooltip
- [test] alarms-list-row-instructions-icon > Tooltip content > Shows withoutAnyExecution tooltip
- [test] alarms-list-row-instructions-icon > Tooltip content > Shows hasManualInstruction tooltip when not withoutAnyExecution
- [test] alarms-list-row-instructions-icon > Tooltip content > Combines multiple instruction tooltips

#### components/widgets/alarm/partials/alarms-list-row.spec.js
Status: passed 1, failed 5, skipped 0

- [test] alarms-list-row > Alarm selected after trigger checkbox
- [test] alarms-list-row > Listeners from feature called after trigger
- [test] alarms-list-row > Row expanded after trigger expand button with hidden groups and without filtered alarms
- [test] alarms-list-row > Row closed after trigger expand button with expanded: true
- [test] alarms-list-row > Click state emitted after trigger click state event
- [test] alarms-list-row > Expand panel opened after double click on row when openExpandPanelByDoubleClick enabled

#### components/widgets/alarm/partials/alarms-list-table.spec.js
Status: passed 18, failed 0, skipped 0

- [test] alarms-list-table > Alarms selected after trigger table
- [test] alarms-list-table > Pagination update event emitted after trigger update pagination
- [test] alarms-list-table > Resize listener added after mount and removed after destroy
- [test] alarms-list-table > Timer cleared after disable sticky
- [test] alarms-list-table > Component adds and removes the same count listeners
- [test] alarms-list-table > Header position changed after trigger scroll
- [test] alarms-list-table > Header hidden after trigger start scroll
- [test] alarms-list-table > Expanded elements works correctly
- [test] alarms-list-table > Root cause diagram opened after trigger click state event
- [test] alarms-list-table > Table pagination emits input event correctly
- [test] alarms-list-table > Fetches entity info properties list on mount
- [test] alarms-list-table > Shows loading state when entityInfoPropertyPending is true
- [test] alarms-list-table > clearSelected with force=true clears selection regardless of keepSelectedAfterAction
- [test] alarms-list-table > clearSelected with force=false clears selection when keepSelectedAfterAction is false
- [test] alarms-list-table > clearSelected with force=false does not clear selection when keepSelectedAfterAction is true
- [test] alarms-list-table > keepSelectedAfterAction initialized from widget.parameters.keepSelectedAfterAction
- [test] alarms-list-table > updateKeepSelectedAfterAction updates keepSelectedAfterAction value
- [test] alarms-list-table > c-mass-actions-panel clear:selected clears selection

#### components/widgets/alarm/timeline/alarm-timeline-step.spec.js
Status: passed 0, failed 1, skipped 0

- [test] alarm-timeline-step > Expand button emits event

#### components/widgets/alarm/timeline/alarm-timeline-steps.spec.js
Status: passed 0, failed 2, skipped 0

- [test] alarm-timeline-steps > Toggles the expanded state when expand button is clicked
- [test] alarm-timeline-steps > Toggles the expanded state when expand button is clicked for second step

#### components/widgets/alarm/timeline/alarm-timeline.spec.js
Status: passed 3, failed 1, skipped 0

- [test] alarm-timeline > User can update group value and emit update event
- [test] alarm-timeline > User can enable only comments filter and emit update event
- [test] alarm-timeline > User can update page value and emit update event
- [test] alarm-timeline > Days computed property returns expected value

### components/widgets/availability

#### components/widgets/availability/availability-widget.spec.js
Status: passed 7, failed 2, skipped 0

- [test] availability-widget > Query changed after component mount
- [test] availability-widget > Data fetched with correct query
- [test] availability-widget > Data fetched when correct query when trend enabled
- [test] availability-widget > Filters changed after trigger availability filters
- [test] availability-widget > Interval changed after trigger availability filters
- [test] availability-widget > Trend changed after trigger availability filters
- [test] availability-widget > Type changed after trigger availability filters
- [test] availability-widget > Display parameter changed after trigger availability filters
- [test] availability-widget > Value filter changed after trigger availability filters

#### components/widgets/availability/partials/availability-widget-filters.spec.js
Status: passed 0, failed 10, skipped 0

- [test] availability-widget-filters > Interval changed after trigger quick date interval field
- [test] availability-widget-filters > Search changed after trigger search field
- [test] availability-widget-filters > Filters changed after trigger filters field
- [test] availability-widget-filters > Filters changed after trigger filters field
- [test] availability-widget-filters > Display parameter changed after trigger availability display parameter field
- [test] availability-widget-filters > Show type changed after trigger availability show type field
- [test] availability-widget-filters > Trend changed after trigger enabled field
- [test] availability-widget-filters > Value filter changed after trigger value filter field
- [test] availability-widget-filters > Local value filter changed after update props
- [test] availability-widget-filters > Export emitted after click on button

### components/widgets/chart

#### components/widgets/chart/bar-chart-widget.spec.js
Status: passed 1, failed 1, skipped 0

- [test] bar-chart-widget > Query updated after mount
- [test] bar-chart-widget > Vectors metrics fetched with correct query

#### components/widgets/chart/entity-chart-widget.spec.js
Status: passed 0, failed 2, skipped 0

- [test] entity-chart-widget > Vector metrics fetched after mount
- [test] entity-chart-widget > Aggregated metrics fetched after mount

#### components/widgets/chart/line-chart-widget.spec.js
Status: passed 1, failed 1, skipped 0

- [test] line-chart-widget > Query updated after mount
- [test] line-chart-widget > Vectors metrics fetched with correct query

#### components/widgets/chart/numbers-widget.spec.js
Status: passed 1, failed 1, skipped 0

- [test] numbers-widget > Query updated after mount
- [test] numbers-widget > Vectors metrics fetched with correct query

#### components/widgets/chart/pie-chart-widget.spec.js
Status: passed 1, failed 1, skipped 0

- [test] pie-chart-widget > Query updated after mount
- [test] pie-chart-widget > Vectors metrics fetched with correct query

### components/widgets/map

#### components/widgets/map/map.spec.js
Status: passed 8, failed 0, skipped 0

- [test] map > Register and unregister editing off handler is working
- [test] map > Map state fetched after component mount in editing mode
- [test] map > Map changed after map preview was triggered
- [test] map > Map changed after breadcrumbs was triggered
- [test] map > Maps history cleared after trigger editing mode
- [test] map > Category query updated after trigger category field
- [test] map > Alarms list modal showed after trigger map preview
- [test] map > Opened alarms fetched correct after trigger fetch list from modal

#### components/widgets/map/partial/map-breadcrumbs.spec.js
Status: passed 1, failed 0, skipped 0

- [test] map-breadcrumbs > Click emitted after button click was triggered

### components/widgets/partials

#### components/widgets/partials/kpi-widget-filters.spec.js
Status: passed 2, failed 1, skipped 0

- [test] kpi-widget-filters > Date interval changed after trigger quick date interval field
- [test] kpi-widget-filters > Sampling changed after trigger sampling field
- [test] kpi-widget-filters > Filters changed after trigger filters selector

### components/widgets/service-weather

#### components/widgets/service-weather/service-pause-event-form.spec.js
Status: passed 2, failed 0, skipped 0

- [test] service-pause-event-form > Reason changed after trigger reason field
- [test] service-pause-event-form > Comment changed after trigger textarea

#### components/widgets/service-weather/service-weather-item.spec.js
Status: passed 3, failed 1, skipped 0

- [test] service-weather-item > Alarms list modal showed after click on button
- [test] service-weather-item > Main information modal showed after click on card
- [test] service-weather-item > Show root cause emitted after trigger click on show root cause button
- [test] service-weather-item > Modal doesn't show after click on link

#### components/widgets/service-weather/service-weather.spec.js
Status: passed 0, failed 10, skipped 0

- [test] service-weather > Query updated after mount
- [test] service-weather > Services list fetched with correct query
- [test] service-weather > Services list fetched with custom limit query
- [test] service-weather > Category query updated after trigger category field
- [test] service-weather > Filter updated after trigger filter field
- [test] service-weather > Gray filter updated after trigger enabled field
- [test] service-weather > Alarms list modal showed after click on button
- [test] service-weather > Alarms list modal showed after click on card
- [test] service-weather > Root cause modal showed after trigger show:root-cause event
- [test] service-weather > Main information modal showed after click on card

### helpers/date/date.spec.js

#### helpers/date/date.spec.js
Status: passed 7, failed 0, skipped 0

- [test] convertDateToString > Date as timestamp with timestamp format
- [test] convertDateToString > Convert date to formats
- [test] convertDateToString > Date as unix
- [test] convertDateToString > Date as native Date object
- [test] convertDateToString > Date as ISO string
- [test] convertDateToString > Invalid date without default value
- [test] convertDateToString > Invalid date with default value

### helpers/date/duration.spec.js

#### helpers/date/duration.spec.js
Status: passed 9, failed 0, skipped 0

- [test] convertUnit > convertUnit from seconds
- [test] convertUnit > convertUnit from minute
- [test] convertUnit > convertUnit from hour
- [test] convertUnit > convertUnit from day
- [test] convertUnit > convertUnit from week
- [test] convertUnit > convertUnit from month
- [test] convertUnit > convertUnit from year
- [test] convertUnit > convertDurationToString in default locale
- [test] convertUnit > convertDurationToString in fr locale with default arguments

### helpers/entities/alarm

#### helpers/entities/alarm/query.spec.js
Status: passed 4, failed 0, skipped 0

- [test] convertAlarmUserPreferenceToQuery > should convert user preferences to a query when all fields are provided
- [test] convertAlarmUserPreferenceToQuery > should handle empty content object gracefully
- [test] convertAlarmWidgetToQuery > should convert widget with 'AlarmsList' type to a query object with default parameters and isCorrelationEnabled
- [test] convertAlarmWidgetToQuery > should handle empty widget.parameters gracefully

### helpers/entities/link

#### helpers/entities/link/list.spec.js
Status: passed 8, failed 0, skipped 0

- [test] entities links list helper > Correct harmonize category links correct arguments
- [test] entities links list helper > Correct harmonize category links without arguments
- [test] entities links list helper > Correct harmonize links with correct arguments
- [test] entities links list helper > Correct harmonize alarms links with alarms without intersection in links
- [test] entities links list helper > Correct harmonize alarms links with alarms with intersection in links
- [test] entities links list helper > Correct harmonize alarms links with one alarm with only old links
- [test] entities links list helper > Correct harmonize alarms links with one alarm
- [test] entities links list helper > Correct harmonize links without arguments

### helpers/entities/remediation

#### helpers/entities/remediation/instruction-execution/form.spec.js
Status: passed 29, failed 0, skipped 0

- [test] helpers/entities/remediation/instruction-execution/form > hasInstructionWithoutAnyExecution > Returns true for withoutAnyExecution icon
- [test] helpers/entities/remediation/instruction-execution/form > hasInstructionWithoutAnyExecution > Returns false for other icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconManualInProgress > Returns true for manual in progress icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconManualInProgress > Returns false for non-manual in progress icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconAutoInProgress > Returns true for auto in progress icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconAutoInProgress > Returns false for non-auto in progress icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconInProgress > Returns true for any in progress icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconInProgress > Returns false for non-in progress icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconFailed > Returns true for failed icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconFailed > Returns false for non-failed icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconSuccess > Returns true for successful icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionIconSuccess > Returns false for non-successful icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionManual > Returns true for manual icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionManual > Returns false for non-manual icons
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionCompleted > Returns true for completed status
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionCompleted > Returns false for non-completed statuses
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionCompleted > Returns false for undefined object
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionFailed > Returns true for failed status
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionFailed > Returns false for non-failed statuses
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionFailed > Returns false for undefined object
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionAborted > Returns true for aborted status
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionAborted > Returns false for non-aborted statuses
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionAborted > Returns false for undefined object
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionPaused > Returns true for paused status
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionPaused > Returns false for non-paused statuses
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionPaused > Returns false for undefined object
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionRunning > Returns true for running status
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionRunning > Returns false for non-running statuses
- [test] helpers/entities/remediation/instruction-execution/form > isInstructionExecutionRunning > Returns false for undefined object

### helpers/entities/widget

#### helpers/entities/widget/query.spec.js
Status: passed 4, failed 0, skipped 0

- [test] prepareWidgetQuery > should convert widget and userPreference to query objects and merge them
- [test] prepareWidgetQuery > should convert widget and userPreference (with isCorrelationEnabled) to query objects and merge them
- [test] prepareWidgetQuery > should convert widget (with isCorrelationEnabled) and userPreference to query objects and merge them
- [test] prepareWidgetQuery > should convert widget (with isCorrelationEnabled) and userPreference (with disabled isCorrelationEnabled) to query objects and merge them

### helpers/forms/pattern.spec.js

#### helpers/forms/pattern.spec.js
Status: passed 200, failed 0, skipped 0

- [test] pattern form converters > should be converted to form and back to pattern with `equal` operator
- [test] pattern form converters > should be converted to form and back to pattern with `not equal` operator
- [test] pattern form converters > should be converted to form and back to pattern with `contains` operator
- [test] pattern form converters > should be converted to form and back to pattern with `not contains` operator
- [test] pattern form converters > should be converted to form and back to pattern with `regexp` operator
- [test] pattern form converters > should be converted to form and back to pattern with `acked` operator
- [test] pattern form converters > should be converted to form and back to pattern with `not acked` operator
- [test] pattern form converters > should be converted to form and back to pattern with `snoozed` operator
- [test] pattern form converters > should be converted to form and back to pattern with `not snoozed` operator
- [test] pattern form converters > should be converted to form and back to pattern with `canceled` operator
- [test] pattern form converters > should be converted to form and back to pattern with `not canceled` operator
- [test] pattern form converters > should be converted to form and back to pattern with `ticket associated` operator
- [test] pattern form converters > should be converted to form and back to pattern with `ticket not associated` operator
- [test] pattern form converters > should be converted to form and back to pattern with `ticket creation failed` operator
- [test] pattern form converters > should be converted to form and back to pattern with `begins with` operator
- [test] pattern form converters > should be converted to form and back to pattern with `not begin with` operator
- [test] pattern form converters > should be converted to form and back to pattern with `ends with` operator
- [test] pattern form converters > should be converted to form and back to pattern with `not end with` operator
- [test] pattern form converters > should be converted to form and back to pattern with alarm infos field and name `exist` operator
- [test] pattern form converters > should be converted to form and back to pattern with alarm infos field and name `not exist` operator
- [test] pattern form converters > should be converted to form and back to pattern with entity infos field and name `exist` operator
- [test] pattern form converters > should be converted to form and back to pattern with entity infos field and name `not exist` operator
- [test] pattern form converters > should be converted to form and back to pattern with `has every` operator
- [test] pattern form converters > should be converted to form and back to pattern with `has one of` operator
- [test] pattern form converters > should keep correct value after converting `%s` field with `%s` condition to form and back
- [test] pattern form converters > <dynamic test name>
- [test] pattern form converters > should keep correct value after converting `%s` condition for an infos value field to form and back
- [test] pattern form converters > <dynamic test name>
- [test] pattern form converters > should keep correct value after converting `%s` condition for an extra infos value field to form and back
- [test] pattern form converters > <dynamic test name>
- [test] pattern form converters > should keep correct value after converting `%s` condition for an object field to form and back
- [test] pattern form converters > <dynamic test name>
- [test] pattern form converters > should include one-of operators for string pattern field `%s`
- [test] pattern form converters > <dynamic test name>
- [test] pattern form converters > should include one-of operators for %s value fields
- [test] pattern form converters > <dynamic test name>
- [test] pattern form converters > should not include one-of operators for %s
- [test] pattern form converters > <dynamic test name>
- [test] pattern form converters > should be converted to form and back to pattern with `has not` operator
- [test] pattern form converters > should be converted to form and back to pattern with `is empty` operator
- [test] pattern form converters > should be converted to form and back to pattern with `is not empty` operator
- [test] pattern form converters > should be converted to form and back to pattern with `higher than` operator
- [test] pattern form converters > should be converted to form and back to pattern with `lower than` operator
- [test] pattern form converters > should be converted to form and back to pattern with `longer` operator
- [test] pattern form converters > should be converted to form and back to pattern with `shorter` operator
- [test] pattern form converters > should be converted to form and back to pattern with `relative time` condition with `to` value
- [test] pattern form converters > should be converted to form and back to pattern with `absolute time` condition
- [test] pattern form converters > should be converted to form and back to pattern with infos and number value
- [test] pattern form converters > should be converted to form and back to pattern with infos and string value
- [test] pattern form converters > should be converted to form and back to pattern with infos and boolean value
- [test] pattern form converters > should be converted to form and back to pattern with infos and string array value
- [test] pattern form converters > should be converted to form and back to pattern with activated
- [test] pattern form converters > should be converted to form and back to pattern with inactive
- [test] pattern form converters > should be converted to form and back to pattern with is meta alarm
- [test] pattern form converters > should be converted to form and back to pattern with inactive
- [test] pattern form converters > should be converted to form and back to pattern with `relative time` condition for activation date
- [test] pattern form converters > should be converted to form and back to pattern with `absolute time` condition for activation date
- [test] pattern form converters > should be converted to form and back to pattern with `relative time` condition and custom duration for `within` operator
- [test] pattern form converters > should be converted to form and back to pattern with `relative time` condition and custom duration for `olderThan` operator
- [test] pattern form converters > should be converted to form and back to pattern with `relative time` condition and custom durations for `inRangePeriod` operator with custom `from` range
- [test] pattern form converters > should be converted to form and back to pattern with `relative time` condition and custom durations for `inRangePeriod` operator with custom `to` range
- [test] pattern form converters > should be converted to form and back to pattern with `relative time` condition and both custom durations for `inRangePeriod` operator

### helpers/handlebars.spec.js

#### helpers/handlebars.spec.js
Status: passed 7, failed 0, skipped 0

- [test] getTemplateVariables > Return empty variables without template
- [test] getTemplateVariables > Return variable with simple template
- [test] getTemplateVariables > Return variable with list template
- [test] getTemplateVariables > Return variable with compare helper
- [test] getTemplateVariables > Return variable with request helper
- [test] getTemplateVariables > Return variable with internal-link helper
- [test] getTemplateVariables > Return variable with copy helper

### helpers/handlebars/helpers.spec.js

#### helpers/handlebars/helpers.spec.js
Status: passed 8, failed 0, skipped 0

- [test] requestHelper > Throws error when url is missing
- [test] requestHelper > Returns empty string when options.fn is not a function
- [test] requestHelper > Creates helper data in window._handlebarsRequestHelper
- [test] requestHelper > Initializes window._handlebarsRequestHelper if it does not exist
- [test] requestHelper > Preserves existing window._handlebarsRequestHelper data
- [test] requestHelper > Returns SafeString with correct component markup
- [test] requestHelper > Generates unique id for each helper call
- [test] requestHelper > Stores all hash options in helper data

### helpers/handlebars/normalize.spec.js

#### helpers/handlebars/normalize.spec.js
Status: passed 17, failed 0, skipped 0

- [test] normalizeHandlebarsNbsp > Returns empty string when no argument provided
- [test] normalizeHandlebarsNbsp > Returns empty string when empty string provided
- [test] normalizeHandlebarsNbsp > Replaces &nbsp; in simple handlebars block
- [test] normalizeHandlebarsNbsp > Replaces \u00A0 in simple handlebars block
- [test] normalizeHandlebarsNbsp > Replaces &nbsp; in triple handlebars block
- [test] normalizeHandlebarsNbsp > Replaces \u00A0 in triple handlebars block
- [test] normalizeHandlebarsNbsp > Replaces multiple &nbsp; in handlebars block
- [test] normalizeHandlebarsNbsp > Replaces multiple \u00A0 in handlebars block
- [test] normalizeHandlebarsNbsp > Replaces mixed &nbsp; and \u00A0 in handlebars block
- [test] normalizeHandlebarsNbsp > Replaces NBSP in multiple handlebars blocks
- [test] normalizeHandlebarsNbsp > Does not replace &nbsp; outside handlebars blocks
- [test] normalizeHandlebarsNbsp > Does not replace \u00A0 outside handlebars blocks
- [test] normalizeHandlebarsNbsp > Replaces NBSP only inside handlebars blocks
- [test] normalizeHandlebarsNbsp > Handles complex template with multiple block types
- [test] normalizeHandlebarsNbsp > Handles template with nested handlebars blocks
- [test] normalizeHandlebarsNbsp > Handles template without any handlebars blocks
- [test] normalizeHandlebarsNbsp > Handles template with only regular spaces in handlebars blocks

### helpers/json.spec.js

#### helpers/json.spec.js
Status: passed 6, failed 0, skipped 0

- [test] stringifyJson > Valid json with default arguments
- [test] stringifyJson > Valid json string with default arguments
- [test] stringifyJson > Undefined json with default arguments
- [test] stringifyJson > Valid json with special indents
- [test] stringifyJson > Undefined json with special defaultValue
- [test] stringifyJson > Invalid json with default arguments

### helpers/metrics.spec.js

#### helpers/metrics.spec.js
Status: passed 13, failed 0, skipped 0

- [test] getDateLabelBySampling > Valid string with sampling by hour
- [test] getDateLabelBySampling > Valid string with sampling by day
- [test] getDateLabelBySampling > Valid string with sampling by week
- [test] getDateLabelBySampling > Valid string with sampling by month
- [test] convertDataSizeValueToTickString > should return 1B
- [test] convertDataSizeValueToTickString > should return 1KB
- [test] convertDataSizeValueToTickString > should return 512KB
- [test] convertDataSizeValueToTickString > should return 1MB
- [test] convertDataSizeValueToTickString > should return 512KB
- [test] convertDataSizeValueToTickString > should return 3GB
- [test] convertDataSizeValueToTickString > should return 512MB
- [test] convertDataSizeValueToTickString > should return 4TB
- [test] convertDataSizeValueToTickString > should return 512GB

### helpers/vue-base.spec.js

#### helpers/vue-base.spec.js
Status: passed 3, failed 0, skipped 0

- [test] mergeReactiveChangedProperties > shouldn't change link to empty object
- [test] mergeReactiveChangedProperties > shouldn't change link to object
- [test] mergeReactiveChangedProperties > should return null when receive null as new object

### i18n/messages.spec.js

#### i18n/messages.spec.js
Status: passed 2, failed 0, skipped 0

- [test] I18n messages > EN translation keys DOESN'T exists inside FR
- [test] I18n messages > FR translation keys DOESN'T exists inside EN

### mixins/entities/corporate-pattern.spec.js

#### mixins/entities/corporate-pattern.spec.js
Status: passed 5, failed 0, skipped 0

- [test] Entities corporate pattern mixin > Initial values applied
- [test] Entities corporate pattern mixin > Corporate patterns list fetched
- [test] Entities corporate pattern mixin > Corporate patterns list fetched with params
- [test] Entities corporate pattern mixin > Corporate patterns list fetched with error
- [test] Entities corporate pattern mixin > Corporate patterns list fetched with previous params

### mixins/entities/event-filter.spec.js

#### mixins/entities/event-filter.spec.js
Status: passed 8, failed 0, skipped 0

- [test] Entities event filter mixin > Initial values applied
- [test] Entities event filter mixin > Event filter rules list fetched
- [test] Entities event filter mixin > Event filter rules list fetched with params
- [test] Entities event filter mixin > Event filter rules list fetched with error
- [test] Entities event filter mixin > Event filter rules list fetched with previous params
- [test] Entities event filter mixin > Event filter rule created
- [test] Entities event filter mixin > Event filter rule updated
- [test] Entities event filter mixin > Event filter rule removed

### mixins/entities/map.spec.js

#### mixins/entities/map.spec.js
Status: passed 11, failed 0, skipped 0

- [test] Entities map mixin > Initial values applied
- [test] Entities map mixin > Event maps list fetched
- [test] Entities map mixin > Event maps list fetched with params
- [test] Entities map mixin > Event maps list fetched without store
- [test] Entities map mixin > Maps list fetched with error
- [test] Entities map mixin > Map fetched without store
- [test] Entities map mixin > Map state fetched without store
- [test] Entities map mixin > Map created
- [test] Entities map mixin > Map updated
- [test] Entities map mixin > Map removed
- [test] Entities map mixin > Map removed

### mixins/entities/pattern.spec.js

#### mixins/entities/pattern.spec.js
Status: passed 5, failed 0, skipped 0

- [test] Entities pattern mixin > Initial values applied
- [test] Entities pattern mixin > Patterns list fetched
- [test] Entities pattern mixin > Patterns list fetched with params
- [test] Entities pattern mixin > Patterns list fetched with error
- [test] Entities pattern mixin > Patterns list fetched with previous params

### mixins/widget/columns

#### mixins/widget/columns/alarm.spec.js
Status: passed 28, failed 0, skipped 0

- [test] widget-columns-alarm > data > Should return initial data with empty columnsFilters and false columnsFiltersPending
- [test] widget-columns-alarm > computed > infoPopupsMap > Should create a map of info popups by column
- [test] widget-columns-alarm > computed > infoPopupsMap > Should return empty object when infoPopups is undefined
- [test] widget-columns-alarm > computed > infoPopupsMap > Should return empty object when widget parameters is undefined
- [test] widget-columns-alarm > computed > columnsFiltersMap > Should create a map of columns filters with filter functions
- [test] widget-columns-alarm > computed > columnsFiltersMap > Should handle columns filters without attributes
- [test] widget-columns-alarm > computed > columnsFiltersMap > Should return empty object when columnsFilters is undefined
- [test] widget-columns-alarm > computed > columnsFiltersMap > Should return identity function when filter does not exist
- [test] widget-columns-alarm > computed > preparedColumns > Should prepare columns with all properties
- [test] widget-columns-alarm > computed > preparedColumns > Should use alias when column value matches entity info property
- [test] widget-columns-alarm > computed > preparedColumns > Should use label when provided instead of alias
- [test] widget-columns-alarm > computed > preparedColumns > Should use original text when no label or alias
- [test] widget-columns-alarm > computed > preparedColumns > Should enable color indicator when colorIndicator is valid
- [test] widget-columns-alarm > computed > preparedColumns > Should disable color indicator when colorIndicator is invalid
- [test] widget-columns-alarm > computed > preparedColumns > Should handle alarm.* prefix for popup template
- [test] widget-columns-alarm > computed > preparedColumns > Should use truncate behavior when specified
- [test] widget-columns-alarm > computed > preparedColumns > Should handle empty columns array
- [test] widget-columns-alarm > computed > preparedColumns > Should handle undefined columns
- [test] widget-columns-alarm > computed > preparedColumns > Should return undefined filter when locale is not set
- [test] widget-columns-alarm > mounted > Should fetch column filters and entity info properties on mount
- [test] widget-columns-alarm > methods > getColumnFilter > Should return filter from columnsFiltersMap when it exists
- [test] widget-columns-alarm > methods > getColumnFilter > Should return default filter when column not in columnsFiltersMap
- [test] widget-columns-alarm > methods > getFilter > Should return filter function with attributes
- [test] widget-columns-alarm > methods > getFilter > Should return filter function without attributes
- [test] widget-columns-alarm > methods > getFilter > Should return identity function when filter does not exist
- [test] widget-columns-alarm > methods > fetchColumnFilters > Should fetch column filters and update state
- [test] widget-columns-alarm > methods > fetchColumnFilters > Should set pending state during fetch
- [test] widget-columns-alarm > methods > fetchColumnFilters > Should handle fetch error gracefully

#### mixins/widget/columns/entity-info-property.spec.js
Status: passed 14, failed 0, skipped 0

- [test] widget-columns-entity-info-property > mounted > Should fetch entity info properties list without pagination on mount
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should return alias when column value matches entity info property with alias
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should return alias when column value matches entity info property without prefix
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should return undefined when column value matches entity info property without alias
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should return undefined when column value matches entity info property with empty alias
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should return undefined when column value does not match any entity info property
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should return undefined when column value format is incorrect
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should return undefined when entityInfoProperties is empty
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should handle column value with different prefix
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should handle column value with null prefix
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should handle column value with undefined prefix
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should handle column value with empty string prefix
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should correctly filter Boolean values in join
- [test] widget-columns-entity-info-property > methods > findAliasByColumnValue > Should not match partial property names

### mixins/widget/periodic-refresh.spec.js

#### mixins/widget/periodic-refresh.spec.js
Status: passed 17, failed 0, skipped 0

- [test] widget-periodic-refresh > computed > Should return periodicRefreshEnabled correctly
- [test] widget-periodic-refresh > computed > Should return periodicRefreshSeconds correctly
- [test] widget-periodic-refresh > methods > Should start periodic refresh correctly
- [test] widget-periodic-refresh > methods > Should stop existing interval before starting new one
- [test] widget-periodic-refresh > methods > Should call fetchList when interval triggers
- [test] widget-periodic-refresh > methods > Should pause periodic refresh correctly
- [test] widget-periodic-refresh > methods > Should resume periodic refresh correctly
- [test] widget-periodic-refresh > methods > Should resume periodic refresh with zero delay when no delay set
- [test] widget-periodic-refresh > methods > Should stop periodic refresh correctly
- [test] widget-periodic-refresh > watchers > Should start periodic refresh when enabled and has value
- [test] widget-periodic-refresh > watchers > Should not start periodic refresh when paused
- [test] widget-periodic-refresh > watchers > Should stop periodic refresh when disabled
- [test] widget-periodic-refresh > watchers > Should restart periodic refresh when value changes
- [test] widget-periodic-refresh > watchers > Should not restart periodic refresh when periodicRefreshSeconds is zero
- [test] widget-periodic-refresh > watchers > Should pause periodic refresh when activeViewPeriodicRefreshPaused becomes true
- [test] widget-periodic-refresh > watchers > Should resume periodic refresh when activeViewPeriodicRefreshPaused becomes false
- [test] widget-periodic-refresh > beforeUnmount > Should stop periodic refresh on component destroy

### plugins/sidebar/components

#### plugins/sidebar/components/sidebar-wrapper.spec.js
Status: passed 0, failed 3, skipped 0

- [test] sidebar-wrapper > Sidebar hidden trigger drawer
- [test] sidebar-wrapper > Sidebar hidden happened after click on the close button
- [test] sidebar-wrapper > Sidebar hidden after click on the close button with close condition

#### plugins/sidebar/components/the-sidebar.spec.js
Status: passed 0, failed 1, skipped 0

- [test] the-sidebar > Renders `the-sidebars` with type: alarmSettings

### plugins/sidebar/store

#### plugins/sidebar/store/index.spec.js
Status: passed 0, failed 9, skipped 0

- [test] Sidebar plugin store module > Mutate state after commit SHOW
- [test] Sidebar plugin store module > Mutate state after commit SHOW without config
- [test] Sidebar plugin store module > Mutate state after commit HIDE
- [test] Sidebar plugin store module > Mutate state after commit HIDE_COMPLETED
- [test] Sidebar plugin store module > Show sidebar. Action: show
- [test] Sidebar plugin store module > Show sidebar with existing id still commits SHOW. Action: show
- [test] Sidebar plugin store module > Hide sidebar. Action: hide
- [test] Sidebar plugin store module > Hide sidebar skips when id unknown. Action: hide
- [test] Sidebar plugin store module > Getter: sidebars

### plugins/vuetify/components

#### plugins/vuetify/components/v-data-table/v-data-table.spec.js
Status: passed 0, failed 10, skipped 0

- [test] v-data-table > Column sorted by DESC after click on the header
- [test] v-data-table > Other column sorted by ASC after click on the header
- [test] v-data-table > Column sort reset after click on the header
- [test] v-data-table > Column sorted by DESC after click on the header with must sort
- [test] v-data-table > First column sorted by DESC after click on the header with multi sort
- [test] v-data-table > First column sorted by ASC after click on the header with multi sort
- [test] v-data-table > Second column sorted by DESC after click on the header with multi sort
- [test] v-data-table > Second column sorted by ASC after click on the header with multi sort
- [test] v-data-table > Second column sort reset after click on the header with multi sort
- [test] v-data-table > Column sort reset after click on the header with multi sort

### store/modules/active-view

#### store/modules/active-view/active-widgets.spec.js
Status: passed 12, failed 0, skipped 0

- [test] Active widgets module > Mutate state after commit FETCH_LIST
- [test] Active widgets module > Mutate state after commit FETCH_LIST_COMPLETED
- [test] Active widgets module > Mutate state after commit FETCH_LIST_FAILED
- [test] Active widgets module > Mutate state after commit CLEAR
- [test] Active widgets module > Get allIds data by widget id. Getter: GET_ALL_IDS_BY_WIDGET_ID
- [test] Active widgets module > Get meta data by widget id. Getter: GET_META_BY_WIDGET_ID
- [test] Active widgets module > Get pending data by widget id. Getter: GET_PENDING_BY_WIDGET_ID
- [test] Active widgets module > Get params data by widget id. Getter: GET_FETCHING_PARAMS_BY_WIDGET_ID
- [test] Active widgets module > Get allIds data by widget id without data. Getter: GET_ALL_IDS_BY_WIDGET_ID
- [test] Active widgets module > Get meta data by widget id without data. Getter: GET_META_BY_WIDGET_ID
- [test] Active widgets module > Get pending data by widget id without data. Getter: GET_PENDING_BY_WIDGET_ID
- [test] Active widgets module > Get fetching params data by widget id without data. Getter: GET_FETCHING_PARAMS_BY_WIDGET_ID

### store/modules/entities

#### store/modules/entities/alarm.spec.js
Status: passed 7, failed 0, skipped 0

- [test] Entities alarm module > Get alarms data by widget id. Getter: getListByWidgetId
- [test] Entities alarm module > Get alarms data by ids. Getter: getList
- [test] Entities alarm module > Get alarm data by alarm id. Getter: getItem
- [test] Entities alarm module > Get meta data by widget id. Getter: getMetaByWidgetId
- [test] Entities alarm module > Get pending data by widget id. Getter: getPendingByWidgetId
- [test] Entities alarm module > Fetch alarms without saving in store. Action: fetchListWithoutStore
- [test] Entities alarm module > Fetch alarm by id. Action: fetchItem

#### store/modules/entities/llm.spec.js
Status: passed 5, failed 0, skipped 0

- [test] Entities llm module > Fetch LLM history without store. Action: fetchLlmHistoryWithoutStore
- [test] Entities llm module > Fetch LLM chats without store. Action: fetchLlmChatsWithoutStore
- [test] Entities llm module > Fetch LLM users without store. Action: fetchLlmUsersWithoutStore
- [test] Entities llm module > Fetch LLM messages without store. Action: fetchLlmMessagesWithoutStore
- [test] Entities llm module > Link LLM history in bulk without store. Action: bulkLinkLlmHistory

#### store/modules/entities/pattern.spec.js
Status: passed 16, failed 0, skipped 0

- [test] Entities pattern module > Fetch list without store. Action: fetchListWithoutStore
- [test] Entities pattern module > Bulk remove. Action: bulkRemove

### store/modules/query

#### store/modules/query/query.spec.js
Status: passed 15, failed 0, skipped 0

- [test] Query module > Mutate state after commit UPDATE
- [test] Query module > Mutate state after commit MERGE
- [test] Query module > Mutate state after commit REMOVE
- [test] Query module > Mutate state after commit FORCE_UPDATE
- [test] Query module > Mutate state after commit UPDATE_LOCKED
- [test] Query module > Mutate state after commit REMOVE_LOCKED
- [test] Query module > Update query. Action: update
- [test] Query module > Merge query. Action: merge
- [test] Query module > Remove query. Action: remove
- [test] Query module > Force update query. Action: forceUpdate
- [test] Query module > Update locked query. Action: updateLocked
- [test] Query module > Remove locked query. Action: removeLocked
- [test] Query module > Get query by id. Getter: getQueryById
- [test] Query module > Get query nonce by id. Getter: getQueryNonceById
- [test] Query module > Get query nonce by id without nonce. Getter: getQueryNonceById

### store/plugins/entities

#### store/plugins/entities/getters.spec.js
Status: passed 1, failed 0, skipped 0

- [test] Entities plugin > Test

### views/admin/maps.spec.js

#### views/admin/maps.spec.js
Status: passed 0, failed 8, skipped 0

- [test] maps > Maps fetched after mounted
- [test] maps > Maps re-fetched after trigger refresh button
- [test] maps > Choose map type modal showed after trigger create button
- [test] maps > Edit %s map modal showed after trigger edit button
- [test] maps > <dynamic test name>
- [test] maps > Duplicate %s map modal showed after trigger edit button
- [test] maps > <dynamic test name>
- [test] maps > Confirmation modal showed after trigger remove map button

### views/admin/tags.spec.js

#### views/admin/tags.spec.js
Status: passed 0, failed 7, skipped 0

- [test] tags > Tags fetched after mounted
- [test] tags > Tags re-fetched after trigger refresh button
- [test] tags > Create tag modal showed after trigger create button
- [test] tags > Update tag modal showed after trigger edit button
- [test] tags > Duplicate tag modal showed after trigger duplicate button
- [test] tags > Confirmation modal showed after trigger remove tag button
- [test] tags > Confirmation modal showed after trigger remove selected tags button

### views/profile/patterns.spec.js

#### views/profile/patterns.spec.js
Status: passed 0, failed 9, skipped 0

- [test] patterns > Patterns refreshed after trigger expand button
- [test] patterns > Corporate patterns refreshed after trigger expand button
- [test] patterns > Edit pattern modal showed after trigger edit on patterns
- [test] patterns > Edit corporate pattern modal showed after trigger edit on patterns
- [test] patterns > Create alarm pattern modal showed after trigger edit on patterns
- [test] patterns > Create entity pattern modal showed after trigger edit on patterns
- [test] patterns > Create pbehavior pattern modal showed after trigger edit on patterns
- [test] patterns > Confirmation delete pattern modal showed after trigger delete on patterns
- [test] patterns > Confirmation delete selected patterns modal showed after trigger delete on patterns
