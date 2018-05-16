<template lang="pug">
  div
    div(v-if="fetchComplete")
      basic-list(:items="items")
        tr.container(slot="header")
            th.box(v-for="columnName in Object.keys(alarmProperty)")
              span {{ columnName }}
              alarm-list-sorting(:columnToSort="alarmProperty[columnName]")
            th.box
        tr.container(slot="row" slot-scope="item")
            td.box(v-for="property in Object.values(alarmProperty)") {{ getProp(item.props, property) }}
            td.box
              actions-panel.actions
        tr.container(slot="expandedRow" slot-scope="item")
            td.box {{ item.props.infos }}
      alarm-list-pagination(:meta="meta", :limit="limit", v-if="fetchComplete")
    loader(v-else)
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import getProp from 'lodash/get';
import omit from 'lodash/omit';

import { PAGINATION_LIMIT } from '@/config';

import BasicList from '../BasicComponent/basic-list.vue';
import ActionsPanel from '../BasicComponent/actions-panel.vue';
import Loader from '../loaders/alarm-list-loader.vue';
import AlarmListPagination from './alarm-list-pagination.vue';
import AlarmListSorting from './alarm-list-sorting.vue';


const { mapGetters, mapActions } = createNamespacedHelpers('entities/alarm');

export default {
/**
 * The Alarm List Component
 *         Props :
 *          - alarmProperty ( Object ) : Object that describe the columns names and the alarms attributes corresponding
 *            e.g : { ColumnName : 'att1.att2', Connector : 'v.connector' }
 *            Default : {}
 *          - itemsPerPage ( Integer ) : Number of Alarm to display per page
 *            Default : 5
 */
  name: 'AlarmList',
  components: {
    AlarmListSorting,
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
  mounted() {
    this.fetchList(this.fetchingParams);
  },
  computed: {
    ...mapGetters([
      'items',
      'meta',
      'fetchComplete',
    ]),
  },
  methods: {
    getProp,
    ...mapActions({
      fetchListAction: 'fetchList',
    }),
    /**
     * TODO: move this function into mixin or helper
     *
     * @returns {*}
     */
    getQuery() {
      const query = omit(this.$route.query, ['page']);

      query.limit = this.limit;
      query.skip = ((this.$route.query.page - 1) * this.limit) || 0;

      return query;
    },
    fetchList() {
      this.fetchListAction({
        params: this.getQuery(),
      });
    },
  },
  watch: {
    $route() {
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
