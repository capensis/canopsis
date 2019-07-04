<template lang="pug">
  div
    v-layout.white(row, wrap, justify-space-between, align-center)
      v-flex
        alarm-list-search(:query.sync="query", :columns="columns")
      v-flex
        pagination(
        v-if="hasColumns",
        :page="query.page",
        :limit="query.limit",
        :total="alarmsMeta.total",
        type="top",
        @input="updateQueryPage"
        )
      v-flex
        filter-selector(
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
        v-if="activeRange",
        close,
        label,
        @input="removeHistoryFilter"
        ) {{ $t(`settings.statsDateInterval.quickRanges.${activeRange.value}`) }}
        v-btn(@click="showEditLiveReportModal", icon, small)
          v-icon(:color="activeRange ? 'primary' : 'black'") schedule
      v-flex.px-3(v-show="selected.length", xs12)
        mass-actions-panel(:itemsIds="selectedIds", :widget="widget")
    no-columns-table(v-if="!hasColumns")
    div(v-else)
      v-data-table.alarms-list-table(
      :class="vDataTableClass",
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
              v-checkbox-functional(v-model="props.selected", primary, hide-details)
            td(
            v-for="column in columns",
            @click="props.expanded = !props.expanded"
            )
              alarm-column-value(:alarm="props.item", :column="column", :widget="widget")
            td
              actions-panel(:item="props.item", :widget="widget", :isEditingMode="isEditingMode")
        template(slot="expand", slot-scope="props")
          time-line(:alarm="props.item", :isHTMLEnabled="widget.parameters.isHtmlEnabledOnTimeLine")
      v-layout.white(align-center)
        v-flex(xs10)
          pagination(
          :page="query.page",
          :limit="query.limit",
          :total="alarmsMeta.total",
          @input="updateQueryPage"
          )
        v-spacer
        v-flex(xs2)
          records-per-page(:value="query.limit", @input="updateRecordsPerPage")
</template>

<script>
import { omit, pick, isEmpty } from 'lodash';

import { MODALS, USERS_RIGHTS } from '@/constants';

import { findRange } from '@/helpers/date-intervals';

import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm/actions/mass-actions-panel.vue';
import TimeLine from '@/components/other/alarm/time-line/time-line.vue';
import AlarmListSearch from '@/components/other/alarm/search/alarm-list-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';
import NoColumnsTable from '@/components/tables/no-columns.vue';
import FilterSelector from '@/components/other/filter/selector/filter-selector.vue';

import authMixin from '@/mixins/auth';
import modalMixin from '@/mixins/modal';
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

    selectedIds() {
      return this.selected.map(item => item._id);
    },

    headers() {
      if (this.hasColumns) {
        return [...this.columns, { text: '', sortable: false }];
      }

      return [];
    },

    vDataTableClass() {
      const columnsLength = this.headers.length;
      const COLUMNS_SIZES_VALUES = {
        sm: { min: 0, max: 10, label: 'sm' },
        md: { min: 11, max: 12, label: 'md' },
        lg: { min: 13, max: Number.MAX_VALUE, label: 'lg' },
      };

      const { label = COLUMNS_SIZES_VALUES.sm.label } = Object.values(COLUMNS_SIZES_VALUES)
        .find(({ min, max }) => columnsLength >= min && columnsLength <= max);

      return {
        [`columns-${label}`]: true,
      };
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
      this.query = omit(this.query, ['interval', 'tstart', 'tstop']);
    },

    showEditLiveReportModal() {
      this.showModal({
        name: MODALS.editLiveReporting,
        config: {
          ...pick(this.query, ['interval', 'tstart', 'tstop']),
          action: params => this.query = { ...omit(this.query, ['tstart', 'tstop']), ...params },
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

<style lang="scss">
  .alarms-list-table {
    &.columns-lg {
      table.v-table {
        tbody, thead {
          td, th {
            padding: 0 8px;
          }

          @media screen and (max-width: 1600px) {
            td, th {
              padding: 0 4px;
            }
          }

          @media screen and (max-width: 1450px) {
            td, th {
              font-size: 0.85em;
            }

            .badge {
              font-size: inherit;
            }
          }
        }
      }
    }

    &.columns-md {
      table.v-table {
        tbody, thead {
          @media screen and (max-width: 1700px) {
            td, th {
              padding: 0 12px;
            }
          }

          @media screen and (max-width: 1250px) {
            td, th {
              padding: 0 8px;
            }
          }

          @media screen and (max-width: 1150px) {
            td, th {
              font-size: 0.85em;
              padding: 0 4px;
            }

            .badge {
              font-size: inherit;
            }
          }
        }
      }
    }
  }
</style>
