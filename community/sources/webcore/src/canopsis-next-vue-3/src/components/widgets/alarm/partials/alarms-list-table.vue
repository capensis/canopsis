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
        class="pa-2 position-relative"
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
            @click:state="openRootCauseDiagram"
            @select:tag="$emit('select:tag', $event)"
            @remove:tag="$emit('remove:tag', $event)"
          />

          <c-table-expand-row
            :expanded="isAlarmExpanded(getTableItem(item))"
            :colspan="headersColspan"
            row-class="alarms-list-table__expand-row"
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
          </c-table-expand-row>
        </template>
      </template>
    </v-data-table-server>

    <c-horizontal-scrollbar
      v-if="stickyHorizontalScroll"
      ref="horizontalScrollbar"
      :class="{ 'alarms-list-table__horizontal-scrollbar--hidden': isHorizontalScrollbarHidden }"
      class="alarms-list-table__horizontal-scrollbar"
    />
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
  MODALS,
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
} from '@/helpers/render';

import { activeViewMixin } from '@/mixins/active-view';
import { handlebarsTagsHelperMixin } from '@/mixins/handlebars/tags-helper';
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

import AlarmsListRow from './alarms-list-row.vue';

const INITIAL_BOOTED_ROWS_LIMIT = 20;

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
    AlarmsListRow,
    MassActionsPanel,
  },

  mixins: [
    activeViewMixin,
    handlebarsTagsHelperMixin,
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
      virtualScrollColumnsWidth: {},
      recalculatingHorizontalScrollbar: false,
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

    isHorizontalScrollbarHidden() {
      return this.isColumnResizing || this.recalculatingHorizontalScrollbar;
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
      const hasVirtualScrollColumnsWidth = Object.keys(this.virtualScrollColumnsWidth).length > 0;

      if (!this.resizableColumn && !hasVirtualScrollColumnsWidth) {
        return this.headers;
      }

      return this.headers.map((header, index) => {
        const virtualScrollWidth = this.virtualScrollColumnsWidth[header.value];
        const configuredWidth = this.getColumnWidthByField(header.value);
        const columnWidth = this.resizableColumn
          ? header.width || (configuredWidth && `${configuredWidth}px`)
          : (virtualScrollWidth && `${virtualScrollWidth}px`) || header.width;
        const preparedColumnWidth = isNumber(columnWidth) ? `${columnWidth}px` : columnWidth;
        const renderedColumnWidth = this.resizableColumn && preparedColumnWidth
          ? `var(${this.getColumnWidthCssVariable(index)}, ${preparedColumnWidth})`
          : preparedColumnWidth;

        return {
          ...header,
          width: renderedColumnWidth,
          minWidth: renderedColumnWidth,
        };
      });
    },

    headersColspan() {
      return Math.max(this.headersWithWidth.length, 1);
    },

    tableStyle() {
      const style = {};

      if (this.resizableColumn) {
        style['--alarms-list-table-width'] = this.sumOfColumnsWidth
          ? `${this.sumOfColumnsWidth}px`
          : '100%';
      } else if (Object.keys(this.virtualScrollColumnsWidth).length) {
        const tableWidth = Object.values(this.virtualScrollColumnsWidth).reduce(
          (sum, width) => sum + width,
          0,
        );

        style['--alarms-list-table-width'] = `${tableWidth}px`;
      }

      return style;
    },

    tableClasses() {
      return {
        table__grid: this.isColumnsChanging,
        'alarms-list-table__table--fixed':
          this.resizableColumn ||
          this.draggableColumn ||
          Object.keys(this.virtualScrollColumnsWidth).length > 0,
        'alarms-list-table__table--resizing': this.isHorizontalScrollbarHidden,
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
      this.scheduleVirtualScrollColumnsWidthCapture();
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
        if (newResizableColumn) {
          this.virtualScrollColumnsWidth = {};

          if (this.isEmptyColumnsWidthByField) {
            this.$nextTick(() => this.calculateColumnsWidths());
          }

          return;
        }

        if (this.preparedWidget.parameters.isVirtualScrollEnabled) {
          this.scheduleVirtualScrollColumnsWidthCapture();
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

          if (
            this.alarms.length &&
            this.alarms.every(({ _id: id }) => this.bootedRows[id])
          ) {
            this.scheduleVirtualScrollColumnsWidthCapture();
          }

          return;
        }

        this.$intersectionObserver.disconnect();
        this.virtualScrollColumnsWidth = {};
      },
    },
  },

  beforeCreate() {
    this.$intersectionObserver = new TableIntersectionObserver();
    this.$cancelBootedRowsRendering = null;
  },

  mounted() {
    this.$tbodyEl = this.$el.querySelector('tbody');

    if (this.resizableColumn && this.isEmptyColumnsWidthByField) {
      this.$nextTick(() => this.calculateColumnsWidths());
    }

    if (this.eager) {
      this.scheduleActionsPanelBooting();
      this.scheduleVirtualScrollColumnsWidthCapture();
    } else {
      this.setBootedRows(this.alarms);
    }

    this.scheduleStickyScrollHandler();
  },

  beforeUnmount() {
    this.$cancelBootedRowsRendering?.();
    this.asyncBootingActionsPanel.clear();
    this.$intersectionObserver.destroy();
  },

  methods: {
    refreshHorizontalScrollbarAfterColumnResize() {
      if (!this.stickyHorizontalScroll) {
        return;
      }

      this.recalculatingHorizontalScrollbar = true;

      this.$nextTick(() => {
        window.requestAnimationFrame(() => {
          this.setHorizontalScrollbarWidth();
          this.calculateScrollPosition();
          this.setHorizontalScrollPosition();
          this.recalculatingHorizontalScrollbar = false;
        });
      });
    },

    captureVirtualScrollColumnsWidth() {
      if (
        !this.preparedWidget.parameters.isVirtualScrollEnabled ||
        Object.keys(this.virtualScrollColumnsWidth).length ||
        this.resizableColumn ||
        this.draggableColumn ||
        !this.headerCells.length
      ) {
        return;
      }

      const columnsWidth = [...this.headerCells].reduce((acc, headerElement) => {
        const { value } = headerElement.dataset;
        const width = headerElement.getBoundingClientRect().width || headerElement.offsetWidth;

        if (value && width) {
          acc[value] = width;
        }

        return acc;
      }, {});

      if (Object.keys(columnsWidth).length === this.headers.length) {
        this.virtualScrollColumnsWidth = columnsWidth;
        this.$nextTick(() => this.$intersectionObserver.connect());
      }
    },

    scheduleVirtualScrollColumnsWidthCapture() {
      this.$nextTick(() => {
        window.requestAnimationFrame(this.captureVirtualScrollColumnsWidth);
      });
    },

    scheduleActionsPanelBooting() {
      this.$nextTick(() => {
        window.requestAnimationFrame(() => this.asyncBootingActionsPanel.run());
      });
    },

    setBootedRows(items, duration) {
      this.$cancelBootedRowsRendering?.();
      this.$cancelBootedRowsRendering = null;

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
      const initialRows = nearest.slice(0, INITIAL_BOOTED_ROWS_LIMIT);
      const deferredRows = [...nearest.slice(INITIAL_BOOTED_ROWS_LIMIT), ...farthest];

      this.bootedRows = initialRows.reduce((acc, id) => {
        acc[id] = true;

        return acc;
      }, {});

      if (!deferredRows.length) {
        window.requestAnimationFrame(() => {
          this.captureVirtualScrollColumnsWidth();

          if (this.stickyHorizontalScroll) {
            this.setHorizontalScrollbarWidth();
            this.calculateScrollPosition();
            this.setHorizontalScrollPosition();
          }

          this.asyncBootingActionsPanel.run();
        });

        return;
      }

      this.$cancelBootedRowsRendering = recursiveRaf(
        deferredRows,
        async (id, index) => {
          this.bootedRows[id] = true;

          await this.$nextTick();

          if (!index && this.stickyHorizontalScroll) {
            window.requestAnimationFrame(() => {
              this.setHorizontalScrollbarWidth();
              this.calculateScrollPosition();
              this.setHorizontalScrollPosition();
            });
          }
        },
        {
          duration,
          onComplete: () => {
            this.$cancelBootedRowsRendering = null;
            window.requestAnimationFrame(() => {
              this.captureVirtualScrollColumnsWidth();
              this.asyncBootingActionsPanel.run();
            });
          },
        },
      );
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

      if (!expanded && this.expandedAlarmIds.includes(alarm._id)) {
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

    openRootCauseDiagram(entity) {
      this.$modals.show({
        name: MODALS.entitiesRootCauseDiagram,
        config: {
          entity,
          colorIndicator: this.preparedWidget.parameters.rootCauseColorIndicator,
        },
      });
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

  &.alarms-list-table__table--resizing {
    > .v-table__wrapper {
      scrollbar-color: transparent transparent;

      &::-webkit-scrollbar-thumb,
      &::-webkit-scrollbar-track {
        visibility: hidden;
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

  &.alarms-list-table__table--truncated {
    .color-indicator {
      max-width: 100%;
    }

    .alarms-column-cell__layout .alarm-column-cell__text {
      display: grid;
    }

    .alarms-list-row__cell {
      .alarm-column-cell__text,
      .alarm-column-cell__text > span,
      .alarm-column-value {
        display: block;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}

.alarms-list-table__horizontal-scrollbar--hidden {
  pointer-events: none;
  visibility: hidden;
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
