<template lang="pug">
  v-layout.alarm-counters(column, justify-space-around)
    alarm-counter(:count="counters.paused", icon="pause")
    alarm-counter(:count="counters.active", icon="build")
    alarm-counter(:count="counters.inactive", icon="brightness_3")
    alarm-counter(:count="counters.maintenanced", icon="more_horiz")
</template>

<script>
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
};
</script>

<style lang="scss">
  .alarm-counters {
    position: relative;
    background-color: rgba(0, 0, 0, .2);
  }
</style>
