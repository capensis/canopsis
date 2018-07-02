<template lang="pug">
  v-container
    v-layout(justify-space-between, align-center)
      v-flex.ml-4(xs4)
        mass-actions-panel(v-show="selected.length", :itemsIds="selectedIds")
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
        v-btn(icon, @click="$emit('openSettings')")
          v-icon settings
    v-layout.my-2(wrap, justify-space-between, align-center)
      v-flex(xs12 md5)
        alarm-list-search
      v-flex(xs4)
        pagination(:meta="meta", :limit="limit", type="top")
    v-data-table(
      v-model="selected"
      :items="items",
      :headers="alarmProperties",
      item-key="_id",
      :total-items="meta.total",
      :pagination.sync="pagination",
      select-all,
      hide-actions,
    )
      template(slot="headerCell", slot-scope="props")
          span(
          ) {{ props.header.text }}
      template(slot="items", slot-scope="props")
        td
          v-checkbox(primary, hide-details, v-model="props.selected")
        td(
          v-for="prop in alarmProperties",
          @click="props.expanded = !props.expanded"
        )
          alarm-column-value(:alarm="props.item", :property="prop")
        td
          actions-panel(:item="props.item")
      template(slot="expand", slot-scope="props")
        time-line(:alarmProps="props.item", @click="props.expanded = !props.expanded")
    v-layout(wrap)
      v-flex(xs12, md7)
      pagination(:meta="meta", :limit="limit")
      //records-per-page
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import omit from 'lodash/omit';

import ListSorting from '@/components/tables/list-sorting.vue';
import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm/actions/mass-actions-panel.vue';
import TimeLine from '@/components/other/alarm/timeline/time-line.vue';
import Loader from '@/components/other/alarm/loader/alarm-list-loader.vue';
import AlarmListSearch from '@/components/other/alarm/search/alarm-list-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';
import FilterSelector from '@/components/other/filter/filter-selector.vue';
import modalMixin from '@/mixins/modal/modal';
import paginationMixin from '@/mixins/pagination';

const { mapActions: alarmMapActions, mapGetters: alarmMapGetters } = createNamespacedHelpers('alarm');

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
    AlarmListSearch,
    RecordsPerPage,
    ListSorting,
    TimeLine,
    MassActionsPanel,
    ActionsPanel,
    Loader,
    AlarmColumnValue,
    FilterSelector,
  },
  mixins: [paginationMixin, modalMixin],
  props: {
    alarmProperties: {
      type: Array,
      default: () => ([]),
    },
  },
  data() {
    return {
      selected: [],
      pagination: {},
    };
  },
  computed: {
    ...alarmMapGetters([
      'items',
      'meta',
      'pending',
    ]),
    selectedIds() {
      return this.selected.map(item => item._id);
    },
  },
  watch: {
    pagination: {
      handler(e) {
        this.$router.push({
          query: {
            page: e.page || '1',
            sort_key: e.sortBy || '',
            sort_dir: e.descending ? 'DESC' : 'ASC',
          },
        });
      },
    },
  },
  methods: {
    ...alarmMapActions({
      fetchListAction: 'fetchList',
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
</style>
