<template>
  <c-layout align="center" justify="space-around" v-bind="$attrs">
    <v-icon :color="color" size="18">
      {{ icon }}
    </v-icon>
    <span class="text-caption text-white">
      {{ count }}
    </span>
  </c-layout>
</template>

<script>
export default {
  inheritAttrs: false,

  props: {
    icon: {
      type: String,
      required: true,
    },
    count: {
      type: Number,
      default: 0,
    },
    color: {
      type: String,
      default: 'white',
    },
  },
};
</script>
