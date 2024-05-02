<template>
  <v-layout>
    <v-flex
      class="pr-3"
      xs7
    >
      <c-number-field
        v-field="duration.value"
        :label="label || $t('common.duration')"
        :disabled="disabled"
        :name="valueFieldName"
        :required="isRequired"
        :min="min"
        :max="max"
        :hide-details="hideDetails"
      />
    </v-flex>
    <v-flex xs5>
      <v-select
        v-validate="unitValidateRules"
        :value="duration.unit"
        :items="availableUnits"
        :disabled="disabled"
        :label="unitsLabel"
        :error-messages="errors.collect(unitFieldName)"
        :name="unitFieldName"
        :clearable="clearable"
        :hide-details="hideDetails"
        @change="updateUnit"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { isNumber } from 'lodash';

import { AVAILABLE_TIME_UNITS, SHORT_AVAILABLE_TIME_UNITS } from '@/constants';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'duration',
    event: 'input',
  },
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
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
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

      const units = this.long
        ? AVAILABLE_TIME_UNITS
        : SHORT_AVAILABLE_TIME_UNITS;

      return Object.values(units).map(({ value, text }) => ({
        value,
        text: this.$tc(text, this.duration.value || 0),
      }));
    },

    isRequired() {
      return this.required || isNumber(this.duration.value) || Boolean(this.duration.unit);
    },

    unitValidateRules() {
      return { required: this.isRequired };
    },
  },
  methods: {
    updateUnit(unit) {
      if (unit) {
        this.updateField('unit', unit);
      } else {
        this.updateModel({
          ...this.duration,
          unit: undefined,
          value: undefined,
        });
      }
    },
  },
};
</script>
