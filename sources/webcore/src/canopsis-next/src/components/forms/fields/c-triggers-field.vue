<template lang="pug">
  v-select(
    v-validate="'required'",
    v-field="value",
    :items="availableTriggers",
    :disabled="disabled",
    :label="label || $t('common.triggers')",
    :error-messages="errors.collect(name)",
    :name="name",
    multiple,
    chips
  )
</template>

<script>
import { SCENARIO_ACTION_TYPES } from '@/constants';

export default {
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Array,
      required: true,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'triggers',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableTriggers() {
      return Object.values(SCENARIO_ACTION_TYPES);
    },
  },
};
</script>
