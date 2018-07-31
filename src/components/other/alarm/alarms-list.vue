<template lang="pug">
  v-container
    v-layout.white(wrap, justify-space-between, align-center)
      v-flex(xs12 md3)
        alarm-list-search(:query.sync="query")
      v-flex(xs2)
        pagination(:meta="alarmsMeta", :query.sync="query", type="top")
      v-flex.ml-4(xs3)
        mass-actions-panel(v-show="selected.length", :itemsIds="selectedIds")
      v-flex(xs3)
        v-chip(
          v-if="query.interval",
          @input="removeHistoryFilter",
          close,
          label,
          color="blue darken-4 white--text"
        ) {{ $t(`modals.liveReporting.${query.interval}`) }}
        v-btn(@click="showEditLiveReportModal", icon, small)
          v-icon(:color="query.interval ? 'blue' : 'black'") schedule
        v-btn(icon, @click="$emit('openSettings')")
          v-icon settings
    transition(name="fade", mode="out-in")
      loader.loader(v-if="alarmsPending")
      div(v-else)
          v-data-table(
            v-model="selected",
            :items="alarms",
            :headers="properties",
            item-key="_id",
            :total-items="alarmsMeta.total",
            :pagination.sync="vDataTablePagination"
            select-all,
            hide-actions,
          )
            template(slot="headerCell", slot-scope="props")
                span {{ props.header.text }}
            template(slot="items", slot-scope="props")
              td
                v-checkbox(primary, hide-details, v-model="props.selected")
              td(
                v-for="prop in properties",
                @click="props.expanded = !props.expanded"
              )
                alarm-column-value(:alarm="props.item", :property="prop", :widget="widget")
              td
                actions-panel(:item="props.item", :widget="widget")
            template(slot="expand", slot-scope="props")
              time-line(:alarmProps="props.item", @click="props.expanded = !props.expanded")
          v-layout.white(align-center)
            v-flex(xs10)
              pagination(:meta="alarmsMeta", :query.sync="query")
            v-spacer
            v-flex(xs2)
              records-per-page(:query.sync="query")
</template>

<script>
import omit from 'lodash/omit';

import ActionsPanel from '@/components/other/alarm/actions/actions-panel.vue';
import MassActionsPanel from '@/components/other/alarm/actions/mass-actions-panel.vue';
import TimeLine from '@/components/other/alarm/timeline/time-line.vue';
import Loader from '@/components/other/alarm/loader/alarm-list-loader.vue';
import AlarmListSearch from '@/components/other/alarm/search/alarm-list-search.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';
import AlarmColumnValue from '@/components/other/alarm/columns-formatting/alarm-column-value.vue';
import modalMixin from '@/mixins/modal/modal';
import queryMixin from '@/mixins/query';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

/**
 * Alarm-list component
 *
 * @module alarm
 *
 * @prop {Object} widget - Object representing the widget
 * @prop {Object} properties - Object that describe the columns names and the alarms attributes corresponding
 *            e.g : { ColumnName : 'att1.att2', Connector : 'v.connector' }
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
  },
  mixins: [
    queryMixin,
    modalMixin,
    entitiesAlarmMixin,
    entitiesUserPreferenceMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    properties: {
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
    selectedIds() {
      return this.selected.map(item => item._id);
    },
  },
  watch: {
    userPreference() {
      this.fetchList(); // TODO: check requests count
    },
  },
  async mounted() {
    await this.fetchUserPreferenceByWidgetId({ widgetId: this.widget.id });

    return this.fetchList();
  },
  methods: {
    removeHistoryFilter() {
      this.query = omit(this.query, ['interval', 'tstart', 'tstop']);
    },
    fetchList() {
      this.fetchAlarmsList({
        params: this.getQuery(),
        widgetId: this.widget.id,
      });
    },
    showEditLiveReportModal() {
      this.showModal({
        name: 'edit-live-reporting',
        config: {
          updateQuery: params => this.query = { ...this.query, ...params },
        },
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
