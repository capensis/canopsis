<template lang="pug">
  div
    v-card.position-relative
      c-progress-overlay(:pending="pending")
      c-alert-overlay(
        :value="hasError",
        :message="serverErrorMessage"
      )
      v-data-table(
        :items="stats",
        :headers="tableHeaders",
        :pagination.sync="pagination"
      )
        template(
          slot="items",
          slot-scope="{ item }",
          xs12
        )
          td {{ item.entity.name }}
          td
            v-layout(align-center)
              v-chip.px-1(:style="{ backgroundColor: getChipColor(item[query.stat.title].value) }", color="white--text")
                div.body-1.font-weight-bold {{ getChipText(item[query.stat.title].value) }}
              div.caption
                div(v-if="hasTrend(item[query.stat.title])")
                  sub.ml-2
                    v-icon.caption(
                      small,
                      :color="item[query.stat.title].trend | trendColor"
                    ) {{ item[query.stat.title].trend | trendIcon }}
                  sub {{ item[query.stat.title].trend | formatValue(widget.parameters.stat.stat.value) }}
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';
import { STATS_DISPLAY_MODE, STATS_CRITICITY, SORT_ORDERS } from '@/constants';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetFetchQueryMixin from '@/mixins/widget/fetch-query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetStatsQueryMixin from '@/mixins/widget/stats/stats-query';
import widgetStatsWrapperMixin from '@/mixins/widget/stats/stats-wrapper';
import widgetStatsTableWrapperMixin from '@/mixins/widget/stats/stats-table-wrapper';

export default {
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
      serverErrorMessage: null,
      stats: [],
      pagination: {
        page: 1,
        sortBy: 'value',
        descending: true,
        totalItems: 0,
        rowsPerPage: PAGINATION_LIMIT,
      },
    };
  },
  computed: {
    statColumn() {
      const { stat } = this.query;

      if (stat) {
        return `${stat.title}.value`;
      }

      return 'value';
    },

    tableHeaders() {
      return [
        {
          text: this.$t('common.entity'),
          value: 'entity.name',
          sortable: false,
        },
        {
          text: this.$t('common.value'),
          value: this.statColumn,
        },
      ];
    },

    getChipColor() {
      return (value) => {
        const { colors, criticityLevels } = this.widget.parameters.displayMode.parameters;

        if (value < criticityLevels.minor) {
          return colors.ok;
        } else if (value < criticityLevels.major) {
          return colors.minor;
        } else if (value < criticityLevels.critical) {
          return colors.major;
        }

        return colors.critical;
      };
    },
    getChipText() {
      return (value) => {
        const { mode, parameters } = this.widget.parameters.displayMode;
        const { criticityLevels } = parameters;

        if (mode === STATS_DISPLAY_MODE.criticity) {
          if (value < criticityLevels.minor) {
            return STATS_CRITICITY.ok;
          } else if (value < criticityLevels.major) {
            return STATS_CRITICITY.minor;
          } else if (value < criticityLevels.critical) {
            return STATS_CRITICITY.major;
          }
          return STATS_CRITICITY.critical;
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
      try {
        const { limit, sortOrder } = this.query;

        this.pending = true;
        this.serverErrorMessage = null;

        const { values } = await this.fetchStatsListWithoutStore({
          params: this.getQuery(),
        });

        this.stats = values;
        this.pagination = {
          page: 1,
          sortBy: this.statColumn,
          totalItems: values.length,
          rowsPerPage: limit || PAGINATION_LIMIT,
          descending: sortOrder === SORT_ORDERS.desc,
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

<style lang="scss">
  .theme--light.v-datatable .v-datatable__actions {
    display: flex;
    justify-content: center;
  }
</style>
