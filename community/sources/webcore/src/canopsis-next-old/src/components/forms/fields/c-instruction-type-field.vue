<template>
  <v-radio-group
    v-field="value"
    :label="$t('remediation.instruction.type')"
    :disabled="disabled"
    class="mt-0"
    hide-details
    mandatory
  >
    <v-radio
      v-for="type in types"
      :key="type.value"
      :label="type.label"
      :value="type.value"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
import { REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Number,
      default: REMEDIATION_INSTRUCTION_TYPES.manual,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    types() {
      return Object.values(REMEDIATION_INSTRUCTION_TYPES).map(value => ({
        value,
        label: this.$t(`remediation.instruction.types.${value}`),
      }));
    },
  },
};
</script>
