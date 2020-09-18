<template lang="pug">
  v-btn.ma-0(icon, small, @click="$emit('expand')")
    v-icon {{ expanded ? 'keyboard_arrow_up' : 'keyboard_arrow_down' }}
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'expand',
  },
  props: {
    expanded: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
