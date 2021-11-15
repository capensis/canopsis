<template lang="pug">
  v-flex.pa-4(md10, offset-md1, lg8, offset-lg2)
    c-information-block.mb-5(:title="$t('kpi.alarmMetrics')")
    c-information-block.mb-5(:title="$tc('common.rating', 2)")
      kpi-rating
    c-information-block.mb-5(:title="$t('kpi.sli')")
      kpi-sli
</template>

<script>
import KpiSli from './kpi-sli.vue';
import KpiRating from './kpi-rating.vue';

export default {
  components: {
    KpiRating,
    KpiSli,
  },
};
</script>
