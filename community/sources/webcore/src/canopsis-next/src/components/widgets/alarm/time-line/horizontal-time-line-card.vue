<template lang="pug">
  span.c-horizontal-time-line-card
    v-tooltip(:disabled="!step.m", top)
      template(#activator="{ on }")
        v-layout(v-on="on", align-center, column)
          c-alarm-chip(v-if="isStepTypeState", :value="step.val")
          v-icon(v-else, :style="{ color: style.icon }") {{ style.icon }}
          span.c-horizontal-time-line-card__time {{ step.t | date('time') }}
      div.pre-line
        div
          strong {{ stepTitle }}
        div(v-html="step.m")
</template>

<script>
import { formatStep } from '@/helpers/formatting';

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
