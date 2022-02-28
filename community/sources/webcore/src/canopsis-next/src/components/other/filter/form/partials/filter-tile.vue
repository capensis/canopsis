<template lang="pug">
  v-list-tile.pa-0(:class="{ 'cursor-move': editable }")
    v-layout
      v-flex(xs12)
        v-list-tile-content {{ filter.title }}
      v-list-tile-action(v-if="editable")
        v-layout
          c-action-btn(type="edit", @click="$emit('edit')")
          c-action-btn(type="delete", @click="$emit('delete')")
</template>

<script>
export default {
  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },
    editable: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
