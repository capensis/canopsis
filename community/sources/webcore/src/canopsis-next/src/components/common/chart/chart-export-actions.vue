<template>
  <v-layout
    class="col-gap-6 row-gap-3"
    justify-end
    wrap
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
    <v-btn
      color="primary"
      small
      @click="exportChartAsPng"
    >
      <v-icon
        small
        left
      >
        file_download
      </v-icon>
      <span>{{ $t('common.downloadAsPng') }}</span>
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
