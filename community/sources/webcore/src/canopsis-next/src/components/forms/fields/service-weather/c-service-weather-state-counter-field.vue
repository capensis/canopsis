<template>
  <v-select
    v-field="value"
    v-validate="rules"
    :label="$t('common.state')"
    :items="availableCounters"
    :name="name"
    :error-messages="errors.collect(name)"
    :item-disabled="isItemDisabled"
    :disabled="disabled"
    chips
    deletable-chips
    small-chips
    multiple
  />
</template>

<script>
import { SERVICE_WEATHER_STATE_COUNTERS, SERVICE_WEATHER_MAX_STATE_COUNTERS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },

    availableCounters() {
      const messages = this.$t('serviceWeather.stateCounters');

      return Object.values(SERVICE_WEATHER_STATE_COUNTERS).map(value => ({
        value,
        text: messages[value],
      }));
    },
  },
  methods: {
    isItemDisabled({ value }) {
      return this.value.includes(value) ? false : this.value.length >= this.max;
    },
  },
};
</script>
