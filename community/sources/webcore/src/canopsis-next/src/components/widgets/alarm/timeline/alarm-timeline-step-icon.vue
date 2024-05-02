<template>
  <div v-if="!isStateCounter">
    <c-alarm-state-chip
      v-if="isState"
      :value="step.val"
      :type="step._t"
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

import { formatAlarmStatus } from '@/helpers/entities/alarm/formatting';
import { getAlarmStepIcon, getAlarmStepColor } from '@/helpers/entities/alarm/step/formatting';
import {
  isStateCounterStepType,
  isChangeStatusStepType,
  isResolveStepType,
  isActivateStepType,
  isCommentStepType,
  isAssocTicketStepType,
  isAckStepType,
  isSnoozeStepType,
  isStateStepType,
  isJunitStepType,
  isPbehaviorStepType,
  isInstructionStepType,
  isAutoInstructionStepType,
  isDeclareTicketStepType,
  isWebhookStepType,
  isMetaAlarmStepType,
} from '@/helpers/entities/alarm/step/entity';

export default {
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const isStateCounter = computed(() => isStateCounterStepType(props.step._t));
    const isState = computed(() => isStateStepType(props.step._t));
    const isChangeStatus = computed(() => isChangeStatusStepType(props.step._t));
    const isJunit = computed(() => isJunitStepType(props.step._t));
    const isPbehavior = computed(() => isPbehaviorStepType(props.step._t));

    const isDefaultStepIcon = computed(() => (
      isResolveStepType(props.step._t)
      || isActivateStepType(props.step._t)
      || isAckStepType(props.step._t)
      || isSnoozeStepType(props.step._t)
      || isCommentStepType(props.step._t)
      || isMetaAlarmStepType(props.step._t)
      || isAssocTicketStepType(props.step._t)
      || isWebhookStepType(props.step._t)
      || isDeclareTicketStepType(props.step._t)
      || isInstructionStepType(props.step._t)
      || isAutoInstructionStepType(props.step._t)
    ));

    const status = computed(() => formatAlarmStatus(props.step.val));
    const stepIcon = computed(() => getAlarmStepIcon(props.step._t));
    const stepColor = computed(() => getAlarmStepColor(props.step._t));

    return {
      isStateCounter,
      isState,
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
