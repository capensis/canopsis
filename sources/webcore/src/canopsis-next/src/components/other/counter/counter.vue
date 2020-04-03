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
          v-flex(v-for="(item, index) in counters", :key="item._id", :class="flexSize")
            counter-item.weatherItem(
              :counter="item",
              :widget="widget",
              :template="widget.parameters.blockTemplate",
              :filter="widget.parameters.viewFilters[index]"
            )
</template>

<script>
import widgetPeriodicRefreshMixin from '@/mixins/widget/periodic-refresh';
import entitiesCounterMixin from '@/mixins/entities/counter';
import widgetQueryMixin from '@/mixins/widget/query';

import CounterItem from './counter-item.vue';

export default {
  components: {
    CounterItem,
  },
  mixins: [
    widgetPeriodicRefreshMixin,
    entitiesCounterMixin,
    widgetQueryMixin,
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
  },
  methods: {
    getQuery() {
      return this.query;
    },

    fetchList() {
      this.fetchCountersList({
        filters: this.widget.parameters.viewFilters,
        params: this.getQuery(),
        widgetId: this.widget._id,
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
