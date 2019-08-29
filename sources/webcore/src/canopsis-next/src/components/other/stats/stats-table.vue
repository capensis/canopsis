<template lang="pug">
  div
    v-card.position-relative
      progress-overlay(:pending="pending")
      stats-alert-overlay(:value="hasError", :message="serverErrorMessage")
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
                  div {{ item[key].value | formatValue(property.stat.value) }}
                  div(v-if="hasTrend(item[key])")
                    sub.ml-2
                      v-icon.caption(
                        small,
                        :color="item[key].trend | trendColor"
                      ) {{ item[key].trend | trendIcon }}
                    sub {{ item[key].trend | formatValue(property.stat.value) }}
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
import widgetStatsTableWrapperMixin from '@/mixins/widget/stats/stats-table-wrapper';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';
import AlarmChips from '@/components/other/alarm/alarm-chips.vue';

import StatsAlertOverlay from './partials/stats-alert-overlay.vue';

export default {
  components: {
    ProgressOverlay,
    AlarmChips,
    StatsAlertOverlay,
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
    widgetStatsTableWrapperMixin,
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
      hasError: false,
      serverErrorMessage: null,
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
      const { stats: widgetStats } = this.widget.parameters;
      const statsOrderedColumns = Object.keys(widgetStats)
        .sort((a, b) => widgetStats[a].position - widgetStats[b].position)
        .map(item => ({
          text: item,
          value: this.$options.filters.statValue(item),
        }));

      return [
        {
          text: this.$t('common.entity'),
          value: 'entity.name',
          sortable: false,
        },

        ...statsOrderedColumns,
      ];
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
      try {
        const { sort = {} } = this.widget.parameters;

        this.pending = true;
        this.hasError = false;
        this.serverErrorMessage = null;

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
      } catch (err) {
        this.hasError = true;
        this.serverErrorMessage = err.description || null;
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
