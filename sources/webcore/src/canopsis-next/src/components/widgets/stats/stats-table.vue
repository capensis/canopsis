<template lang="pug">
  div
    c-progress-overlay(:pending="pending")
    c-alert-overlay(
      :value="hasError",
      :message="serverErrorMessage"
    )
    v-data-table(
      :items="stats",
      :headers="columns",
      :pagination.sync="pagination",
      :custom-sort="customSort"
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

import { dataTableCustomSortWithNullIgnoring } from '@/helpers/sort';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetFetchQueryMixin from '@/mixins/widget/fetch-query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetStatsQueryMixin from '@/mixins/widget/stats/stats-query';
import widgetStatsWrapperMixin from '@/mixins/widget/stats/stats-wrapper';
import widgetStatsTableWrapperMixin from '@/mixins/widget/stats/stats-table-wrapper';

import AlarmChips from '@/components/widgets/alarm/alarm-chips.vue';

export default {
  components: {
    AlarmChips,
  },
  filters: {
    statValue(name) {
      return `${name}.value`;
    },
  },
  mixins: [
    entitiesStatsMixin,
    widgetFetchQueryMixin,
    entitiesUserPreferenceMixin,
    widgetStatsQueryMixin,
    widgetStatsWrapperMixin,
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
    customSort: dataTableCustomSortWithNullIgnoring,

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
        this.serverErrorMessage = err.description || this.$t('errors.statsRequestProblem');
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
