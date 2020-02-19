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
        v-btn(data-test="alarmsDateInterval", @click="showEditLiveReportModal", icon, small)
          v-icon(:color="activeRange ? 'primary' : 'black'") schedule
    alarms-table(
      :widget="widget",
      :alarms="alarms",
      :totalItems="alarmsMeta.total",
      :pagination.sync="vDataTablePagination",
      :loading="alarmsPending",
      :isEditingMode="isEditingMode",
      :hasColumns="hasColumns",
      :columns="columns",
      ref="alarmsTable"
    )
    v-layout.white(align-center)
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
</template>

<script>
import { omit, pick, isEmpty } from 'lodash';

import { MODALS, USERS_RIGHTS } from '@/constants';

import { findRange } from '@/helpers/date-intervals';
import AlarmsTable from '@/components/other/alarm/alarms-table.vue';
import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm/actions/mass-actions-panel.vue';
import AlarmsExpandPanel from '@/components/other/alarm/partials/alarms-expand-panel.vue';
import AlarmListSearch from '@/components/other/alarm/search/alarm-list-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';
import NoColumnsTable from '@/components/tables/no-columns.vue';
import FilterSelector from '@/components/other/filter/selector/filter-selector.vue';
import AlarmsListRow from '@/components/other/alarm/partials/alarms-list-row.vue';

import authMixin from '@/mixins/auth';
import widgetQueryMixin from '@/mixins/widget/query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetPaginationMixin from '@/mixins/widget/pagination';
import widgetFilterSelectMixin from '@/mixins/widget/filter-select';
import widgetRecordsPerPageMixin from '@/mixins/widget/records-per-page';
import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesAlarmMixin from '@/mixins/entities/alarm';

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
    AlarmsTable,
    AlarmListSearch,
    AlarmsListRow,
    RecordsPerPage,
    AlarmsExpandPanel,
    MassActionsPanel,
    ActionsPanel,
    AlarmColumnValue,
    NoColumnsTable,
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
