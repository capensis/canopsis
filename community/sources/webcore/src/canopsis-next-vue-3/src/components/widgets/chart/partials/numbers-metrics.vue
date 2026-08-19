<template>
  <c-layout class="kpi-widget" align="center" column>
    <h4 class="kpi-widget__title">
      {{ title }}
    </h4>
    <c-layout class="kpi-widget__list w-100" wrap>
      <numbers-metrics-item
        v-for="metric in metrics"
        :key="metric.title"
        :metric="metric"
        :show-trend="showTrend"
        :value-font-size="fontSize"
      />
    </c-layout>
    <c-layout
      v-if="!hideActions"
      class="kpi-widget__actions mt-4"
      justify="end"
    >
      <v-btn
        :loading="downloading"
        color="primary"
        size="small"
        @click="$emit('export:csv')"
      >
        <v-icon class="mr-1" size="small">
          file_download
        </v-icon>
        <span>{{ $t('common.exportAsCsv') }}</span>
      </v-btn>
    </c-layout>
  </c-layout>
</template>

<script>
import NumbersMetricsItem from './numbers-metrics-item.vue';

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
    hideActions: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['export:csv'],
};
</script>
