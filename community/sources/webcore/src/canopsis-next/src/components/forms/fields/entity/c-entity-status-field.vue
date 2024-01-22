<template>
  <v-select
    v-field="value"
    :items="alarmStates"
    :disabled="disabled"
    :label="label || $t('common.status')"
    :name="name"
    hide-details
  />
</template>

<script>
import { ENTITIES_STATUSES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Number],
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'status',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    alarmStates() {
      return Object.values(ENTITIES_STATUSES).map(value => ({
        value,
        text: this.$t(`common.statusTypes.${value}`),
      }));
    },
  },
};
</script>
