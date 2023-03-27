<template lang="pug">
  v-layout.numbers-metrics(column, align-center)
    h4.numbers-metrics__title {{ title }}
    v-layout.numbers-metrics__list(row, wrap)
      numbers-metrics-item(
        v-for="metric in metrics",
        :key="metric.title",
        :metric="metric",
        :show-trend="showTrend"
      )
</template>

<script>
import NumbersMetricsItem from '@/components/widgets/chart/partials/numbers-metrics-item.vue';

export default {
  components: { NumbersMetricsItem },
  props: {
    metrics: {
      type: Array,
      default: () => [],
    },
    title: {
      type: String,
      required: false,
    },
    showTrend: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
.numbers-metrics {
  &__title {
    font-size: 18px;
    font-weight: 500;
    margin-bottom: 16px;
  }

  &__list {
    gap: 50px;
  }
}
</style>
