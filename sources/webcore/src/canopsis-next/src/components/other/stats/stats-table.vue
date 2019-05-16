<template lang="pug">
  div
    progress-overlay(:pending="pending")
    v-data-table(
      :items="stats",
      :headers="columns",
      :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES"
    )
      template(slot="items", slot-scope="{ item }")
        td {{ item.entity.name }}
        td(v-for="(property, key) in widget.parameters.stats")
          template(v-if="isStatNotEmpty(item[key])")
            td(v-if="property.stat.value === $constants.STATS_TYPES.currentState.value")
              alarm-chips(:type="$constants.ENTITY_INFOS_TYPE.state", :value="item[key].value")
            td(v-else)
              div {{ getFormattedValue(item[key].value, property.stat.value) }}
                sub {{ item[key].trend }}
          div(v-else) {{ $t('tables.noData') }}
</template>

<script>
import { isUndefined, isNull } from 'lodash';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetStatsQueryMixin from '@/mixins/widget/stats/stats-query';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';
import AlarmChips from '@/components/other/alarm/alarm-chips.vue';

export default {
  components: {
    ProgressOverlay,
    AlarmChips,
  },
  mixins: [
    entitiesStatsMixin,
    widgetQueryMixin,
    entitiesUserPreferenceMixin,
    widgetStatsQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: true,
      stats: [],
    };
  },
  computed: {
    isStatNotEmpty() {
      return stat => stat && !isUndefined(stat.value) && !isNull(stat.value);
    },
    columns() {
      return [
        {
          text: this.$t('common.entity'),
          value: 'entity.name',
        },

        ...Object.keys(this.widget.parameters.stats).map(item => ({
          text: item,
          value: this.widget.parameters.stats[item].stat.value,
        })),
      ];
    },
    getFormattedValue() {
      const PROPERTIES_FILTERS_MAP = {
        state_rate: value => this.$options.filters.percentage(value),
        ack_time_sla: value => this.$options.filters.percentage(value),
        resolve_time_sla: value => this.$options.filters.percentage(value),
        time_in_state: value => this.$options.filters.duration(value),
        mtbf: value => this.$options.filters.duration(value),
      };

      return (value, columnValue) => {
        if (PROPERTIES_FILTERS_MAP[columnValue]) {
          return PROPERTIES_FILTERS_MAP[columnValue](value);
        }

        return value;
      };
    },
  },
  methods: {
    getQuery() {
      const {
        stats,
        mfilter,
        tstop,
        periodUnit,
        tstart,
      } = this.getStatsQuery();

      const durationValue = tstop.diff(tstart, periodUnit);

      return {
        stats,
        mfilter,

        duration: `${durationValue}${periodUnit.toLowerCase()}`,
        tstop: tstop.startOf('h').unix(),
      };
    },

    async fetchList() {
      this.pending = true;

      const { values } = await this.fetchStatsListWithoutStore({
        params: this.getQuery(),
      });

      this.stats = values;
      this.pending = false;
    },
  },
};
</script>
