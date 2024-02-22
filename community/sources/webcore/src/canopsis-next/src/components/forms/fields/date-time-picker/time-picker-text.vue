<template>
  <v-text-field
    :value="value"
    :label="label"
    class="time-picker__select"
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
