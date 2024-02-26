<template>
  <v-tooltip top>
    <template #activator="{ on }">
      <v-layout
        class="alarm-counters fill-height"
        column
        v-on="on"
      >
        <alarm-counter
          v-for="counter in preparedCounters"
          :key="counter.name"
          :count="counter.count"
          :color="counter.color"
          :icon="counter.icon"
        />
      </v-layout>
    </template>
    <v-layout column>
      <span
        v-for="(counter, index) in preparedCounters"
        :key="counter.name"
      >{{ countersMessages[index] }}</span>
    </v-layout>
  </v-tooltip>
</template>

<script>
import { get } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import { ENTITIES_STATES, SERVICE_WEATHER_STATE_COUNTERS, SERVICE_WEATHER_STATE_COUNTERS_ICONS } from '@/constants';

import { getEntityStateColor } from '@/helpers/entities/entity/color';

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
    countersColorsByType() {
      return {
        [SERVICE_WEATHER_STATE_COUNTERS.all]: CSS_COLORS_VARS.error,
        [SERVICE_WEATHER_STATE_COUNTERS.active]: CSS_COLORS_VARS.error,
        [SERVICE_WEATHER_STATE_COUNTERS.minor]: getEntityStateColor(ENTITIES_STATES.minor),
        [SERVICE_WEATHER_STATE_COUNTERS.major]: getEntityStateColor(ENTITIES_STATES.major),
        [SERVICE_WEATHER_STATE_COUNTERS.critical]: getEntityStateColor(ENTITIES_STATES.critical),
      };
    },

    preparedCounters() {
      return this.types.map(type => ({
        type,
        icon: SERVICE_WEATHER_STATE_COUNTERS_ICONS[type],
        color: this.countersColorsByType[type],
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
