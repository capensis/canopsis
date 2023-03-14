<template lang="pug">
  v-select(
    v-field="value",
    v-validate="rules",
    :items="availableParameters",
    :label="label",
    :name="name",
    :multiple="isMultiple",
    :hide-details="hideDetails"
  )
    template(v-if="isMultiple", #selection="{ item, index }")
      span(v-if="!index") {{ getSelectionLabel(item) }}
</template>

<script>
import { isArray } from 'lodash';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Array],
      required: true,
    },
    name: {
      type: String,
      default: 'parameters',
    },
    label: {
      type: String,
      required: false,
    },
    min: {
      type: Number,
      default: 1,
    },
    required: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    isMultiple() {
      return isArray(this.value);
    },

    isMinValueLength() {
      return this.value.length === this.min;
    },

    availableParameters() {
      return Object.values(ALARM_METRIC_PARAMETERS).map(value => ({
        value,
        disabled: this.isMinValueLength && this.value.includes(value),
        text: this.$t(`alarm.metrics.${value}`),
      }));
    },

    rules() {
      return {
        required: this.required,
      };
    },
  },
  methods: {
    getSelectionLabel(item) {
      if (this.isMinValueLength) {
        return item.text;
      }

      return this.$t('common.parametersToDisplay', { count: this.value.length });
    },
  },
};
</script>
