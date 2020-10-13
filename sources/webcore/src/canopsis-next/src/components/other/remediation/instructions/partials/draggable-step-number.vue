<template lang="pug">
  v-layout(align-center, justify-space-around)
    v-icon.handler.draggable(:class="dragClass") drag_indicator
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
  },
};
</script>

<style lang="scss">
  .handler {
    width: 18px !important;
  }
</style>
