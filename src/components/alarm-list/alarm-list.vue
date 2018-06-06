<template lang="pug">
  div
    v-layout(justify-end, align-center)
      v-chip(
        v-if="liveReportingFilter"
        @input="handleRemoveLiveReportingFilter",
        close,
        label,
        color="blue darken-4 white--text"
      ) {{ liveReportingFilter.text }}
      v-btn(@click="handleLiveReportingClick", icon, small)
        v-icon(:color="liveReportingFilter ? 'blue' : 'black'") schedule
    v-layout
      v-flex(xs5)
        alarm-list-searching
      v-btn(icon, @click="openSettingsPanel")
        v-icon settings
    div(v-if="!pending")
      basic-list(:items="items")
        tr.container(slot="header")
          th.box(v-for="column in alarmProperties")
            span {{ column.text }}
            list-sorting(:column="column.value")
            th.box
        tr.container(slot="row" slot-scope="item")
            td.box(v-for="property in alarmProperties") {{ item.props | get(property.value) }}
            td.box
              actions-panel.actions(:item="item.props")
        tr.container(slot="expandedRow", slot-scope="item")
            time-line(:alarmProps="item.props")
      pagination(:meta="meta", :limit="limit")
    loader(v-else)
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import BasicList from '@/components/basic-component/basic-list.vue';
import ActionsPanel from '@/components/basic-component/actions-panel.vue';
import Loader from '@/components/loaders/alarm-list-loader.vue';
import AlarmListSearching from '@/components/alarm-list/alarm-list-searching.vue';
import TimeLine from '@/components/alarm-list/time-line.vue';
import ListSorting from '@/components/basic-component/list-sorting.vue';
import PaginationMixin from '@/mixins/pagination';

const { mapActions: alarmsListMapActions, mapGetters: alarmsListMapGetters } = createNamespacedHelpers('alarmsList');
const { mapActions: alarmMapActions, mapGetters: alarmMapGetters } = createNamespacedHelpers('alarm');
const { mapActions: settingsMapActions } = createNamespacedHelpers('alarmsListSettings');
const { mapActions: modalsMapActions } = createNamespacedHelpers('modal');

/**
 * Alarm-list component.
 *
 * @module components/alarm-list
 * @param {object} alarmProperties - Object that describe the columns names and the alarms attributes corresponding
 *            e.g : { ColumnName : 'att1.att2', Connector : 'v.connector' }
 * @param {integer} [itemsPerPage=5] - Number of Alarm to display per page
 */
export default {
  components: {
    ListSorting,
    TimeLine,
    AlarmListSearching,
    ActionsPanel,
    BasicList,
    Loader,
  },
  mixins: [PaginationMixin],
  props: {
    alarmProperties: {
      type: Array,
      default: () => ([]),
    },
  },
  computed: {
    ...alarmMapGetters([
      'items',
      'meta',
      'pending',
    ]),
    ...alarmsListMapGetters(['liveReportingFilter']),
  },
  methods: {
    ...alarmMapActions({
      fetchListAction: 'fetchList',
    }),

    ...settingsMapActions({
      openSettingsPanel: 'openPanel',
    }),

    ...modalsMapActions({
      showModal: 'show',
    }),

    ...alarmsListMapActions(['removeLiveReportingFilter']),

    handleLiveReportingClick() {
      this.showModal({ name: 'edit-live-reporting' });
    },

    handleRemoveLiveReportingFilter() {
      this.removeLiveReportingFilter();
      this.fetchList();
    },
  },
};
</script>

<style scoped>
  th {
    overflow: hidden;
    text-overflow: ellipsis;
  }
  td {
    overflow-wrap: break-word;
  }
  .container {
    display: flex;
  }
  .box{
    width: 10%;
    flex: 1;
  }
</style>
