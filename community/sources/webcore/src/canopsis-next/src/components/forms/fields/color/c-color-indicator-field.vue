<template>
  <v-radio-group
    v-field="value"
    :name="name"
    :disabled="disabled"
  >
    <v-radio
      :value="$constants.COLOR_INDICATOR_TYPES.state"
      :label="$t('settings.colorIndicator.fields.displayAsSeverity')"
      color="primary"
    />
    <v-radio
      :value="$constants.COLOR_INDICATOR_TYPES.impactState"
      :label="$t('settings.colorIndicator.fields.displayAsPriority')"
      color="primary"
    />
  </v-radio-group>
</template>

<script>
import { COLOR_INDICATOR_TYPES } from '@/constants';

export default {
  props: {
    value: {
      type: String,
      default: COLOR_INDICATOR_TYPES.state,
    },
    name: {
      type: String,
      default: 'colorIndicator',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
