<template lang="pug">
  v-list-tile.pa-0(:class="{ 'filter-item': hasAccessToEdit }")
    v-layout
      v-flex(xs12)
        v-list-tile-content {{ filter.title }}
      v-list-tile-action(v-if="hasAccessToEdit")
        v-layout
          c-action-btn(type="edit", @click="$emit('edit')")
          c-action-btn(type="delete", @click="$emit('delete')")
</template>

<script>
export default {
  model: {
    prop: 'filter',
    event: 'input',
  },
  props: {
    filter: {
      type: Object,
      default: () => ({}),
    },
    hasAccessToEdit: {
      type: Boolean,
      default: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.filter-item {
  cursor: move;
}
</style>
