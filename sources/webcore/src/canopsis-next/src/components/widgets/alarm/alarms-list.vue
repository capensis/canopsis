<template lang="pug">
  div(data-test="tableWidget")
    v-layout.white(row, wrap, justify-space-between, align-center)
      v-flex
        advanced-search(
          :query.sync="query",
          :columns="columns",
          :tooltip="$t('search.alarmAdvancedSearch')"
        )
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
        v-tooltip(bottom)
          v-btn(
            slot="activator",
            data-test="alarmsDateInterval",
            icon,
            small,
            @click="showEditLiveReportModal"
          )
            v-icon(:color="activeRange ? 'primary' : 'black'") schedule
          span {{ $t('liveReporting.button') }}
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
      expandable,
      ref="alarmsTable"
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

import { MODALS, TOURS } from '@/constants';

import { findRange } from '@/helpers/date/date-intervals';

import FilterSelector from '@/components/other/filter/filter-selector.vue';
import AdvancedSearch from '@/components/common/search/advanced-search.vue';

import authMixin from '@/mixins/auth';
import widgetFetchQueryMixin from '@/mixins/widget/fetch-query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetFilterSelectMixin from '@/mixins/widget/filter-select';
import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import widgetRemediationInstructionsFilterMixin from '@/mixins/widget/remediation-instructions-filter-select';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import alarmColumnFilters from '@/mixins/entities/alarm-column-filters';
import rightsWidgetsAlarmsListCorrelation from '@/mixins/rights/widgets/alarms-list/correlation';
import rightsWidgetsAlarmsListFilters from '@/mixins/rights/widgets/alarms-list/filters';
import rightsWidgetsAlarmsListRemediationInstructionsFilters
  from '@/mixins/rights/widgets/alarms-list/remediation-instructions-filters';

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
    AdvancedSearch,
    AlarmsListTable,
    AlarmsExpandPanelTour,
    AlarmsListRemediationInstructionsFilters,
  },
  mixins: [
    authMixin,
    alarmColumnFilters,
    widgetFetchQueryMixin,
    widgetColumnsMixin,
    widgetFilterSelectMixin,
    widgetPeriodicRefreshMixin,
    widgetRemediationInstructionsFilterMixin,
    entitiesAlarmMixin,
    rightsWidgetsAlarmsListCorrelation,
    rightsWidgetsAlarmsListFilters,
    rightsWidgetsAlarmsListRemediationInstructionsFilters,
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
  },
  mounted() {
    this.fetchAlarmColumnFilters();
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
  },
};
</script>
