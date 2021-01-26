<template lang="pug">
  v-layout.alarm-counters(column, justify-space-around)
    v-tooltip(
      v-for="counter in alarmCounters.selected",
      :key="counter.key",
      top
    )
      alarm-counter(
        :count="counter.count",
        :icon="counter.icon",
        slot="activator"
      )
      span {{ counter.name }}
    v-tooltip(top, :disabled="!otherCountersValue")
      alarm-counter(
        :count="otherCountersValue",
        icon="more_horiz",
        slot="activator"
      )
      div(v-for="otherCounter in alarmCounters.other", :key="otherCounter.key")
        strong {{ otherCounter.name }}
        span : {{ otherCounter.count }}
</template>

<script>
import AlarmCounter from './alarm-counter.vue';

export default {
  components: {
    AlarmCounter,
  },
  props: {
    selectedTypes: {
      type: Array,
      required: true,
    },
    counters: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    alarmCounters() {
      return this.counters.reduce((acc, { count, type }) => {
        if (this.selectedTypes.includes(type._id)) {
          acc.selected.push({
            key: type._id,
            name: type.name,
            icon: type.icon_name,
            count,
          });
        } else {
          acc.other.push({
            key: type._id,
            name: type.name,
            count,
          });
        }

        return acc;
      }, {
        selected: [],
        other: [],
      });
    },

    otherCountersValue() {
      return this.alarmCounters.other.reduce((acc, { count }) => acc + count, 0);
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
