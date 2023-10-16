<template>
  <v-radio-group
    class="mt-0"
    v-field="value"
    :label="label"
    :disabled="disabled"
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
import { IDLE_RULE_ALARM_CONDITIONS } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: IDLE_RULE_ALARM_CONDITIONS.lastEvent,
    },
    label: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    types() {
      return Object.values(IDLE_RULE_ALARM_CONDITIONS).map(value => ({
        value,
        label: this.$t(`idleRules.alarmConditions.${value}`),
      }));
    },
  },
};
</script>
