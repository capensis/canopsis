<template lang="pug">
  div(data-test="tableWidget")
    v-layout.white(row, wrap, justify-space-between, align-center)
      v-flex
        c-advanced-search(
          :query.sync="query",
          :columns="columns",
          :tooltip="$t('search.alarmAdvancedSearch')"
        )
      v-flex(v-if="hasAccessToCategory")
        c-entity-category-field.mr-3(:category="query.category", @input="updateCategory")
      v-flex(v-if="hasAccessToCorrelation")
        v-switch(
          :value="query.correlation",
          :label="$t('common.correlation')",
          color="primary",
          @change="updateCorrelation"
        )
      v-flex
        filter-selector(
          data-test="tableFilterSelector",
          :label="$t('settings.selectAFilter')",
          :filters="viewFilters",
          :locked-filters="widgetViewFilters",
          :value="mainFilter",
          :condition="mainFilterCondition",
          :has-access-to-edit-filter="hasAccessToEditFilter",
          :has-access-to-user-filter="hasAccessToUserFilter",
          :has-access-to-list-filters="hasAccessToListFilters",
          @input="updateSelectedFilter",
          @update:condition="updateSelectedCondition",
          @update:filters="updateFilters"
        )
      v-flex
        alarms-list-remediation-instructions-filters(
          :filters.sync="remediationInstructionsFilters",
          :locked-filters.sync="widgetRemediationInstructionsFilters",
          :has-access-to-edit-filter="hasAccessToEditRemediationInstructionsFilter",
          :has-access-to-user-filter="hasAccessToUserRemediationInstructionsFilter",
          :has-access-to-list-filters="hasAccessToListRemediationInstructionsFilters"
        )
      v-flex
        v-chip.primary.white--text(
          data-test="resetAlarmsDateInterval",
          v-if="activeRange",
          close,
          label,
          @input="removeHistoryFilter"
        ) {{ $t(`settings.statsDateInterval.quickRanges.${activeRange.value}`) }}
        c-action-btn(
          :tooltip="$t('liveReporting.button')",
          :color="activeRange ? 'primary' : 'black'",
          icon="schedule",
          @click="showEditLiveReportModal"
        )
      v-flex(v-if="hasAccessToExportAsCsv")
        c-action-btn(
          :loading="!!alarmsExportPending",
          :tooltip="$t('settings.exportAsCsv')",
          icon="cloud_download",
          color="black",
          @click="exportAlarmsList"
        )
    v-layout(row, wrap, align-center)
      c-pagination(
        data-test="topPagination",
        v-if="hasColumns",
        :page="query.page",
        :limit="query.limit",
        :total="alarmsMeta.total_count",
        type="top",
        @input="updateQueryPage"
      )
    alarms-list-table(
      ref="alarmsTable",
      :widget="widget",
      :alarms="alarms",
      :total-items="alarmsMeta.total_count",
      :pagination.sync="vDataTablePagination",
      :loading="alarmsPending",
      :is-tour-enabled="isTourEnabled",
      :hide-groups="!query.correlation",
      :has-columns="hasColumns",
      :columns="columns",
      selectable,
      expandable
    )
      c-table-pagination(
        :total-items="alarmsMeta.total_count",
        :rows-per-page="query.limit",
        :page="query.page",
        @update:page="updateQueryPage",
        @update:rows-per-page="updateRecordsPerPage"
      )
    alarms-expand-panel-tour(v-if="isTourEnabled", :callbacks="tourCallbacks")
</template>

<script>
import { omit, pick, isEmpty } from 'lodash';

import { MODALS, TOURS, USERS_PERMISSIONS } from '@/constants';

import { findRange } from '@/helpers/date/date-intervals';

import FilterSelector from '@/components/other/filter/filter-selector.vue';

