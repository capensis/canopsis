<template lang="pug">
  v-radio-group.mt-0(
    v-field="value",
    :label="label",
    hide-details,
    mandatory,
    color="primary"
  )
    v-radio(
      :label="stopLabel || $t('common.stop')",
      :value="stopValue",
      color="primary"
    )
    v-radio(
      :label="continueLabel || $t('common.continue')",
      :value="continueValue",
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
      type: Boolean,
      default: WORKFLOW_TYPES.stop,
    },
    label: {
      type: String,
      required: true,
    },
    stopLabel: {
      type: String,
      required: false,
    },
    continueLabel: {
      type: String,
      required: false,
    },
    stopValue: {
      type: Boolean,
      default: WORKFLOW_TYPES.stop,
    },
    continueValue: {
      type: Boolean,
      default: WORKFLOW_TYPES.continue,
    },
  },
};
</script>
