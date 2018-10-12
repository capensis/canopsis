<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
    v-container
      v-switch(:input-value="value", @change="$emit('input', $event)")
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      default: '',
    },
    value: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
