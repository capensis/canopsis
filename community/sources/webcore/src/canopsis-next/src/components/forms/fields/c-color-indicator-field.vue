<template lang="pug">
  v-radio-group(v-field="value")
    v-radio(
      :value="$constants.COLOR_INDICATOR_TYPES.state",
      :label="$t('settings.colorIndicator.fields.displayAsSeverity')",
      color="primary"
    )
    v-radio(
      :value="$constants.COLOR_INDICATOR_TYPES.impactState",
      :label="$t('settings.colorIndicator.fields.displayAsPriority')",
      color="primary"
    )
</template>

<script>
import { COLOR_INDICATOR_TYPES } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      default: COLOR_INDICATOR_TYPES.state,
    },
  },
};
</script>
