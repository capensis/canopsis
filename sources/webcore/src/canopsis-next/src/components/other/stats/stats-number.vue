<template lang="pug">
  div
    v-card(v-if="!pending")
      v-data-table(
        :items="stats",
        :headers="tableHeaders",
        :pagination.sync="pagination",
        :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
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
                template(v-if="item[query.stat.title].trend >= 0") + {{ item[query.stat.title].trend }}
    v-layout(v-else, justify-center)
      v-progress-circular(
      indeterminate,
      color="primary",
      )
</template>

<script>
import { PAGINATION_LIMIT } from '@/config';
import { STATS_DISPLAY_MODE, STATS_CRITICITY, SORT_ORDERS } from '@/constants';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetStatsQueryMixin from '@/mixins/widget/stats/stats-query';


import Ellipsis from '@/components/tables/ellipsis.vue';
import RecordsPerPage from '@/components/tables/records-per-page.vue';

export default {
  components: {
    Ellipsis,
    RecordsPerPage,
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
      pending: false,
      stats: [],
      pagination: {
        sortBy: 'value',
        descending: true,
        rowsPerPage: PAGINATION_LIMIT,
      },
      tableHeaders: [
        {
          text: this.$t('common.entity'),
          value: 'entity.name',
          sortable: true,
        },
        {
          text: this.$t('common.value'),
          value: 'value',
          sortable: true,
        },
      ],
    };
  },
  computed: {
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
      const { sortOrder, stat, limit = PAGINATION_LIMIT } = this.query;
      const {
        stats,
        mfilter,
        tstop,
        duration,
      } = this.getStatsQuery();

      return {
        duration,
        stats,
        mfilter,
        limit,

        tstop: tstop.startOf('h').unix(),
        sort_column: stat.title,
        sort_order: sortOrder ? sortOrder.toLowerCase() : SORT_ORDERS.desc.toLowerCase(),
      };
    },

    async fetchList() {
      this.pending = true;

      this.stats = await this.fetchStatValuesWithoutStore({
        params: this.getQuery(),
      });

      this.pending = false;
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
