<template>
  <v-layout
    class="kpi-widget"
    column
    align-center
  >
    <h4 class="kpi-widget__title">
      {{ title }}
    </h4>
    <v-layout
      class="kpi-widget__list"
      wrap
    >
      <numbers-metrics-item
        v-for="metric in metrics"
        :key="metric.title"
        :metric="metric"
        :show-trend="showTrend"
        :value-font-size="fontSize"
      />
    </v-layout>
    <v-layout
      class="kpi-widget__actions mt-4"
      justify-end
    >
      <v-btn
        :loading="downloading"
        color="primary"
        small
        @click="$emit('export:csv')"
      >
        <v-icon
          small
          left
        >
          file_download
        </v-icon>
        <span>{{ $t('common.exportAsCsv') }}</span>
      </v-btn>
    </v-layout>
  </v-layout>
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
