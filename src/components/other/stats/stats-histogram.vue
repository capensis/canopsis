<template lang="pug">
  v-btn
    v-icon settings
</template>


<script>
import { Bar } from 'vue-chartjs';

export default {
  extends: Bar,
  props: {
    widget: {
      type: Object,
      required: true,
    },
    rowId: {
      type: String,
      required: true,
    },
  },
  mounted() {
    const data = { ...this.widget.parameters.data };

    this.renderChart(data);
  },
};
</script>
