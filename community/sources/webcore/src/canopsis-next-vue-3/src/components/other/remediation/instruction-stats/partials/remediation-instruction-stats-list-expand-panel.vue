<template>
  <div>
    <v-tabs
      v-model="activeTab"
      bg-color="secondary-lighten-1"
      color="primary"
      align-tabs="center"
    >
      <v-tab value="summary">{{ t('common.summary') }}</v-tab>
      <v-tab value="statistics">{{ t('common.statistics') }}</v-tab>
      <v-tab v-if="remediationInstructionStatsItem.has_executions" value="alarmsTimeline">
        {{ t('remediation.instructionStat.alarmsTimeline') }}
      </v-tab>
      <v-tab value="rating">{{ t('common.rating') }}</v-tab>
    </v-tabs>

    <v-window v-model="activeTab" class="bg-secondary-lighten-2">
      <c-window-expand-item value="summary" cols="8">
        <remediation-instruction-stats-summary-tab :remediation-instruction="remediationInstruction" />
      </c-window-expand-item>
      <c-window-expand-item value="statistics" cols="8">
        <remediation-instruction-stats-statistics-tab
          :remediation-instruction="remediationInstruction"
          :interval="interval"
        />
      </c-window-expand-item>
      <c-window-expand-item
        v-if="remediationInstructionStatsItem.has_executions"
        value="alarmsTimeline"
      >
        <remediation-instruction-stats-alarms-timeline-tab
          :remediation-instruction="remediationInstruction"
          :interval="interval"
        />
      </c-window-expand-item>
      <c-window-expand-item value="rating" cols="8">
        <remediation-instruction-stats-rating-tab
          :remediation-instruction="remediationInstruction"
          :interval="interval"
        />
      </c-window-expand-item>
    </v-window>
  </div>
</template>

<script>
import { onMounted, ref, watch } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useRemediationInstructionStats } from '@/hooks/store/modules/remediation-instruction-stats';

import RemediationInstructionStatsAlarmsTimelineTab from './remediation-instruction-stats-alarms-timeline-tab.vue';
import RemediationInstructionStatsRatingTab from './remediation-instruction-stats-rating-tab.vue';
import RemediationInstructionStatsStatisticsTab from './remediation-instruction-stats-statistics-tab.vue';
import RemediationInstructionStatsSummaryTab from './remediation-instruction-stats-summary-tab.vue';

export default {
  components: {
    RemediationInstructionStatsAlarmsTimelineTab,
    RemediationInstructionStatsRatingTab,
    RemediationInstructionStatsStatisticsTab,
    RemediationInstructionStatsSummaryTab,
  },

  props: {
    remediationInstructionStatsItem: {
      type: Object,
      required: true,
    },
    interval: {
      type: Object,
      required: true,
    },
  },

  setup(props) {
    const { t } = useI18n();
    const { fetchRemediationInstructionStatsSummaryWithoutStore } =
      useRemediationInstructionStats();
    const activeTab = ref('summary');
    const remediationInstruction = ref({});

    const fetchRemediationInstructionStatsSummary = async () => {
      remediationInstruction.value = await fetchRemediationInstructionStatsSummaryWithoutStore({
        id: props.remediationInstructionStatsItem._id,
      });
    };

    watch(() => props.remediationInstructionStatsItem, fetchRemediationInstructionStatsSummary);
    onMounted(fetchRemediationInstructionStatsSummary);

    return {
      t,
      activeTab,
      remediationInstruction,
    };
  },
};
</script>
