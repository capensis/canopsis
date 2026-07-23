<template>
  <v-autocomplete
    :model-value="value"
    :items="availableParameters"
    :label="label || $t('kpi.metrics.parameter')"
    :name="name"
    :multiple="isMultiple"
    :hide-details="hideDetails"
    :return-object="false"
    item-title="text"
    item-value="value"
    @update:model-value="$emit('update:value', $event)"
  >
    <template #selection="{ item, index }">
      <template v-if="isMultiple">
        <span v-if="!index">
          {{ getSelectionLabel(item) }}
        </span>
      </template>
      <template v-else>
        {{ getItemText(item) }}
      </template>
    </template>
  </v-autocomplete>
</template>

<script>
import { isArray, omit } from 'lodash';

import { ALARM_METRIC_PARAMETERS } from '@/constants';

export default {
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
      default: '',
    },
    min: {
      type: Number,
      default: 1,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    parameters: {
      type: Array,
      default: () => Object.values(omit(ALARM_METRIC_PARAMETERS, ['timeToAck', 'timeToResolve'])),
    },
    disabledParameters: {
      type: Array,
      default: () => [],
    },
  },

  emits: ['update:value'],

  computed: {
    isMultiple() {
      return isArray(this.value);
    },

    isMinValueLength() {
      return this.value?.length === this.min;
    },

    availableParameters() {
      return this.parameters.map((value) => ({
        value,
        disabled:
          (this.disabledParameters.includes(value) && !this.isActiveValue(value)) ||
          (this.isMinValueLength && this.isActiveValue(value)),
        text: this.$t(`alarm.metrics.${value}`),
      }));
    },
  },

  methods: {
    isActiveValue(value) {
      return this.isMultiple ? this.value.includes(value) : this.value === value;
    },

    getSelectionLabel(item) {
      if (this.isMinValueLength) {
        return this.getItemText(item);
      }

      return this.$t('common.parametersToDisplay', { count: this.value.length });
    },

    getItemText(item) {
      const preparedItem = item?.raw ?? item;
      const value = preparedItem?.value ?? preparedItem?.props?.value ?? preparedItem;
      const parameter = this.availableParameters.find(
        ({ value: parameterValue }) => parameterValue === value,
      );

      return parameter?.text ?? preparedItem?.text ?? preparedItem?.title ?? '';
    },
  },
};
</script>
