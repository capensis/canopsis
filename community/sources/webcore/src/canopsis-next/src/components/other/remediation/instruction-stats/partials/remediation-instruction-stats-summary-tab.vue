<template lang="pug">
  v-layout.remediation-instruction-stats-summary(column)
    remediation-instruction-stats-summary-row(
      :label="$t('common.created')",
      :value="remediationInstruction.created | date('long', true)"
    )
    remediation-instruction-stats-summary-row(
      :label="$t('remediationInstructionStats.lastModifiedOn')",
      :value="remediationInstruction.last_modified | date('long', true)"
    )
    remediation-instruction-stats-summary-row(
      :label="$t('remediationInstructionStats.executionCount')",
      :value="remediationInstruction.execution_count"
    )
    remediation-instruction-stats-summary-row(
      :label="$t('remediationInstructionStats.lastExecutedOn')",
      :value="remediationInstruction.last_executed_on | date('long', true)"
    )
    remediation-instruction-stats-summary-row(:label="$t('remediationInstructionStats.alarmStates')")
      affect-alarm-states.remediation-instruction-stats-summary--alarm-states(
        v-if="remediationInstruction.alarm_states",
        :alarm-states="remediationInstruction.alarm_states"
      )
    remediation-instruction-stats-summary-row(:label="$t('remediationInstructionStats.okAlarmStates')")
      span.c-state-count-changes-chip.primary {{ remediationInstruction.ok_alarm_states }}
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
