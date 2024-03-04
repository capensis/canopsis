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
          :tooltip="$t('alarm.advancedSearch')"
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
            :label="$t('settings.selectAFilter')",
            :filters="userPreference.filters",
            :locked-filters="widget.filters",
            :locked-value="lockedFilter",
            :value="mainFilter",
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
      :is-tour-enabled="isTourEnabled",
      :hide-children="!query.correlation",
      :columns="widget.parameters.widgetColumns",
      :sticky-header="widget.parameters.sticky_header",
      :dense="dense",
      :refresh-alarms-list="fetchList",
      :selected-tag="query.tag",
      :search="query.search",
      :selectable="!hideMassSelection",
      :hide-actions="hideActions",
      expandable,
      densable,
      @select:tag="selectTag",
      @update:dense="updateDense",
      @update:page="updateQueryPage",
      @update:rows-per-page="updateRecordsPerPage",
      @clear:tag="clearTag"
    )
    alarms-expand-panel-tour(v-if="isTourEnabled", :callbacks="tourCallbacks")
</template>

<script>
import { omit, pick, isObject, isEqual } from 'lodash';

import { API_HOST, API_ROUTES } from '@/config';

import { MODALS, TOURS, USERS_PERMISSIONS } from '@/constants';

import { findQuickRangeValue } from '@/helpers/date/date-intervals';
import { getPageForNewRecordsPerPage } from '@/helpers/pagination';

import { authMixin } from '@/mixins/auth';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import { exportMixinCreator } from '@/mixins/widget/export';
import { widgetFilterSelectMixin } from '@/mixins/widget/filter-select';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import { widgetRemediationInstructionsFilterMixin } from '@/mixins/widget/remediation-instructions-filter-select';
import { entitiesAlarmMixin } from '@/mixins/entities/alarm';
import { entitiesAlarmTagMixin } from '@/mixins/entities/alarm-tag';
import { entitiesAlarmDetailsMixin } from '@/mixins/entities/alarm/details';
import { permissionsWidgetsAlarmsListCorrelation } from '@/mixins/permissions/widgets/alarms-list/correlation';
import { permissionsWidgetsAlarmsListCategory } from '@/mixins/permissions/widgets/alarms-list/category';
import { permissionsWidgetsAlarmsListFilters } from '@/mixins/permissions/widgets/alarms-list/filters';
import { permissionsWidgetsAlarmsListRemediationInstructionsFilters }
from '@/mixins/permissions/widgets/alarms-list/remediation-instructions-filters';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import FiltersListBtn from '@/components/other/filter/filters-list-btn.vue';

import AlarmsListTable from './partials/alarms-list-table.vue';
import AlarmsExpandPanelTour from './expand-panel/alarms-expand-panel-tour.vue';
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
    AlarmsExpandPanelTour,
    AlarmsListRemediationInstructionsFilters,
  },
  mixins: [
    authMixin,
    widgetFetchQueryMixin,
    widgetFilterSelectMixin,
    widgetPeriodicRefreshMixin,
    widgetRemediationInstructionsFilterMixin,
    entitiesAlarmMixin,
    entitiesAlarmTagMixin,
    entitiesAlarmDetailsMixin,
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
    tourCallbacks() {
      return {
        onPreviousStep: this.onTourPreviousStep,
        onNextStep: this.onTourNextStep,
      };
    },

    isTourEnabled() {
      return this.checkIsTourEnabled(TOURS.alarmsExpandPanel)
        && !!this.alarms.length;
    },

    activeRange() {
      const { tstart, tstop } = this.query;

      if (tstart || tstop) {
        return findQuickRangeValue(tstart, tstop);
      }

      return null;
    },

    firstAlarmExpanded() {
      const [alarm] = this.alarms;

      return alarm && this.$refs.alarmsTable.expanded[alarm._id];
    },

    hasAccessToExportAsCsv() {
      return this.checkAccess(USERS_PERMISSIONS.business.alarmsList.actions.exportAsCsv);
    },

    dense() {
      return this.userPreference.content.dense ?? this.widget.parameters.dense;
    },
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

    updateCorrelation(correlation) {
      this.updateContentInUserPreference({
        isCorrelationEnabled: correlation,
      });

      this.query = {
        ...this.query,

        page: 1,
        correlation,
      };
    },

    updateCategory(category) {
      const categoryId = category && category._id;

      this.updateContentInUserPreference({
        category: categoryId,
      });

      this.query = {
        ...this.query,

        page: 1,
        category: categoryId,
      };
    },

    updateRecordsPerPage(limit) {
      this.updateContentInUserPreference({
        itemsPerPage: limit,
      });

      this.query = {
        ...this.query,

        limit,
        page: getPageForNewRecordsPerPage(limit, this.query.limit, this.query.page),
      };
    },

    updateDense(dense) {
      this.updateContentInUserPreference({
        dense,
      });
    },

    expandFirstAlarm() {
      if (!this.firstAlarmExpanded) {
        this.$set(this.$refs.alarmsTable.expanded, this.alarms[0]._id, true);
      }
    },

    onTourPreviousStep(currentStep) {
      if (currentStep !== 1) {
        this.expandFirstAlarm();
      }

      return this.$nextTick();
    },

    onTourNextStep() {
      this.expandFirstAlarm();

      return this.$nextTick();
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

    getExportQuery() {
      const query = this.getQuery();
      const {
        widgetExportColumns,
        widgetColumns,
        exportCsvSeparator,
        exportCsvDatetimeFormat,
      } = this.widget.parameters;
      const columns = widgetExportColumns?.length ? widgetExportColumns : widgetColumns;

      return {
        ...pick(query, ['search', 'category', 'correlation', 'opened', 'tstart', 'tstop']),

        fields: columns.map(({ value, text }) => ({ name: value, label: text })),
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

        this.downloadFile(`${API_HOST}${API_ROUTES.alarmListExport}/${fileData._id}/download`);
      } catch (err) {
        this.$popups.error({ text: this.$t('alarm.popups.exportFailed') });
      } finally {
        this.downloading = false;
      }
    },
  },
};
</script>
