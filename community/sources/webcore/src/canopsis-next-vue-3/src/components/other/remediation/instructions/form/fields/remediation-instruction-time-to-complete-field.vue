<template>
  <v-layout>
    <v-flex xs8>
      <c-number-field
        :value="duration.value"
        :label="$t('remediation.instruction.timeToComplete')"
        :name="durationFieldName"
        :disabled="disabled"
        :min="min"
        required
        filled
        @update:value="updateDurationField('value', $event)"
      />
    </v-flex>
    <v-flex class="pl-3" xs4>
      <v-select
        v-validate="'required'"
        :model-value="duration.unit"
        :label="$t('common.unit')"
        :items="availableUnits"
        :name="unitFieldName"
        :disabled="disabled"
        class="time-complete-unit"
        item-title="text"
        item-value="value"
        hide-details
        @update:model-value="updateDurationField('unit', $event)"
      />
    </v-flex>
  </v-layout>
</template>

<script>
import { SHORT_AVAILABLE_TIME_UNITS } from '@/constants';

export default {
  inject: ['$validator'],
  props: {
    duration: {
      type: Object,
      default: () => ({}),
    },
    name: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      default: 1,
    },
  },
  emits: ['update:duration'],
  computed: {
    durationFieldName() {
      return `${this.name}.duration`;
    },

    unitFieldName() {
      return `${this.name}.unit`;
    },

    availableUnits() {
      return Object.values(SHORT_AVAILABLE_TIME_UNITS).map(({ value, text }) => ({
        value,
        text: this.$t(text, this.duration.value),
      }));
    },
  },
  methods: {
    updateDurationField(field, value) {
      this.$emit('update:duration', {
        ...this.duration,
        [field]: value,
      });
    },
  },
};
</script>

<style lang="scss">
.time-complete-unit .v-input__slot {
  &:before,
  &:after {
    content: none !important;
  }
}
</style>
