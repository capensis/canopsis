<template lang="pug" functional>
  v-layout(align-center, justify-end)
    v-icon.handler.draggable(:class="props.dragClass", v-show="props.draggable") drag_indicator
    v-avatar.white--text(color="primary", size="32")
      slot
</template>

<script>
export default {
  props: {
    draggable: {
      type: Boolean,
      default: false,
    },
    dragClass: {
      type: String,
      default: '',
    },
  },
};
</script>

<style lang="scss">
  .handler {
    width: 20px !important;
  }
</style>
