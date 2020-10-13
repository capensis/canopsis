<template lang="pug">
  v-layout(align-center, justify-space-around)
    span.handler
      v-icon.draggable(v-show="!disabled", :class="dragClass") drag_indicator
    v-avatar.white--text(color="primary", size="32")
      slot
</template>

<script>
export default {
  props: {
    dragClass: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss">
  .handler {
    width: 18px !important;
  }
</style>