import { authMixin } from '@/mixins/auth';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetExportMixinCreator from '@/mixins/widget/export';
import widgetFilterSelectMixin from '@/mixins/widget/filter-select';
import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import widgetRemediationInstructionsFilterMixin from '@/mixins/widget/remediation-instructions-filter-select';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import { permissionsWidgetsAlarmsListCorrelation } from '@/mixins/permissions/widgets/alarms-list/correlation';
import { permissionsWidgetsAlarmsListCategory } from '@/mixins/permissions/widgets/alarms-list/category';
import { permissionsWidgetsAlarmsListFilters } from '@/mixins/permissions/widgets/alarms-list/filters';
import { permissionsWidgetsAlarmsListRemediationInstructionsFilters }
  from '@/mixins/permissions/widgets/alarms-list/remediation-instructions-filters';

import AlarmsListTable from './partials/alarms-list-table.vue';
import AlarmsExpandPanelTour from './partials/alarms-expand-panel-tour.vue';
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
    AlarmsListTable,
    AlarmsExpandPanelTour,
    AlarmsListRemediationInstructionsFilters,
  },
  mixins: [
    authMixin,
    widgetFetchQueryMixin,
    widgetColumnsMixin,
    widgetFilterSelectMixin,
    widgetPeriodicRefreshMixin,
    widgetRemediationInstructionsFilterMixin,
    entitiesAlarmMixin,
    permissionsWidgetsAlarmsListCategory,
    permissionsWidgetsAlarmsListCorrelation,
    permissionsWidgetsAlarmsListFilters,
    permissionsWidgetsAlarmsListRemediationInstructionsFilters,
    widgetExportMixinCreator({
      createExport: 'createAlarmsListExport',
      fetchExport: 'fetchAlarmsListExport',
      fetchExportFile: 'fetchAlarmsListCsvFile',
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
  },
  data() {
    return {
      selected: [],
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
      return this.checkIsTourEnabled(TOURS.alarmsExpandPanel) && !!this.alarms.length;
    },

    activeRange() {
      const { tstart, tstop } = this.query;

      if (tstart || tstop) {
        return findRange(tstart, tstop);
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
  },
  methods: {
    updateCorrelation(correlation) {
      this.updateWidgetPreferencesInUserPreference({
        ...this.userPreference.widget_preferences,

        isCorrelationEnabled: correlation,
      });

      this.query = {
        ...this.query,

        correlation,
      };
    },

    updateCategory(category) {
      const categoryId = category && category._id;

      this.updateWidgetPreferencesInUserPreference({
        ...this.userPreference.widget_preferences,

        category: categoryId,
      });

      this.query = {
        ...this.query,

        category: categoryId,
      };
    },

    updateRecordsPerPage(limit) {
      this.updateWidgetPreferencesInUserPreference({
        ...this.userPreference.widget_preferences,

        itemsPerPage: limit,
      });

      this.query = {
        ...this.query,

        limit,
      };
    },

    expandFirstAlarm() {
      if (this.alarms[0] && !this.firstAlarmExpanded) {
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
      this.query = omit(this.query, ['tstart', 'tstop']);
    },

    showEditLiveReportModal() {
      this.$modals.show({
        name: MODALS.editLiveReporting,
        config: {
          ...pick(this.query, ['tstart', 'tstop']),
          action: params => this.query = { ...this.query, ...params },
        },
      });
    },

    fetchList({ isPeriodicRefresh } = {}) {
      if (this.hasColumns) {
        const query = this.getQuery();

        if (isPeriodicRefresh && !isEmpty(this.$refs.alarmsTable.expanded)) {
          query.with_steps = true;
        }

        this.fetchAlarmsList({
          widgetId: this.widget._id,
          params: query,
        });
      }
    },

    exportAlarmsList() {
      const query = this.getQuery();

      this.exportWidgetAsCsv({
        name: `${this.widget._id}-${new Date().toLocaleString()}`,
        params: {
          ...pick(query, ['search', 'filter', 'correlation', 'opened', 'resolved', 'active_columns']),

          separator: this.widget.parameters.exportCsvSeparator,
          time_format: this.widget.parameters.exportCsvDatetimeFormat,
        },
      });
    },
  },
};
</script>