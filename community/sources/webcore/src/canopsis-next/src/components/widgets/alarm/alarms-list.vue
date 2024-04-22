<template lang="pug">
  div
    v-layout(
      v-if="!hideToolbar",
      :class="{ 'mb-4': !dense }",
      row,
      wrap,
      justify-space-between,
      align-end
    )
      v-flex
        c-advanced-search-field(
          :query.sync="query",
          :columns="widget.parameters.widgetColumns",
          :tooltip="$t('alarm.advancedSearch')",
          :items="searches",
          combobox,
          @submit="updateSearchesInUserPreferences",
          @toggle-pin="togglePinSearchInUserPreferences",
          @remove="removeSearchFromUserPreferences"
        )
      v-flex(v-if="hasAccessToCategory")
        c-entity-category-field.mr-3.mt-0(:category="query.category", hide-details, @input="updateCategory")
      v-flex(v-if="hasAccessToCorrelation")
        v-switch.mt-0(
          :value="query.correlation",
          :label="$t('common.correlation')",
          color="primary",
          hide-details,
          @change="updateCorrelation"
        )
      v-flex
        v-layout(v-if="hasAccessToUserFilter", row, align-end)
          filter-selector(
            :value="query.filter",
            :locked-value="query.lockedFilter",
            :label="$t('settings.selectAFilter')",
            :filters="userPreference.filters",
            :locked-filters="widget.filters",
            :disabled="!hasAccessToListFilters",
            :clearable="!widget.parameters.clearFilterDisabled",
            hide-details,
            @input="updateSelectedFilter"
          )
          filters-list-btn(
            v-if="hasAccessToAddFilter || hasAccessToEditFilter",
            :widget-id="widget._id",
            :addable="hasAccessToAddFilter",
            :editable="hasAccessToEditFilter",
            private,
            with-alarm,
            with-entity,
            with-pbehavior
          )
      v-flex(v-if="hasAccessToFilterByBookmark")
        v-switch.mt-0(
          :value="query.only_bookmarks",
          :label="$t('alarm.filterByBookmark')",
          color="primary",
          hide-details,
          @change="updateOnlyBookmarks"
        )
      v-flex
        alarms-list-remediation-instructions-filters(
          :filters.sync="remediationInstructionsFilters",
          :locked-filters.sync="widgetRemediationInstructionsFilters",
          :editable="hasAccessToEditRemediationInstructionsFilter",
          :addable="hasAccessToUserRemediationInstructionsFilter",
          :has-access-to-list-filters="hasAccessToListRemediationInstructionsFilters"
        )
      v-flex
        v-chip.primary.white--text(
          v-if="activeRange",
          close,
          label,
          @input="removeHistoryFilter"
        ) {{ $t(`quickRanges.types.${activeRange.value}`) }}
        c-action-btn(
          :tooltip="$t('alarm.liveReporting')",
          :color="activeRange ? 'primary' : ''",
          icon="schedule",
          @click="showEditLiveReportModal"
        )
      v-flex(v-if="hasAccessToExportAsCsv")
        c-action-btn(
          :loading="downloading",
          :tooltip="$t('settings.exportAsCsv')",
          icon="cloud_download",
          @click="exportAlarmsList"
        )
    alarms-list-table.mt-1(
      ref="alarmsTable",
      :widget="widget",
      :alarms="alarms",
      :total-items="alarmsMeta.total_count",
      :pagination.sync="pagination",
      :loading="alarmsPending",
      :hide-children="!query.correlation",
      :columns="widget.parameters.widgetColumns",
      :sticky-header="widget.parameters.sticky_header",
      :dense="dense",
      :refresh-alarms-list="fetchList",
      :selected-tag="query.tag",
      :search="query.search",
      :selectable="!hideMassSelection",
      :hide-actions="hideActions",
      :resizable-column="resizableColumn",
      :draggable-column="draggableColumn",
      :cells-content-behavior="cellsContentBehavior",
      :columns-settings="columnsSettings",
      expandable,
      densable,
      @select:tag="selectTag",
      @update:dense="updateDense",
      @update:page="updateQueryPage",
      @update:rows-per-page="updateRecordsPerPage",
      @update:columns-settings="updateColumnsSettings",
      @clear:tag="clearTag"
    )
