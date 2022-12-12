<template lang="pug">
  div.pa-2
    v-fade-transition(v-if="countersPending", key="progress", mode="out-in")
      v-progress-linear.progress-linear-absolute--top(height="2", indeterminate)
    v-layout.fill-height(key="content", wrap)
      v-alert(v-if="hasNoData", type="info", :value="true") {{ $t('common.noData') }}
      template(v-else)
        v-flex(v-for="counter in countersWithFilters", :key="counter.filter.title", :class="flexSize")
          counter-item.weatherItem(
            :counter="counter",
            :widget="widget",
            :query="queryWithoutFilters"
          )
</template>

<script>
import { omit } from 'lodash';

import { widgetPeriodicRefreshMixin } from '@/mixins/widget/periodic-refresh';
import entitiesCounterMixin from '@/mixins/entities/counter';
import { widgetFetchQueryMixin } from '@/mixins/widget/fetch-query';

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
      const columnsCount = {
        m: this.widget.parameters.columnMobile,
        t: this.widget.parameters.columnTablet,
        l: this.widget.parameters.columnDesktop,
        xl: this.widget.parameters.columnDesktop,
      }[this.$mq];

      return `xs${12 / columnsCount}`;
    },

    hasNoData() {
      return this.counters.length === 0;
    },

    countersWithFilters() {
      const { filters } = this.widget;

      return this.counters.map((counter, index) => ({ ...counter, filter: filters[index] }));
    },

    queryWithoutFilters() {
      return omit(this.query, ['filters']);
    },
  },
  methods: {
    fetchList() {
      this.fetchCountersList({
        widgetId: this.widget._id,
        filters: this.query.filters,
        params: this.queryWithoutFilters,
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
