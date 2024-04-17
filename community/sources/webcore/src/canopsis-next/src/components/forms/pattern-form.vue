<template>
  <v-layout column>
    <v-name-field
      v-field="form.title"
      :label="$t('common.title')"
      name="title"
      required
    />
    <c-alarm-patterns-field
      v-if="isAlarmPattern"
      v-field="form"
      :name="$constants.PATTERNS_FIELDS.alarm"
      :check-count-name="$constants.PATTERNS_FIELDS.alarm"
    />
    <c-entity-patterns-field
      v-else-if="isEntityPattern"
      v-field="form"
      :name="$constants.PATTERNS_FIELDS.entity"
      :check-count-name="$constants.PATTERNS_FIELDS.entity"
    />
    <c-pbehavior-patterns-field
      v-else-if="isPbehaviorPattern"
      v-field="form"
      :name="$constants.PATTERNS_FIELDS.pbehavior"
      :check-count-name="$constants.PATTERNS_FIELDS.pbehavior"
    />
    <c-service-weather-patterns-field
      v-else-if="isServiceWeatherPattern"
      v-field="form"
      :name="$constants.PATTERNS_FIELDS.serviceWeather"
      :check-count-name="$constants.PATTERNS_FIELDS.serviceWeather"
    />
  </v-layout>
</template>

<script>
import { computed } from 'vue';

import { PATTERN_TYPES } from '@/constants';

export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const isAlarmPattern = computed(() => props.form.type === PATTERN_TYPES.alarm);
    const isEntityPattern = computed(() => props.form.type === PATTERN_TYPES.entity);
    const isPbehaviorPattern = computed(() => props.form.type === PATTERN_TYPES.pbehavior);
    const isServiceWeatherPattern = computed(() => props.form.type === PATTERN_TYPES.serviceWeather);

    return {
      isAlarmPattern,
      isEntityPattern,
      isPbehaviorPattern,
      isServiceWeatherPattern,
    };
  },
};
</script>
