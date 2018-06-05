<template lang="pug">
  div
    v-layout(justify-space-between)
      mass-actions
      v-flex(xs5)
        alarm-list-searching
      v-btn(icon, @click="openSettingsPanel")
        v-icon settings
    div(v-if="!pending")
      basic-list(:items="items")
        tr.container(slot="header")
          v-checkbox.checkbox.box( @click.stop, v-model="idSelectedItems", :value="allIds(items)")
          th.box(v-for="columnName in Object.keys(alarmProperty)")
            span {{ columnName }}
            list-sorting(:column="alarmProperty[columnName]")
            th.box.actions
        tr.container(slot="row" slot-scope="item")
          v-checkbox.checkbox(
          @click.stop,
          v-model="idSelectedItems", :value="item.props._id", hide-details class="pa-0")
          td.box(v-for="property in Object.values(alarmProperty)") {{ getProp(item.props, property) }}
          td.box.actions
            actions-panel
        tr.container(slot="expandedRow", slot-scope="item")
          time-line(:alarmProps="item.props")
      alarm-list-pagination(:meta="meta", :limit="limit")
    loader(v-else)
    p {{ idSelectedItems }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import getProp from 'lodash/get';

import { PAGINATION_LIMIT } from '@/config';
import getQuery from '@/helpers/pagination';

import BasicList from '@/components/basic-component/basic-list.vue';
import ActionsPanel from '@/components/basic-component/actions-panel.vue';
import Loader from '@/components/loaders/alarm-list-loader.vue';
import MassActions from '@/components/alarm-list/mass-actions.vue';
import AlarmListPagination from '@/components/alarm-list/alarm-list-pagination.vue';
import AlarmListSearching from '@/components/alarm-list/alarm-list-searching.vue';
import TimeLine from '@/components/alarm-list/time-line.vue';
import ListSorting from '@/components/basic-component/list-sorting.vue';

const { mapActions, mapGetters } = createNamespacedHelpers('alarm');
const { mapActions: settingsMapActions } = createNamespacedHelpers('alarmsListSettings');

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
    ListSorting,
    TimeLine,
    MassActions,
    AlarmListSearching,
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
  data() {
    return {
      // alarm's ids selected by the checkbox
      idSelectedItems: [],
    };
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
    allIds(items) {
      const a = [];
      items.forEach((item) => {
        a.push(item._id);
      });
      return a;
    },
    ...mapActions({
      fetchListAction: 'fetchList',
    }),
    ...settingsMapActions({
      openSettingsPanel: 'openPanel',
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
  .checkbox {
    flex: 0.2;
  }
  .actions {
    flex: 0.6;
  }
</style>
