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
      :value="$constants.WORKFLOW_TYPES.stop",
      color="primary"
    )
    v-radio(
      :label="$t('remediationInstructions.remainingStep')",
      :value="$constants.WORKFLOW_TYPES.continue",
      color="primary"
    )
</template>

<script>
import { WORKFLOW_TYPES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: WORKFLOW_TYPES.stop,
    },
  },
};
</script>
