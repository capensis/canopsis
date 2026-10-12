<template>
  <c-tooltip :theme="$system.dark ? 'dark' : 'light'" location="top">
    <template #activator="{ props: tooltipProps }">
      <div class="alarm-counters fill-height d-flex flex-column" v-bind="tooltipProps">
        <alarm-counter
          v-for="counter in preparedCounters"
          :key="counter.type"
          :count="counter.count"
          :color="counter.color"
          :icon="counter.icon"
        />
      </div>
    </template>
    <div class="d-flex flex-column">
      <span v-for="(message, index) in countersMessages" :key="preparedCounters[index].type">
        {{ message }}
      </span>
    </div>
  </c-tooltip>
</template>

<script>
import { get } from 'lodash';

import { CSS_COLORS_VARS } from '@/config';
import {
  ALARM_STATES,
  SERVICE_WEATHER_STATE_COUNTERS,
  SERVICE_WEATHER_STATE_COUNTERS_ICONS,
} from '@/constants';

import { getAlarmStateColor } from '@/helpers/entities/alarm/color';

import AlarmCounter from './alarm-counter.vue';

export default {
  inject: ['$system'],

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
        [SERVICE_WEATHER_STATE_COUNTERS.minor]: getAlarmStateColor(ALARM_STATES.minor),
        [SERVICE_WEATHER_STATE_COUNTERS.major]: getAlarmStateColor(ALARM_STATES.major),
        [SERVICE_WEATHER_STATE_COUNTERS.critical]: getAlarmStateColor(ALARM_STATES.critical),
        [SERVICE_WEATHER_STATE_COUNTERS.paused]: CSS_COLORS_VARS.secondary,
        [SERVICE_WEATHER_STATE_COUNTERS.acknowledgedInactive]: CSS_COLORS_VARS.secondary,
        [SERVICE_WEATHER_STATE_COUNTERS.inactiveDepends]: CSS_COLORS_VARS.secondary,
      };
    },

    preparedCounters() {
      return this.types
        .map(type => ({
          type,
          icon: SERVICE_WEATHER_STATE_COUNTERS_ICONS[type],
          color: this.countersColorsByType[type],
          count: get(this.counters, type, 0),
        }))
        .sort((a, b) => b.count - a.count);
    },

    countersMessages() {
      const messages = this.$tm('serviceWeather.stateCountersTooltips');

      return this.preparedCounters.map(({ type, count }) => `${count} ${messages[type]}`);
    },
  },
};
</script>
