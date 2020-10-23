<template lang="pug">
  v-layout.alarm-counters(column, justify-space-around)
    alarm-counter(:count="alarmCounters.pause", icon="pause")
    alarm-counter(:count="alarmCounters.build", icon="build")
    alarm-counter(:count="alarmCounters.brightness", icon="brightness_3")
    alarm-counter(:count="alarmCounters.other", icon="more_horiz")
</template>

<script>
import { sum } from 'lodash';

import AlarmCounter from './alarm-counter.vue';

export default {
  components: {
    AlarmCounter,
  },
  props: {
    counters: {
      type: Object,
      required: true,
    },
    orientation: {
      type: String,
      default: 'vertical',
    },
  },
  computed: {
    alarmCounters() {
      const {
        paused: pause = 0,
        maintenanced: build = 0,
        inactive: brightness = 0,
        ...restCounters
      } = this.counters;

      return {
        pause,
        build,
        brightness,
        other: sum(Object.values(restCounters)),
      };
    },
  },
};
</script>

<style lang="scss">
  .alarm-counters {
    position: relative;
    background-color: rgba(0, 0, 0, .2);
  }
</style>
