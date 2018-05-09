<template lang="pug">
  div( v-if="fetchComplete" )
    basic-list( :items="items" )
      tr.container( slot="header" )
          th.box( v-for="columnName in Object.keys(alarmProperty)" @click="sortAlarms(columnName)" ) {{ columnName }}
          th.box
      tr.container( slot="row" slot-scope="item" )
          td.box( v-for="property in Object.values(alarmProperty)" ) {{ getProp(item.props, property) }}
          td.box
            actions-panel.actions
      tr.container(slot="expandedRow" slot-scope="item")
          time-line(:idAlarm="item.props.d")
    alarm-list-pagination( :itemsPerPage="itemsPerPage" @changedPage="changePage" v-if="fetchComplete" )
    loader(v-else)
</template>

<script>
import getProp from 'lodash/get';
import merge from 'lodash/merge';

import { createNamespacedHelpers } from 'vuex';
import BasicList from '../BasicComponent/basic-list.vue';
import ActionsPanel from '../BasicComponent/actions-panel.vue';
import Loader from '../loaders/alarm-list-loader.vue';
import AlarmListPagination from './alarm-list-pagination.vue';
import TimeLine from '../time-line.vue';


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
    TimeLine,
    AlarmListPagination,
    ActionsPanel,
    BasicList,
    Loader,
  },
  mounted() {
    this.fetchList(this.fetchingParams);
  },
  props: {
    alarmProperty: {
      type: Object,
      default() {
        return {};
      },
    },
    itemsPerPage: {
      type: Number,
      default: 5,
    },
  },
  data() {
    return {
      columnSorted: 'opened',
      fetchingParams: {
        params: {
          limit: this.itemsPerPage,
        },
      },
    };
  },
  computed: {
    ...mapGetters([
      'items',
      'meta',
      'fetchComplete',
    ]),
  },
  methods: {
    ...mapActions(['fetchList']),
    changePage(newFetchingParameters) {
      merge(this.fetchingParams, newFetchingParameters);
      this.fetchList(this.fetchingParams);
    },
    getProp,
    sortAlarms(columnToSort) {
      this.fetchingParams.sort_key = this.alarmProperty[columnToSort];
      this.fetchingParams.sort_dir = 'DESC';
      this.fetchList(this.fetchingParams);
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
