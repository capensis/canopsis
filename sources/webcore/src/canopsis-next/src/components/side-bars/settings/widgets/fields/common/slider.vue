<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
    v-container
      v-slider(
      :value="value",
      :max="max",
      :min="min",
      ticks="always"
      @input="$emit('input', $event)"
      always-dirty,
      thumb-label,
      )
</template>

<script>
export default {
  props: {
    title: {
      type: String,
      required: true,
    },
    value: {
      type: Number,
      required: true,
    },
    min: {
      type: Number,
      default: 0,
    },
    max: {
      type: Number,
      default: 5,
    },
  },
};
</script>

