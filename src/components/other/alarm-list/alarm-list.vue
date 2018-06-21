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
      v-btn(icon, @click="openSettingsPanel")
        v-icon settings
  v-layout.my-2(wrap, justify-space-between, align-center)
    v-flex(xs12 md5)
      alarm-list-searching
    v-flex(xs4)
      pagination(:meta="meta", :limit="limit", type="top")
  basic-list(:items="items", :pending="pending", @update:selected="selected = $event", expanded)
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
    pagination(:meta="meta", :limit="limit")
    records-per-page
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import intersectionWith from 'lodash/intersectionWith';
import omit from 'lodash/omit';

// TABLE
import BasicList from '@/components/tables/basic-list.vue';
import ListSorting from '@/components/tables/list-sorting.vue';
// ACTIONS
import ActionsPanel from '@/components/other/alarm-list/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm-list/actions/mass-actions-panel.vue';
// TIMELINE
import TimeLine from '@/components/other/alarm-list/timeline/time-line.vue';
// LOADER
import Loader from '@/components/other/alarm-list/loader/alarm-list-loader.vue';
// SEARCHING
import AlarmListSearching from '@/components/other/alarm-list/searching/alarm-list-searching.vue';
// PAGINATION
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import PaginationMixin from '@/mixins/pagination';
// COLUMNS FORMATTING
import AlarmColumnValue from '@/components/other/alarm-list/columns-formatting/alarm-column-value.vue';
// FILTER SELECTOR
import FilterSelector from '@/components/other/filter/filter-selector.vue';
// MODAL
import ModalMixin from '@/mixins/modal/modal';

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
    MassActionsPanel,
    AlarmListSearching,
    ActionsPanel,
    BasicList,
    Loader,
    AlarmColumnValue,
    FilterSelector,
  },
  mixins: [PaginationMixin, ModalMixin],
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
  watch: {
    items(items) {
      this.selected = intersectionWith(
        this.selected,
        items,
        (selectedItemId, item) => selectedItemId === item._id,
      );
    },
  },
  methods: {
    ...alarmMapActions({
      fetchListAction: 'fetchList',
    }),
    ...settingsMapActions({
      openSettingsPanel: 'openPanel',
    }),

    removeHistoryFilter() {
      const query = omit(this.$route.query, ['interval']);
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
