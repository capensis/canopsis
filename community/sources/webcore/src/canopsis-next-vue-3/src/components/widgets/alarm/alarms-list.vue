<template>
  <div class="alarms-list-widget pa-2">
    <c-layout
      v-if="isToolbarShown"
      :class="['alarms-list-widget__toolbar px-4', { 'mb-4': !dense }]"
      align="end"
      justify="start"
      wrap
    >
      <div class="alarms-list-widget__toolbar-cell alarms-list-widget__toolbar-cell--advanced-search">
        <c-alarm-advanced-search
          :searches="searches"
          @submit="updateSearch"
          @reset="resetSearch"
          @toggle-pin:search="togglePinSearch"
          @remove:search="removeSearch"
        />
      </div>

      <div v-if="hasAccessToCategory" class="alarms-list-widget__toolbar-cell">
        <c-entity-category-field
          :category="query.category"
          class="ma-0"
          hide-details
          @update:category="updateCategory"
        />
      </div>

      <div v-if="hasAccessToCorrelation" class="alarms-list-widget__toolbar-cell">
        <c-enabled-field
          :model-value="query.correlation"
          :label="$t('common.correlation')"
          class="mt-0"
          color="primary"
          hide-details
          inline
          no-margin
          @update:model-value="updateCorrelation"
        />
      </div>

      <c-layout
        v-if="hasAccessToUserFilter"
        class="alarms-list-widget__toolbar-cell alarms-list-widget__toolbar-cell--filter"
        align="end"
      >
        <filter-selector
          :model-value="query.filter"
          :locked-value="query.lockedFilter"
          :label="$t('settings.selectAFilter')"
          :filters="userPreferenceFilters"
          :locked-filters="widget.filters"
          :clearable="preparedWidget.parameters.clearFilterEnabled"
          hide-details
          @update:model-value="updateSelectedFilter"
        />

        <filters-list-btn
          :widget-id="widget._id"
          addable
          editable
          private
          with-alarm
          with-entity
          with-pbehavior
        />
      </c-layout>

      <div v-if="hasAccessToFilterByBookmark" class="alarms-list-widget__toolbar-cell">
        <c-enabled-field
          :model-value="query.only_bookmarks"
          :label="$t('alarm.filterByBookmark')"
          class="mt-0"
          color="primary"
          hide-details
          inline
          no-margin
          @update:model-value="updateOnlyBookmarks"
        />
      </div>

      <div class="alarms-list-widget__toolbar-cell alarms-list-widget__toolbar-cell--tags">
        <c-alarm-tag-field
          :value="query.tags"
          :label="$t('common.tag')"
          :show-count="2"
          multiple
          combobox
          @update:value="updateTags"
        />
      </div>

      <div
        v-if="hasAccessToUserRemediationInstructionsFilter"
        class="alarms-list-widget__toolbar-cell alarms-list-widget__toolbar-cell--instructions"
      >
        <alarms-list-remediation-instructions-filter
          :model-value="query.instructionsFilter"
          @update:model-value="updateRemediationInstructionsFilters"
        />
      </div>

      <c-layout class="alarms-list-widget__toolbar-cell alarms-list-widget__toolbar-cell--live-reporting">
        <c-chip
          v-if="activeRange"
          color="primary"
          title=""
          class="px-2"
          closable
          @close="removeHistoryFilter"
        >
          {{ $t(`quickRanges.types.${activeRange.value}`) }}
        </c-chip>

        <c-action-btn
          :tooltip="$t('alarm.liveReporting')"
          :color="activeRange ? 'primary' : ''"
          icon="schedule"
          @click="showEditLiveReportModal"
        />
      </c-layout>

      <div
        v-if="hasAccessToExportAsCsv"
        class="alarms-list-widget__toolbar-cell alarms-list-widget__toolbar-cell--export"
      >
        <c-action-btn
          :loading="downloading"
          :tooltip="$t('settings.exportAsCsv')"
          icon="cloud_download"
          @click="exportAlarmsList"
        />
      </div>
    </c-layout>

    <c-layout
      v-if="shownTopPagination"
      class="alarms-list-widget__top-pagination px-4 py-2 position-relative"
      align="center"
    >
      <div class="alarms-list-widget__top-pagination-left">
        <c-density-btn-toggle
          v-if="isTableControlsShown"
          :value="dense"
          button-size="x-small"
          density="compact"
          icon-size="16"
          @change="updateDense"
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
              :refresh-alarms-list="fetchList"
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
        class="alarms-list-widget__top-pagination-center"
        type="top"
        @update:page="updatePage"
      />

      <c-grid-btns
        v-if="hasColumnEditButtons"
        :editing="isColumnsChanging"
        class="alarms-list-widget__top-pagination-right"
        @toggle="toggleColumnEditingMode"
        @reset="resetColumnsSettings"
      />
    </c-layout>

    <v-data-table-server
      ref="dataTable"
      :headers="headersWithWidth"
      :hide-default-footer="true"
      :items="alarms"
      :items-length="totalItems"
      :items-per-page="options.itemsPerPage"
      :loading="alarmsPending"
      :page="options.page"
      :sort-by="options.sortBy"
      :disable-sort="isColumnsChanging"
      :density="tableDensity"
      :class="tableClasses"
      :style="tableStyle"
      class="alarms-list-widget__table"
      multi-sort
      @update:options="updateOptions"
    >
      <template
        v-for="header in headersWithWidth"
        :key="`header.${header.key}`"
        #[`header.${header.key}`]="{ column, isSorted, sortBy, getSortIcon }"
      >
        <div v-if="header.key === rowActionsColumnKey" class="alarms-list-widget__select-all">
          <div class="alarms-list-widget__select-all-checkbox">
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
          v-if="draggingMode && header.key !== rowActionsColumnKey"
          :data-value="header.value"
          class="table__dragging-handler"
          @click.stop=""
        />
        <span
          v-if="resizingMode && header.key !== rowActionsColumnKey"
          class="table__resize-handler"
          @mousedown.stop.prevent="startColumnResize(header.value)"
          @click.stop=""
        />
      </template>

      <template #[`item.row-icons`]="{ item }">
        <c-layout class="alarms-list-widget__row-icons" align="center">
          <c-layout class="alarms-list-widget__row-checkbox" align="center" justify="center">
            <c-light-checkbox
              v-if="selectable"
              :model-value="isAlarmSelected(getTableItem(item))"
              :disabled="!isAlarmSelectable(getTableItem(item))"
              @update:model-value="toggleAlarmSelection(getTableItem(item), $event)"
            />
          </c-layout>

          <c-layout v-if="hasAlarmInstruction(getTableItem(item))" align="center">
            <alarms-list-row-instructions-icon :alarm="getTableItem(item)" />
          </c-layout>

          <c-layout v-if="getTableItem(item).bookmark" align="center">
            <alarms-list-row-bookmark-icon />
          </c-layout>
        </c-layout>
      </template>

      <template
        v-for="column in columns"
        :key="`item.${column.value}`"
        #[`item.${column.value}`]="{ item }"
      >
        <alarm-column-value
          :alarm="getTableItem(item)"
          :widget="preparedWidget"
          :column="column"
          :selected-tags="query.tags"
          :small="isSmallDensity"
          @select:tag="selectTag"
          @remove:tag="removeTag"
        />
        <span
          v-if="resizingMode"
          class="table__resize-handler"
          @mousedown.stop.prevent="startColumnResize(column.value)"
          @click.stop=""
        />
      </template>

      <template #[`item.actions`]="{ item }">
        <c-layout justify="end" align="center">
          <actions-panel
            v-if="!hideActions"
            :item="getTableItem(item)"
            :widget="preparedWidget"
            :parent-alarm="parentAlarm"
            :refresh-alarms-list="fetchList"
            :small="isSmallDensity"
            :inline-count="actionsInlineCount"
            :ignore-media-query="resizableColumn"
          />
        </c-layout>
        <span
          v-if="resizingMode"
          class="table__resize-handler"
          @mousedown.stop.prevent="startColumnResize('actions')"
          @click.stop=""
        />
      </template>

      <template #no-data>
        <tr>
          <td :colspan="headersColspan" class="text-center text-medium-emphasis py-6">
            {{ $t('common.noData') }}
          </td>
        </tr>
      </template>
    </v-data-table-server>

    <c-table-pagination
      :items="paginationItems"
      :total-items="totalItems"
      :items-per-page="options.itemsPerPage"
      :page="options.page"
      class="mt-2"
      @update:page="updatePage"
      @update:items-per-page="updateItemsPerPage"
      @update:pagination="updatePaginationOptions"
    />
  </div>
