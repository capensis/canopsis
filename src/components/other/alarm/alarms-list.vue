<template lang="pug">
  div
    v-layout.white(row, wrap, justify-space-between, align-center)
      v-flex
        alarm-list-search(:query.sync="query")
      v-flex
        pagination(v-if="hasColumns", :meta="alarmsMeta", :query.sync="query", type="top")
      v-flex(v-if="hasAccessToListFilters")
        filter-selector(
        :label="$t('settings.selectAFilter')",
        :items="viewFilters",
        :value="mainFilter",
        :condition="mainFilterCondition",
        @update:condition="updateSelectedCondition",
        @input="updateSelectedFilter"
        )
      v-flex
        v-chip.primary.white--text(
        v-if="query.interval",
        @input="removeHistoryFilter",
        close,
        label,
        ) {{ $t(`modals.liveReporting.${query.interval}`) }}
        v-btn(@click="showEditLiveReportModal", icon, small)
          v-icon(:color="query.interval ? 'primary' : 'black'") schedule
      v-flex.px-3(v-show="selected.length", xs12)
        mass-actions-panel(:itemsIds="selectedIds")
    no-columns-table(v-if="!hasColumns")
    div(v-else)
      v-data-table(
      v-model="selected",
      :items="alarms",
      :headers="headers",
      :total-items="alarmsMeta.total",
      :pagination.sync="vDataTablePagination",
      :loading="alarmsPending",
      ref="dataTable",
      item-key="_id",
      hide-actions,
      select-all,
      expand
      )
        template(slot="progress")
          v-fade-transition
            v-progress-linear(height="2", indeterminate, color="primary")
        template(slot="headerCell", slot-scope="props")
          span {{ props.header.text }}
        template(slot="items", slot-scope="props")
          tr
            td
              v-checkbox(primary, hide-details, v-model="props.selected")
            td(
            v-for="column in columns",
            @click="props.expanded = !props.expanded"
            )
              alarm-column-value(:alarm="props.item", :column="column", :widget="widget")
            td
              actions-panel(:item="props.item", :widget="widget")
        template(slot="expand", slot-scope="props")
          time-line(:alarmProps="props.item")
      v-layout.white(align-center)
        v-flex(xs10)
          pagination(:meta="alarmsMeta", :query.sync="query")
        v-spacer
        v-flex(xs2)
          records-per-page(:query.sync="query")
</template>

<script>
import omit from 'lodash/omit';
import pick from 'lodash/pick';
import isEmpty from 'lodash/isEmpty';

import { MODALS, USERS_RIGHTS } from '@/constants';

import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm/actions/mass-actions-panel.vue';
import TimeLine from '@/components/other/alarm/timeline/time-line.vue';
import AlarmListSearch from '@/components/other/alarm/search/alarm-list-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';
import NoColumnsTable from '@/components/tables/no-columns.vue';
import FilterSelector from '@/components/other/filter/selector/filter-selector.vue';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal/modal';
import widgetQueryMixin from '@/mixins/widget/query';
import widgetColumnsMixin from '@/mixins/widget/columns';
import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import filterSelectMixin from '@/mixins/filter-select';

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
    TimeLine,
    MassActionsPanel,
    ActionsPanel,
    AlarmColumnValue,
    NoColumnsTable,
    FilterSelector,
  },
  mixins: [
    authMixin,
    modalMixin,
    widgetQueryMixin,
    widgetColumnsMixin,
    widgetPeriodicRefreshMixin,
    entitiesAlarmMixin,
    filterSelectMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      selected: [],
    };
  },
  computed: {
    selectedIds() {
      return this.selected.map(item => item._id);
    },

    headers() {
      if (this.hasColumns) {
        return [...this.columns, { text: '', sortable: false }];
      }

      return [];
    },

    hasAccessToListFilters() {
      return this.checkAccess(USERS_RIGHTS.business.alarmList.actions.listFilters);
    },

    hasAccessToEditFilter() {
      return this.checkAccess(USERS_RIGHTS.business.alarmList.actions.editFilter);
    },
  },
  methods: {
    removeHistoryFilter() {
      this.query = omit(this.query, ['interval', 'tstart', 'tstop']);
    },

    showEditLiveReportModal() {
      this.showModal({
        name: MODALS.editLiveReporting,
        config: {
          ...pick(this.query, ['interval', 'tstart', 'tstop']),
          action: params => this.query = { ...this.query, ...params },
        },
      });
    },

    fetchList({ isPeriodicRefresh } = {}) {
      if (this.hasColumns) {
        const query = this.getQuery();

        if (isPeriodicRefresh && !isEmpty(this.$refs.dataTable.expanded)) {
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
