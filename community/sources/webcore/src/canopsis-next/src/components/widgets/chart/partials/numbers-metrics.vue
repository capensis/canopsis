<template lang="pug">
  v-layout.numbers-metrics(column, align-center)
    h4.numbers-metrics__title {{ title }}
    v-layout.numbers-metrics__list(row, wrap)
      numbers-metrics-item(
        v-for="metric in metrics",
        :key="metric.title",
        :metric="metric",
        :show-trend="showTrend",
        :value-font-size="fontSize"
      )
    v-layout.numbers-metrics__actions.mt-4(row, justify-end)
      v-btn.ma-0(:loading="downloading", color="primary", small, @click="$emit('export:csv')")
        v-icon(small, left) file_download
        span {{ $t('common.exportAsCsv') }}
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
    downloading: {
      type: Boolean,
      default: false,
    },
    fontSize: {
      type: Number,
      required: false,
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

  &__actions {
    width: 100%;
  }
}
</style>
