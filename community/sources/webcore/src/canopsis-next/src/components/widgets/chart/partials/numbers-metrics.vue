<template lang="pug">
  v-layout.kpi-widget(column, align-center)
    h4.kpi-widget__title {{ title }}
    v-layout.kpi-widget__list(row, wrap)
      numbers-metrics-item(
        v-for="metric in metrics",
        :key="metric.title",
        :metric="metric",
        :show-trend="showTrend",
        :value-font-size="fontSize"
      )
    v-layout.kpi-widget__actions.mt-4(row, justify-end)
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
