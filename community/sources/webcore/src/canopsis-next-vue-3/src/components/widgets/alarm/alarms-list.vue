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
      </div>

      <c-pagination
        :page="options.page"
        :limit="options.itemsPerPage"
        :total="totalItems"
        class="alarms-list-widget__top-pagination-center"
        type="top"
        @update:page="updatePage"
      />
    </c-layout>

    <v-data-table-server
      ref="alarmsTable"
      :headers="headers"
      :hide-default-footer="true"
      :items="alarms"
      :items-length="totalItems"
      :items-per-page="options.itemsPerPage"
      :loading="alarmsPending"
      :page="options.page"
      :sort-by="options.sortBy"
      :density="tableDensity"
      class="alarms-list-widget__table"
      @update:options="updateOptions"
    >
      <template
        v-for="header in headers"
        :key="`header.${header.key}`"
        #[`header.${header.key}`]
      >
        <alarm-header-cell :header="header" />
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
      </template>

      <template #item.actions>
        <c-layout justify="end" align="center" />
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
import { isEqual, isObject, omit, pick } from 'lodash';
import { mapActions, mapStores } from 'pinia';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';
import {
  COLOR_INDICATOR_TYPES,
  DENSE_TYPES,
  LIVE_REPORTING_QUICK_RANGES,
  MODALS,
  RESIZING_CELLS_CONTENTS_BEHAVIORS,
  USER_PERMISSIONS,
} from '@/constants';
import { findQuickRangeValue } from '@/helpers/date/date-intervals';
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
import { exportMixinCreator } from '@/mixins/widget/export';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { useAlarmStore } from '@/stores/modules/alarm';
import { useWidgetStore } from '@/stores/modules/widget';

import AlarmColumnValue from './columns-formatting/alarm-column-value.vue';
import AlarmHeaderCell from './headers-formatting/alarm-header-cell.vue';
import AlarmsListRemediationInstructionsFilter from './partials/alarms-list-remediation-instructions-filter.vue';

/**
 * Alarm-list widget runtime foundation.
 *
 * The dedicated table, alarm query state and actions are migrated in the next
 * AlarmsList stages; this component already consumes prepared widget
 * parameters so create/edit/duplicate flows use the final settings payload.
 */
export default {
  components: {
    AlarmColumnValue,
    AlarmHeaderCell,
    AlarmsListRemediationInstructionsFilter,
    FilterSelector,
    FiltersListBtn,
  },

  mixins: [
    widgetFetchQueryMixin,
    widgetAdvancedSearchSavedItemsMixin,
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
    };
  },
  computed: {
    ...mapStores(useAlarmStore),

    preparedWidget() {
      return prepareAlarmListWidget(this.widget);
    },

    columns() {
      const {
        showRootCauseByStateClick = true,
        columns: columnsSettings = {},
      } = this.preparedWidget.parameters ?? {};

      return (this.preparedWidget.parameters?.widgetColumns ?? []).map(column => ({
        ...column,
        popupTemplate: this.infoPopupsMap[column.value] ?? this.infoPopupsMap[`alarm.${column.value}`],
        filter: getAlarmsListWidgetColumnValueFilter(column.value),
        getComponent: getAlarmsListWidgetColumnComponentGetter(column, {
          showRootCauseByStateClick,
          isCellContentTruncated:
            columnsSettings.cells_content_behavior
            === RESIZING_CELLS_CONTENTS_BEHAVIORS.truncate,
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

    headers() {
      const headers = this.columns.map(column => ({
        ...column,
        key: column.value,
        title: column.text,
      }));

      if (!this.hideActions) {
        const actionsText = this.$t('common.actionsLabel');

        headers.push({
          key: 'actions',
          value: 'actions',
          title: actionsText,
          text: actionsText,
          sortable: false,
        });
      }

      return headers;
    },

    headersColspan() {
      return Math.max(this.headers.length, 1);
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
  created() {
    this.actualizeUsedProperties();
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
      const selected = this.$refs.alarmsTable?.selected ?? [];

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
        search: selectedToQuery(selected, query.search),
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
  padding-right: 80px;
}

.alarms-list-widget__top-pagination-center {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
}

.alarms-list-widget__table {
  min-height: 100%;
}
</style>
