<template>
  <span class="c-horizontal-time-line-card">
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
            v-if="isStepTypeState"
            :value="step.val"
          />
          <v-icon
            v-else
            :style="{ color: style.icon }"
          >
            {{ style.icon }}
          </v-icon>
          <span class="c-horizontal-time-line-card__time">{{ step.t | date('time') }}</span>
        </v-layout>
      </template>
      <div class="pre-line">
        <div><strong>{{ stepTitle }}</strong></div>
        <div v-html="sanitizedStepMessage" />
      </div>
    </v-tooltip>
  </span>
</template>

<script>
import { sanitizeHtml, linkifyHtml } from '@/helpers/html';
import { formatStep } from '@/helpers/entities/entity/formatting';

import { widgetExpandPanelAlarmTimelineCard } from '@/mixins/widget/expand-panel/alarm/timeline-card';

export default {
  mixins: [widgetExpandPanelAlarmTimelineCard],
  props: {
    step: {
      type: Object,
      required: true,
    },
  },
  computed: {
    style() {
      return formatStep(this.step);
    },

    sanitizedStepMessage() {
      return sanitizeHtml(linkifyHtml(String(this.step?.m ?? '')));
    },
  },
};
</script>

<style lang="scss" scoped>
.c-horizontal-time-line-card {
  display: inline-flex;
  align-items: center;

  &__time {
    font-size: 10px;
  }
}
</style>
