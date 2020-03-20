<template lang="pug">
  div(data-test="tableWidget")
    v-layout.white(row, wrap, justify-space-between, align-center)
      v-flex
        alarm-list-search(:query.sync="query", :columns="columns")
      v-flex
        pagination(
          data-test="topPagination",
          v-if="hasColumns",
          :page="query.page",
          :limit="query.limit",
          :total="alarmsMeta.total",
          type="top",
          @input="updateQueryPage"
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
    alarms-list-table(
      :widget="widget",
      :alarms="alarms",
      :totalItems="alarmsMeta.total",
      :pagination.sync="vDataTablePagination",
      :loading="alarmsPending",
      :isEditingMode="isEditingMode",
      :isTourEnabled="isTourEnabled",
      :hasColumns="hasColumns",
      :columns="columns",
      ref="alarmsTable"
    )
      v-layout.white(v-show="alarmsMeta.total", align-center)
        v-flex(xs10)
          pagination(
            data-test="bottomPagination",
            :page="query.page",
            :limit="query.limit",
            :total="alarmsMeta.total",
            @input="updateQueryPage"
          )
        v-spacer
        v-flex(xs2, data-test="itemsPerPage")
          records-per-page(:value="query.limit", @input="updateRecordsPerPage")
    alarms-expand-panel-tour(v-if="isTourEnabled", :callbacks="tourCallbacks")
</template>

<script>
import { omit, pick, isEmpty } from 'lodash';

import { MODALS, USERS_RIGHTS, TOURS } from '@/constants';

import { findRange } from '@/helpers/date-intervals';

import RecordsPerPage from '@/components/tables/records-per-page.vue';
import FilterSelector from '@/components/other/filter/selector/filter-selector.vue';

import authMixin from '@/mixins/auth';
import widgetQueryMixin from '@/mixins/widget/query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetPaginationMixin from '@/mixins/widget/pagination';
import widgetFilterSelectMixin from '@/mixins/widget/filter-select';
import widgetRecordsPerPageMixin from '@/mixins/widget/records-per-page';
import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import alarmColumnFilters from '@/mixins/entities/alarm-column-filters';

import AlarmListSearch from './search/alarm-list-search.vue';
import AlarmsExpandPanelTour from './partials/alarms-expand-panel-tour.vue';
import AlarmsListTable from './partials/alarms-list-table.vue';

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
    AlarmListSearch,
    RecordsPerPage,
    AlarmsListTable,
    AlarmsExpandPanelTour,
    FilterSelector,
  },
  mixins: [
    authMixin,
    widgetQueryMixin,
    widgetColumnsMixin,
    widgetPaginationMixin,
    widgetFilterSelectMixin,
    widgetRecordsPerPageMixin,
    widgetPeriodicRefreshMixin,
    entitiesAlarmMixin,
    alarmColumnFilters,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
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
        onNextStep: this.onTourNextStep,
      };
    },

    isTourEnabled() {
      return this.checkIsTourEnabled(TOURS.alarmsExpandPanel) && this.alarms.length;
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
  },
  methods: {
    onTourNextStep(currentStep) {
      if (currentStep === 0) {
        this.$set(this.$refs.alarmsTable.expanded, this.alarms[0]._id, true);
      }

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
