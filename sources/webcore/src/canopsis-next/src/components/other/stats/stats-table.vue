<template lang="pug">
  div
    v-data-table(
      :items="stats",
      :headers="columns",
      :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES"
    )
      v-progress-linear(slot="progress", color="primary", indeterminate)
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

import StatsNumber from './stats-number.vue';

export default {
  components: {
    StatsNumber,
  },
  mixins: [entitiesStatsMixin, widgetQueryMixin, entitiesUserPreferenceMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      stats: [],
    };
  },
  computed: {
    columns() {
      return Object.keys(this.widget.parameters.stats).map(item => ({ value: item }));
    },
  },
  methods: {
    async fetchList() {
      const query = { ...this.query };

      const stats = await this.fetchStatsListWithoutStore({
        params: query,
      });

      this.stats = stats.values;
    },
  },
};
</script>
