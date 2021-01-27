<template lang="pug">
  v-select(
    v-validate="'required'",
    v-field="value",
    :items="availableTriggers",
    :disabled="disabled",
    :label="label || $t('scenarios.fields.triggers')",
    :error-messages="errors.collect('triggers')",
    name="triggers",
    multiple,
    chips
  )
</template>

<script>
import { SCENARIO_TRIGGERS } from '@/constants';

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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    availableTriggers() {
      return Object.values(SCENARIO_TRIGGERS);
    },
  },
};
</script>
