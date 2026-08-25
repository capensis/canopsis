<template>
  <div
    class="alarms-list-table"
    v-on="wrapperListeners"
  >
    <c-empty-data-table-columns v-if="!preparedColumns.length" />
    <template v-else>
    <div
      v-if="shownTopPagination"
      ref="actions"
      class="alarms-list-table__top-pagination"
    >
      <c-layout
        class="px-4 py-2 position-relative"
        align="center"
      >
        <div class="alarms-list-table__top-pagination-left">
          <c-density-btn-toggle
            v-if="densable && isTableControlsShown"
            :value="dense"
            button-size="x-small"
            density="compact"
            icon-size="16"
            @change="$emit('update:dense', $event)"
          />

          <c-mass-actions-panel
            v-if="selectable && unresolvedSelected.length"
            v-model="keepSelectedAfterAction"
            v-bind="widgetGridParameters"
            :selected="unresolvedSelected"
            :default-value="defaultKeepSelectedAfterAction"
            @clear:selected="clearSelected"
          >
            <template #actions>
              <mass-actions-panel
                :items="unresolvedSelected"
                :widget="preparedWidget"
                :refresh-alarms-list="refreshAlarmsList"
                :small="isSmallDensity"
                @clear:items="clearSelectedAfterAction"
              />
            </template>
          </c-mass-actions-panel>
        </div>

        <c-pagination
          :page="options.page"
          :limit="options.itemsPerPage"
          :total="totalItems"
          class="alarms-list-table__top-pagination-center"
          type="top"
          @update:page="updatePage"
        />

        <c-grid-btns
          v-if="hasColumnEditButtons"
          :editing="isColumnsChanging"
          class="alarms-list-table__top-pagination-right"
          @toggle="toggleColumnEditingMode"
          @reset="resetColumnsSettings"
        />
      </c-layout>
    </div>

    <v-data-table-server
      ref="dataTable"
      :headers="headersWithWidth"
      :hide-default-footer="true"
      :items="alarms"
      :items-length="totalItems"
      :items-per-page="options.itemsPerPage"
      :loading="isTableLoading"
      :page="options.page"
      :sort-by="options.sortBy"
      :disable-sort="isColumnsChanging"
      :density="tableDensity"
      :class="tableClasses"
      :style="tableStyle"
      class="alarms-list-table__table"
      multi-sort
      @update:options="updateOptions"
    >
      <template
        v-for="header in headersWithWidth"
        :key="`header.${header.key}`"
        #[`header.${header.key}`]="{ column, isSorted, sortBy, getSortIcon }"
      >
        <div v-if="header.key === rowIconsColumnKey" class="alarms-list-table__select-all">
          <div class="alarms-list-table__select-all-checkbox">
            <c-light-checkbox
              v-if="selectable"
              :model-value="isAllSelectableSelected"
              :indeterminate="isSomeSelectableSelected && !isAllSelectableSelected"
              :disabled="!selectableAlarms.length"
              @update:model-value="toggleAllSelectableAlarms"
            />
          </div>
        </div>
        <div v-else class="v-data-table-header__content">
          <alarm-header-cell
            :header="header"
            :resizing="resizingMode"
            :ellipsis-headers="isCellContentTruncated"
          />
          <template v-if="column.sortable && !isColumnsChanging">
            <v-icon
              :icon="getSortIcon(column)"
              class="v-data-table-header__sort-icon"
            />
            <div
              v-if="isSorted(column)"
              class="v-data-table-header__sort-badge"
            >
              {{ getHeaderSortIndex(column, sortBy) }}
            </div>
          </template>
        </div>
        <span
          v-if="draggingMode && header.key !== rowIconsColumnKey"
          :data-value="header.value"
          class="table__dragging-handler"
          @click.stop=""
        />
        <span
          v-if="resizingMode && header.key !== rowIconsColumnKey"
          class="table__resize-handler"
          @pointerdown.stop.prevent="startColumnResize(header.value)"
          @mousedown.stop.prevent=""
          @click.stop=""
        />
      </template>

      <template #body="{ items }">
        <tr v-if="!items.length && !isTableLoading" class="v-data-table-rows-no-data">
          <td :colspan="headersColspan">
            {{ $t('common.noDataAvailable') }}
          </td>
        </tr>
        <template
          v-for="item in items"
          :key="getTableItem(item)._id"
        >
          <alarms-list-row
            :ref="`row${getTableItem(item)._id}`"
            :alarm="getTableItem(item)"
            :widget="preparedWidget"
            :headers="headersWithWidth"
            :row-icons-column-key="rowIconsColumnKey"
            :parent-alarm="parentAlarm"
            :selected-tags="selectedTags"
            :selectable="selectable"
            :selectable-alarm="isAlarmSelectable(getTableItem(item))"
            :selected="isAlarmSelected(getTableItem(item))"
            :selecting="selecting"
            :small="isSmallDensity"
            :resizing="resizingMode"
            :hide-actions="hideActions"
            :expandable="expandable"
            :expanded="isAlarmExpanded(getTableItem(item))"
            :search="query.search"
            :actions-inline-count="actionsInlineCount"
            :actions-ignore-media-query="resizableColumn"
            :refresh-alarms-list="refreshAlarmsList"
            :virtual-scroll="preparedWidget.parameters.isVirtualScrollEnabled"
            :booted="bootedRows[getTableItem(item)._id]"
            :eager="eager"
            @update:selected="toggleAlarmSelection(getTableItem(item), $event)"
            @update:expanded="toggleAlarmExpanded(getTableItem(item), $event)"
            @start:resize="startColumnResize"
            @select:tag="$emit('select:tag', $event)"
            @remove:tag="$emit('remove:tag', $event)"
          />

          <alarms-list-expand-row
            :expanded="isAlarmExpanded(getTableItem(item))"
            :colspan="headersColspan"
          >
            <alarms-expand-panel
              :alarm="getExpandedAlarm(getTableItem(item))"
              :widget="preparedWidget"
              :hide-children="hideChildren"
              :search="query.search"
              :selected-tags="selectedTags"
              @select:tag="$emit('select:tag', $event)"
              @remove:tag="$emit('remove:tag', $event)"
            />
          </alarms-list-expand-row>
        </template>
      </template>
    </v-data-table-server>

    <c-horizontal-scrollbar v-if="stickyHorizontalScroll" ref="horizontalScrollbar" />
    </template>
  </div>
