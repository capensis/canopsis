<template>
  <c-layout class="remediation-instruction-stats-summary" column>
    <remediation-instruction-stats-summary-row
      :label="t('common.created')"
      :value="dateToString(remediationInstruction.created)"
    />
    <remediation-instruction-stats-summary-row
      :label="t('common.lastModifiedOn')"
      :value="dateToString(remediationInstruction.last_modified)"
    />
    <remediation-instruction-stats-summary-row
      :label="t('remediation.instructionStat.totalExecutions')"
      :value="remediationInstruction.execution_count"
    />
    <remediation-instruction-stats-summary-row
      :label="t('remediation.instructionStat.successfulExecutions')"
      :value="remediationInstruction.successful"
    />
    <remediation-instruction-stats-summary-row
      :label="t('remediation.instructionStat.lastExecutedOn')"
      :value="dateToString(remediationInstruction.last_executed_on)"
    />
    <remediation-instruction-stats-summary-row :label="t('remediation.instructionStat.alarmStates')">
      <affect-alarm-states
        v-if="remediationInstruction.alarm_states"
        :alarm-states="remediationInstruction.alarm_states"
        class="remediation-instruction-stats-summary--alarm-states"
      />
    </remediation-instruction-stats-summary-row>
    <remediation-instruction-stats-summary-row :label="t('remediation.instructionStat.okAlarmStates')">
      <c-state-count-changes-chip>{{ remediationInstruction.ok_alarm_states }}</c-state-count-changes-chip>
    </remediation-instruction-stats-summary-row>
  </c-layout>
</template>

<script>
import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

import AffectAlarmStates from './affect-alarm-states.vue';
import RemediationInstructionStatsSummaryRow from './remediation-instruction-stats-summary-row.vue';

export default {
  components: { AffectAlarmStates, RemediationInstructionStatsSummaryRow },

  props: {
    remediationInstruction: {
      type: Object,
      default: () => ({}),
    },
  },

  setup() {
    const { t } = useI18n();
    const dateToString = value => convertDateToString(value);

    return {
      t,
      dateToString,
    };
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
