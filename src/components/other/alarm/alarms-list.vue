<template lang="pug">
  v-container
    v-layout.white(wrap, justify-space-between, align-center)
      v-flex(xs12 md3)
        alarm-list-search(:query.sync="query")
      v-flex(xs2)
        pagination(:meta="alarmsMeta", :query.sync="query", type="top")
      v-flex.ml-4(xs3)
        mass-actions-panel(v-show="selected.length", :itemsIds="selectedIds")
      v-flex(xs3)
        v-chip(
        v-if="query.interval",
        @input="removeHistoryFilter",
        close,
        label,
        color="blue darken-4 white--text"
        ) {{ $t(`modals.liveReporting.${query.interval}`) }}
        v-btn(@click="showEditLiveReportModal", icon, small)
          v-icon(:color="query.interval ? 'blue' : 'black'") schedule
        v-btn(icon, @click="showSettings")
          v-icon settings
    div
      v-data-table(
      v-model="selected",
      :items="alarms",
      :headers="headers",
      :total-items="alarmsMeta.total",
      :pagination.sync="vDataTablePagination",
      :loading="alarmsPending",
      item-key="_id",
      select-all,
      hide-actions,
      )
        template(slot="headerCell", slot-scope="props")
          span {{ props.header.text }}
        template(slot="items", slot-scope="props")
          td
            v-checkbox(primary, hide-details, v-model="props.selected")
          td(
          v-for="prop in properties",
          @click="props.expanded = !props.expanded"
          )
            alarm-column-value(:alarm="props.item", :property="prop", :widget="widget")
          td
            actions-panel(:item="props.item", :widget="widget")
        template(slot="expand", slot-scope="props")
          time-line(:alarmProps="props.item", @click="props.expanded = !props.expanded")
      v-layout.white(align-center)
        v-flex(xs10)
          pagination(:meta="alarmsMeta", :query.sync="query")
        v-spacer
        v-flex(xs2)
          records-per-page(:query.sync="query")
</template>

<script>
import omit from 'lodash/omit';

import { MODALS, SIDE_BARS } from '@/constants';

import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm/actions/mass-actions-panel.vue';
import TimeLine from '@/components/other/alarm/timeline/time-line.vue';
import AlarmListSearch from '@/components/other/alarm/search/alarm-list-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';

import modalMixin from '@/mixins/modal/modal';
import sideBarMixin from '@/mixins/side-bar/side-bar';
import widgetQueryMixin from '@/mixins/widget/query';
import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

/**
 * Alarm-list component
 *
 * @module alarm
 *
 * @prop {Object} widget - Object representing the widget
 * @prop {Object} properties - Object that describe the columns names and the alarms attributes corresponding
 *            e.g : { ColumnName : 'att1.att2', Connector : 'v.connector' }
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
  },
  mixins: [
    modalMixin,
    sideBarMixin,
    widgetQueryMixin,
    widgetPeriodicRefreshMixin,
    entitiesAlarmMixin,
    entitiesUserPreferenceMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    properties: {
      type: Array,
      default: () => ([]),
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
      return [...this.properties, { text: '', sortable: false }];
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
          updateQuery: params => this.query = { ...this.query, ...params },
        },
      });
    },
    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.alarmSettings,
        config: {
          widget: this.widget,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  th {
    overflow: hidden;
    text-overflow: ellipsis;
  }

  td {
    overflow-wrap: break-word;
  }
  .fade-enter-active, .fade-leave-active {
    transition: opacity .5s;
  }
  .fade-enter, .fade-leave-to {
    opacity: 0;
  }
  .loader {
    top: 15%;
    position: absolute;
  }
</style>
