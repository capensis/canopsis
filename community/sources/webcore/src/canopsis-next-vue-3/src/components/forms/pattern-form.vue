<template>
  <c-layout column>
    <c-name-field
      v-model="title"
      :label="$t('common.title')"
      name="title"
      autofocus
      required
    />
    <c-alarm-patterns-field
      v-if="isAlarmPattern"
      :patterns="formModel"
      :name="PATTERNS_FIELDS.alarm"
      :check-count-name="PATTERNS_FIELDS.alarm"
      required
      @update:patterns="formModel = $event"
    />
    <c-entity-patterns-field
      v-else-if="isEntityPattern"
      :patterns="formModel"
      :name="PATTERNS_FIELDS.entity"
      :check-count-name="PATTERNS_FIELDS.entity"
      required
      @update:patterns="formModel = $event"
    />
    <c-pbehavior-patterns-field
      v-else-if="isPbehaviorPattern"
      :patterns="formModel"
      :name="PATTERNS_FIELDS.pbehavior"
      :check-count-name="PATTERNS_FIELDS.pbehavior"
      required
      @update:patterns="formModel = $event"
    />
    <c-service-weather-patterns-field
      v-else-if="isServiceWeatherPattern"
      :patterns="formModel"
      :name="PATTERNS_FIELDS.serviceWeather"
      :check-count-name="PATTERNS_FIELDS.serviceWeather"
      required
      @update:patterns="formModel = $event"
    />
  </c-layout>
</template>

<script>
import { computed } from 'vue';

import { PATTERNS_FIELDS, PATTERN_TYPES } from '@/constants';
import { useFormContext, useFormModels } from '@/hooks/form';

export default {
  setup() {
    const form = useFormContext();
    const { title, type } = useFormModels(['title', 'type']);

    const formModel = computed({
      get: () => form.values.value,
      set: (value) => {
        form.values.value = value;
      },
    });

    const isAlarmPattern = computed(() => type.value === PATTERN_TYPES.alarm);
    const isEntityPattern = computed(() => type.value === PATTERN_TYPES.entity);
    const isPbehaviorPattern = computed(() => type.value === PATTERN_TYPES.pbehavior);
    const isServiceWeatherPattern = computed(() => type.value === PATTERN_TYPES.serviceWeather);

    return {
      PATTERNS_FIELDS,
      title,
      formModel,
      isAlarmPattern,
      isEntityPattern,
      isPbehaviorPattern,
      isServiceWeatherPattern,
    };
  },
};
</script>
