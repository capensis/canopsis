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
          searching
        v-flex(xs4)
          pagination(:meta="meta", :limit="limit", type="top")
        v-flex(xs2)
          records-per-page
    div(v-if="!pending")
      basic-list(:items="items", @update:selected="selected = $event")
        tr.container.header.pa-0(slot="header")
          th.box(v-for="column in alarmProperties")
            span {{ column.text }}
            list-sorting(:column="column.value", class="blue--text")
          th.box
        tr.container(slot="row" slot-scope="item")
            td.box(v-for="property in alarmProperties")
              alarm-column-value(:alarm="item.props", :pathToProperty="property.value", :filter="property.filter")
            td.box
              actions-panel.actions(:item="item.props")
        tr.container(slot="expandedRow", slot-scope="item")
          time-line(:alarmProps="item.props")
      v-layout(wrap)
        v-flex(xs12, md7)
        pagination(:meta="meta", :limit="limit")
        records-per-page
    loader(v-else)
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

// TABLE
import BasicList from '@/components/tables/basic-list.vue';
import ListSorting from '@/components/tables/list-sorting.vue';
// ACTIONS
import ActionsPanel from '@/components/other/alarm-list/actions/actions-panel.vue';
import MassActions from '@/components/other/alarm-list/actions/mass-actions.vue';
// TIMELINE
import TimeLine from '@/components/other/alarm-list/timeline/time-line.vue';
// LOADER
import Loader from '@/components/other/alarm-list/loader/alarm-list-loader.vue';
// SEARCHING
import Searching from '@/components/other/alarm-list/searching/alarm-list-searching.vue';
// PAGINATION
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import PaginationMixin from '@/mixins/pagination';
// COLUMNS FORMATTING
import AlarmColumnValue from '@/components/other/alarm-list/columns-formatting/alarm-column-value.vue';
// FILTER SELECTOR
import FilterSelector from '@/components/other/filter/filter-selector.vue';

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
    RecordsPerPage,
    ListSorting,
    TimeLine,
    MassActions,
    Searching,
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
      selected: [],
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