</template>

<script>
import { get, isEqual, isNumber } from 'lodash';
import { mapState } from 'pinia';

import { PAGINATION_LIMIT, VUETIFY_ANIMATION_DELAY } from '@/config';
import {
  ALARM_ACTION_BUTTON_MARGINS,
  ALARM_ACTION_BUTTON_WIDTHS,
  ALARM_ACTIONS_PADDINGS,
  DEFAULT_ALARM_ACTIONS_INLINE_COUNT,
  DENSE_TYPES,
  RESIZING_CELLS_CONTENTS_BEHAVIORS,
} from '@/constants';

import { AsyncBooting } from '@/services/async-booting';
import { TableIntersectionObserver } from '@/services/table-intersection-observer';

import { useAlarmStore } from '@/stores/modules/alarm';

import { mapIds } from '@/helpers/array';
import { isActionAvailableForAlarm } from '@/helpers/entities/alarm/form';
import { calculateAlarmLinksColumnWidth } from '@/helpers/entities/alarm/list';
import {
  convertDataTableSortToQuery,
  convertQuerySortToDataTableSort,
} from '@/helpers/entities/shared/query';
import { prepareAlarmListWidget } from '@/helpers/entities/widget/forms/alarm';
import {
  getNearestAndFarthestIndexes,
  getNearestViewportIndexesBound,
  recursiveRaf,
  splitIdsToChunk,
} from '@/helpers/render';

import { activeViewMixin } from '@/mixins/active-view';
import {
  widgetColumnDraggingAlarmMixin,
  widgetColumnResizingAlarmMixin,
  widgetColumnsAlarmMixin,
} from '@/mixins/widget/columns';
import { widgetRowsSelectingAlarmMixin } from '@/mixins/widget/rows/alarm-selecting';
import { widgetStickyAlarmMixin } from '@/mixins/widget/rows/alarm-sticky';

