<template lang="pug">
  v-tooltip(top)
    template(#activator="{ on }")
      v-layout.alarm-counters.fill-height(v-on="on", column)
        alarm-counter(
          v-for="counter in preparedCounters",
          :key="counter.name",
          :count="counter.count",
          :color="counter.color",
          :icon="counter.icon"
        )
    v-layout(column)
      span(
        v-for="(counter, index) in preparedCounters",
        :key="counter.name"
      ) {{ countersMessages[index] }}
</template>

<script>
import { get } from 'lodash';

import { SERVICE_WEATHER_STATE_COUNTERS_COLORS, SERVICE_WEATHER_STATE_COUNTERS_ICONS } from '@/constants';

import AlarmCounter from './alarm-counter.vue';

export default {
  components: { AlarmCounter },
  props: {
    counters: {
      type: Object,
      default: () => ({}),
    },
    types: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    preparedCounters() {
      return this.types.map(type => ({
        type,
        icon: SERVICE_WEATHER_STATE_COUNTERS_ICONS[type],
        color: SERVICE_WEATHER_STATE_COUNTERS_COLORS[type],
        count: get(this.counters, type, 0),
      })).sort((a, b) => b.count - a.count);
    },

    countersMessages() {
      const messages = this.$t('serviceWeather.stateCountersTooltips');

      return this.preparedCounters.map(({ type, count }) => `${count} ${messages[type]}`);
    },
  },
};
</script>
