<template>
  <compact
    v-field="value"
    v-bind="$attrs"
  />
</template>

<script>
import { Compact } from 'vue-color';

export default {
  components: {
    Compact,
  },
  inheritAttrs: false,
  props: {
    value: {
      type: [Object, String],
      required: false,
    },
  },
};
</script>
