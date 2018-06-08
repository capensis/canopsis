<template lang="pug">
  div
    v-layout(justify-space-between)
      v-flex
        alarm-list-searching
    v-layout(justify-space-between, align-center)
      v-flex(xs4, class="ml-4")
        mass-actions
      v-flex(xs5)
        pagination(:meta="meta", :limit="limit", type="top")
      v-btn(icon, @click="openSettingsPanel")
        v-icon settings
    div(v-if="!pending")
      basic-list(:items="items")
        tr.container(slot="header", class="pa-0")
          v-checkbox.checkbox.box( @click.stop="selectAll(items)", v-model="allSelected", hide-details)
          th.box(v-for="column in alarmProperties")
            span {{ column.text }}
            list-sorting(:column="column.value")
            th.box
        tr.container(slot="row" slot-scope="item")
            v-checkbox.checkbox(@click.stop="select", v-model="selected", :value="item.props._id", hide-details)
            td.box(v-for="property in alarmProperties") {{ item.props | get(property.value, property.filter) }}
            td.box
              actions-panel.actions(:item="item.props")
        tr.container(slot="expandedRow", slot-scope="item")
          time-line(:alarmProps="item.props")
      v-layout
        pagination(:meta="meta", :limit="limit")
        v-flex
          page-iterator
    loader(v-else)
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import BasicList from '@/components/basic-component/basic-list.vue';
import ActionsPanel from '@/components/basic-component/actions-panel.vue';
import Loader from '@/components/loaders/alarm-list-loader.vue';
import MassActions from '@/components/alarm-list/mass-actions.vue';
import AlarmListSearching from '@/components/alarm-list/alarm-list-searching.vue';
import TimeLine from '@/components/alarm-list/time-line.vue';
import ListSorting from '@/components/basic-component/list-sorting.vue';
import PageIterator from '@/components/basic-component/records-per-page.vue';
import PaginationMixin from '@/mixins/pagination';

const { mapActions: alarmMapActions, mapGetters: alarmMapGetters } = createNamespacedHelpers('alarm');
const { mapActions: settingsMapActions } = createNamespacedHelpers('alarmsListSettings');

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
    PageIterator,
    ListSorting,
    TimeLine,
    MassActions,
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
  data() {
    return {
      // alarm's ids selected by the checkboxes
      selected: [],
      allSelected: false,
    };
  },
  computed: {
    ...alarmMapGetters([
      'items',
      'meta',
      'pending',
    ]),
  },
  methods: {
    selectAll(items) {
      this.selected = [];
      if (!this.allSelected) {
        items.forEach((item) => {
          this.selected.push(item._id);
        });
      }
      this.allSelected = !this.allSelected;
    },
    select() {
      this.allSelected = false;
    },
    ...alarmMapActions({
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
  .box{
    flex: 1;
    padding: 1px;
  }
  .bottomToolbox {
    display: flex;
    flex-flow: row wrap;
  }
  .checkbox {
    flex: 0.5;
    padding: 0 0.5em;
  }
</style>
