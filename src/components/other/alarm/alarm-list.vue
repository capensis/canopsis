<template lang="pug">
.white
  v-layout(justify-space-between, align-center)
    v-flex.ml-4(xs4)
      mass-actions-panel(v-show="selected.length", :itemsIds="selected")
    v-flex(xs2)
      v-chip(
        v-if="$route.query.interval",
        @input="removeHistoryFilter",
        close,
        label,
        color="blue darken-4 white--text"
      ) {{ $t(`modals.liveReporting.${$route.query.interval}`) }}
      v-btn(@click="showModal({ name: 'edit-live-reporting' })", icon, small)
        v-icon(:color="$route.query.interval ? 'blue' : 'black'") schedule
      v-btn(icon, @click="$emit('openSettings')")
        v-icon settings
  v-layout.my-2(wrap, justify-space-between, align-center)
    v-flex(xs12 md5)
      alarm-list-search
    v-flex(xs4)
      pagination(:meta="meta", :limit="limit", type="top")
  basic-list(:items="items", :pending="pending", :selected.sync="selected", expanded)
    loader(slot="loader")
    tr.container.header.pa-0(slot="header")
      th.box(v-for="column in alarmProperties")
        span {{ column.text }}
        list-sorting(:column="column.value", class="blue--text")
      th.box
    tr.container(slot="row" slot-scope="item")
        td.box(v-for="property in alarmProperties")
          alarm-column-value(:alarm="item.props", :property="property")
        td.box
          actions-panel.actions(:item="item.props")
    tr.container(slot="expandedRow", slot-scope="item")
      time-line(:alarmProps="item.props")
  v-layout(wrap)
    v-flex(xs12, md7)
    pagination(:meta="meta", :limit="limit", :first="first", :last="last")
    records-per-page
</template>

<script>
import omit from 'lodash/omit';

import BasicList from '@/components/tables/basic-list.vue';
import ListSorting from '@/components/tables/list-sorting.vue';
import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm/actions/mass-actions-panel.vue';
import TimeLine from '@/components/other/alarm/timeline/time-line.vue';
import Loader from '@/components/other/alarm/loader/alarm-list-loader.vue';
import AlarmListSearch from '@/components/other/alarm/search/alarm-list-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';
import FilterSelector from '@/components/other/filter/filter-selector.vue';
import modalMixin from '@/mixins/modal/modal';
import paginationMixin from '@/mixins/pagination';
import alarmsMixin from '@/mixins/alarms';

/**
 * Alarm-list component
 *
 * @module alarm
 *
 * @prop {object} alarmProperties - Object that describe the columns names and the alarms attributes corresponding
 *            e.g : { ColumnName : 'att1.att2', Connector : 'v.connector' }
 * @prop {integer} [itemsPerPage=5] - Number of Alarm to display per page
 *
 * @event openSettings#click
 */
export default {
  components: {
    AlarmListSearch,
    RecordsPerPage,
    ListSorting,
    TimeLine,
    MassActionsPanel,
    ActionsPanel,
    BasicList,
    Loader,
    AlarmColumnValue,
    FilterSelector,
  },
  mixins: [alarmsMixin, paginationMixin, modalMixin],
  props: {
    alarmProperties: {
      type: Array,
      default: () => ([]),
    },
  },
  data() {
    return {
      selected: [],
    };
  },
  methods: {
    removeHistoryFilter() {
      const query = omit(this.$route.query, ['interval', 'tstart', 'tstop']);
      this.$router.push({ query });
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

  .container {
    padding: 0;
    display: flex;
    align-items: center;
  }

  .header {
    border: 1px solid gray;
  }

  .box{
    flex: 1;
    padding: 1px;
  }
</style>
