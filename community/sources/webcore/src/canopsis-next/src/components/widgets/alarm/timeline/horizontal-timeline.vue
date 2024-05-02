<template>
  <div class="c-horizontal-time-line">
    <div
      v-for="(group, groupIndex) in groupedSteps"
      :key="group.day"
      class="c-horizontal-time-line__groups"
    >
      <span class="c-horizontal-time-line__day">{{ group.day }}</span>
      <v-divider
        class="grey mr-2"
        vertical
      />
      <div class="c-horizontal-time-line__cards">
        <template v-for="(step, stepIndex) in group.steps">
          <horizontal-timeline-card
            :key="stepIndex"
            :step="step"
          />
          <v-icon
            v-if="groupIndex !== groupedSteps.length - 1 || stepIndex !== group.steps.length - 1"
            :key="`arrow-${stepIndex}`"
            class="mx-2"
            size="16"
          >
            arrow_forward
          </v-icon>
        </template>
      </div>
    </div>
  </div>
</template>

<script>
import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';

import HorizontalTimelineCard from './horizontal-timeline-card.vue';

export default {
  components: { HorizontalTimelineCard },
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  computed: {
    groupedSteps() {
      return Object.entries(groupAlarmSteps(this.alarm.v.steps)).map(([day, steps]) => ({
        day,
        steps,
      }));
    },
  },
};
</script>

<style lang="scss" scoped>
  $timeFontSize: 10px;
  $dayPadding: 5px;

  .c-horizontal-time-line {
    display: inline-flex;
    position: relative;
    padding-top: $timeFontSize + $dayPadding;

    &__day {
      font-size: $timeFontSize;
      position: absolute;
      top: 0;
    }

    &__groups {
      display: flex;
    }

    &__cards {
      display: flex;
    }
  }
</style>
