<template>
  <v-row class="ga-2 flex-nowrap" no-gutters>
    <v-col :cols="sameWidth ? 6 : 7">
      <c-number-field
        :value="duration.value"
        :label="label || $t('common.duration')"
        :disabled="disabled"
        :name="valueFieldName"
        :required="isRequired"
        :min="min"
        :max="max"
        :hide-details="hideDetails"
        :autofocus="autofocus"
        @update:value="updateValue"
      >
        <template v-if="$slots.append" #append>
          <slot name="append" />
        </template>
      </c-number-field>
    </v-col>

    <v-col :cols="sameWidth ? 6 : 5">
      <c-select
        :model-value="duration.unit"
        :items="availableUnits"
        :disabled="disabled"
        :label="unitsLabel"
        :name="unitFieldName"
        :required="isRequired"
        :clearable="clearable"
        :hide-details="hideDetails"
        item-title="text"
        item-value="value"
        @update:model-value="updateUnit"
      />
    </v-col>
  </v-row>
</template>

<script>
import { isNumber } from 'lodash';

import {
  AVAILABLE_TIME_UNITS,
  MEDIUM_AVAILABLE_TIME_UNITS,
  SHORT_AVAILABLE_TIME_UNITS,
} from '@/constants';

export default {
  props: {
    duration: {
      type: Object,
      default: () => ({
        value: 1,
        unit: AVAILABLE_TIME_UNITS.minute.value,
      }),
    },
    label: {
      type: String,
      default: null,
    },
    unitsLabel: {
      type: String,
      default: null,
    },
    units: {
      type: Array,
      default: null,
    },
    name: {
      type: String,
      default: 'duration',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    clearable: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      default: 1,
    },
    max: {
      type: Number,
      required: false,
    },
    long: {
      type: Boolean,
      default: false,
    },
    medium: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    autofocus: {
      type: Boolean,
      default: false,
    },
    sameWidth: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:duration'],

  computed: {
    valueFieldName() {
      return `${this.name}.value`;
    },

    unitFieldName() {
      return `${this.name}.unit`;
    },

    availableUnits() {
      if (this.units) {
        return this.units;
      }

      const units = {
        true: SHORT_AVAILABLE_TIME_UNITS,
        [this.long]: AVAILABLE_TIME_UNITS,
        [this.medium]: MEDIUM_AVAILABLE_TIME_UNITS,
      }.true;

      return Object.values(units).map(({ value, text }) => ({
        value,
        text: this.$t(text, this.duration.value || 0),
      }));
    },

    isRequired() {
      return this.required || isNumber(this.duration.value) || Boolean(this.duration.unit);
    },
  },

  methods: {
    updateModel(duration) {
      this.$emit('update:duration', duration);
    },

    updateField(field, value) {
      this.updateModel({
        ...this.duration,
        [field]: value,
      });
    },

    updateValue(value) {
      this.updateField('value', value);
    },

    updateUnit(unit) {
      if (unit) {
        this.updateField('unit', unit);

        return;
      }

      this.updateModel({
        ...this.duration,
        unit: undefined,
        value: undefined,
      });
    },
  },
};
</script>
