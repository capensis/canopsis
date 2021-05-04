<template lang="pug">
  v-radio-group.mt-0(
    v-field="value",
    :label="$t('remediationInstructions.type')",
    hide-details,
    mandatory,
    color="primary"
  )
    v-radio(
      :label="$t('remediationInstructions.types.manual')",
      :value="$constants.REMEDIATION_INSTRUCTION_TYPES.manual",
      color="primary"
    )
    v-radio(
      :label="$t('remediationInstructions.types.automatic')",
      :value="$constants.REMEDIATION_INSTRUCTION_TYPES.automatic",
      color="primary"
    )
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
  },
};
</script>
