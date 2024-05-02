<template>
  <span class="horizontal-time-line-card">
    <v-tooltip
      :disabled="!step.m"
      top
    >
      <template #activator="{ on }">
        <v-layout
          align-center
          column
          v-on="on"
        >
          <c-alarm-chip
            v-if="isChangeState"
            :value="step.val"
          />
          <v-icon
            v-else
            :color="style.color"
          >
            {{ style.icon }}
          </v-icon>
          <span class="horizontal-time-line-card__time">{{ step.t | date('time') }}</span>
        </v-layout>
      </template>
      <div class="pre-line">
        <alarm-timeline-step-title :step="step" class="text-subtitle-1" />
        <div v-html="sanitizedStepMessage" />
      </div>
    </v-tooltip>
  </span>
</template>

<script>
import { computed } from 'vue';

import { sanitizeHtml, linkifyHtml } from '@/helpers/html';
import { formatNotificationAlarmStep } from '@/helpers/entities/alarm/step/formatting';
import { isStateStepType } from '@/helpers/entities/alarm/step/entity';

import AlarmTimelineStepTitle from './alarm-timeline-step-title.vue';

export default {
  components: { AlarmTimelineStepTitle },
  props: {
    step: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const isChangeState = computed(() => isStateStepType(props.step._t));
    const style = computed(() => formatNotificationAlarmStep(props.step));
    const sanitizedStepMessage = computed(() => sanitizeHtml(linkifyHtml(String(props.step?.m ?? ''))));

    return {
      style,
      isChangeState,
      sanitizedStepMessage,
    };
  },
};
</script>

<style lang="scss" scoped>
.horizontal-time-line-card {
  display: inline-flex;
  align-items: center;

  &__time {
    font-size: 10px;
  }
}
</style>
