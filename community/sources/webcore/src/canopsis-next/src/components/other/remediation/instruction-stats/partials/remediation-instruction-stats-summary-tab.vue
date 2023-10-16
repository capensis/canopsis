<template>
  <v-layout
    class="remediation-instruction-stats-summary"
    column
  >
    <remediation-instruction-stats-summary-row
      :label="$t('common.created')"
      :value="remediationInstruction.created | date"
    />
    <remediation-instruction-stats-summary-row
      :label="$t('common.lastModifiedOn')"
      :value="remediationInstruction.last_modified | date"
    />
    <remediation-instruction-stats-summary-row
      :label="$t('remediation.instructionStat.totalExecutions')"
      :value="remediationInstruction.execution_count"
    />
    <remediation-instruction-stats-summary-row
      :label="$t('remediation.instructionStat.successfulExecutions')"
      :value="remediationInstruction.successful"
    />
    <remediation-instruction-stats-summary-row
      :label="$t('remediation.instructionStat.lastExecutedOn')"
      :value="remediationInstruction.last_executed_on | date"
    />
    <remediation-instruction-stats-summary-row :label="$t('remediation.instructionStat.alarmStates')">
      <affect-alarm-states
        class="remediation-instruction-stats-summary--alarm-states"
        v-if="remediationInstruction.alarm_states"
        :alarm-states="remediationInstruction.alarm_states"
      />
    </remediation-instruction-stats-summary-row>
    <remediation-instruction-stats-summary-row :label="$t('remediation.instructionStat.okAlarmStates')">
      <c-state-count-changes-chip>{{ remediationInstruction.ok_alarm_states }}</c-state-count-changes-chip>
    </remediation-instruction-stats-summary-row>
  </v-layout>
</template>

<script>
import RemediationInstructionStatsSummaryRow from './remediation-instruction-stats-summary-row.vue';
import AffectAlarmStates from './affect-alarm-states.vue';

export default {
  components: { RemediationInstructionStatsSummaryRow, AffectAlarmStates },
  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

<style lang="scss">
.remediation-instruction-stats-summary {
  &--alarm-states {
    width: 250px;
  }
}
</style>
