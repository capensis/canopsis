<template lang="pug">
  div(data-test="tableWidget")
    v-layout.white(row, wrap, justify-space-between, align-center)
      v-flex
        advanced-search(
          :query.sync="query",
          :columns="columns",
          :tooltip="$t('search.alarmAdvancedSearch')"
        )
      v-flex(v-if="hasAccessToCorrelationSwitcher")
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
          :lockedFilters="widgetViewFilters",
          :value="mainFilter",
          :condition="mainFilterCondition",
          :hasAccessToEditFilter="hasAccessToEditFilter",
          :hasAccessToUserFilter="hasAccessToUserFilter",
          :hasAccessToListFilter="hasAccessToListFilter",
          @input="updateSelectedFilter",
          @update:condition="updateSelectedCondition",
          @update:filters="updateFilters"
        )
      v-flex
        alarms-list-remediation-instructions-filters(
          :filters="remediationInstructionsFilters",
          @input="updateRemediationInstructionsFilters"
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
      pagination(
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
      :totalItems="alarmsMeta.total_count",
      :pagination.sync="vDataTablePagination",
      :loading="alarmsPending",
      :isTourEnabled="isTourEnabled",
      :hideGroups="!query.correlation",
      :hasColumns="hasColumns",
      :columns="columns",
      selectable,
      expandable,
      ref="alarmsTable"
    )
      v-layout.white(v-show="alarmsMeta.total_count", align-center)
        v-flex(xs10)
          pagination(
            data-test="bottomPagination",
            :page="query.page",
            :limit="query.limit",
            :total="alarmsMeta.total_count",
            @input="updateQueryPage"
          )
        v-spacer
        v-flex(xs2, data-test="itemsPerPage")
          records-per-page.py-4(:value="query.limit", @input="updateRecordsPerPage")
    alarms-expand-panel-tour(v-if="isTourEnabled", :callbacks="tourCallbacks")
</template>

<script>
import { omit, pick, isEmpty } from 'lodash';

import { MODALS, USERS_RIGHTS, TOURS } from '@/constants';

import { findRange } from '@/helpers/date-intervals';

import RecordsPerPage from '@/components/tables/records-per-page.vue';
import FilterSelector from '@/components/other/filter/selector/filter-selector.vue';
import AdvancedSearch from '@/components/other/shared/search/advanced-search.vue';

import authMixin from '@/mixins/auth';
import widgetFetchQueryMixin from '@/mixins/widget/fetch-query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetPaginationMixin from '@/mixins/widget/pagination';
import widgetFilterSelectMixin from '@/mixins/widget/filter-select';
import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import alarmColumnFilters from '@/mixins/entities/alarm-column-filters';

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
    RecordsPerPage,
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
    widgetPaginationMixin,
    widgetFilterSelectMixin,
    widgetPeriodicRefreshMixin,
    entitiesAlarmMixin,
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
    remediationInstructionsFilters() {
      return this.userPreference.widget_preferences.remediationInstructionsFilters || [];
    },

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

    hasAccessToListFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.editFilter);
    },

    hasAccessToUserFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.userFilter);
    },

    hasAccessToCorrelationSwitcher() {
      return this.checkAccess(USERS_RIGHTS.business.alarmsList.actions.correlation);
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
    updateRemediationInstructionsFilters(remediationInstructionsFilters = []) {
      this.updateWidgetPreferencesInUserPreference({
        ...this.userPreference.widget_preferences,

        remediationInstructionsFilters,
      });
    },

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
