<template lang="pug">
  v-flex.pa-4(v-else, md10, offset-md1, lg8, offset-lg2)
    c-information-block.mb-5(:title="$t('kpi.alarmMetrics')")
    c-information-block.mb-5(:title="$tc('common.rating', 2)")
    c-information-block.mb-5(:title="$t('kpi.sli')")
      kpi-sli-graph
</template>

<script>
import KpiSliGraph from './partials/kpi-sli-graph.vue';

export default {
  components: {
    KpiSliGraph,
  },
};
</script>
