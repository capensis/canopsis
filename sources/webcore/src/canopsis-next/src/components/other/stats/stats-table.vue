<template lang="pug">
  div
    progress-overlay(:pending="pending")
    v-data-table(
      :items="stats",
      :headers="columns",
      :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES"
    )
      template(slot="headers", slot-scope="{ headers }")
        th {{ $t('common.entity') }}
        th(v-for="header in headers", :key="header.value") {{ header.value }}
      template(slot="items", slot-scope="{ item }")
        td {{ item.entity.name }}
        td(v-for="(property, key) in widget.parameters.stats")
          template(
          v-if="item[key] && item[key].value !== undefined && item[key].value !== null"
          )
            td
              div {{ item[key].value }}
                sub {{ item[key].trend }}
          div(v-else) {{ $t('tables.noData') }}
</template>

<script>
import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';
import widgetStatsQueryMixin from '@/mixins/widget/stats/stats-query';

import ProgressOverlay from '@/components/layout/progress/progress-overlay.vue';

export default {
  components: {
    ProgressOverlay,
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
    columns() {
      return Object.keys(this.widget.parameters.stats).map(item => ({ value: item }));
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
      this.pending = true;

      const stats = await this.fetchStatsListWithoutStore({
        params: this.getQuery(),
      });

      this.stats = stats.values;
      this.pending = false;
    },
  },
};
</script>
