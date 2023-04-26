<template lang="pug">
  v-layout(column)
    entity-chart-widget(
      v-for="(chart, index) in charts",
      :key="index",
      :widget="chart",
      :entity="entity"
    )
</template>

<script>
import EntityChartWidget from '@/components/widgets/chart/entity-chart-widget.vue';

export default {
  components: { EntityChartWidget },
  props: {
    charts: {
      type: Array,
      required: true,
    },
    entity: {
      type: Object,
      required: true,
    },
  },
};
</script>
