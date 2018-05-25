<template lang="pug">
  div
    div(v-if="!pending")
      basic-list(:items="items")
        tr.container(slot="header")
            th.box(v-for="columnName in Object.keys(alarmProperty)") {{ columnName }}
            th.box
        tr.container(slot="row", slot-scope="item")
            td.box(v-for="property in Object.values(alarmProperty)") {{ getProp(item.props, property) }}
            td.box
              actions-panel.actions
        tr.container(slot="expandedRow", slot-scope="item")
            time-line(:alarmId="item.props.d")
      alarm-list-pagination(:meta="meta", :limit="limit")
    loader(v-else)
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import getProp from 'lodash/get';

import { PAGINATION_LIMIT } from '@/config';
import getQuery from '@/helpers/pagination';

import AlarmListPagination from '@/components/AlarmList/alarm-list-pagination.vue';
import BasicList from '@/components/BasicComponent/basic-list.vue';
import ActionsPanel from '@/components/BasicComponent/actions-panel.vue';
import Loader from '@/components/loaders/alarm-list-loader.vue';
import TimeLine from '@/components/AlarmList/time-line.vue';


const { mapActions, mapGetters } = createNamespacedHelpers('alarm');

/**
 * Alarm-list component.
 *
 * @module components/alarm-list
 * @param {object} alarmProperty - Object that describe the columns names and the alarms attributes corresponding
 *            e.g : { ColumnName : 'att1.att2', Connector : 'v.connector' }
 * @param {integer} [itemsPerPage=5] - Number of Alarm to display per page
 */
export default {
  name: 'AlarmList',
  components: {
    TimeLine,
    AlarmListPagination,
    ActionsPanel,
    BasicList,
    Loader,
  },
  props: {
    alarmProperty: {
      type: Object,
      default() {
        return {};
      },
    },
    limit: {
      type: Number,
      default: PAGINATION_LIMIT,
    },
  },
  computed: {
    ...mapGetters([
      'items',
      'meta',
      'pending',
    ]),
  },
  watch: {
    $route() {
      this.fetchList();
    },
  },
  mounted() {
    this.fetchList(this.fetchingParams);
  },
  methods: {
    getProp,
    getQuery,
    ...mapActions({
      fetchListAction: 'fetchList',
    }),
    fetchList() {
      this.fetchListAction({
        params: this.getQuery(),
      });
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
