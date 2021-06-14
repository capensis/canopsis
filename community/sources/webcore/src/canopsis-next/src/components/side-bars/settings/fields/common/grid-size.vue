<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
    v-container
      v-slider(
        :value="value",
        :max="max",
        :min="min",
        ticks="always",
        @input="$emit('input', $event)",
        always-dirty,
        thumb-label
      )
</template>

<script>
export default {
  props: {
    value: {
      type: Number,
    },
    title: {
      type: String,
      required: true,
    },
    max: {
      type: Number,
      default: 12,
    },
    min: {
      type: Number,
      default: 0,
    },
  },
};
</script>

