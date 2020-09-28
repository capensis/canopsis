<template lang="pug" functional>
  v-layout(align-center, justify-end)
    v-icon.mr-1.draggable(v-if="props.draggable") drag_indicator
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
  },
};
</script>
