<template lang="pug">
  v-select(
    v-field="value",
    :items="availableParameters",
    :name="name",
    hide-details,
    multiple
  )
    template(#selection="{ item, index }")
      span(v-if="!index") {{ getSelectionLabel(item) }}
</template>

<script>
import { ALARM_METRIC_PARAMETERS } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      required: true,
    },
    name: {
      type: String,
      default: 'parameters',
    },
    min: {
      type: Number,
      default: 1,
    },
  },
  computed: {
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
