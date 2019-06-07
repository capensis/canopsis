<template lang="pug">
  div
    v-card
      progress-overlay(:pending="pending")
      v-data-table(
      :items="stats",
      :headers="columns",
      :pagination.sync="pagination"
      )
        template(slot="items", slot-scope="{ item }")
          td {{ item.entity.name }}
          td(v-for="(property, key) in widget.parameters.stats")
            template(v-if="isStatNotEmpty(item[key])")
              td(v-if="property.stat.value === $constants.STATS_TYPES.currentState.value")
                alarm-chips(:type="$constants.ENTITY_INFOS_TYPE.state", :value="item[key].value")
              td(v-else)
                v-layout(align-center)
                  div {{ getFormattedValue(item[key].value, property.stat.value) }}
                  div(v-if="item[key].trend !== undefined && item[key].trend !== null")
                    sub.ml-2
                      v-icon.caption(
                      small,
                      :color="trendFormat(item[key].value).color"
                      ) {{ trendFormat(item[key].value).icon }}
                    sub {{ getFormattedValue(item[key].trend, property.stat.value) }}
            div(v-else) {{ $t('tables.noData') }}
</template>

<script>
import { isUndefined, isNull } from 'lodash';

import { PAGINATION_LIMIT } from '@/config';
import { SORT_ORDERS } from '@/constants';

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
  filters: {
    statValue(name) {
      return `${name}.value`;
    },
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
      page: 1,
      pagination: {
        page: 1,
        sortBy: null,
        descending: true,
        totalItems: 0,
        rowsPerPage: PAGINATION_LIMIT,
      },
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
          sortable: false,
        },

        ...Object.keys(this.widget.parameters.stats).map(item => ({
          text: item,
          value: this.$options.filters.statValue(item),
        })),
      ];
    },
    getFormattedValue() {
      const PROPERTIES_FILTERS_MAP = {
        state_rate: value => this.$options.filters.percentage(value),
        ack_time_sla: value => this.$options.filters.percentage(value),
        resolve_time_sla: value => this.$options.filters.percentage(value),
        time_in_state: value => this.$options.filters.duration({ value, locale: this.$i18n.locale }),
        mtbf: value => this.$options.filters.duration({ value, locale: this.$i18n.locale }),
      };

      return (value, columnValue) => {
        if (PROPERTIES_FILTERS_MAP[columnValue]) {
          return PROPERTIES_FILTERS_MAP[columnValue](value);
        }

        return value;
      };
    },
    trendFormat() {
      return (value) => {
        if (value > 0) {
          return {
            icon: 'trending_up',
            color: 'primary',
          };
        } else if (value < 0) {
          return {
            icon: 'trending_down',
            color: 'error',
          };
        }

        return {
          icon: 'trending_flat',
        };
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
      const { sort = {} } = this.widget.parameters;

      this.pending = true;

      const { values } = await this.fetchStatsListWithoutStore({
        params: this.getQuery(),
      });

      this.stats = values;

      this.pagination = {
        page: 1,
        sortBy: sort.column ? this.$options.filters.statValue(sort.column) : null,
        totalItems: values.length,
        rowsPerPage: PAGINATION_LIMIT,
        descending: sort.order === SORT_ORDERS.desc,
      };

      this.pending = false;
    },
  },
};
</script>