</template>

<script>
import { get, isEqual, isNumber, isObject, omit, pick } from 'lodash';
import { mapActions, mapStores } from 'pinia';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';
import {
  ALARM_ACTION_BUTTON_MARGINS,
  ALARM_ACTION_BUTTON_WIDTHS,
  ALARM_ACTIONS_PADDINGS,
  COLOR_INDICATOR_TYPES,
  DEFAULT_ALARM_ACTIONS_INLINE_COUNT,
  DENSE_TYPES,
  LIVE_REPORTING_QUICK_RANGES,
  MODALS,
  RESIZING_CELLS_CONTENTS_BEHAVIORS,
  USER_PERMISSIONS,
  WIDGET_GRID_SIZES_KEYS,
} from '@/constants';
import { findQuickRangeValue } from '@/helpers/date/date-intervals';
import { isActionAvailableForAlarm } from '@/helpers/entities/alarm/form';
import { calculateAlarmLinksColumnWidth } from '@/helpers/entities/alarm/list';
import { convertAlarmWidgetParametersToActiveColumns } from '@/helpers/entities/alarm/query';
import { getAlarmListExportDownloadFileUrl } from '@/helpers/entities/alarm/url';
import {
  getAlarmsListWidgetColumnComponentGetter,
  getAlarmsListWidgetColumnValueFilter,
  prepareAlarmListWidget,
} from '@/helpers/entities/widget/forms/alarm';
import { selectedToQuery, widgetToExportQueryColumns } from '@/helpers/entities/widget/query';
import { setSeveralFields } from '@/helpers/immutable';
import { permissionsWidgetsAlarmsListBookmark } from '@/mixins/permissions/widgets/alarms-list/bookmark';
import { permissionsWidgetsAlarmsListCategory } from '@/mixins/permissions/widgets/alarms-list/category';
import { permissionsWidgetsAlarmsListCorrelation } from '@/mixins/permissions/widgets/alarms-list/correlation';
import { permissionsWidgetsAlarmsListFilters } from '@/mixins/permissions/widgets/alarms-list/filters';
import {
  permissionsWidgetsAlarmsListRemediationInstructionsFilters,
} from '@/mixins/permissions/widgets/alarms-list/remediation-instructions-filters';
import { widgetAdvancedSearchSavedItemsMixin } from '@/mixins/widget/advanced-search/saved-items';
import {
  widgetColumnDraggingAlarmMixin,
  widgetColumnResizingAlarmMixin,
} from '@/mixins/widget/columns';
import { exportMixinCreator } from '@/mixins/widget/export';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { useAlarmStore } from '@/stores/modules/alarm';
import { useWidgetStore } from '@/stores/modules/widget';

