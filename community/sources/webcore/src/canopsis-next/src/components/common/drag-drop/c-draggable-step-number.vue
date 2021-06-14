<template lang="pug">
  v-layout(align-center, justify-space-around)
    span.handler
      v-icon.draggable(v-show="!disabled", :class="dragClass") drag_indicator
    v-avatar.white--text(:color="color", size="32")
      slot
</template>

<script>
export default {
  props: {
    dragClass: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'primary',
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
    height: 24px !important;
    width: 18px !important;
  }
</style>
