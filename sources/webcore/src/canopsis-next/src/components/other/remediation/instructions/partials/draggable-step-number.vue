<template lang="pug">
  v-layout(align-center, justify-space-around)
    v-icon.handler.draggable(:class="dragClass") drag_indicator
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
  },
};
</script>

<style lang="scss">
  .handler {
    width: 18px !important;
  }
</style>
