<template lang="pug">
  v-layout.remediation-instruction-stats-information(column)
    remediation-instruction-stats-information-row(
      :label="$t('common.created')",
      :value="remediationInstruction.created | date('long', true)"
    )
    remediation-instruction-stats-information-row(
      :label="$t('remediationInstructionStats.lastModifiedOn')",
      :value="remediationInstruction.last_modified | date('long', true)"
    )
    remediation-instruction-stats-information-row(
      :label="$t('remediationInstructionStats.executionCount')",
      :value="remediationInstruction.execution_count"
    )
    remediation-instruction-stats-information-row(
      :label="$t('remediationInstructionStats.lastExecutedOn')",
      :value="remediationInstruction.last_executed_on | date('long', true)"
    )
    remediation-instruction-stats-information-row(:label="$t('remediationInstructionStats.alarmStates')")
      affect-alarm-states.remediation-instruction-stats-information--alarm-states(
        v-if="remediationInstruction.alarm_states",
        :alarm-states="remediationInstruction.alarm_states"
      )
    remediation-instruction-stats-information-row(:label="$t('remediationInstructionStats.okAlarmStates')")
      span.c-state-count-changes-chip.primary {{ remediationInstruction.ok_alarm_states }}
</template>

<script>
import RemediationInstructionStatsInformationRow from './remediation-instruction-stats-information-row.vue';
import AffectAlarmStates from './affect-alarm-states.vue';

export default {
  components: { RemediationInstructionStatsInformationRow, AffectAlarmStates },
  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

<style lang="scss">
.remediation-instruction-stats-information {
  &--alarm-states {
    width: 250px;
  }
}
</style>