</template>

<script>
import { omit, pick, isObject, isEqual } from 'lodash';

import { MODALS, USERS_PERMISSIONS, LIVE_REPORTING_QUICK_RANGES } from '@/constants';

import { findQuickRangeValue } from '@/helpers/date/date-intervals';
import { getAlarmListExportDownloadFileUrl } from '@/helpers/entities/alarm/url';
import { setSeveralFields } from '@/helpers/immutable';
import { getPageForNewRecordsPerPage } from '@/helpers/pagination';

import { authMixin } from '@/mixins/auth';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { exportMixinCreator } from '@/mixins/widget/export';
import { widgetSearchMixin } from '@/mixins/widget/search';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetAlarmsSocketMixin } from '@/mixins/widget/alarms-socket';
import { widgetRemediationInstructionsFilterMixin } from '@/mixins/widget/remediation-instructions-filter-select';
import { entitiesAlarmMixin } from '@/mixins/entities/alarm';
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';
import { entitiesAlarmDetailsMixin } from '@/mixins/entities/alarm/details';
import { permissionsWidgetsAlarmsListCorrelation } from '@/mixins/permissions/widgets/alarms-list/correlation';
import { permissionsWidgetsAlarmsListBookmark } from '@/mixins/permissions/widgets/alarms-list/bookmark';
import { permissionsWidgetsAlarmsListCategory } from '@/mixins/permissions/widgets/alarms-list/category';
import { permissionsWidgetsAlarmsListFilters } from '@/mixins/permissions/widgets/alarms-list/filters';
import {
  permissionsWidgetsAlarmsListRemediationInstructionsFilters,
} from '@/mixins/permissions/widgets/alarms-list/remediation-instructions-filters';
import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

import FilterSelector from '@/components/other/filter/partials/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/partials/filters-list-btn.vue';

import AlarmsListTable from './partials/alarms-list-table.vue';
import AlarmsListRemediationInstructionsFilters from './partials/alarms-list-remediation-instructions-filters.vue';

/**
 * Alarm-list component
 *
 * @module alarm
 *
 * @prop {Object} widget - Object representing the widget
 *
 * @event openSettings#click
 */
