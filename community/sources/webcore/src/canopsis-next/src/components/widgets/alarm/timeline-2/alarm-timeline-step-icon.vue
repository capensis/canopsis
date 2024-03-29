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
    <v-chip v-else-if="isJunit" :color="stepColor" small>
      {{ $t('alarm.timeline.junit') }}
    </v-chip>
    <v-chip v-else-if="isPbehavior" :color="pbehaviorColor" small>
      <v-icon :color="pbehaviorIconColor" />
    </v-chip>
    <v-icon
      v-else-if="isDefaultStepIcon"
      :color="stepColor"
    >
      {{ stepIcon }}
    </v-icon>
  </div>
</template>

<script>
import { ALARM_LIST_STEPS } from '@/constants';

import { formatAlarmStatus } from '@/helpers/entities/alarm/formatting';
import { getAlarmStepIcon, getAlarmStepColor } from '@/helpers/entities/alarm/step/formatting';
import { getPbehaviorColor } from '@/helpers/entities/pbehavior/form';
import { getMostReadableTextColor } from '@/helpers/color';

export default {
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    isStateCounter() {
      return this.step._t === ALARM_LIST_STEPS.stateCounter;
    },

    isChangeState() {
      return this.step._t === ALARM_LIST_STEPS.changeState;
    },

    isChangeStatus() {
      return [ALARM_LIST_STEPS.statusinc, ALARM_LIST_STEPS.statusdec].includes(this.step._t);
    },

    isJunit() {
      return [ALARM_LIST_STEPS.junitTestSuiteUpdate, ALARM_LIST_STEPS.junitTestCaseUpdate].includes(this.step._t);
    },

    status() {
      return formatAlarmStatus(this.step.val);
    },

    isDefaultStepIcon() {
      return [
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
        ALARM_LIST_STEPS.ticketDeclarationRuleProgress,
        ALARM_LIST_STEPS.ticketDeclarationRuleCompleted,
        ALARM_LIST_STEPS.ticketDeclarationRuleFailed,
        ALARM_LIST_STEPS.instructionStart,
        ALARM_LIST_STEPS.instructionPause,
        ALARM_LIST_STEPS.instructionResume,
        ALARM_LIST_STEPS.instructionComplete,
        ALARM_LIST_STEPS.instructionAbort,
        ALARM_LIST_STEPS.instructionFail,
        ALARM_LIST_STEPS.autoInstructionStart,
        ALARM_LIST_STEPS.autoInstructionComplete,
        ALARM_LIST_STEPS.autoInstructionFail,
      ].includes(this.step._t);
    },

    isPbehavior() {
      return [ALARM_LIST_STEPS.pbhenter, ALARM_LIST_STEPS.pbhleave].includes(this.step._t);
    },

    stepIcon() {
      return getAlarmStepIcon(this.step._t);
    },

    stepColor() {
      return getAlarmStepColor(this.step._t);
    },

    pbehaviorColor() {
      return getPbehaviorColor(this.pbehaviorInfo);
    },

    pbehaviorIconColor() {
      return getMostReadableTextColor(this.color, { level: 'AA', size: 'large' });
    },
  },
};
</script>
