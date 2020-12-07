<template lang="pug">
  v-list-tile.filter-item.pa-0
    v-layout(:data-test="`filterItem-${filter.title}`")
      v-flex(xs12)
        v-list-tile-content {{ filter.title }}
      v-list-tile-action(v-if="hasAccessToEdit")
        v-layout
          v-btn.ma-1(
            icon,
            :data-test="`editFilter-${filter.title}`",
            @click="$emit('edit')"
          )
            v-icon edit
          v-btn.ma-1(
            icon,
            :data-test="`deleteFilter-${filter.title}`",
            @click="$emit('delete')"
          )
            v-icon delete
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