export default {
  components: {
    FilterSelector,
    FiltersListBtn,
    AlarmsListTable,
    AlarmsListRemediationInstructionsFilters,
  },
  mixins: [
    authMixin,
    widgetFetchQueryMixin,
    widgetSearchMixin,
    widgetFilterSelectMixin,
    widgetPeriodicRefreshMixin,
    widgetAlarmsSocketMixin,
    widgetRemediationInstructionsFilterMixin,
    entitiesWidgetMixin,
    entitiesAlarmMixin,
    entitiesAlarmTagMixin,
    entitiesAlarmDetailsMixin,
    permissionsWidgetsAlarmsListCategory,
    permissionsWidgetsAlarmsListCorrelation,
    permissionsWidgetsAlarmsListBookmark,
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
    tabId: {
      type: String,
      default: '',
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
  },
  data() {
    return {
      downloading: false,
    };
  },
  computed: {
    activeRange() {
      const { tstart, tstop } = this.query;

      if (tstart || tstop) {
        return findQuickRangeValue(tstart, tstop, LIVE_REPORTING_QUICK_RANGES, LIVE_REPORTING_QUICK_RANGES.custom);
      }

      return null;
    },

    hasAccessToExportAsCsv() {
      return this.checkAccess(USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv);
    },

    dense() {
      return this.userPreference.content.dense ?? this.widget.parameters.dense;
    },

    columnsSettings() {
      return this.userPreference.content.columns_settings;
    },

    resizableColumn() {
      return !!this.widget.parameters?.columns?.resizable;
    },

    cellsContentBehavior() {
      return this.widget.parameters?.columns?.cells_content_behavior;
    },

    draggableColumn() {
      return !!this.widget.parameters?.columns?.draggable;
    },
  },
  created() {
    this.actualizeUsedProperties();
  },
  methods: {
    refreshExpanded() {
      if (this.$refs.alarmsTable?.expanded) {
        Object.entries(this.$refs.alarmsTable.expanded).forEach(([id, expanded]) => {
          if (expanded && !this.alarms.some(alarm => alarm._id === id)) {
            this.$set(this.$refs.alarmsTable.expanded, id, false);
          }
        });
      }
    },

    selectTag(tag) {
      this.query = {
        ...this.query,

        page: 1,
        tag,
      };
    },

    clearTag() {
      const newQuery = omit(this.query, ['tag']);

      newQuery.page = 1;

      this.query = newQuery;
    },

    updateColumnsSettings(columnsSettings) {
      this.updateContentInUserPreference({ columns_settings: columnsSettings });
    },

    updateCorrelation(correlation) {
      this.updateContentInUserPreference({ isCorrelationEnabled: correlation });

      this.query = {
        ...this.query,

        page: 1,
        correlation,
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
      const categoryId = category && category._id;

      this.updateContentInUserPreference({ category: categoryId });

      this.query = {
        ...this.query,

        page: 1,
        category: categoryId,
      };
    },

    updateRecordsPerPage(limit) {
      this.updateContentInUserPreference({ itemsPerPage: limit });

      this.query = {
        ...this.query,

        limit,
        page: getPageForNewRecordsPerPage(limit, this.query.limit, this.query.page),
      };
    },

    updateDense(dense) {
      this.updateContentInUserPreference({ dense });
    },

    removeHistoryFilter() {
      const newQuery = omit(this.query, ['tstart', 'tstop']);

      newQuery.page = 1;

      this.query = newQuery;
    },

    showEditLiveReportModal() {
      this.$modals.show({
        name: MODALS.editLiveReporting,
        config: {
          ...pick(this.query, ['tstart', 'tstop', 'time_field']),
          action: params => this.query = {
            ...this.query,
            ...params,

            page: 1,
          },
        },
      });
    },

    async fetchList() {
      if (this.widget.parameters.widgetColumns.length) {
        const params = this.getQuery();

        this.fetchAlarmsDetailsList({ widgetId: this.widget._id });

        if (!this.alarmTagsPending) {
          this.fetchAlarmTagsList({ params: { paginate: false } });
        }

        if (!this.alarmsPending || !isEqual(params, this.alarmsFetchingParams)) {
          await this.fetchAlarmsList({
            widgetId: this.widget._id,
            params,
          });

          this.refreshExpanded();
        }
      }
    },

    getExportQueryColumns() {
      const {
        widgetExportColumns,
        widgetColumns,
      } = this.widget.parameters;

      const hasExportColumns = !!widgetExportColumns?.length;
      const columns = hasExportColumns ? widgetExportColumns : widgetColumns;

      return columns.map(({ value, text, template }) => ({
        name: value,
        label: text,
        template: hasExportColumns ? template : undefined,
      }));
    },

    getExportQuery() {
      const query = this.getQuery();
      const {
        exportCsvSeparator,
        exportCsvDatetimeFormat,
      } = this.widget.parameters;

      return {
        ...pick(query, [
          'search',
          'category',
          'correlation',
          'opened',
          'tstart',
          'tstop',
          'only_bookmarks',
        ]),

        fields: this.getExportQueryColumns(),
        filters: query.filters,
        separator: exportCsvSeparator,
        /**
         * @link https://git.canopsis.net/canopsis/canopsis-pro/-/issues/3997
         */
        time_format: isObject(exportCsvDatetimeFormat)
          ? exportCsvDatetimeFormat.value
          : exportCsvDatetimeFormat,
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
      } catch (err) {
        this.$popups.error({ text: this.$t('alarm.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },

    actualizeUsedProperties() {
      const unwatch = this.$watch(() => this.query.active_columns, (activeColumns) => {
        if (!isEqual(activeColumns, this.widget.parameters.usedAlarmProperties)) {
          this.updateWidget({
            id: this.widget._id,
            data: setSeveralFields(this.widget, {
              'parameters.usedAlarmProperties': activeColumns,
            }),
          });
        }

        unwatch();
      });
    },
  },
};
</script>
