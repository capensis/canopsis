<template lang="pug">
  div.pa-2
    v-fade-transition(mode="out-in")
      v-layout(v-if="countersPending", key="progress", column)
        v-flex(xs12)
          v-layout(justify-center)
            v-progress-circular(indeterminate, color="primary")
      v-layout.fill-height(v-else, key="content", wrap)
        v-alert(v-if="hasNoData", type="info", :value="true") {{ $t('tables.noData') }}
        template(v-else)
          v-flex(v-for="counter in countersWithFilters", :key="counter.filter.title", :class="flexSize")
            counter-item.weatherItem(
              :counter="counter",
              :widget="widget"
            )
</template>

<script>
import { omit } from 'lodash';

import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesCounterMixin from '@/mixins/entities/counter';
import widgetFetchQueryMixin from '@/mixins/widget/fetch-query';

import CounterItem from './counter-item.vue';

export default {
  components: {
    CounterItem,
  },
  mixins: [
    widgetPeriodicRefreshMixin,
    entitiesCounterMixin,
    widgetFetchQueryMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    flexSize() {
      return [
        `xs${this.widget.parameters.columnSM}`,
        `md${this.widget.parameters.columnMD}`,
        `lg${this.widget.parameters.columnLG}`,
      ];
    },

    hasNoData() {
      return this.counters.length === 0;
    },

    countersWithFilters() {
      const { viewFilters } = this.widget.parameters;

      return this.counters.map((counter, index) => ({ ...counter, filter: viewFilters[index] }));
    },
  },
  methods: {
    getQuery() {
      return omit(this.query, ['filters']);
    },

    fetchList() {
      this.fetchCountersList({
        widgetId: this.widget._id,
        filters: this.query.filters,
        params: this.getQuery(),
      });
    },
  },
};
</script>

<style lang="scss" scoped>
  .weatherItem {
    height: 100%;
  }
</style>
