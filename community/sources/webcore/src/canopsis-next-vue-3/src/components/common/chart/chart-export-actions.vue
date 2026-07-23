<template>
  <c-layout class="col-gap-6 row-gap-3" justify="end" wrap>
    <v-btn
      v-if="!hideCsv"
      :loading="downloading"
      :disabled="!chart"
      color="primary"
      size="small"
      @click="$emit('export:csv')"
    >
      <v-icon class="mr-1" size="small"> file_download </v-icon>
      <span>{{ $t('common.exportAsCsv') }}</span>
    </v-btn>
    <v-btn :disabled="!chart" color="primary" size="small" @click="exportChartAsPng">
      <v-icon class="mr-1" size="small"> file_download </v-icon>
      <span>{{ $t('common.downloadAsPng') }}</span>
    </v-btn>
  </c-layout>
</template>

<script>
import { canvasToBlob } from '@/helpers/charts/canvas';

export default {
  props: {
    chart: {
      type: Object,
      required: false,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
    hideCsv: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['export:csv', 'export:png'],

  methods: {
    async exportChartAsPng() {
      if (!this.chart) {
        return;
      }

      this.$emit('export:png', await canvasToBlob(this.chart.canvas));
    },
  },
};
</script>
