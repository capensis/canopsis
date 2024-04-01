<template>
  <div v-if="!isStateCounter">
    <c-alarm-chip
      v-if="isChangeState"
      :value="step.val"
    />
    <v-icon
      v-else-if="isChangeStatus"
      :color="status.color"
    >
      {{ status.icon }}
    </v-icon>
    <c-alarm-pbehavior-chip
      v-else-if="isPbehavior"
      :color="step.color"
      :icon="step.icon_name"
    />
    <v-icon
      v-else-if="isDefaultStepIcon"
      :color="stepColor"
    >
      {{ stepIcon }}
    </v-icon>
    <c-alarm-extra-details-chip
      v-else-if="isJunit"
      :color="stepColor"
    >
      <span class="white--text">{{ $t('alarm.timeline.junit') }}</span>
    </c-alarm-extra-details-chip>
  </div>
</template>

<script>
import { computed } from 'vue';

import { ALARM_LIST_STEPS } from '@/constants';

import { formatAlarmStatus } from '@/helpers/entities/alarm/formatting';
import { getAlarmStepIcon, getAlarmStepColor } from '@/helpers/entities/alarm/step/formatting';

export default {
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const isStateCounter = computed(() => props.step._t === ALARM_LIST_STEPS.stateCounter);
    const isChangeState = computed(() => props.step._t === ALARM_LIST_STEPS.changeState);

    const isChangeStatus = computed(() => [
      ALARM_LIST_STEPS.statusinc,
      ALARM_LIST_STEPS.statusdec,
    ].includes(props.step._t));

    const isJunit = computed(() => [
      ALARM_LIST_STEPS.junitTestSuiteUpdate,
      ALARM_LIST_STEPS.junitTestCaseUpdate,
    ].includes(props.step._t));

    const isPbehavior = computed(() => [
      ALARM_LIST_STEPS.pbhenter,
      ALARM_LIST_STEPS.pbhleave,
    ].includes(props.step._t));

    const isDefaultStepIcon = computed(() => [
      ALARM_LIST_STEPS.resolve,
      ALARM_LIST_STEPS.activate,
      ALARM_LIST_STEPS.ack,
      ALARM_LIST_STEPS.ackRemove,
      ALARM_LIST_STEPS.snooze,
      ALARM_LIST_STEPS.unsnooze,
      ALARM_LIST_STEPS.comment,
      ALARM_LIST_STEPS.metaalarmattach,
      ALARM_LIST_STEPS.assocTicket,
      ALARM_LIST_STEPS.webhookStart,
      ALARM_LIST_STEPS.webhookProgress,
      ALARM_LIST_STEPS.webhookComplete,
      ALARM_LIST_STEPS.webhookFail,
      ALARM_LIST_STEPS.declareTicket,
      ALARM_LIST_STEPS.declareTicketFail,
      ALARM_LIST_STEPS.declareTicketRuleProgress,
      ALARM_LIST_STEPS.declareTicketRuleComplete,
      ALARM_LIST_STEPS.declareTicketRuleFailed,
      ALARM_LIST_STEPS.instructionStart,
      ALARM_LIST_STEPS.instructionPause,
      ALARM_LIST_STEPS.instructionResume,
      ALARM_LIST_STEPS.instructionComplete,
      ALARM_LIST_STEPS.instructionAbort,
      ALARM_LIST_STEPS.instructionFail,
      ALARM_LIST_STEPS.autoInstructionStart,
      ALARM_LIST_STEPS.autoInstructionComplete,
      ALARM_LIST_STEPS.autoInstructionFail,
    ].includes(props.step._t));

    const status = computed(() => formatAlarmStatus(props.step.val));
    const stepIcon = computed(() => getAlarmStepIcon(props.step._t));
    const stepColor = computed(() => getAlarmStepColor(props.step._t));

    return {
      isStateCounter,
      isChangeState,
      isChangeStatus,
      isJunit,
      isPbehavior,
      isDefaultStepIcon,
      status,
      stepIcon,
      stepColor,
    };
  },
};
</script>
