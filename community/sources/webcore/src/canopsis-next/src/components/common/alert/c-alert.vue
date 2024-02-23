<template>
  <v-alert
    v-bind="$attrs"
    :value="value"
    class="c-alert"
  >
    <slot />
  </v-alert>
</template>

<script>
export default {
  inheritAttrs: false,
  props: {
    value: {
      type: Boolean,
      default: true,
    },
  },
};
</script>

<style>
.c-alert {
  width: 100%;
}
</style>
