<template lang="pug">
  div
    div.white
      v-layout(justify-space-between, align-center)
        v-flex.ml-4(xs4)
          mass-actions
        v-btn(icon, @click="openSettingsPanel")
          v-icon settings
      v-layout.my-2(wrap, justify-space-between, align-center)
        v-flex(xs12 md5)
          alarm-list-searching
        v-flex(xs4)
          pagination(:meta="meta", :limit="limit", type="top")
        v-flex(xs2)
          page-iterator
    div(v-if="!pending")
      basic-list(:items="items")
        tr.container.header.pa-0(slot="header")
          v-checkbox.checkbox.box( @click.stop="selectAll(items)", v-model="allSelected", hide-details)
          th.box(v-for="column in alarmProperties")
            span {{ column.text }}
            list-sorting(:column="column.value", class="blue--text")
            th.box
        tr.container(slot="row" slot-scope="item")
            v-checkbox.checkbox(@click.stop="select", v-model="selected", :value="item.props._id", hide-details)
            td.box(v-for="property in alarmProperties")
              alarm-column-value(:alarm="item.props", :pathToProperty="property.value", :filter="property.filter")
            td.box
              actions-panel.actions(:item="item.props")
        tr.container(slot="expandedRow", slot-scope="item")
          time-line(:alarmProps="item.props")
      v-layout(wrap)
        v-flex(xs12, md7)
        pagination(:meta="meta", :limit="limit")
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
import AlarmColumnValue from '@/components/alarm-list/alarm-column-value.vue';
import FilterSelector from '@/components/other/filter/selector.vue';

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
    AlarmColumnValue,
    FilterSelector,
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

  .bottomToolbox {
    display: flex;
    flex-flow: row wrap;
  }

  .checkbox {
    flex: 0.5;
    padding: 0 0.5em;
  }
</style>
