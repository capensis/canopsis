<template>
  <v-autocomplete
    v-validate="rules"
    v-field="value"
    :items="availableParameters"
    :label="label"
    :name="name"
    :multiple="isMultiple"
    :hide-details="hideDetails"
    :return-object="false"
    :error-messages="errors.collect(name)"
  >
    <template #selection="{ item, index }">
      <template v-if="isMultiple">
        <span v-if="!index">{{ getSelectionLabel(item) }}</span>
      </template>
      <template v-else>
        {{ item.text }}
      </template>
    </template>
    <template #item="{ parent, item, attrs, on }">
      <v-list-item
        v-bind="attrs"
        v-on="on"
      >
        <v-list-item-action v-if="isMultiple">
          <v-checkbox
            :input-value="attrs.value"
            :color="parent.color"
            :disabled="attrs.disabled"
          />
        </v-list-item-action>
        <v-list-item-content>
          <v-list-item-title>{{ item.text }}</v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </template>
  </v-autocomplete>
</template>

<script>
import { isArray, omit } from 'lodash';

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
    parameters: {
      type: Array,
      default: () => Object.values(omit(ALARM_METRIC_PARAMETERS, ['timeToAck', 'timeToResolve'])),
    },
    disabledParameters: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    isMultiple() {
      return isArray(this.value);
    },

    isMinValueLength() {
      return this.value?.length === this.min;
    },

    availableParameters() {
      return this.parameters.map(value => ({
        value,
        disabled: (this.disabledParameters.includes(value) && !this.isActiveValue(value))
          || (this.isMinValueLength && this.isActiveValue(value)),
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
    isActiveValue(value) {
      return this.isMultiple ? this.value.includes(value) : this.value === value;
    },

    getSelectionLabel(item) {
      if (this.isMinValueLength) {
        return item.text;
      }

      return this.$t('common.parametersToDisplay', { count: this.value.length });
    },
  },
};
</script>
