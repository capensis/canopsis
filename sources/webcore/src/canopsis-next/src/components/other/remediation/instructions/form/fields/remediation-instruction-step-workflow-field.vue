<template lang="pug">
  v-radio-group.mt-0(
    v-field="value",
    :label="$t('remediationInstructions.workflow')",
    hide-details,
    mandatory,
    color="primary"
  )
    v-radio(
      :label="$t('remediationInstructions.stop')",
      :value="$constants.REMEDIATION_WORKFLOW_TYPES.stop",
      color="primary"
    )
    v-radio(
      :label="$t('remediationInstructions.remainingStep')",
      :value="$constants.REMEDIATION_WORKFLOW_TYPES.continue",
      color="primary"
    )
</template>

<script>
import { REMEDIATION_WORKFLOW_TYPES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Boolean,
      default: REMEDIATION_WORKFLOW_TYPES.stop,
    },
  },
};
</script>
