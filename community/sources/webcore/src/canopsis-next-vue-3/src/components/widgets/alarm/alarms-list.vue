<template>
  <div class="alarms-list-widget">
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

    <alarms-list-table-with-pagination
      :widget="preparedWidget"
      :alarms="alarms"
      :columns="columns"
      :columns-settings="columnsSettings"
      :meta="alarmsMeta"
      :query="query"
      :parent-alarm="parentAlarm"
      :loading="alarmsPending"
      :hide-actions="hideActions"
      :selected-tags="query.tags"
      :densable="isTableControlsShown"
      :dense="dense"
      :sticky-header="preparedWidget.parameters.sticky_header"
      :sticky-horizontal-scroll="preparedWidget.parameters.sticky_horizontal_scroll"
      :widget-grid-parameters="widgetGridParameters"
      :default-keep-selected-after-action="defaultKeepSelectedAfterAction"
      :refresh-alarms-list="fetchList"
      :selectable="selectable"
      class="mt-2"
      expandable
      @update:query="updateTableQuery"
      @update:columns-settings="updateColumnsSettings"
      @update:dense="updateDense"
      @select:tag="selectTag"
      @remove:tag="removeTag"
    />
  </div>
</template>

<script>
import { isEqual, isObject, omit, pick } from 'lodash';
import { mapActions, mapStores } from 'pinia';

import {
  LIVE_REPORTING_QUICK_RANGES,
  MODALS,
  USER_PERMISSIONS,
  DISPLAY_NAME_TO_WIDGET_GRID_SIZES_KEYS_MAP,
  WIDGET_GRID_SIZES_KEYS,
} from '@/constants';

import { useAlarmStore } from '@/stores/modules/alarm';
import { useWidgetStore } from '@/stores/modules/widget';

import { findQuickRangeValue } from '@/helpers/date/date-intervals';
import { convertAlarmWidgetParametersToActiveColumns } from '@/helpers/entities/alarm/query';
import { getAlarmListExportDownloadFileUrl } from '@/helpers/entities/alarm/url';
import { prepareAlarmListWidget } from '@/helpers/entities/widget/forms/alarm';
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
import { widgetAlarmsSocketMixin } from '@/mixins/widget/alarms-socket';
import { exportMixinCreator } from '@/mixins/widget/export';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';

import AlarmsListRemediationInstructionsFilter from './partials/alarms-list-remediation-instructions-filter.vue';
import AlarmsListTableWithPagination from './partials/alarms-list-table-with-pagination.vue';

/**
 * Alarm-list widget runtime foundation.
 *
 * The dedicated table, alarm query state and actions are migrated in the next
 * AlarmsList stages; this component already consumes prepared widget
 * parameters so create/edit/duplicate flows use the final settings payload.
 */
export default {
  components: {
    AlarmsListRemediationInstructionsFilter,
    AlarmsListTableWithPagination,
    FilterSelector,
    FiltersListBtn,
  },

  mixins: [
    widgetFetchQueryMixin,
    widgetAdvancedSearchSavedItemsMixin,
    widgetPeriodicRefreshMixin,
    widgetAlarmsSocketMixin,
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
      return (this.preparedWidget.parameters?.widgetColumns ?? []).map(column => ({
        ...column,
        popupTemplate: this.infoPopupsMap[column.value] ?? this.infoPopupsMap[`alarm.${column.value}`],
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

    columnsSettings() {
      return this.userPreference.content?.columns_settings;
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

    selectedAlarms() {
      return [];
    },

    defaultKeepSelectedAfterAction() {
      return this.preparedWidget.parameters.keepSelectedAfterAction ?? false;
    },

    displayGridSize() {
      return DISPLAY_NAME_TO_WIDGET_GRID_SIZES_KEYS_MAP[this.$vuetify.display.name] ?? WIDGET_GRID_SIZES_KEYS.desktop;
    },

    widgetGridParameters() {
      const gridParameters = this.widget.grid_parameters ?? {};

      return pick(gridParameters[this.displayGridSize] ?? gridParameters, ['x', 'y', 'w', 'h']);
    },

    alarmsFetchingParams() {
      return this.alarmStore.getFetchingParamsByWidgetId(this.widget._id);
    },

    dense() {
      return this.userPreference.content.dense ?? this.preparedWidget.parameters.dense;
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
      fetchAlarmsDetailsList: 'fetchAlarmDetailsList',
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

      this.fetchAlarmsDetailsList({ widgetId: this.widget._id });

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

    updateTableQuery(query) {
      if (query.itemsPerPage !== this.query.itemsPerPage) {
        this.updateContentInUserPreference({ itemsPerPage: query.itemsPerPage });
      }

      this.query = {
        ...query,
      };
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
</style>
