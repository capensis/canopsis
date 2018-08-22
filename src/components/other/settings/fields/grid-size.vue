<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") Columns - Large
    v-container
      v-slider(
      :value="value",
      max=12,
      always-dirty,
      thumb-label,
      @input="$emit('update:value', $event)"
      )
</template>

<script>
export default {
  props: {
    value: {
      type: Number,
    },
  },
};
</script>

