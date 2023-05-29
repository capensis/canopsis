<template lang="pug">
  span.alarm-header-actions
    slot
    v-btn(small, icon, @click="$emit('resizing')")
      v-icon(small) {{ resizing ? 'lock_open' : 'lock_outline' }}
</template>

<script>
export default {
  props: {
    resizing: {
      type: Boolean,
      default: false,
    },
  },
};
</script>

<style lang="scss" scoped>
.alarm-header-actions {
  position: relative;
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