import ActionsPanel from './actions/actions-panel.vue';
import MassActionsPanel from './actions/mass-actions-panel.vue';
import AlarmColumnValue from './columns-formatting/alarm-column-value.vue';
import AlarmHeaderCell from './headers-formatting/alarm-header-cell.vue';
import AlarmsListRemediationInstructionsFilter from './partials/alarms-list-remediation-instructions-filter.vue';
import AlarmsListRowBookmarkIcon from './partials/alarms-list-row-bookmark-icon.vue';
import AlarmsListRowInstructionsIcon from './partials/alarms-list-row-instructions-icon.vue';

const DISPLAY_NAME_TO_WIDGET_GRID_SIZE = {
  xs: WIDGET_GRID_SIZES_KEYS.mobile,
  sm: WIDGET_GRID_SIZES_KEYS.mobile,
  md: WIDGET_GRID_SIZES_KEYS.tablet,
  lg: WIDGET_GRID_SIZES_KEYS.desktop,
  xl: WIDGET_GRID_SIZES_KEYS.desktop,
  xxl: WIDGET_GRID_SIZES_KEYS.desktop,
};

/**
 * Alarm-list widget runtime foundation.
 *
 * The dedicated table, alarm query state and actions are migrated in the next
 * AlarmsList stages; this component already consumes prepared widget
 * parameters so create/edit/duplicate flows use the final settings payload.
 */
