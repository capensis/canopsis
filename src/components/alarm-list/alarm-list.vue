<template lang="pug">
  div
    alarm-list-searching
    div(v-if="!pending")
      basic-list(:items="items")
        tr.container(slot="header")
          th.box(v-for="columnName in Object.keys(alarmProperties)")
            span {{ columnName }}
            list-sorting(:column="alarmProperties[columnName]")
            th.box
        tr.container(slot="row" slot-scope="item")
            td.box(v-for="property in Object.values(alarmProperties)") {{ getProp(item.props, property) }}
            td.box
              actions-panel.actions
        tr.container(slot="expandedRow", slot-scope="item")
            time-line(:alarmProps="item.props")
      pagination(:meta="meta", :limit="limit")
    loader(v-else)
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import getProp from 'lodash/get';

import BasicList from '@/components/basic-component/basic-list.vue';
import ActionsPanel from '@/components/basic-component/actions-panel.vue';
import Loader from '@/components/loaders/alarm-list-loader.vue';
import AlarmListSearching from '@/components/alarm-list/alarm-list-searching.vue';
import TimeLine from '@/components/alarm-list/time-line.vue';
import ListSorting from '@/components/basic-component/list-sorting.vue';
import PaginationMixin from '@/mixins/pagination';

const { mapActions, mapGetters } = createNamespacedHelpers('alarm');

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
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    ...mapGetters([
      'items',
      'meta',
      'pending',
    ]),
  },
  methods: {
    getProp,
    ...mapActions({
      fetchListAction: 'fetchList',
    }),
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
