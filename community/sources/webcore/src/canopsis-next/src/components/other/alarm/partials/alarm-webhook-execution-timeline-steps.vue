<template>
  <div class="alarm-webhook-execution-timeline-steps">
    <ul
      v-for="(groupSteps, day) in groupedSteps"
      :key="day"
    >
      <li
        v-for="(step, index) in groupSteps"
        :key="index"
      >
        <div
          v-show="index === 0"
          class="alarm-webhook-execution-timeline-steps-date text--disabled"
        >
          <div class="date">
            {{ day }}
          </div>
        </div>
        <div class="alarm-webhook-execution-timeline-steps-card">
          <div class="time text--disabled">
            {{ step.t | date('time') }}
          </div>
          <alarm-webhook-execution-timeline-flag
            :step="step"
            class="flag"
          />
          <slot
            :step="step"
            name="card"
          />
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import { groupAlarmSteps } from '@/helpers/entities/alarm/step/list';

import AlarmWebhookExecutionTimelineFlag from './alarm-webhook-execution-timeline-flag.vue';

export default {
  components: { AlarmWebhookExecutionTimelineFlag },
  props: {
    steps: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    groupedSteps() {
      return groupAlarmSteps(this.steps);
    },
  },
};
</script>

<style lang="scss" scoped>
$borderLine: #DDDDE0;

.alarm-webhook-execution-timeline-steps {
  ul {
    list-style: none;
  }

  &-date {
    padding: 3em 2em 0;
    position: relative;

    .date {
      top: 4px;
      left: -11px;
      position: absolute;
    }

    &:before, &:after {
      content: '';
      position: absolute;
      left: -2px;
      width: 2px;
      background-color: $borderLine;
    }

    &:before {
      top: 0;
      height: 4px;
    }

    &:after {
      top: 24px;
      bottom: 0;
    }
  }

  &-card {
    padding: 3em 2em 0;
    position: relative;

    .time {
      position: absolute;
      left: 2em;
      top: 9px;
      display: block;
      font-size: 11px;
    }

    .flag {
      height: 30px;
      top: 0;
      left: -13px;
      position: absolute;
      display: flex;
      align-items: center;
    }

    &:before, &:after {
      content: '';
      position: absolute;
      left: -2px;
      width: 2px;
      background-color: $borderLine;
    }

    &:after {
      top: 30px;
      bottom: 0;
    }
  }

  ul:last-of-type li:last-child &-card:after {
    background-color: unset;
    background-image: linear-gradient(
        to bottom,
        $borderLine 60%,
        transparent
    );
  }
}
</style>
