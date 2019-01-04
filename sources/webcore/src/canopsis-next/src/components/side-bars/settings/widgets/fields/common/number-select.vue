<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{ title }}
    v-container
      v-text-field(type="number", :value="value", @input="$emit('input', Number($event))")
</template>

<script>
export default {
  props: {
    value: {
      type: Number,
      required: true,
    },
    title: {
      type: String,
      required: true,
    },
  },
};
</script>
