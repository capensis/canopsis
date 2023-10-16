<template>
  <v-text-field
    class="time-picker__select"
    :value="value"
    :label="label"
    placeholder="−−:−−"
    hide-details
    disabled
  />
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: null,
    },
    label: {
      type: String,
      default: null,
    },
  },
};
</script>
