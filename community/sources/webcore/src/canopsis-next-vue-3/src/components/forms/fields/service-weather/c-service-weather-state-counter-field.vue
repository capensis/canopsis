<template>
  <c-select
    :model-value="modelValue"
    :items="availableCounters"
    :label="t('common.state')"
    :name="name"
    :disabled="disabled"
    :required="required"
    :item-props="getItemProps"
    chips
    closable-chips
    multiple
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { computed } from 'vue';

import {
  SERVICE_WEATHER_MAX_STATE_COUNTERS,
  SERVICE_WEATHER_STATE_COUNTERS,
} from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'counter',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    max: {
      type: Number,
      default: SERVICE_WEATHER_MAX_STATE_COUNTERS,
    },
  },

  emits: ['update:modelValue'],

  setup(props) {
    const { t, tm } = useI18n();

    const availableCounters = computed(() => {
      const messages = tm('serviceWeather.stateCounters');

      return Object.values(SERVICE_WEATHER_STATE_COUNTERS).map(value => ({
        value,
        title: messages[value],
      }));
    });

    const getItemProps = item => ({
      disabled: props.modelValue.includes(item.value)
        ? false
        : props.modelValue.length >= props.max,
    });

    return {
      t,
      availableCounters,
      getItemProps,
    };
  },
};
</script>
