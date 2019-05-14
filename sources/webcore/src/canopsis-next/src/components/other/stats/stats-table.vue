<template lang="pug">
  div
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
            td
              div {{ item[key].value }}
                sub {{ item[key].trend }}
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

export default {
  components: {
    ProgressOverlay,
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
  },
  methods: {
    getQuery() {
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

        tstop: tstop.startOf('h').unix(),
      };
    },

    async fetchList() {
      const { sort } = this.widget.parameters;

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
