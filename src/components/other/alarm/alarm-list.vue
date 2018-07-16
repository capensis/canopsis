<template lang="pug">
  v-container
    v-layout.white(wrap, justify-space-between, align-center)
      v-flex(xs12 md3)
        alarm-list-search
      v-flex(xs2)
        pagination(:meta="meta", :limit="limit", type="top")
      v-flex.ml-4(xs3)
        mass-actions-panel(v-show="selected.length", :itemsIds="selectedIds")
      v-flex(xs3)
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
    transition(name="fade", mode="out-in")
      loader.loader(v-if="pending")
      div(v-else)
          v-data-table(
            v-model="selected",
            :items="items",
            :headers="alarmProperties",
            item-key="_id",
            :total-items="meta.total",
            :pagination.sync="pagination",
            select-all,
            hide-actions,
          )
            template(slot="headerCell", slot-scope="props")
                span {{ props.header.text }}
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
          v-layout.white(align-center)
            v-flex(xs10)
              pagination(:meta="meta", :limit="limit")
            v-spacer
            v-flex(xs2)
              records-per-page
</template>

<script>
import { createNamespacedHelpers } from 'vuex';
import omit from 'lodash/omit';

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
import dateIntervals from '@/helpers/date-intervals';
import alarmsMixin from '@/mixins/alarms';

const { mapGetters: alarmMapGetters } = createNamespacedHelpers('alarm');

/**
 * Alarm-list component
 *
 * @module alarm
 *
 * @prop {object} alarmProperties - Object that describe the columns names and the alarms attributes corresponding
 *            e.g : { ColumnName : 'att1.att2', Connector : 'v.connector' }
 * @prop {integer} [itemsPerPage=5] - Number of Alarm to display per page
 *
 * @event openSettings#click
 */
export default {
  components: {
    AlarmListSearch,
    RecordsPerPage,
    TimeLine,
    MassActionsPanel,
    ActionsPanel,
    Loader,
    AlarmColumnValue,
    FilterSelector,
  },
  mixins: [alarmsMixin, paginationMixin, modalMixin],
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
    selectedIds() {
      return this.selected.map(item => item._id);
    },
  },
  methods: {
    removeHistoryFilter() {
      const query = omit(this.$route.query, ['interval']);
      this.$router.push({ query });
    },

    getQuery() {
      const query = omit(this.$route.query, ['page', 'interval']);

      if (this.$route.query.interval && this.$route.query.interval !== 'custom') {
        try {
          const { tstart, tstop } = dateIntervals[this.$route.query.interval]();
          query.tstart = tstart;
          query.tstop = tstop;
        } catch (err) {
          console.warn(err);
        }
      }
      query.limit = this.limit;
      query.skip = ((this.$route.query.page - 1) * this.limit) || 0;

      return query;
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
  .fade-enter-active, .fade-leave-active {
    transition: opacity .5s;
  }
  .fade-enter, .fade-leave-to {
    opacity: 0;
  }
  .loader {
    top: 15%;
    position: absolute;
  }
</style>
