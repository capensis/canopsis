<template>
  <v-layout
    class="kpi-chart-export-actions"
    justify-end="justify-end"
    wrap="wrap"
  >
    <v-btn
      class="ma-0"
      :loading="downloading"
      color="primary"
      small="small"
      @click="$emit('export:csv')"
    >
      <v-icon
        small="small"
        left="left"
      >
        file_download
      </v-icon><span>{{ $t('common.exportAsCsv') }}</span>
    </v-btn>
    <v-btn
      class="ma-0"
      color="primary"
      small="small"
      @click="exportChartAsPng"
    >
      <v-icon
        small="small"
        left="left"
      >
        file_download
      </v-icon><span>{{ $t('common.downloadAsPng') }}</span>
    </v-btn>
  </v-layout>
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

<style scoped lang="scss">
.kpi-chart-export-actions {
  column-gap: 24px;
  row-gap: 12px;
}
</style>
