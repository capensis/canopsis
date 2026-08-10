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
      <v-tab value="rating">{{ t('common.rating') }}</v-tab>
    </v-tabs>

    <c-layout class="pa-3 bg-secondary-lighten-2">
      <v-card class="pa-3 w-100">
        <v-window v-model="activeTab">
          <v-window-item value="summary">
            <v-row justify="center">
              <v-col cols="8">
                <remediation-instruction-stats-summary-tab :remediation-instruction="remediationInstruction"/>
              </v-col>
            </v-row>
          </v-window-item>
          <v-window-item value="statistics">
            <v-row justify="center">
              <v-col cols="8">
                <remediation-instruction-stats-statistics-tab
                  :remediation-instruction="remediationInstruction"
                  :interval="interval"
                />
              </v-col>
            </v-row>
          </v-window-item>
          <v-window-item value="rating">
            <v-row justify="center">
              <v-col cols="8">
                <remediation-instruction-stats-rating-tab
                  :remediation-instruction="remediationInstruction"
                  :interval="interval"
                />
              </v-col>
            </v-row>
          </v-window-item>
        </v-window>
      </v-card>
    </c-layout>
  </div>
</template>

<script>
import { onMounted, ref, watch } from 'vue';

import { useI18n } from '@/hooks/i18n';
import { useRemediationInstructionStats } from '@/hooks/store/modules/remediation-instruction-stats';

import RemediationInstructionStatsRatingTab from './remediation-instruction-stats-rating-tab.vue';
import RemediationInstructionStatsStatisticsTab from './remediation-instruction-stats-statistics-tab.vue';
import RemediationInstructionStatsSummaryTab from './remediation-instruction-stats-summary-tab.vue';

export default {
  components: {
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
