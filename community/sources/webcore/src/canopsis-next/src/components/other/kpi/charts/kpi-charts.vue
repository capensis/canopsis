<template lang="pug">
  v-layout.pa-4(column, align-center)
    div.kpi-charts
      c-information-block.mb-5(:title="$t('kpi.alarmMetrics')")
        kpi-alarms(:unavailable="unavailable")
      c-information-block.mb-5(:title="$tc('common.rating', 2)")
        kpi-rating(:unavailable="unavailable")
      c-information-block.mb-5(:title="$t('kpi.sli')")
        kpi-sli(:unavailable="unavailable")
</template>

<script>
import KpiAlarms from './kpi-alarms.vue';
import KpiSli from './kpi-sli.vue';
import KpiRating from './kpi-rating.vue';

export default {
  components: {
    KpiAlarms,
    KpiRating,
    KpiSli,
  },
  props: {
    unavailable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style scoped lang="scss">
.kpi-charts {
  max-width: 1500px;
  width: 100%;
}
</style>
