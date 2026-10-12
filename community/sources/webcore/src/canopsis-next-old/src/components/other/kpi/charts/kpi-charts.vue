<template>
  <v-layout
    class="pa-4"
    column
    align-center
  >
    <div class="kpi-charts">
      <c-information-block
        :title="$t('kpi.alarmMetrics')"
        class="mb-5"
      >
        <kpi-alarms :unavailable="unavailable" />
      </c-information-block>
      <c-information-block
        :title="$tc('common.rating', 2)"
        class="mb-5"
      >
        <kpi-rating :unavailable="unavailable" />
      </c-information-block>
      <c-information-block
        :title="$t('kpi.sli')"
        class="mb-5"
      >
        <kpi-sli :unavailable="unavailable" />
      </c-information-block>
    </div>
  </v-layout>
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