import MassActionsPanel from '../actions/mass-actions-panel.vue';
import AlarmsExpandPanel from '../expand-panel/alarms-expand-panel.vue';
import AlarmHeaderCell from '../headers-formatting/alarm-header-cell.vue';

import AlarmsListExpandRow from './alarms-list-expand-row.vue';
import AlarmsListRow from './alarms-list-row.vue';

export default {
  provide() {
    return {
      $asyncBootingActionsPanel: this.asyncBootingActionsPanel,
      $intersectionObserver: this.$intersectionObserver,
    };
  },

  components: {
    AlarmHeaderCell,
    AlarmsExpandPanel,
    AlarmsListExpandRow,
    AlarmsListRow,
    MassActionsPanel,
  },

  mixins: [
    activeViewMixin,
    widgetColumnsAlarmMixin,
    widgetColumnDraggingAlarmMixin,
    widgetColumnResizingAlarmMixin,
    widgetRowsSelectingAlarmMixin,
    widgetStickyAlarmMixin,
  ],

  props: {
    widget: {
      type: Object,
      required: true,
    },
    alarms: {
      type: Array,
      required: true,
    },
    columns: {
      type: Array,
      default: () => [],
    },
    columnsSettings: {
      type: Object,
      default: () => ({}),
    },
    meta: {
      type: Object,
      required: true,
    },
    query: {
      type: Object,
      required: true,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
    loading: {
      type: Boolean,
      default: false,
    },
    expandable: {
      type: Boolean,
      default: false,
    },
    selectable: {
      type: Boolean,
      default: false,
    },
    hideChildren: {
      type: Boolean,
      default: false,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    hidePagination: {
      type: Boolean,
      default: false,
    },
    densable: {
      type: Boolean,
      default: false,
    },
    denseValue: {
      type: Number,
      default: undefined,
    },
    resizableColumn: {
      type: Boolean,
      default: false,
    },
    draggableColumn: {
      type: Boolean,
      default: false,
    },
    widgetGridParameters: {
      type: Object,
      default: () => ({}),
    },
    defaultKeepSelectedAfterAction: {
      type: Boolean,
      default: false,
    },
    selectedTags: {
      type: Array,
      default: () => [],
    },
    refreshAlarmsList: {
      type: Function,
      default: () => {},
    },
    eager: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:query', 'update:columns-settings', 'update:dense', 'select:tag', 'remove:tag'],

  data() {
    return {
      keepSelectedAfterAction: this.defaultKeepSelectedAfterAction,
      rowIconsColumnKey: 'row-icons',
      rowIconsColumnClass: 'alarms-list-table__row-icons-column',
      selectedAlarmIds: [],
      expandedAlarmIds: [],
      asyncBootingActionsPanel: new AsyncBooting(),
      bootedRows: {},
    };
  },

  computed: {
    preparedWidget() {
      return prepareAlarmListWidget(this.widget);
    },

    ...mapState(useAlarmStore, {
      alarmDetailsById(store) {
        return store.getDetailsByWidgetId(this.preparedWidget._id);
      },
    }),

    widgetColumnValues() {
      return this.preparedColumns.map(({ value }) => value).join('\0');
    },

    columnsSettingsHeaderValues() {
      return this.baseHeaders.map(({ value }) => value).filter(Boolean);
    },

    hasRowActions() {
      return (
        this.expandable ||
        this.selectable ||
        this.alarms.some(alarm => this.hasAlarmInstruction(alarm) || alarm.bookmark)
      );
    },

    cellsContentBehavior() {
      return this.preparedWidget.parameters?.columns?.cells_content_behavior;
    },

    isCellContentWrapped() {
      return this.cellsContentBehavior === RESIZING_CELLS_CONTENTS_BEHAVIORS.wrap;
    },

    isCellContentTruncated() {
      return this.cellsContentBehavior === RESIZING_CELLS_CONTENTS_BEHAVIORS.truncate;
    },

    hasColumnEditButtons() {
      return this.resizableColumn || this.draggableColumn;
    },

    isTableControlsShown() {
      return !!this.preparedColumns.length;
    },

    shownTopPagination() {
      return !this.hidePagination && !!this.totalItems;
    },

    isTableLoading() {
      return this.loading || this.columnsFiltersPending || this.entityInfoPropertyPending;
    },

    isColumnsChanging() {
      return this.resizingMode || this.draggingMode;
    },

    baseHeaders() {
      const headers = this.preparedColumns.map((column) => {
        const header = {
          ...column,
          key: column.value,
          title: column.text,
          class: this.draggableClass,
          headerProps: {
            'data-value': column.value,
            class: this.draggableClass,
          },
        };

        if (column.linksInRowCount) {
          const linksCounts = this.alarms.map(
            alarm => Object.values(get(alarm, column.value, {})).flat().length,
          );
          const maxLinksCount = Math.max(...linksCounts);
          const actualInlineLinksCount = maxLinksCount > column.inlineLinksCount
            ? column.inlineLinksCount + 1
            : maxLinksCount;

          header.width = calculateAlarmLinksColumnWidth(
            this.dense,
            Math.max(Math.min(actualInlineLinksCount, column.linksInRowCount), 1),
          );
        }

        return header;
      });

      if (this.hasRowActions) {
        headers.unshift({
          key: this.rowIconsColumnKey,
          value: this.rowIconsColumnKey,
          title: '',
          text: '',
          sortable: false,
          width: 100,
          headerProps: {
            'data-value': this.rowIconsColumnKey,
            class: this.rowIconsColumnClass,
          },
        });
      }

      if (!this.hideActions) {
        const actionsText = this.$t('common.actionsLabel');

        headers.push({
          key: 'actions',
          value: 'actions',
          title: actionsText,
          text: actionsText,
          sortable: false,
          class: this.draggableClass,
          headerProps: {
            'data-value': 'actions',
            class: this.draggableClass,
          },
        });
      }

      return headers;
    },

    headers() {
      if (!this.draggableColumn) {
        return this.baseHeaders;
      }

      const baseHeadersPositions = new Map(
        this.baseHeaders.map(({ value }, index) => [value, index]),
      );
      const rowActionsHeaders = this.baseHeaders.filter(
        ({ value }) => value === this.rowIconsColumnKey,
      );
      const draggableHeaders = this.baseHeaders.filter(
        ({ value }) => value !== this.rowIconsColumnKey,
      );

      return [
        ...rowActionsHeaders,
        ...draggableHeaders.sort(
          (a, b) => this.getColumnSortPosition(a, baseHeadersPositions)
            - this.getColumnSortPosition(b, baseHeadersPositions),
        ),
      ];
    },

    headersWithWidth() {
      if (!this.resizableColumn) {
        return this.headers;
      }

      return this.headers.map((header) => {
        const width = this.getColumnWidthByField(header.value);
        const columnWidth = header.width ? header.width : width && `${width}px`;

        return {
          ...header,
          width: columnWidth,
          minWidth: columnWidth,
        };
      });
    },

    headersColspan() {
      return Math.max(this.headersWithWidth.length, 1);
    },

    tableStyle() {
      if (this.resizableColumn) {
        return {
          '--alarms-list-table-width': this.sumOfColumnsWidth
            ? `${this.sumOfColumnsWidth}px`
            : '100%',
        };
      }

      return {};
    },

    tableClasses() {
      return {
        table__grid: this.isColumnsChanging,
        'alarms-list-table__table--fixed': this.resizableColumn || this.draggableColumn,
        'alarms-list-table__table--sticky-header': this.stickyHeader,
        'alarms-list-table__table--sticky-horizontal-scroll': this.stickyHorizontalScroll,
        'alarms-list-table__table--wrapped': this.isCellContentWrapped,
        'alarms-list-table__table--truncated': this.isCellContentTruncated,
        'alarms-list-table__table--dragging': this.draggingMode,
        'alarms-list-table__table--selecting': this.selecting,
        'alarms-list-table__table--selecting-text-unselectable': this.selectingMousePressed,
      };
    },

    selectableAlarms() {
      return this.alarms.filter(alarm => this.isAlarmSelectable(alarm));
    },

    selectedAlarms() {
      const selectedIds = new Set(this.selectedAlarmIds);

      return this.alarms.filter(alarm => selectedIds.has(alarm._id));
    },

    unresolvedSelected() {
      return this.selectedAlarms.filter(alarm => isActionAvailableForAlarm(alarm, this.preparedWidget));
    },

    isAllSelectableSelected() {
      return (
        !!this.selectableAlarms.length &&
        this.selectableAlarms.every(alarm => this.selectedAlarmIds.includes(alarm._id))
      );
    },

    isSomeSelectableSelected() {
      return this.selectableAlarms.some(alarm => this.selectedAlarmIds.includes(alarm._id));
    },

    totalItems() {
      return this.meta.total_count ?? 0;
    },

    dense() {
      return this.denseValue ?? this.preparedWidget.parameters.dense;
    },

    tableDensity() {
      return {
        [DENSE_TYPES.medium]: 'comfortable',
        [DENSE_TYPES.small]: 'compact',
      }[this.dense] ?? 'default';
    },

    isSmallDensity() {
      return this.dense === DENSE_TYPES.small;
    },

    isMediumDensity() {
      return this.dense === DENSE_TYPES.medium;
    },

    actionsInlineCount() {
      if (!this.resizableColumn) {
        return DEFAULT_ALARM_ACTIONS_INLINE_COUNT;
      }

      const denseType = {
        [this.isSmallDensity]: DENSE_TYPES.small,
        [this.isMediumDensity]: DENSE_TYPES.medium,
      }.true ?? DENSE_TYPES.large;
      const width = this.getColumnWidthByField('actions') ?? 0;
      const widthWithoutPaddings = width - (ALARM_ACTIONS_PADDINGS[denseType] * 2);
      const actionWidth = ALARM_ACTION_BUTTON_WIDTHS[denseType] + ALARM_ACTION_BUTTON_MARGINS[denseType];

      return Math.floor(widthWithoutPaddings / actionWidth) || 1;
    },

    options: {
      get() {
        const {
          page = 1,
          itemsPerPage = PAGINATION_LIMIT,
          sortBy = [],
          sortDesc = [],
        } = this.query;

        return {
          page,
          itemsPerPage,
          sortBy: convertQuerySortToDataTableSort(sortBy, sortDesc),
        };
      },

      set(newOptions = {}) {
        const {
          page = 1,
          itemsPerPage = PAGINATION_LIMIT,
          sortBy = [],
        } = newOptions;
        const sort = convertDataTableSortToQuery(sortBy);
        const oldItemsPerPage = this.query.itemsPerPage ?? PAGINATION_LIMIT;
        const convertedOptions = {
          page: itemsPerPage <= oldItemsPerPage ? page : 1,
          itemsPerPage,
          sortBy: sort.sortBy,
          sortDesc: sort.sortDesc,
        };
        const oldOptions = {
          page: this.query.page ?? 1,
          itemsPerPage: oldItemsPerPage,
          sortBy: this.query.sortBy ?? [],
          sortDesc: this.query.sortDesc ?? [],
        };

        if (isEqual(convertedOptions, oldOptions)) {
          return;
        }

        this.$emit('update:query', {
          ...this.query,
          ...convertedOptions,
        });
      },
    },

  },

  watch: {
    alarms() {
      this.synchronizeSelectedAlarms();
      this.synchronizeExpandedAlarms();
      this.scheduleStickyScrollHandler();

      if (!this.eager) {
        this.setBootedRows(this.alarms);

        return;
      }

      this.scheduleActionsPanelBooting();
    },

    defaultKeepSelectedAfterAction(value) {
      this.keepSelectedAfterAction = value;
    },

    widgetColumnValues(newValues, oldValues) {
      if (!oldValues || newValues === oldValues || !this.isColumnsChanging) {
        return;
      }

      this.$nextTick(() => {
        if (this.draggingMode) {
          this.calculateColumnsSortingOrderByHeaders(this.baseHeaders);
        }

        if (this.resizingMode) {
          this.calculateColumnsWidths();
        }

        this.updateColumnsSettings();
        this.disableDraggingMode();
        this.disableResizingMode();
      });
    },

    columnsSettings: {
      immediate: true,
      deep: true,
      handler() {
        if (!this.draggingMode) {
          this.setColumnsPosition(
            this.normalizeColumnsPosition(this.columnsSettings?.columns_position),
          );
        }

        if (!this.resizingMode) {
          const columnsWidth = this.normalizeColumnsWidth(this.columnsSettings?.columns_width);

          this.setColumnsWidth(columnsWidth);

          if (!this.hasWidthForAllColumns(columnsWidth)) {
            this.$nextTick(() => this.calculateColumnsWidths());
          }
        }
      },
    },

    resizableColumn: {
      immediate: true,
      handler(newResizableColumn) {
        if (newResizableColumn && this.isEmptyColumnsWidthByField) {
          this.$nextTick(() => this.calculateColumnsWidths());
        }
      },
    },

    dense() {
      if (!this.resizableColumn) {
        this.$nextTick(() => this.calculateColumnsWidths());
      }

      this.scheduleStickyScrollHandler();
    },

    expandedAlarmIds(expandedAlarmIds) {
      if (this.preparedWidget.parameters.pausePeriodicRefreshOnExpandPanel) {
        this.checkExpandedForPausePeriodicRefresh(expandedAlarmIds);
      }
    },

    'widget.parameters.pausePeriodicRefreshOnExpandPanel': {
      handler(pausePeriodicRefreshOnExpandPanel) {
        this.checkExpandedForPausePeriodicRefresh(
          pausePeriodicRefreshOnExpandPanel ? this.expandedAlarmIds : [],
        );
      },
    },

    'widget.parameters.isVirtualScrollEnabled': {
      immediate: true,
      handler(isVirtualScrollEnabled) {
        if (isVirtualScrollEnabled) {
          this.$intersectionObserver.connect();

          return;
        }

        this.$intersectionObserver.disconnect();
      },
    },
  },

  beforeCreate() {
    this.$intersectionObserver = new TableIntersectionObserver();
  },

  mounted() {
    this.$tbodyEl = this.$el.querySelector('tbody');

    if (this.resizableColumn && this.isEmptyColumnsWidthByField) {
      this.$nextTick(() => this.calculateColumnsWidths());
    }

    if (this.eager) {
      this.scheduleActionsPanelBooting();
    } else {
      this.setBootedRows(this.alarms);
    }

    this.scheduleStickyScrollHandler();
  },

  beforeUnmount() {
    this.asyncBootingActionsPanel.clear();
    this.$intersectionObserver.destroy();
  },

  methods: {
    scheduleActionsPanelBooting() {
      this.$nextTick(() => {
        window.requestAnimationFrame(() => this.asyncBootingActionsPanel.run());
      });
    },

    setBootedRows(items, itemsPerRender = 10) {
      if (!this.$tbodyEl || !items.length) {
        return;
      }

      const { start, end } = getNearestViewportIndexesBound({
        wrapperEl: this.$tbodyEl,
        length: items.length,
      });
      const { nearest, farthest } = getNearestAndFarthestIndexes({
        start,
        end,
        ids: mapIds(items),
      });

      this.bootedRows = nearest.reduce((acc, id) => {
        acc[id] = true;

        return acc;
      }, {});

      const chunks = splitIdsToChunk(farthest, itemsPerRender);

      if (!chunks.length) {
        window.requestAnimationFrame(() => {
          if (this.stickyHorizontalScroll) {
            this.setHorizontalScrollbarWidth();
            this.calculateScrollPosition();
            this.setHorizontalScrollPosition();
          }

          this.asyncBootingActionsPanel.run();
        });

        return;
      }

      chunks.forEach((idsChunk, index) => {
        recursiveRaf(() => {
          idsChunk.forEach((id) => {
            this.bootedRows[id] = true;
          });

          if (!index && this.stickyHorizontalScroll) {
            window.requestAnimationFrame(() => {
              this.setHorizontalScrollbarWidth();
              this.calculateScrollPosition();
              this.setHorizontalScrollPosition();
            });
          }

          if (index === chunks.length - 1) {
            window.requestAnimationFrame(() => this.asyncBootingActionsPanel.run());
          }
        }, index + 1);
      });
    },

    updateOptions(options) {
      if (this.isColumnsChanging) {
        return;
      }

      this.options = options;
    },

    updatePage(page) {
      this.$emit('update:query', {
        ...this.query,
        page,
      });
    },

    updateColumnsSettings() {
      const settings = {};

      if (this.resizingMode) {
        settings.columns_width = this.normalizeColumnsWidth(this.columnsWidthByField);
      }

      if (this.draggingMode) {
        settings.columns_position = this.normalizeColumnsPosition(this.columnsPositionByField);
      }

      this.$emit('update:columns-settings', settings);
    },

    toggleColumnEditingMode() {
      if (this.isColumnsChanging) {
        this.updateColumnsSettings();
      }

      if (this.resizableColumn) {
        this.toggleResizingMode();
      }

      if (this.draggableColumn) {
        this.toggleDraggingMode();
      }
    },

    resetColumnsSettings() {
      if (this.draggableColumn) {
        this.setColumnsPosition({});
      }

      if (this.resizableColumn) {
        this.setColumnsWidth({});
        this.$nextTick(this.calculateColumnsWidths);
      }
    },

    getColumnSortPosition(header, baseHeadersPositions) {
      const position = this.getColumnPositionByField(header.value);

      return Number.isFinite(position) ? position : baseHeadersPositions.get(header.value);
    },

    normalizeColumnsPosition(columnsPosition = {}) {
      const columnsPositionFields = new Set(Object.keys(columnsPosition));
      const orderedFields = this.columnsSettingsHeaderValues
        .filter(value => columnsPositionFields.has(value))
        .sort((a, b) => columnsPosition[a] - columnsPosition[b]);
      const missingFields = this.columnsSettingsHeaderValues.filter(
        value => !columnsPositionFields.has(value),
      );

      return [...orderedFields, ...missingFields].reduce((acc, value, index) => {
        acc[value] = index;

        return acc;
      }, {});
    },

    normalizeColumnsWidth(columnsWidth = {}) {
      const headerValues = new Set(this.columnsSettingsHeaderValues);

      return Object.entries(columnsWidth).reduce((acc, [value, width]) => {
        if (headerValues.has(value)) {
          acc[value] = width;
        }

        return acc;
      }, {});
    },

    hasWidthForAllColumns(columnsWidth) {
      return this.columnsSettingsHeaderValues.every(value => columnsWidth[value]);
    },

    getHeaderSortIndex(column, sortBy = []) {
      return sortBy.findIndex(sortItem => sortItem.key === column.key) + 1;
    },

    getTableItem(item) {
      return item?.raw ?? item;
    },

    hasAlarmInstruction(alarm) {
      return !!alarm.assigned_instructions?.length || isNumber(alarm.instruction_execution_icon);
    },

    isAlarmSelectable(alarm) {
      return this.selectable && isActionAvailableForAlarm(alarm, this.preparedWidget);
    },

    isAlarmSelected(alarm) {
      return this.selectedAlarmIds.includes(alarm._id);
    },

    isAlarmExpanded(alarm) {
      return this.expandedAlarmIds.includes(alarm._id);
    },

    toggleAlarmSelection(alarm, selected) {
      if (!this.isAlarmSelectable(alarm)) {
        return;
      }

      if (selected) {
        this.selectedAlarmIds = [...new Set([...this.selectedAlarmIds, alarm._id])];
      } else {
        this.selectedAlarmIds = this.selectedAlarmIds.filter(id => id !== alarm._id);
      }
    },

    toggleAllSelectableAlarms(selected) {
      const selectableIds = this.selectableAlarms.map(({ _id }) => _id);

      if (selected) {
        this.selectedAlarmIds = [...new Set([...this.selectedAlarmIds, ...selectableIds])];
      } else {
        this.selectedAlarmIds = this.selectedAlarmIds.filter(id => !selectableIds.includes(id));
      }
    },

    toggleAlarmExpanded(alarm, expanded) {
      if (expanded && !this.expandedAlarmIds.includes(alarm._id)) {
        this.expandedAlarmIds = [...this.expandedAlarmIds, alarm._id];

        return;
      }

      if (!expanded) {
        this.expandedAlarmIds = this.expandedAlarmIds.filter(id => id !== alarm._id);
      }
    },

    getExpandedAlarm(alarm) {
      const alarmDetails = this.alarmDetailsById[alarm._id] ?? {};

      return {
        ...alarm,
        ...(alarmDetails.data ?? alarmDetails),
      };
    },

    clearSelected() {
      this.selectedAlarmIds = [];
    },

    clearSelectedAfterAction() {
      if (!this.keepSelectedAfterAction) {
        this.clearSelected();
      }
    },

    synchronizeSelectedAlarms() {
      const alarmIds = new Set(this.alarms.map(({ _id }) => _id));

      this.selectedAlarmIds = this.selectedAlarmIds.filter(id => alarmIds.has(id));
    },

    synchronizeExpandedAlarms() {
      const alarmIds = new Set(this.alarms.map(({ _id }) => _id));

      this.expandedAlarmIds = this.expandedAlarmIds.filter(id => alarmIds.has(id));
    },

    checkExpandedForPausePeriodicRefresh(expanded) {
      if (expanded.length) {
        this.pausePeriodicRefresh();

        return;
      }

      this.resumePeriodicRefresh();

      setTimeout(() => this.scheduleStickyScrollHandler(), VUETIFY_ANIMATION_DELAY);
    },
  },
};
</script>

<style lang="scss">
.alarms-list-table__top-pagination {
  position: relative;
  min-height: 48px;
  background: rgb(var(--v-theme-surface));
  transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
  transition-property: opacity, background-color;
  z-index: 2;

  &:after {
    content: ' ';
    width: 100%;
    height: 2px;
    background: inherit;
    position: absolute;
    left: 0;
    right: 0;
    bottom: -1px;
  }
}

.alarms-list-table__top-pagination-left {
  align-items: center;
  display: flex;
  gap: 8px;
  padding-right: 80px;
}

.alarms-list-table__top-pagination-center {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

.alarms-list-table__top-pagination-right {
  position: absolute;
  right: 16px;
}

.alarms-list-table__table {
  > .v-table__wrapper {
    overflow-x: auto;
  }

  &.alarms-list-table__table--sticky-horizontal-scroll {
    > .v-table__wrapper {
      -ms-overflow-style: none;
      scrollbar-width: none;

      &::-webkit-scrollbar {
        display: none;
      }
    }
  }

  &.alarms-list-table__table--sticky-header {
    thead {
      position: relative;
      transition: 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
      transition-property: opacity, background-color;
      width: 100%;
      z-index: 2;

      &.head-shadow {
        tr:first-child {
          box-shadow: 0 1px 10px 0 rgba(0, 0, 0, 0.12) !important;
        }
      }

      tr {
        background: rgb(var(--v-theme-surface));
        box-shadow: 0 1px 0 rgba(var(--v-border-color), var(--v-border-opacity));
        transition: background-color 0.3s cubic-bezier(0.25, 0.8, 0.5, 1);
        width: 100%;

        th {
          border-bottom: none;
        }
      }
    }
  }

  tr:not(.v-data-table__expanded) th:first-of-type {
    width: 120px !important;
  }

  .alarms-list-table__row-icons-column {
    padding-left: 6px !important;
    padding-right: 6px !important;
  }

  &.alarms-list-table__table--fixed {
    > .v-table__wrapper > table {
      table-layout: fixed;
      width: var(--alarms-list-table-width);
      max-width: unset;
      min-width: 100%;
    }
  }

  &.alarms-list-table__table--selecting {
    .alarms-list-row:after {
      opacity: 1;
      pointer-events: auto;
    }
  }

  &.alarms-list-table__table--selecting-text-unselectable {
    * {
      user-select: none;
    }
  }
}

.alarms-list-table__select-all {
  display: flex;
  justify-content: flex-start;
  width: 100%;
}

.alarms-list-table__select-all-checkbox {
  display: flex;
  justify-content: center;
  width: 24px;
}
</style>
