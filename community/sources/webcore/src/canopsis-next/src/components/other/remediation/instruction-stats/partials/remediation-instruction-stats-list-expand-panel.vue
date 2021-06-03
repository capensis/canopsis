<template lang="pug">
  div
    v-tabs(v-model="activeTab", color="secondary lighten-1", dark, centered, slider-color="primary")
      v-tab {{ $t('common.summary') }}
      v-tab {{ $t('common.statistics') }}
      v-tab {{ $t('remediationInstructionStats.alarmsTimeline') }}
      v-tab {{ $t('remediationInstructionStats.rating') }}
    v-layout.pa-3.secondary.lighten-2
      v-flex(xs12)
        v-card.pa-3
          v-tabs-items.pt-2(v-model="activeTab")
            v-tab-item(lazy)
              v-flex(offset-xs2, xs8)
                remediation-instruction-stats-information-tab(:remediation-instruction="remediationInstruction")
            v-tab-item(lazy)
              remediation-instruction-stats-statistics-tab(:remediation-instruction="remediationInstruction")
            v-tab-item(lazy)
              remediation-instruction-stats-alarms-timeline-tab(:remediation-instruction="remediationInstruction")
            v-tab-item(lazy)
              remediation-instruction-stats-rating-tab(:remediation-instruction="remediationInstruction")
</template>

<script>
import RemediationInstructionStatsInformationTab from './remediation-instruction-stats-information-tab.vue';
import RemediationInstructionStatsStatisticsTab from './remediation-instruction-stats-statistics-tab.vue';
import RemediationInstructionStatsAlarmsTimelineTab from './remediation-instruction-stats-alarms-timeline-tab.vue';
import RemediationInstructionStatsRatingTab from './remediation-instruction-stats-rating-tab.vue';

export default {
  components: {
    RemediationInstructionStatsInformationTab,
    RemediationInstructionStatsStatisticsTab,
    RemediationInstructionStatsAlarmsTimelineTab,
    RemediationInstructionStatsRatingTab,
  },
  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
  },
  data() {
    return {
      activeTab: 0,
    };
  },
};
</script>
