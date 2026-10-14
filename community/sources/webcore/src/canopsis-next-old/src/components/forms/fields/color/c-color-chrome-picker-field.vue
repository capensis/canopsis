<template>
  <chrome
    v-field="value"
    v-bind="$attrs"
  />
</template>

<script>
import { Chrome } from 'vue-color';

export default {
  components: {
    Chrome,
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
