<template>
  <div>
    <v-tabs
      v-model="activeTab"
      background-color="secondary lighten-1"
      slider-color="primary"
      dark
      centered
    >
      <v-tab>{{ $t('common.summary') }}</v-tab>
      <v-tab>{{ $t('common.statistics') }}</v-tab>
      <v-tab v-if="remediationInstructionStatsItem.has_executions">
        {{ $t('remediation.instructionStat.alarmsTimeline') }}
      </v-tab>
      <v-tab>{{ $tc('common.rating') }}</v-tab>
    </v-tabs>
    <v-layout class="pa-3 secondary lighten-2">
      <v-flex xs12>
        <v-card class="pa-3">
          <v-tabs-items
            v-model="activeTab"
            class="pt-2"
          >
            <v-tab-item>
              <v-flex
                offset-xs2
                xs8
              >
                <remediation-instruction-stats-summary-tab :remediation-instruction="remediationInstruction" />
              </v-flex>
            </v-tab-item>
            <v-tab-item>
              <v-flex
                offset-xs2
                xs8
              >
                <remediation-instruction-stats-statistics-tab
                  :remediation-instruction="remediationInstruction"
                  :interval="interval"
                />
              </v-flex>
            </v-tab-item>
            <v-tab-item v-if="remediationInstructionStatsItem.has_executions">
              <remediation-instruction-stats-alarms-timeline-tab
                :remediation-instruction="remediationInstruction"
                :interval="interval"
              />
            </v-tab-item>
            <v-tab-item>
              <v-flex
                offset-xs2
                xs8
              >
                <remediation-instruction-stats-rating-tab
                  :remediation-instruction="remediationInstruction"
                  :interval="interval"
                />
              </v-flex>
            </v-tab-item>
          </v-tabs-items>
        </v-card>
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import { entitiesRemediationInstructionStatsMixin } from '@/mixins/entities/remediation/instruction-stats';

import RemediationInstructionStatsSummaryTab from './remediation-instruction-stats-summary-tab.vue';
import RemediationInstructionStatsStatisticsTab from './remediation-instruction-stats-statistics-tab.vue';
import RemediationInstructionStatsAlarmsTimelineTab from './remediation-instruction-stats-alarms-timeline-tab.vue';
import RemediationInstructionStatsRatingTab from './remediation-instruction-stats-rating-tab.vue';

export default {
  components: {
    RemediationInstructionStatsSummaryTab,
    RemediationInstructionStatsStatisticsTab,
    RemediationInstructionStatsAlarmsTimelineTab,
    RemediationInstructionStatsRatingTab,
  },
  mixins: [entitiesRemediationInstructionStatsMixin],
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
  data() {
    return {
      activeTab: 0,
      remediationInstruction: {},
    };
  },
  watch: {
    remediationInstructionStatsItem: 'fetchRemediationInstructionStatsSummary',
  },
  mounted() {
    this.fetchRemediationInstructionStatsSummary();
  },
  methods: {
    async fetchRemediationInstructionStatsSummary() {
      this.remediationInstruction = await this.fetchRemediationInstructionStatsSummaryWithoutStore({
        id: this.remediationInstructionStatsItem._id,
      });
    },
  },
};
</script>
