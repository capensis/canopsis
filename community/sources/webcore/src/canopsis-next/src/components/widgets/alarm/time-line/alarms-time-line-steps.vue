<template>
  <div class="alarms-timeline-steps">
    <ul
      v-for="(groupSteps, day) in groupedSteps"
      :key="day"
    >
      <li
        v-for="(step, index) in groupSteps"
        :key="index"
      >
        <div
          class="alarms-timeline-steps-date text--disabled"
          v-show="index === 0"
        >
          <div class="date">
            {{ day }}
          </div>
        </div>
        <div class="alarms-timeline-steps-card">
          <div class="time text--disabled">
            {{ step.t | date('time') }}
          </div>
          <alarms-time-line-flag
            class="flag"
            :step="step"
          />
          <slot
            name="card"
            :step="step"
          />
        </div>
      </li>
    </ul>
  </div>
</template>

<script>
import { groupAlarmSteps } from '@/helpers/entities/alarm/list';

import AlarmsTimeLineFlag from './alarms-time-line-flag.vue';

export default {
  components: { AlarmsTimeLineFlag },
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

.alarms-timeline-steps {
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