export default {
  components: {
    ActionsPanel,
    AlarmColumnValue,
    AlarmHeaderCell,
    AlarmsListRowBookmarkIcon,
    AlarmsListRowInstructionsIcon,
    AlarmsListRemediationInstructionsFilter,
    FilterSelector,
    FiltersListBtn,
    MassActionsPanel,
  },

  mixins: [
    widgetFetchQueryMixin,
    widgetAdvancedSearchSavedItemsMixin,
    widgetColumnResizingAlarmMixin,
    widgetColumnDraggingAlarmMixin,
    permissionsWidgetsAlarmsListBookmark,
    permissionsWidgetsAlarmsListCategory,
    permissionsWidgetsAlarmsListCorrelation,
    permissionsWidgetsAlarmsListFilters,
    permissionsWidgetsAlarmsListRemediationInstructionsFilters,
    exportMixinCreator({
      createExport: 'createAlarmsListExport',
      fetchExport: 'fetchAlarmsListExport',
    }),
  ],

  props: {
    widget: {
      type: Object,
      required: true,
    },
    hideActions: {
      type: Boolean,
      default: false,
    },
    hideMassSelection: {
      type: Boolean,
      default: false,
    },
    hideToolbar: {
      type: Boolean,
      default: false,
    },
    parentAlarm: {
      type: Object,
      default: null,
    },
  },

  data() {
    return {
      downloading: false,
      keepSelectedAfterAction: this.widget.parameters?.keepSelectedAfterAction ?? false,
      rowActionsColumnKey: 'row-icons',
      selectedAlarmIds: [],
    };
  },
  computed: {
    ...mapStores(useAlarmStore),

    preparedWidget() {
      return prepareAlarmListWidget(this.widget);
    },

    columns() {
      const { showRootCauseByStateClick = true } = this.preparedWidget.parameters ?? {};

      return (this.preparedWidget.parameters?.widgetColumns ?? []).map(column => ({
        ...column,
        popupTemplate: this.infoPopupsMap[column.value] ?? this.infoPopupsMap[`alarm.${column.value}`],
        filter: getAlarmsListWidgetColumnValueFilter(column.value),
        getComponent: getAlarmsListWidgetColumnComponentGetter(column, {
          showRootCauseByStateClick,
          isCellContentTruncated: this.isCellContentTruncated,
        }),
        colorIndicatorEnabled: Object.values(COLOR_INDICATOR_TYPES).includes(column.colorIndicator),
      }));
    },

    infoPopupsMap() {
      return (this.preparedWidget.parameters?.infoPopups ?? []).reduce((acc, { column, template }) => {
        acc[column] = template;

        return acc;
      }, {});
    },

    selectable() {
      return !this.hideMassSelection && !this.preparedWidget.parameters.kiosk?.hideMassSelection;
    },

    hasRowActions() {
      return this.selectable || this.alarms.some(alarm => this.hasAlarmInstruction(alarm) || alarm.bookmark);
    },

    columnsSettings() {
      return this.userPreference.content?.columns_settings;
    },

    resizableColumn() {
      return !!this.preparedWidget.parameters?.columns?.resizable;
    },

    draggableColumn() {
      return !!this.preparedWidget.parameters?.columns?.draggable;
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

    isColumnsChanging() {
      return this.resizingMode || this.draggingMode;
    },

    baseHeaders() {
      const headers = this.columns.map((column) => {
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
          key: this.rowActionsColumnKey,
          value: this.rowActionsColumnKey,
          title: '',
          text: '',
          sortable: false,
          width: 100,
          headerProps: {
            'data-value': this.rowActionsColumnKey,
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
      return this.draggableColumn
        ? [...this.baseHeaders].sort(
          (a, b) => this.getColumnPositionByField(a.value) - this.getColumnPositionByField(b.value),
        )
        : this.baseHeaders;
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
          '--alarms-list-table-width': `${this.sumOfColumnsWidth}px`,
        };
      }

      return {};
    },

    tableClasses() {
      return {
        table__grid: this.isColumnsChanging,
        'alarms-list-widget__table--fixed': this.resizableColumn || this.draggableColumn,
        'alarms-list-widget__table--wrapped': this.isCellContentWrapped,
        'alarms-list-widget__table--truncated': this.isCellContentTruncated,
        'alarms-list-widget__table--dragging': this.draggingMode,
      };
    },

    alarms() {
      return this.alarmStore.getListByWidgetId(this.widget._id);
    },

    alarmsMeta() {
      return this.alarmStore.getMetaByWidgetId(this.widget._id);
    },

    alarmsPending() {
      return this.alarmStore.getPendingByWidgetId(this.widget._id);
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

    defaultKeepSelectedAfterAction() {
      return this.preparedWidget.parameters.keepSelectedAfterAction ?? false;
    },

    displayGridSize() {
      return DISPLAY_NAME_TO_WIDGET_GRID_SIZE[this.$vuetify.display.name] ?? WIDGET_GRID_SIZES_KEYS.desktop;
    },

    widgetGridParameters() {
      const gridParameters = this.widget.grid_parameters ?? {};

      return pick(gridParameters[this.displayGridSize] ?? gridParameters, ['x', 'y', 'w', 'h']);
    },

    alarmsFetchingParams() {
      return this.alarmStore.getFetchingParamsByWidgetId(this.widget._id);
    },

    totalItems() {
      return this.alarmsMeta.total_count ?? 0;
    },

    dense() {
      return this.userPreference.content.dense ?? this.preparedWidget.parameters.dense;
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

    activeRange() {
      const { tstart, tstop } = this.query;

      if (tstart || tstop) {
        return findQuickRangeValue(
          tstart,
          tstop,
          LIVE_REPORTING_QUICK_RANGES,
          LIVE_REPORTING_QUICK_RANGES.custom,
        );
      }

      return null;
    },

    hasAccessToExportAsCsv() {
      return this.checkAccess(USER_PERMISSIONS.business.alarmsList.actions.exportAsCsv);
    },

    isToolbarShown() {
      return !this.hideToolbar;
    },

    isTableControlsShown() {
      return !!this.columns.length;
    },

    shownTopPagination() {
      return !!this.totalItems;
    },

    paginationItems() {
      return [5, 10, 20, 50, 100, 200];
    },

    userPreferenceFilters() {
      return this.userPreference.filters ?? [];
    },
  },
  watch: {
    alarms() {
      this.synchronizeSelectedAlarms();
    },

    defaultKeepSelectedAfterAction(value) {
      this.keepSelectedAfterAction = value;
    },

    columns() {
      if (this.isColumnsChanging) {
        this.persistColumnsSettings();

        this.disableDraggingMode();
        this.disableResizingMode();
      }
    },

    columnsSettings: {
      immediate: true,
      deep: true,
      handler() {
        if (!this.draggingMode && this.columnsSettings?.columns_position) {
          this.setColumnsPosition(this.columnsSettings?.columns_position);
        }

        if (!this.resizingMode && this.columnsSettings?.columns_width) {
          this.setColumnsWidth(this.columnsSettings?.columns_width);
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
    },
  },
  created() {
    this.actualizeUsedProperties();
  },
  mounted() {
    if (this.resizableColumn && this.isEmptyColumnsWidthByField) {
      this.$nextTick(() => this.calculateColumnsWidths());
    }
  },
  methods: {
    ...mapActions(useAlarmStore, {
      fetchAlarmsList: 'fetchList',
      fetchAlarmsListExport: 'fetchAlarmsListExport',
      createAlarmsListExport: 'createAlarmsListExport',
    }),
    ...mapActions(useWidgetStore, {
      updateWidget: 'update',
    }),

    async fetchList() {
      if (!this.columns.length) {
        return;
      }

      const params = this.getQuery();

      if (this.alarmsPending && isEqual(params, this.alarmsFetchingParams)) {
        return;
      }

      await this.fetchAlarmsList({
        widgetId: this.widget._id,
        params,
      });
    },

    actualizeUsedProperties() {
      const activeColumns = convertAlarmWidgetParametersToActiveColumns(
        this.preparedWidget.parameters,
      );

      if (isEqual(activeColumns, this.preparedWidget.parameters.usedAlarmProperties)) {
        return;
      }

      this.updateWidget({
        id: this.widget._id,
        data: setSeveralFields(this.widget, {
          'parameters.usedAlarmProperties': activeColumns,
        }),
      });
    },

    updateCorrelation(correlation) {
      this.updateContentInUserPreference({ isCorrelationEnabled: correlation });

      this.query = {
        ...this.query,
        page: 1,
        correlation,
      };
    },

    updateTags(tags) {
      this.updateContentInUserPreference({ tags });

      this.query = {
        ...this.query,
        page: 1,
        tags,
      };
    },

    selectTag(tag) {
      this.updateTags([...(this.query.tags || []), tag]);
    },

    removeTag(tag) {
      this.updateTags((this.query.tags || []).filter(queryTag => queryTag !== tag));
    },

    updateDense(dense) {
      this.updateContentInUserPreference({ dense });
    },

    updateColumnsSettings(columnsSettings) {
      this.updateContentInUserPreference({ columns_settings: columnsSettings });
    },

    persistColumnsSettings() {
      const settings = {};

      if (this.resizingMode) {
        settings.columns_width = this.columnsWidthByField;
      }

      if (this.draggingMode) {
        settings.columns_position = this.columnsPositionByField;
      }

      this.updateColumnsSettings(settings);
    },

    toggleColumnEditingMode() {
      if (this.isColumnsChanging) {
        this.persistColumnsSettings();
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

    updateOptions(options) {
      if (this.isColumnsChanging) {
        return;
      }

      this.options = options;
    },

    getHeaderSortIndex(column, sortBy = []) {
      return sortBy.findIndex(item => item.key === column.key) + 1;
    },

    updateItemsPerPage(itemsPerPage) {
      this.updateContentInUserPreference({ itemsPerPage });

      this.query = {
        ...this.query,
        itemsPerPage,
      };
    },

    updatePage(page) {
      this.query = {
        ...this.query,
        page,
      };
    },

    updatePaginationOptions({ page, itemsPerPage }) {
      this.updateContentInUserPreference({ itemsPerPage });

      this.query = {
        ...this.query,
        page,
        itemsPerPage,
      };
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

    updateOnlyBookmarks(onlyBookmarks) {
      this.updateContentInUserPreference({ onlyBookmarks });

      this.query = {
        ...this.query,
        page: 1,
        only_bookmarks: onlyBookmarks,
      };
    },

    updateCategory(category) {
      const categoryId = category?._id;

      this.updateContentInUserPreference({ category: categoryId });

      this.query = {
        ...this.query,
        page: 1,
        category: categoryId,
      };
    },

    updateSelectedFilter(filter) {
      this.updateContentInUserPreference({ mainFilter: filter });

      this.query = {
        ...this.query,
        page: 1,
        filter,
      };
    },

    updateRemediationInstructionsFilters(instructionsFilter) {
      this.updateContentInUserPreference({ instructionsFilter });

      this.query = {
        ...this.query,
        page: 1,
        instructionsFilter,
      };
    },

    removeHistoryFilter() {
      this.query = {
        ...omit(this.query, ['tstart', 'tstop']),
        page: 1,
      };
    },

    showEditLiveReportModal() {
      this.$modals.show({
        name: MODALS.editLiveReporting,
        config: {
          ...pick(this.query, ['tstart', 'tstop', 'time_field']),
          action: (params) => {
            this.query = {
              ...this.query,
              ...params,
              page: 1,
            };
          },
        },
      });
    },

    getExportQuery() {
      const query = this.getQuery();
      const { exportCsvSeparator, exportCsvDatetimeFormat } = this.preparedWidget.parameters;

      return {
        ...pick(query, [
          'category',
          'correlation',
          'opened',
          'tstart',
          'tstop',
          'only_bookmarks',
          'tags',
          'instruction_filter_type',
          'instruction_type',
          'instruction_statuses',
          'instruction_ids',
          'alarm_pattern',
          'entity_pattern',
          'pbehavior_pattern',
        ]),

        fields: widgetToExportQueryColumns(this.preparedWidget),
        filters: query.filters,
        separator: exportCsvSeparator,
        /**
         * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/3997
         */
        time_format: isObject(exportCsvDatetimeFormat)
          ? exportCsvDatetimeFormat.value
          : exportCsvDatetimeFormat,
        search: selectedToQuery(this.selectedAlarms, query.search),
      };
    },

    async exportAlarmsList() {
      this.downloading = true;

      try {
        const fileData = await this.generateFile({
          data: this.getExportQuery(),
          widgetId: this.widget._id,
        });

        this.downloadFile(getAlarmListExportDownloadFileUrl(fileData._id));
      } catch {
        this.$popups.error({ text: this.$t('alarm.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.alarms-list-widget {
  height: 100%;
}

.alarms-list-widget__toolbar {
  gap: 16px;
}

.alarms-list-widget__toolbar-cell {
  min-width: 180px;
}

.alarms-list-widget__toolbar-cell--advanced-search {
  flex: 1 1 360px;
  min-width: 320px;
}

.alarms-list-widget__toolbar-cell--filter {
  min-width: 260px;
}

.alarms-list-widget__toolbar-cell--tags {
  flex: 0 1 25%;
  min-width: 280px;
}

.alarms-list-widget__toolbar-cell--instructions {
  min-width: 280px;
}

.alarms-list-widget__toolbar-cell--live-reporting {
  display: flex;
  min-width: auto;
}

.alarms-list-widget__toolbar-cell--export {
  min-width: auto;
}

.alarms-list-widget__top-pagination {
  min-height: 48px;
  background: rgb(var(--v-theme-surface));
  z-index: 2;
}

.alarms-list-widget__top-pagination-left {
  align-items: center;
  display: flex;
  gap: 8px;
  padding-right: 80px;
}

.alarms-list-widget__top-pagination-center {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

.alarms-list-widget__top-pagination-right {
  position: absolute;
  right: 16px;
}

.alarms-list-widget__table {
  min-height: 100%;

  &--fixed {
    :deep(.v-table__wrapper > table) {
      table-layout: fixed;
      width: var(--alarms-list-table-width);
      max-width: unset;
      min-width: var(--alarms-list-table-width);
    }

    :deep(thead > tr > th) {
      word-break: break-all;
      white-space: pre-wrap;
    }
  }

  &--dragging {
    :deep(thead > tr > th) {
      position: relative;
    }

    :deep(.v-data-table-header__content) {
      pointer-events: none;
    }
  }

  &--wrapped {
    :deep(tbody > tr > td:not(:last-of-type)) {
      word-break: break-all;
      word-wrap: break-word;
    }
  }

  &--truncated {
    :deep(.color-indicator) {
      max-width: 100%;
    }

    :deep(.alarms-column-cell__layout .alarm-column-cell__text) {
      display: grid;
      min-width: 0;
    }

    :deep(.alarm-column-cell__text),
    :deep(.alarm-column-cell__text > span),
    :deep(.alarm-column-value) {
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      display: block;
      max-width: 100%;
      min-width: 0;
    }

    :deep(.v-data-table-header__span--ellipsis) {
      display: block;
      min-width: 0;
      overflow: hidden;
      text-overflow: ellipsis;
      white-space: nowrap;
    }
  }
}

.alarms-list-widget__draggable-column {
  position: relative;
}

.alarms-list-widget__select-all {
  align-items: center;
  display: flex;
  justify-content: flex-start;
  width: 100px;
}

.alarms-list-widget__select-all-checkbox {
  align-items: center;
  display: flex;
  height: 24px;
  justify-content: center;
  width: 24px;
}

.alarms-list-widget__row-icons {
  min-width: 100px;
  width: 100px;
}

.alarms-list-widget__row-checkbox {
  height: 24px;
  max-width: 24px;
  width: 24px;
}
</style>
