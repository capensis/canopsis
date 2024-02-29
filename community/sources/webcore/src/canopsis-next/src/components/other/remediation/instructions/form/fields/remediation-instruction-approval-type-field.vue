<template>
  <v-radio-group
    v-field="type"
    class="mt-0"
    hide-details
    mandatory
    row
  >
    <v-radio
      :label="$tc('common.role')"
      :value="$constants.REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role"
      color="primary"
    />
    <v-radio
      :label="$tc('common.user')"
      :value="$constants.REMEDIATION_INSTRUCTION_APPROVAL_TYPES.user"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
import { REMEDIATION_INSTRUCTION_APPROVAL_TYPES } from '@/constants';

export default {
  model: {
    prop: 'type',
    event: 'input',
  },
  props: {
    type: {
      type: Number,
      default: REMEDIATION_INSTRUCTION_APPROVAL_TYPES.role,
    },
  },
};
</script>
