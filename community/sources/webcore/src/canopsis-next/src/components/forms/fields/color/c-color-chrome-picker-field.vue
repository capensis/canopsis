<template lang="pug">
  chrome(v-field="value", v-bind="$attrs")
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
      type: String,
      required: false,
    },
  },
};
</script>
