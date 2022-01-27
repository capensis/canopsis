<template lang="pug">
  v-layout(justify-end)
    v-btn.ma-0.mr-4(:loading="downloading", color="primary", small, @click="$emit('export:csv')")
      v-icon(small, left) file_download
      span {{ $t('common.exportAsCsv') }}
    v-btn.ma-0(color="primary", small, @click="exportChartAsPng")
      v-icon(small, left) file_download
      span {{ $t('common.downloadAsPng') }}
</template>

<script>
import { canvasToBlob } from '@/helpers/charts/canvas';

export default {
  props: {
    chart: {
      type: Object,
      required: true,
    },
    downloading: {
      type: Boolean,
      default: false,
    },
  },
  methods: {
    async exportChartAsPng() {
      this.$emit('export:png', await canvasToBlob(this.chart.canvas));
    },
  },
};
</script>
