<template>
  <c-select
    :model-value="value"
    :items="alarmStatuses"
    :disabled="disabled"
    :label="label || $t('common.status')"
    :name="name"
    item-title="text"
    item-value="value"
    hide-details
    @update:model-value="updateValue"
  />
</template>

<script>
import { ALARM_STATUSES } from '@/constants';

export default {
  emits: ['update:value', 'update:modelValue'],
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
    alarmStatuses() {
      return Object.values(ALARM_STATUSES).map(value => ({
        value,
        text: this.$t(`common.statusTypes.${value}`),
      }));
    },
  },
  methods: {
    updateValue(value) {
      this.$emit('update:value', value);
      this.$emit('update:modelValue', value);
    },
  },
};
</script>
